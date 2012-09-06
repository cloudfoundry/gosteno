package steno

import (
	"code.google.com/p/go.net/websocket"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
)

const (
	HTTP_ROOT_PATH         = "/"
	HTTP_REGEXP_PATH       = "/regexp"
	HTTP_LIST_LOGGERS_PATH = "/loggers"
	HTTP_LOGGER_PATH       = "/logger/"
	WEBSOCKET_TAIL_PATH    = "/ws/tail/"
	HTTP_TAIL_PATH         = "/tail/"
)

type chanSet map[chan []byte]bool

var wsChans = make(map[string]chanSet)
var wsMutex sync.RWMutex

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET":
		page, _ := template.New("index page").Parse(asset("index.html.tp"))

		loggersInfo := make(map[string]string)
		for k, v := range loggers {
			bytes, _ := v.MarshalJSON()
			loggersInfo[k] = string(bytes)
		}

		i := 0
		levels := make([]string, len(LEVELS))
		for k, _ := range LEVELS {
			levels[i] = k
			i++
		}
		sort.Strings(levels)

		page.Execute(w, struct {
			LoggersInfo map[string]string
			Levels      []string
		}{loggersInfo, levels})

	default:
		http.NotFound(w, r)
	}
}

func loggerHandler(w http.ResponseWriter, r *http.Request) {
	loggerName := r.URL.Path[len(HTTP_LOGGER_PATH):]
	logger, ok := loggers[loggerName]
	if !ok {
		http.Error(w, "No logger with the name exists : "+loggerName, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		bytes, _ := logger.MarshalJSON()
		if _, err := w.Write(bytes); err != nil {
			log.Println(err)
		}

	case "PUT":
		var levelParams struct{ Level string }
		jsonData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		if err = json.Unmarshal(jsonData, &levelParams); err != nil {
			http.Error(w, "Can't parse the parameters", http.StatusBadRequest)
			return
		}
		level, ok := LEVELS[levelParams.Level]
		if !ok {
			http.Error(w, "No level with that name exists : "+levelParams.Level, http.StatusBadRequest)
			return
		}
		logger.level = level

	default:
		http.NotFound(w, r)
	}
}

func regExpHandler(w http.ResponseWriter, r *http.Request) {
	type regexpParams struct {
		RegExp string
		Level  string
	}
	switch r.Method {
	case "GET":
		bytes, err := json.Marshal(regexpParams{loggerRegexp.String(), loggerRegexpLevel.name})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if _, err = w.Write(bytes); err != nil {
			log.Println(err)
		}

	case "PUT":
		jsonData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		params := regexpParams{}
		err = json.Unmarshal(jsonData, &params)
		if err != nil {
			http.Error(w, "Can't parse parameters", http.StatusBadRequest)
			return
		}
		level, ok := LEVELS[params.Level]
		if !ok {
			http.Error(w, "No level with the name exists : "+params.Level, http.StatusBadRequest)
			return
		}
		err = SetLoggerRegexp(params.RegExp, level)
		if err != nil {
			http.Error(w, "The parameter of regexp is not correct", http.StatusBadRequest)
			return
		}

	default:
		http.NotFound(w, r)
	}
}

func loggersListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if _, err := w.Write([]byte(loggersInJson())); err != nil {
			log.Println(err)
		}
	default:
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

func checkAuth(req *http.Request, user string, password string) bool {
	// FIXME:websocket client authentication is simply ignored currently for 
	//       there is no suitable solution existing in official protocol yet.
	if req.Header.Get("Upgrade") == "websocket" {
		return true
	}

	if user == "" && password == "" {
		return true
	}

	authParts := strings.Split(req.Header.Get("Authorization"), " ")
	if len(authParts) != 2 || authParts[0] != "Basic" {
		return false
	}
	code, err := base64.StdEncoding.DecodeString(authParts[1])
	if err != nil {
		return false
	}
	userPass := strings.Split(string(code), ":")
	if len(userPass) != 2 || userPass[0] != user || userPass[1] != password {
		return false
	}
	return true
}

type BasicAuth struct {
	handler http.Handler
}

func (a *BasicAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !checkAuth(req, config.User, config.Password) {
		w.Header().Set("WWW-Authenticate", "Basic")
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
	} else {
		a.handler.ServeHTTP(w, req)
	}
}

func initHttpServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc(HTTP_ROOT_PATH, rootHandler)
	mux.HandleFunc(HTTP_REGEXP_PATH, regExpHandler)
	mux.HandleFunc(HTTP_LOGGER_PATH, loggerHandler)
	mux.HandleFunc(HTTP_LIST_LOGGERS_PATH, loggersListHandler)
	mux.HandleFunc(HTTP_TAIL_PATH, tailHandler)
	mux.Handle(WEBSOCKET_TAIL_PATH, websocket.Handler(tailWSServer))

	basicAuth := &BasicAuth{
		handler: mux,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: basicAuth,
	}

	go server.ListenAndServe()
}
