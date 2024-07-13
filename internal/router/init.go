package router

import (
	"errors"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.ada.wf/status/internal/cache"
	"go.ada.wf/status/internal/config"
	"go.ada.wf/status/internal/constant"
	"go.ada.wf/status/internal/log"
	"go.ada.wf/status/internal/router/api"
)

func static(w http.ResponseWriter, r *http.Request) {
	file := "build/" + r.PathValue("file") + ".html"

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
	var db *redis.Client

	if c.Redis.Enabled {
		db = c.Redis.Connect()
		err := cache.Ping(db)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	apiConfig := api.Config{
		Targets: c.Probe,
		Db:      db,
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
