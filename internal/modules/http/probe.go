package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type BadStatusCodeError struct {
	Expected []int
	Got      int
}

func (m *BadStatusCodeError) Error() string {
	return fmt.Sprintf("bad status code, Expected: %d Got: %d", m.Expected, m.Got)
}

type BadBodyError struct {
	Expected string
	Got      string
}

func (m *BadBodyError) Error() string {
	return fmt.Sprintf("bad responde body, Expected: %s Got: %s", m.Expected, m.Got)
}

func (c Config) IsUp() error {
	logrus.Debugf("the check for %s has begin", c.Target)

	r, err := http.Get(c.Target)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(r.Body)

	valid := false

	for _, v := range c.Valid {
		if r.StatusCode == v {
			valid = true
		}
	}

	if !valid {
		return &BadStatusCodeError{
			Expected: c.Valid,
			Got:      r.StatusCode,
		}
	}

	if c.Response != "" {
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), c.Response) {
			return &BadBodyError{
				Expected: c.Response,
				Got:      string(body),
			}
		}
	}

	logrus.Debugf("the check for %s is completed with no error", c.Target)

	return nil
}
