package steno

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

type chanSet map[chan []byte]bool

var wsChans = make(map[string]chanSet)
var wsMutex sync.RWMutex

func handler(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.Method, "GET") {
		io.WriteString(w, loggersInJson())
	} else {
		// TODO: PUT not implemented
		http.NotFound(w, r)
	}
}

func tailServer(rw *websocket.Conn) {
	path := rw.Request().RequestURI
	logName := strings.Split(path, "/")[2]

	if loggers[logName] == nil {
		// TODO: How to gracefully close ws conn
		rw.Write([]byte("logger not found"))
		rw.Close()
		return
	}

	wsMutex.Lock()

	if wsChans[logName] == nil {
		wsChans[logName] = make(chanSet)
	}
	set := wsChans[logName]

	ch := make(chan []byte)
	set[ch] = true

	wsMutex.Unlock()

	for {
		msg := <-ch
		n, err := rw.Write(msg)

		// TODO: Need to confirm: this === remote conn closed
		if n == 0 || err != nil {
			break
		}
	}

	wsMutex.Lock()
	delete(set, ch)
	wsMutex.Unlock()
}

func initHttpServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/loggers", handler)
	mux.Handle("/tail/", websocket.Handler(tailServer))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go server.ListenAndServe()
}
