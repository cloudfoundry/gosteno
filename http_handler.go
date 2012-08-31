package steno

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"sync"
)

const (
	WEBSOCKET_TAIL_PATH    = "/ws/tail/"
	HTTP_TAIL_PATH         = "/tail/"
	HTTP_LIST_LOGGERS_PATH = "/loggers"
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

func tailHandler(w http.ResponseWriter, r *http.Request) {
	path := r.RequestURI
	logName := path[len(HTTP_TAIL_PATH):]

	if loggers[logName] == nil {
		http.NotFound(w, r)
		return
	}

	url := fmt.Sprintf("ws://%s%s%s", r.Host, WEBSOCKET_TAIL_PATH, logName)
	fmt.Println(url)

	t, err := template.New("tail").Parse(asset("tail.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.Execute(w, template.URL(url))
}

func tailWSServer(rw *websocket.Conn) {
	path := rw.Request().RequestURI
	logName := path[len(WEBSOCKET_TAIL_PATH):]

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

	mux.HandleFunc(HTTP_LIST_LOGGERS_PATH, handler)
	mux.HandleFunc(HTTP_TAIL_PATH, tailHandler)

	mux.Handle(WEBSOCKET_TAIL_PATH, websocket.Handler(tailWSServer))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go server.ListenAndServe()
}
