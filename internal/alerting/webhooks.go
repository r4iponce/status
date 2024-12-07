package alerting

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.gnous.eu/ada/status/internal/config"
)

var errDisabledWebhook = errors.New("webhook is disabled")

type webhook struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func SendNotification(c config.Alerting, message string) error {
	if !c.Enabled {
		return errDisabledWebhook
	}

	data, err := json.Marshal(webhook{
		Username: c.Username,
		Content:  message,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(c.Url, "application/json", bytes.NewBuffer(data)) //nolint:bodyclose
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func InitCheck(targets []config.Target) {
	logrus.Debug("initializing alerting")

	for {
		err := checkStatus(targets)
		if err != nil {
			logrus.Error(err)
		}

		time.Sleep(5 * time.Minute) // TODO make it configurable
	}
}

func checkStatus(targets []config.Target) error { // TODO add context to properly stop task
	for _, t := range targets {
		switch t.Module {
		case "http":
			if t.Webhooks.Enabled {
				logrus.Debugf("Verify if %s is up", t.Name)
				err := t.Http.IsUp()
				logrus.Debug(err)
				if err != nil {
					err := SendNotification(t.Webhooks, t.Name+" is down")
					if err != nil {
						return err
					}
				}
			}
		default:
			logrus.Errorf("Invalid module name: %s", t.Module)
		}
	}

	return nil
}
