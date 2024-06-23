package log

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var errCannotOpenLogFile = errors.New("cannot open log file")

func (config Config) Init() error {
	if config.File != "" {
		file, err := os.OpenFile(config.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o640) //nolint:mnd
		if err != nil {
			logrus.Debug(err)

			return errCannotOpenLogFile
		}
		logrus.SetOutput(file)
	}

	switch config.Level {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	default:

		return fmt.Errorf("unknown log level: %s", config.Level) //nolint: err113
	}

	return nil
}
