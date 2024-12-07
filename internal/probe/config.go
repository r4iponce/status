package probe

import (
	"gitlab.gnous.eu/ada/status/internal/config"
)

var c Config

type Config struct {
	Cache   bool
	Targets []config.Target
}

func Init(cache bool, targets []config.Target) {
	c = Config{
		Cache:   cache,
		Targets: targets,
	}
}
