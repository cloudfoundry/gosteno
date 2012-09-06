package steno

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

func checkAuth(req *http.Request, user string, password string) bool {
	log.Printf("Authenticating for request: %s ...", req.RequestURI)

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

	log.Printf("Request authorized : %s", req.RequestURI)
	return true
}
