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

var (
	errDisabledWebhook = errors.New("webhook is disabled")
	last               = make(map[string]lastAction)
)

type webhook struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

type lastAction struct {
	lastStatus string
	time       time.Time
}

func SendNotification(c config.Alerting, message string) error {
	if !c.Enabled {
		return errDisabledWebhook
	}

	logrus.Debugf("sending a webhook with message : %s", message)

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

func InitCheck(targets []config.Target, interval int) {
	logrus.Debug("initializing alerting")

	for {
		err := checkStatus(targets)
		if err != nil {
			logrus.Error(err)
		}

		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

func checkStatus(targets []config.Target) error { // TODO add context to properly stop task
	for _, t := range targets {
		switch t.Module {
		case "http":
			if t.Webhooks.Enabled {
				err := checkStatusHttp(t)
				if err != nil {
					return err
				}
			}
		default:
			logrus.Errorf("Invalid module name: %s", t.Module)
		}
	}

	return nil
}

func checkStatusHttp(t config.Target) error {
	err := t.Http.IsUp()
	if err != nil {
		v := last[t.Name].lastStatus
		if v == "down" {
			return nil
		}

		logrus.Debugf("%s is down", t.Name)
		err = SendNotification(t.Webhooks, t.Name+" is down")
		if err != nil {
			return err
		}

		last[t.Name] = lastAction{
			lastStatus: "down",
			time:       time.Now(),
		}
		logrus.Debugf("l %v", last)

		return nil
	}

	if last[t.Name].lastStatus == "down" {
		logrus.Debugf("%s is up", t.Name)
		err = SendNotification(t.Webhooks, t.Name+" is up")
		if err != nil {
			return err
		}

		delete(last, t.Name)
	}

	return nil
}
