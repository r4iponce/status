package probe

import (
	"gitlab.gnous.eu/ada/status/internal/modules/http"
)

var config Config

type Config struct {
	Cache   bool
	Targets []Target
}

type Target struct {
	Name        string
	Description string
	Module      string
	Http        http.Config
}

func Init(cache bool, targets []Target) {
	config = Config{
		Cache:   cache,
		Targets: targets,
	}
}
