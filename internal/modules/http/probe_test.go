package http

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.gnous.eu/ada/status/internal/constant"
)

type prepare struct {
	body   string
	code   int
	listen string
}

func getRandomPort() int {
	port := rand.IntN(65535-1024) + 1024
	for isPortAvailable(port) == false {
		port = rand.IntN(65535-1024) + 1024
	}

	return port
}

func isPortAvailable(port int) bool {
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return false
	}

	err = l.Close()
	if err != nil {
		logrus.Error(err)
	}

	return true
}

func (c prepare) basicHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(c.code)
	_, err := w.Write([]byte(c.body))
	if err != nil {
		logrus.Error(err)
	}
}

func (c prepare) prepare() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", c.basicHandler)
	httpServer := &http.Server{
		Addr:              c.listen,
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

		port := getRandomPort()

		p := prepare{
			body:   "OK",
			code:   http.StatusOK,
			listen: fmt.Sprintf("localhost:%d", port),
		}

		p.prepare()

		time.Sleep(1 * time.Second)
		c := Config{
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

		port := getRandomPort()

		p := prepare{
			body:   "OK",
			code:   http.StatusNotFound,
			listen: fmt.Sprintf("localhost:%d", port),
		}

		p.prepare()

		time.Sleep(1 * time.Second)
		c := Config{
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

		port := getRandomPort()

		p := prepare{
			body:   "OK",
			code:   http.StatusBadGateway,
			listen: fmt.Sprintf("localhost:%d", port),
		}

		p.prepare()

		time.Sleep(1 * time.Second)
		c := Config{
			Target:   fmt.Sprintf("http://localhost:%d", port),
			Valid:    []int{200},
			Response: "OK",
		}

		got := c.IsUp()
		want := &BadStatusCodeError{}
		if !errors.As(got, &want) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("Down bad body", func(t *testing.T) {
		t.Parallel()

		port := getRandomPort()

		p := prepare{
			body:   "KO",
			code:   http.StatusOK,
			listen: fmt.Sprintf("localhost:%d", port),
		}

		p.prepare()

		time.Sleep(1 * time.Second)
		c := Config{
			Target:   fmt.Sprintf("http://localhost:%d", port),
			Valid:    []int{200},
			Response: "OK",
		}

		got := c.IsUp()
		want := &BadBodyError{}
		if !errors.As(got, &want) {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}
