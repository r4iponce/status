package log

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Next struct {
	Fn func(http.ResponseWriter, *http.Request)
}

func (n Next) HttpLogger(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	n.Fn(w, r)
}
