package steno

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"text/template"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/loggers" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	switch r.Method {
	case "GET":
		io.WriteString(w, loggersInJson())

	case "PUT", "POST":
		level, ok := LEVELS[r.FormValue("level")]
		if !ok {
			http.Error(w, "The parameter of 'level' is not correct:", http.StatusBadRequest)
			return
		}
		regexp := r.FormValue("regexp")
		_, err := SetLoggerRegexp(regexp, level)
		if err != nil {
			http.Error(w, "The parameter of 'regexp' is not correct", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	default:
		http.NotFound(w, r)
	}
}

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

func initHttpServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/loggers", handler)
	mux.HandleFunc("/", index)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go server.ListenAndServe()
}
