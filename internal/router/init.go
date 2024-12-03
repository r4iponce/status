package router

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gitlab.gnous.eu/ada/status/internal/config"
	"gitlab.gnous.eu/ada/status/internal/constant"
	"gitlab.gnous.eu/ada/status/internal/log"
	"gitlab.gnous.eu/ada/status/internal/router/api"
)

func static(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue("file") + ".html"
	file := filepath.Join("build", fileName)

	index, err := os.ReadFile(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "Not found", http.StatusNotFound)

			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error(err)

		return
	}
	_, err = w.Write(index)
	if err != nil {
		logrus.Error(err)
	}
}

func Init(c config.Config) {
	apiConfig := api.Config{
		Targets: c.Probe,
		Cache:   c.Redis.Enabled,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/status", log.Next{Fn: apiConfig.GetAll}.HttpLogger)
	mux.HandleFunc("GET /{file}", static)
	mux.Handle("GET /", http.FileServer(http.Dir("build")))

	server := &http.Server{
		Addr:              c.Listen,
		Handler:           mux,
		ReadHeaderTimeout: constant.ReadHeaderTimeout,
	}

	logrus.Infof("Listen on %s", c.Listen)

	logrus.Fatal(server.ListenAndServe())
}
