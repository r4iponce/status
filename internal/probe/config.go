package probe

import (
	"github.com/redis/go-redis/v9"
	"go.ada.wf/status/internal/modules/http"
)

type Config struct {
	Db      *redis.Client
	Cache   bool
	Targets []Target
}

type Target struct {
	Name        string
	Description string
	Module      string
	Http        http.Config
}
