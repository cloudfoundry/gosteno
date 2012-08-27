package steno

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.Method, "GET") {
		io.WriteString(w, loggersInJson())
		return
	}

	if strings.EqualFold(r.Method, "PUT") {
		level, ok := LEVELS[r.FormValue("level")]
		if !ok {
			http.Error(w, "The parameter of 'level' is not correct:", 400)
			return
		}

		regexp := r.FormValue("regexp")
		_, err := SetLoggerRegexp(regexp, level)
		if err != nil {
			http.Error(w, "The parameter of 'regexp' is not correct", 400)
			return
		}

		io.WriteString(w, "Level changed successful")
		return
	}

	http.NotFound()
}

func initHttpServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/loggers", handler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go server.ListenAndServe()
}
