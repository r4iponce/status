package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.gnous.eu/ada/status/internal/config"
	"gitlab.gnous.eu/ada/status/internal/router"
)

func main() {
	var configPath string

	switch len(os.Args) {
	case 2: //nolint:mnd
		configPath = os.Args[1]
	case 1:
		configPath = "config.toml"
	default:
		logrus.Fatal("Max 1 argument is valid.")
	}

	c, err := config.LoadToml(configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	err = c.Verify()
	if err != nil {
		logrus.Fatal(err)
	}

	err = c.Log.Init()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Debugf("Loaded config : %v", c)

	router.Init(c)
}
