package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type BadStatusCodeError struct {
	expected []int
	got      int
}

func (m *BadStatusCodeError) Error() string {
	return fmt.Sprintf("bad status code, expected: %d got: %d", m.expected, m.got)
}

type BadBodyError struct {
	expected string
	got      string
}

func (m *BadBodyError) Error() string {
	return fmt.Sprintf("bad responde body, expected: %s got: %s", m.expected, m.got)
}

func (c Config) IsUp() error {
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
			expected: c.Valid,
			got:      r.StatusCode,
		}
	}

	if c.Response != "" {
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), c.Response) {
			return &BadBodyError{
				expected: c.Response,
				got:      string(body),
			}
		}
	}

	logrus.Debugf("the check for %s is completed with no error", c.Target)

	return nil
}
