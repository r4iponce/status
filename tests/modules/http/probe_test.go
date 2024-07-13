package http

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"go.ada.wf/status/internal/constant"
	local_http "go.ada.wf/status/internal/modules/http"
	"go.ada.wf/status/tests/utils"
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

func TestVerify(t *testing.T) { //nolint:funlen
	t.Parallel()

	t.Run("UP basic", func(t *testing.T) {
		t.Parallel()

		port := utils.GetRandomPort()
		listen := fmt.Sprintf("localhost:%d", port)

		p := utils.Handler{
			Body: "OK",
			Code: http.StatusOK,
		}

		prepare(listen, p)

		time.Sleep(1 * time.Second)
		c := local_http.Config{
			Target:   fmt.Sprintf("http://localhost:%d", port),
			Valid:    []int{200},
			Response: "OK",
		}

		err := c.IsUp()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("UP multiple code", func(t *testing.T) {
		t.Parallel()

		port := utils.GetRandomPort()
		listen := fmt.Sprintf("localhost:%d", port)

		p := utils.Handler{
			Body: "OK",
			Code: http.StatusNotFound,
		}

		prepare(listen, p)

		time.Sleep(1 * time.Second)
		c := local_http.Config{
			Target:   fmt.Sprintf("http://localhost:%d", port),
			Valid:    []int{200, 404},
			Response: "",
		}

		err := c.IsUp()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Down bad code", func(t *testing.T) {
		t.Parallel()

		port := utils.GetRandomPort()
		listen := fmt.Sprintf("localhost:%d", port)

		p := utils.Handler{
			Body: "OK",
			Code: http.StatusBadGateway,
		}

		prepare(listen, p)

		time.Sleep(1 * time.Second)
		c := local_http.Config{
			Target:   fmt.Sprintf("http://localhost:%d", port),
			Valid:    []int{200},
			Response: "",
		}

		got := c.IsUp()
		want := &local_http.BadStatusCodeError{}
		if !errors.As(got, &want) {
			t.Fatalf("want %v, got %v", want, got)
		}
	})

	t.Run("Down bad body", func(t *testing.T) {
		t.Parallel()

		port := utils.GetRandomPort()
		listen := fmt.Sprintf("localhost:%d", port)

		p := utils.Handler{
			Body: "OK",
			Code: http.StatusOK,
		}

		prepare(listen, p)

		time.Sleep(1 * time.Second)
		c := local_http.Config{
			Target:   fmt.Sprintf("http://localhost:%d", port),
			Valid:    []int{200},
			Response: "KO",
		}

		got := c.IsUp()
		want := &local_http.BadBodyError{}
		if !errors.As(got, &want) {
			t.Fatalf("want %v, got %v", want, got)
		}
	})
}
