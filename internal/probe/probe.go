package probe

import (
	"github.com/sirupsen/logrus"
	"go.ada.wf/status/internal/cache"
	"go.ada.wf/status/internal/models"
)

func RunAll(c Config) []models.Status {
	var statuses []models.Status

	for _, v := range c.Targets {
		switch v.Module {
		case "http":
			s := runHttp(c, v)
			statuses = append(statuses, s)
		default:
			logrus.Errorf("Invalid module name: %s", v.Module)
		}
	}

	return statuses
}

func runHttp(c Config, t Target) models.Status {
	err := t.Http.IsUp()

	var statuses models.Status

	if c.Cache {
		if cache.KeyExist(c.Db, t.Name) {
			statuses, err = cache.GetCacheResult(c.Db, t.Name)
			if err != nil {
				logrus.Error(err)
			}

			return statuses
		}
	}

	if err != nil {
		statuses = models.Status{
			Name:        t.Name,
			Description: t.Description,
			Target:      t.Http.Target,
			Success:     false,
			Error:       err.Error(),
		}
	} else {
		statuses = models.Status{
			Name:        t.Name,
			Description: t.Description,
			Target:      t.Http.Target,
			Success:     true,
		}
	}

	if c.Cache {
		cache.SetCacheResult(c.Db, statuses)
	}

	return statuses
}
