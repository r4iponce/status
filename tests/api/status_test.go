package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.gnous.eu/ada/status/internal/cache"
	"gitlab.gnous.eu/ada/status/internal/constant"
	"gitlab.gnous.eu/ada/status/internal/models"
	local_http "gitlab.gnous.eu/ada/status/internal/modules/http"
	"gitlab.gnous.eu/ada/status/internal/probe"
	"gitlab.gnous.eu/ada/status/tests/utils"
)

func prepare(listen string, handler utils.Handler) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handler.BasicHandler)
	httpServer := &http.Server{
		Addr:              listen,
		Handler:           mux,
		ReadHeaderTimeout: constant.ReadHeaderTimeout,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			logrus.Error(err)

			return
		}
	}()
}

func TestStatusApi(t *testing.T) {
	t.Parallel()

	t.Run("Basic up target", func(t *testing.T) {
		t.Parallel()

		portTest1 := utils.GetRandomPort()

		cacheConfig := cache.Config{
			Enabled: false,
		}

		targets := []probe.Target{{
			Name:        "test1",
			Description: "This is test one",
			Module:      "http",
			Http: local_http.Config{
				Target:   fmt.Sprintf("http://localhost:%d", portTest1),
				Valid:    []int{200},
				Response: "OK",
			},
		}}

		url := fmt.Sprintf("http://%s/api/status", utils.Prepare(targets, cacheConfig))

		prepare(fmt.Sprintf("localhost:%d", portTest1), utils.Handler{
			Body: "OK",
			Code: 200,
		})

		time.Sleep(1 * time.Second)

		r, err := http.Get(url)
		if err != nil {
			t.Fatalf("Hi, %v", err)
		}

		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var got []models.Status

		err = json.Unmarshal(body, &got)
		if err != nil {
			t.Fatal(err)
		}

		want := models.Status{
			Name:        targets[0].Name,
			Description: targets[0].Description,
			Success:     true,
			Error:       "",
			Target:      targets[0].Http.Target,
		}

		if got[0] != want {
			t.Fatalf("want %v, got %v", want, got[0])
		}
	})

	t.Run("Basic down bad status code", func(t *testing.T) {
		t.Parallel()

		portTest2 := utils.GetRandomPort()

		cacheConfig := cache.Config{
			Enabled: false,
		}

		targets := []probe.Target{{
			Name:        "test2",
			Description: "This is test two",
			Module:      "http",
			Http: local_http.Config{
				Target:   fmt.Sprintf("http://localhost:%d", portTest2),
				Valid:    []int{200},
				Response: "",
			},
		}}

		url := fmt.Sprintf("http://%s/api/status", utils.Prepare(targets, cacheConfig))

		prepare(fmt.Sprintf("localhost:%d", portTest2), utils.Handler{
			Body: "OK",
			Code: 503,
		})

		time.Sleep(1 * time.Second)

		r, err := http.Get(url)
		if err != nil {
			t.Fatalf("Hi, %v", err)
		}

		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var got []models.Status

		err = json.Unmarshal(body, &got)
		if err != nil {
			t.Fatal(err)
		}

		wantError := &local_http.BadStatusCodeError{Expected: targets[0].Http.Valid, Got: 503}

		want := models.Status{
			Name:        targets[0].Name,
			Description: targets[0].Description,
			Success:     false,
			Error:       wantError.Error(),
			Target:      targets[0].Http.Target,
		}

		if got[0] != want {
			t.Fatalf("want %v, got %v", want, got[0])
		}
	})

	t.Run("Basic down bad body", func(t *testing.T) {
		t.Parallel()

		portTest3 := utils.GetRandomPort()

		cacheConfig := cache.Config{
			Enabled: false,
		}

		targets := []probe.Target{{
			Name:        "test3",
			Description: "This is test three",
			Module:      "http",
			Http: local_http.Config{
				Target:   fmt.Sprintf("http://localhost:%d", portTest3),
				Valid:    []int{200},
				Response: "OK",
			},
		}}

		url := fmt.Sprintf("http://%s/api/status", utils.Prepare(targets, cacheConfig))

		prepare(fmt.Sprintf("localhost:%d", portTest3), utils.Handler{
			Body: "KO",
			Code: 200,
		})

		time.Sleep(1 * time.Second)

		r, err := http.Get(url)
		if err != nil {
			t.Fatalf("Hi, %v", err)
		}

		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var got []models.Status

		err = json.Unmarshal(body, &got)
		if err != nil {
			t.Fatal(err)
		}

		wantError := &local_http.BadBodyError{Expected: "OK", Got: "KO"}

		want := models.Status{
			Name:        targets[0].Name,
			Description: targets[0].Description,
			Success:     false,
			Error:       wantError.Error(),
			Target:      targets[0].Http.Target,
		}

		if got[0] != want {
			t.Fatalf("want %v, got %v", want, got[0])
		}
	})
}
