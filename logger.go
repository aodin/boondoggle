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
	ip := strings.SplitN(r.Header.Get("X-Real-IP"), ":", 2)[0]
	if ip == "" {
		ip = strings.SplitN(r.RemoteAddr, ":", 2)[0]
	}
	log.Printf(
		`%q %s %q %q`,
		fmt.Sprintf(`%s %s`, r.Method, r.URL),
		ip,
		r.Header.Get("User-Agent"),
		r.Header.Get("Referer"),
	)
}

var defaultLogger = &DefaultLogger{}
