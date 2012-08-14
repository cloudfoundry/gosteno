package steno

import (
	. "launchpad.net/gocheck"
)

type HttpHandlerSuite struct {
}

var _ = Suite(&HttpHandlerSuite{})

func (s *HttpHandlerSuite) TestGenerateToken(c *C) {
	t := generateToken()
	// base64.URLEncoder should escape '/', '+'
	c.Assert(string(t), Matches, "^[^/+=?&]*$")
}
