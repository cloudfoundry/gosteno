package steno

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"text/template"
)

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if strings.EqualFold(r.Method, "GET") {
		page, err := template.New("index page").Parse(index_page_template)
		if err != nil {
			panic(err)
		}

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
	} else {
		http.NotFound(w, r)
	}
}

func handleLoggers(w http.ResponseWriter, r *http.Request) {
	//FIXME:What if the url is not a well form url such as /loggers//test ?
	trimedPath := (strings.Trim(r.URL.Path, "/"))
	if strings.Count(trimedPath, "/") > 1 {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	//FIXME:handling all kinds of error more sophisticatedly
	case "GET":
		if strings.Count(trimedPath, "/") == 0 {
			io.WriteString(w, loggersInJson())
			return
		}

		loggerName := strings.Split(trimedPath, "/")[1]
		logger, ok := loggers[loggerName]
		if !ok {
			http.Error(w, "No logger with the name exist", http.StatusBadRequest)
			return
		}
		bytes, _ := logger.MarshalJSON()
		w.Write(bytes)

	case "PUT":
		if strings.Count(trimedPath, "/") == 0 {
			http.Error(w, "Not implement yet", http.StatusBadRequest)
			return
		}

		loggerName := strings.Split(trimedPath, "/")[1]
		logger, ok := loggers[loggerName]
		if !ok {
			http.Error(w, "No logger with the name exist", http.StatusBadRequest)
			return
		}
		var levelJson struct{ Level string }
		jsonData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		if err = json.Unmarshal(jsonData, &levelJson); err != nil {
			http.Error(w, "The parameter of level is not correct", http.StatusBadRequest)
			return
		}
		level, ok := LEVELS[levelJson.Level]
		if !ok {
			http.Error(w, "No level with that name exist", http.StatusBadRequest)
		}
		logger.level = level
	default:
		http.NotFound(w, r)
	}
}

func handleRexExp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/regexp" {
		http.NotFound(w, r)
		return
	}

	type regexpMsg struct {
		RegExp string
		Level  string
	}
	switch r.Method {
	//FIXME:handling all kinds of error more sophisticatedly
	case "GET":
		bytes, err := json.Marshal(regexpMsg{loggerRegexp.String(), loggerRegexpLevel.name})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			panic(err)
			return
		}
		if _, err = w.Write(bytes); err != nil {
			panic(err)
		}

	case "PUT":
		jsonData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		msg := regexpMsg{}
		err = json.Unmarshal(jsonData, &msg)
		if err != nil {
			panic(err)
		}
		level, ok := LEVELS[msg.Level]
		if !ok {
			http.Error(w, "The parameter of level is not correct", http.StatusBadRequest)
			return
		}
		_, err = SetLoggerRegexp(msg.RegExp, level)
		if err != nil {
			http.Error(w, "The parameter of regexp is not correct", http.StatusBadRequest)
			return
		}
	default:
	}
}

func initHttpServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/regexp", handleRexExp)
	mux.HandleFunc("/loggers/", handleLoggers)
	mux.HandleFunc("/", index)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go server.ListenAndServe()
}
