package steno

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"os"
)

type HttpHandlerSuite struct {
	mux *http.ServeMux
}

var _ = Suite(&HttpHandlerSuite{})

func (s *HttpHandlerSuite) SetUpSuite(c *C) {
	mux := http.NewServeMux()
	mux.HandleFunc(HTTP_REGEXP_PATH, regExpHandler)
	mux.HandleFunc(HTTP_LOGGER_PATH, loggerHandler)
	mux.HandleFunc(HTTP_LIST_LOGGERS_PATH, loggersListHandler)

	s.mux = mux
}

func (s *HttpHandlerSuite) TearDownSuite(c *C) {
	s.mux = nil
}

func (s *HttpHandlerSuite) SetUpTest(c *C) {
	cfg := Config{}
	cfg.Sinks = []Sink{NewIOSink(os.Stdout)}
	Init(&cfg)
	loggers = make(map[string]*BaseLogger)
}

func (s *HttpHandlerSuite) TearDownTest(c *C) {
	config = Config{}
	loggers = nil
	loggerRegexp = nil
	loggerRegexpLevel = nil
}

func (s *HttpHandlerSuite) TestGetEmptyRegexp(c *C) {
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	r, _ := http.Get(fmt.Sprintf("%s%s", ts.URL, HTTP_REGEXP_PATH))
	data, _ := ioutil.ReadAll(r.Body)

	var rp regexpParams
	err := json.Unmarshal(data, &rp)
	c.Assert(err, IsNil)
	c.Assert(rp.RegExp, Equals, "")
	c.Assert(rp.Level, Equals, "")
}

func (s *HttpHandlerSuite) TestGetRegexp(c *C) {
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	SetLoggerRegexp("^test$", LOG_FATAL)

	r, _ := http.Get(fmt.Sprintf("%s%s", ts.URL, HTTP_REGEXP_PATH))
	data, _ := ioutil.ReadAll(r.Body)

	var rp regexpParams
	err := json.Unmarshal(data, &rp)
	c.Assert(err, IsNil)
	c.Assert(rp.RegExp, Equals, "^test$")
	c.Assert(rp.Level, Equals, "fatal")
}

func (s *HttpHandlerSuite) TestPutRegexp(c *C) {
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	NewLogger("foo")
	NewLogger("bar")
	c.Assert(loggerRegexp, IsNil)
	c.Assert(loggerRegexpLevel, IsNil)
	c.Assert(loggers["foo"].level, Equals, LOG_INFO)
	c.Assert(loggers["bar"].level, Equals, LOG_INFO)

	r := regexpParams{"bar", "off"}
	b, _ := json.Marshal(r)
	buf := bytes.NewBuffer(b)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s%s", ts.URL, HTTP_REGEXP_PATH), buf)
	http.DefaultClient.Do(req)

	c.Assert(loggerRegexp.String(), Equals, "bar")
	c.Assert(loggerRegexpLevel, Equals, LOG_OFF)
	c.Assert(loggers["foo"].level, Equals, LOG_INFO)
	c.Assert(loggers["bar"].level, Equals, LOG_OFF)
}

func (s *HttpHandlerSuite) TestGetNotExistingLogger(c *C) {
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	r, _ := http.Get(fmt.Sprintf("%s%s%s", ts.URL, HTTP_LOGGER_PATH, "foobar"))
	c.Assert(r.StatusCode, Equals, http.StatusBadRequest)
}

func (s *HttpHandlerSuite) TestGetLogger(c *C) {
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	NewLogger("foobar")

	r, _ := http.Get(fmt.Sprintf("%s%s%s", ts.URL, HTTP_LOGGER_PATH, "foobar"))
	data, _ := ioutil.ReadAll(r.Body)

	c.Assert(string(data), Matches, `^{.*info.*}$`)
}

func (s *HttpHandlerSuite) TestPutLogger(c *C) {
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	NewLogger("foobar")
	c.Assert(loggers["foobar"].level, Equals, LOG_INFO)

	l := struct{ Level string }{"off"}
	b, _ := json.Marshal(l)
	buf := bytes.NewBuffer(b)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s%s%s", ts.URL, HTTP_LOGGER_PATH, "foobar"), buf)
	http.DefaultClient.Do(req)

	c.Assert(loggers["foobar"].level, Equals, LOG_OFF)
}

func (s *HttpHandlerSuite) TestGetEmptyLoggers(c *C) {
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	r, _ := http.Get(fmt.Sprintf("%s%s", ts.URL, HTTP_LIST_LOGGERS_PATH))
	data, _ := ioutil.ReadAll(r.Body)

	decodedLoggers := make(map[string]interface{})
	err := json.Unmarshal(data, &decodedLoggers)
	c.Assert(err, IsNil)
	c.Assert(decodedLoggers, HasLen, 0)
}

func (s *HttpHandlerSuite) TestGetLoggers(c *C) {
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	NewLogger("foo")
	NewLogger("bar")

	r, _ := http.Get(fmt.Sprintf("%s%s", ts.URL, HTTP_LIST_LOGGERS_PATH))
	data, _ := ioutil.ReadAll(r.Body)

	decodedLoggers := make(map[string]interface{})
	err := json.Unmarshal(data, &decodedLoggers)
	c.Assert(err, IsNil)

	_, ok := decodedLoggers["foo"]
	c.Assert(ok, Equals, true)
	_, ok = decodedLoggers["bar"]
	c.Assert(ok, Equals, true)
}
