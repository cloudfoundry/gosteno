package steno

import (
	"fmt"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
)

type HttpAuthSuite struct {
	basicAuth *BasicAuth
}

var _ = Suite(&HttpAuthSuite{})

func (s *HttpAuthSuite) SetUpSuite(c *C) {
	mux := http.NewServeMux()
	mux.HandleFunc(HTTP_REGEXP_PATH, regExpHandler)
	mux.HandleFunc(HTTP_LOGGER_PATH, loggerHandler)
	mux.HandleFunc(HTTP_LIST_LOGGERS_PATH, loggersListHandler)

	s.basicAuth = &BasicAuth{
		handler: mux,
	}

	cfg := Config{}
	cfg.User = "jeff"
	cfg.Password = "li"
	cfg.Sinks = []Sink{newNullSink()}
	Init(&cfg)
}

func (s *HttpAuthSuite) TearDownSuite(c *C) {
	config = Config{}

	s.basicAuth = nil
}

func (s *HttpAuthSuite) SetUpTest(c *C) {
	loggers = make(map[string]*BaseLogger)
}

func (s *HttpAuthSuite) TearDownTest(c *C) {
	loggerRegexp = nil
	loggerRegexpLevel = nil

	loggers = nil
}

func testGetMethodStatusCode(req *http.Request, c *C) {
	r, _ := http.DefaultClient.Do(req)
	c.Assert(r.StatusCode, Equals, http.StatusUnauthorized)

	req.SetBasicAuth("jeff", "li")
	r, _ = http.DefaultClient.Do(req)
	c.Assert(r.StatusCode, Equals, http.StatusOK)
}

func (s *HttpAuthSuite) TestGetRegexpWithAuth(c *C) {
	ts := httptest.NewServer(s.basicAuth)
	defer ts.Close()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", ts.URL, HTTP_REGEXP_PATH), nil)
	testGetMethodStatusCode(req, c)
}

func (s *HttpAuthSuite) TestGetLoggerWithAuth(c *C) {
	ts := httptest.NewServer(s.basicAuth)
	defer ts.Close()

	NewLogger("foobar")

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s%s", ts.URL, HTTP_LOGGER_PATH, "foobar"), nil)
	testGetMethodStatusCode(req, c)
}

func (s *HttpAuthSuite) TestGetLoggersWithAuth(c *C) {
	ts := httptest.NewServer(s.basicAuth)
	defer ts.Close()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", ts.URL, HTTP_LIST_LOGGERS_PATH), nil)
	testGetMethodStatusCode(req, c)
}
