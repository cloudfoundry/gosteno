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
	} else {
		// TODO: PUT not implemented
		http.NotFound(w, r)
	}
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
