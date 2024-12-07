package config

import (
	"gitlab.gnous.eu/ada/status/internal/cache"
	"gitlab.gnous.eu/ada/status/internal/log"
	"gitlab.gnous.eu/ada/status/internal/modules/http"
)

type Config struct {
	Log    log.Config
	Listen string
	Probe  []Target
	Redis  cache.Config
}

type Target struct {
	Name        string
	Description string
	Module      string
	Http        http.Config
	Webhooks    Alerting
}

type Alerting struct {
	Enabled  bool
	Username string
	Url      string
}

func (c Config) Verify() error {
	allowedValue := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	found := false
	for _, v := range allowedValue {
		if v == c.Log.Level {
			found = true
		}
	}

	if !found {
		return errLogLevel
	}

	return nil
}
