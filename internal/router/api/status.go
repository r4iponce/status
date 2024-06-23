package api

import (
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gitlab.gnous.eu/ada/status/internal/probe"
)

type Config struct {
	Targets []probe.Target
	Db      *redis.Client
	Cache   bool
}

func (c Config) GetAll(w http.ResponseWriter, _ *http.Request) {
	probeConfig := probe.Config{
		Db:      c.Db,
		Cache:   c.Cache,
		Targets: c.Targets,
	}

	s := probe.RunAll(probeConfig)

	response, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		logrus.Error(err)

		return
	}
}
