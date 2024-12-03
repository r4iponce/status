package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.gnous.eu/ada/status/internal/probe"
)

type Config struct {
	Targets []probe.Target
	Cache   bool
}

func (c Config) GetAll(w http.ResponseWriter, _ *http.Request) {
	probe.Init(c.Cache, c.Targets)

	s := probe.RunAll()

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
