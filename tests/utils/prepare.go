package utils

import (
	"fmt"
	"math/rand/v2"
	"net"

	"github.com/sirupsen/logrus"
	"go.ada.wf/status/internal/cache"
	"go.ada.wf/status/internal/config"
	"go.ada.wf/status/internal/log"
	"go.ada.wf/status/internal/probe"
	"go.ada.wf/status/internal/router"
)

func GetRandomPort() int {
	port := rand.IntN(65535-1024) + 1024
	for isPortAvailable(port) == false {
		port = rand.IntN(65535-1024) + 1024
	}

	return port
}

func isPortAvailable(port int) bool {
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return false
	}

	err = l.Close()
	if err != nil {
		logrus.Error(err)
	}

	return true
}

func Prepare(targets []probe.Target, redis cache.Config) string {
	listen := fmt.Sprintf("localhost:%d", GetRandomPort())

	c := config.Config{
		Log: log.Config{
			Level: "DEBUG",
			File:  "",
		},
		Listen: listen,
		Probe:  targets,
		Redis:  redis,
	}

	err := c.Verify()
	if err != nil {
		logrus.Fatal(err)
	}

	err = c.Log.Init()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Debugf("Loaded config : %v", c)

	go router.Init(c)

	logrus.Debug("Status backend started")

	return listen
}
