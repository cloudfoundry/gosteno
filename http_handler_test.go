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

	cfg := Config{}
	cfg.Sinks = []Sink{NewIOSink(os.Stdout)}
	Init(&cfg)
}

func (s *HttpHandlerSuite) TearDownSuite(c *C) {
	config = Config{}

	s.mux = nil
}

func (s *HttpHandlerSuite) SetUpTest(c *C) {
	loggers = make(map[string]*BaseLogger)
}

func (s *HttpHandlerSuite) TearDownTest(c *C) {
	loggerRegexp = nil
	loggerRegexpLevel = nil

	loggers = nil
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

	NewLogger("foobar")
	c.Assert(loggerRegexp, IsNil)
	c.Assert(loggerRegexpLevel, IsNil)
	c.Assert(loggers["foobar"].level, Equals, LOG_INFO)

	r := regexpParams{"^foobar$", "off"}
	b, _ := json.Marshal(r)
	buf := bytes.NewBuffer(b)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s%s", ts.URL, HTTP_REGEXP_PATH), buf)
	http.DefaultClient.Do(req)

	c.Assert(loggerRegexp.String(), Equals, "^foobar$")
	c.Assert(loggerRegexpLevel, Equals, LOG_OFF)
	c.Assert(loggers["foobar"].level, Equals, LOG_OFF)
}

func (s *HttpHandlerSuite) TestPutRegexpWithWrongParams(c *C) {
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	NewLogger("foobar")

	buf := bytes.NewBufferString("whatever")
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s%s", ts.URL, HTTP_REGEXP_PATH), buf)
	r, err := http.DefaultClient.Do(req)
	c.Assert(err, IsNil)

	m, _ := ioutil.ReadAll(r.Body)
	c.Assert(r.StatusCode, Equals, http.StatusBadRequest)
	c.Assert(string(m), Equals, "Can't parse the parameters\n")
}

func (s *HttpHandlerSuite) TestPutRegexpWithWrongParams2(c *C) {
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	NewLogger("foobar")

	rp := regexpParams{"foobar", "NotExistingLevel"}
	b, _ := json.Marshal(rp)
	buf := bytes.NewBuffer(b)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s%s", ts.URL, HTTP_REGEXP_PATH), buf)
	r, err := http.DefaultClient.Do(req)
	c.Assert(err, IsNil)

	m, _ := ioutil.ReadAll(r.Body)
	c.Assert(r.StatusCode, Equals, http.StatusBadRequest)
	c.Assert(string(m), Equals, "No level with that name exists : NotExistingLevel\n")
}

func (s *HttpHandlerSuite) TestPutRegexpWithWrongParams3(c *C) {
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	NewLogger("foobar")

	rp := regexpParams{"[", "off"}
	b, _ := json.Marshal(rp)
	buf := bytes.NewBuffer(b)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s%s", ts.URL, HTTP_REGEXP_PATH), buf)
	r, err := http.DefaultClient.Do(req)
	c.Assert(err, IsNil)

	m, _ := ioutil.ReadAll(r.Body)
	c.Assert(r.StatusCode, Equals, http.StatusBadRequest)
	c.Assert(string(m), Equals, "The parameter is not a valid regular expression\n")
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

func (s *HttpHandlerSuite) TestPutLoggerWithWrongParams(c *C) {
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	NewLogger("foobar")

	buf := bytes.NewBufferString("whatever")
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s%s%s", ts.URL, HTTP_LOGGER_PATH, "foobar"), buf)
	r, err := http.DefaultClient.Do(req)
	c.Assert(err, IsNil)

	m, _ := ioutil.ReadAll(r.Body)
	c.Assert(r.StatusCode, Equals, http.StatusBadRequest)
	c.Assert(string(m), Equals, "Can't parse the parameters\n")
}

func (s *HttpHandlerSuite) TestPutLoggerWithWrongParams2(c *C) {
	ts := httptest.NewServer(s.mux)
	defer ts.Close()

	NewLogger("foobar")

	l := struct{ Level string }{"NotExistingLevel"}
	b, _ := json.Marshal(l)
	buf := bytes.NewBuffer(b)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s%s%s", ts.URL, HTTP_LOGGER_PATH, "foobar"), buf)
	r, err := http.DefaultClient.Do(req)
	c.Assert(err, IsNil)

	m, _ := ioutil.ReadAll(r.Body)
	c.Assert(r.StatusCode, Equals, http.StatusBadRequest)
	c.Assert(string(m), Equals, "No level with that name exists : NotExistingLevel\n")
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
