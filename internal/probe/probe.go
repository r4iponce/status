package probe

import (
	"github.com/sirupsen/logrus"
	"gitlab.gnous.eu/ada/status/internal/cache"
	"gitlab.gnous.eu/ada/status/internal/config"
	"gitlab.gnous.eu/ada/status/internal/models"
)

func RunAll() []models.Status {
	var statuses []models.Status

	for _, v := range c.Targets {
		switch v.Module {
		case "http":
			s := RunHttp(c.Cache, v)
			statuses = append(statuses, s)
		default:
			logrus.Errorf("Invalid module name: %s", v.Module)
		}
	}

	return statuses
}

func RunHttp(cacheEnabled bool, t config.Target) models.Status {
	var statuses models.Status

	if cacheEnabled && cache.KeyExist(t.Name) {
		statuses, err := cache.GetCacheResult(t.Name)
		if err != nil {
			logrus.Error(err)
		}

		return statuses
	}

	err := t.Http.IsUp()

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

	if cacheEnabled {
		cache.SetCacheResult(statuses)
	}

	return statuses
}
