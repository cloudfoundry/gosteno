package steno

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	HTTP_REGEXP_PATH       = "/regexp"
	HTTP_LIST_LOGGERS_PATH = "/loggers"
	HTTP_LOGGER_PATH       = "/logger/"
)

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

	mux.HandleFunc(HTTP_REGEXP_PATH, regExpHandler)
	mux.HandleFunc(HTTP_LOGGER_PATH, loggerHandler)
	mux.HandleFunc(HTTP_LIST_LOGGERS_PATH, loggersListHandler)

	basicAuth := &BasicAuth{
		handler: mux,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: basicAuth,
	}

	go server.ListenAndServe()
}
