package probe

import (
	"github.com/sirupsen/logrus"
	"gitlab.gnous.eu/ada/status/internal/cache"
	"gitlab.gnous.eu/ada/status/internal/models"
)

func RunAll() []models.Status {
	var statuses []models.Status

	for _, v := range config.Targets {
		switch v.Module {
		case "http":
			s := runHttp(config, v)
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
		if cache.KeyExist(t.Name) {
			statuses, err = cache.GetCacheResult(t.Name)
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
		cache.SetCacheResult(statuses)
	}

	return statuses
}
