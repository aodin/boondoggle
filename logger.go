package boondoggle

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type RequestLogger interface {
	Log(*http.Request)
}

type DefaultLogger struct{}

// By default, a timestamp will be written by the logger
func (l *DefaultLogger) Log(r *http.Request) {
	log.Printf(
		`%q %s %q`,
		fmt.Sprintf(`%s %s`, r.Method, r.URL),
		strings.SplitN(r.RemoteAddr, ":", 2)[0],
		r.Header.Get("User-Agent"),
	)
}

var defaultLogger = &DefaultLogger{}
