package steno

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		http.Error(w, fmt.Sprintf("No logger with the name exists : %s", loggerName), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		bytes, _ := logger.MarshalJSON()
		w.Write(bytes)

	case "PUT":
		var levelParams struct{ Level string }
		jsonData, _ := ioutil.ReadAll(r.Body)

		if err := json.Unmarshal(jsonData, &levelParams); err != nil {
			http.Error(w, "Can't parse the parameters", http.StatusBadRequest)
			return
		}
		level, err := GetLogLevel(levelParams.Level)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.level = level

	default:
		http.NotFound(w, r)
	}
}

type regexpParams struct {
	RegExp string
	Level  string
}

func regExpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		regexpInfo := regexpParams{}
		if loggerRegexp != nil {
			regexpInfo.RegExp = loggerRegexp.String()
			regexpInfo.Level = loggerRegexpLevel.String()
		}
		bytes, err := json.Marshal(regexpInfo)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Write(bytes)

	case "PUT":
		jsonData, _ := ioutil.ReadAll(r.Body)
		params := regexpParams{}
		err := json.Unmarshal(jsonData, &params)
		if err != nil {
			http.Error(w, "Can't parse the parameters", http.StatusBadRequest)
			return
		}
		level, err := GetLogLevel(params.Level)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = SetLoggerRegexp(params.RegExp, level)
		if err != nil {
			http.Error(w, "The parameter is not a valid regular expression", http.StatusBadRequest)
			return
		}

	default:
		http.NotFound(w, r)
	}
}

func loggersListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write([]byte(loggersInJson()))
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
		w.WriteHeader(http.StatusUnauthorized)
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
