package utils

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	Body string
	Code int
}

func (c Handler) BasicHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(c.Code)
	_, err := w.Write([]byte(c.Body))
	if err != nil {
		logrus.Error(err)
	}
}
