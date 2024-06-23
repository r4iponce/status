package config

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml/v2"
)

var (
	errLogLevel              = errors.New("log level is invalid, valid list is : \"DEBUG\", \"INFO\", \"WARN\", \"ERROR\", \"FATAL\"") //nolint: lll
	errConfigFileNotReadable = errors.New("config file is not loadable")
)

func LoadToml(file string) (Config, error) {
	var c Config

	source, err := os.ReadFile(file)
	if err != nil {
		return c, errConfigFileNotReadable
	}

	err = toml.Unmarshal(source, &c)
	if err != nil {
		panic(err)
	}

	return c, nil
}
