package config

import (
	"go.ada.wf/status/internal/cache"
	"go.ada.wf/status/internal/log"
	"go.ada.wf/status/internal/probe"
)

type Config struct {
	Log    log.Config
	Listen string
	Probe  []probe.Target
	Redis  cache.Config
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
