package steno

import (
	"io"
	"net/http"
	"fmt"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.Method, "GET") {
		io.WriteString(w, loggersInJson())
	} else {
		// TODO: PUT not implemented
		http.NotFound(w, r)
	}
}

func initHttp(port int) {
	http.HandleFunc("/loggers", handler)

	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
