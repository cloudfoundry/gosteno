package steno

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func checkAuth(req *http.Request, user string, password string) bool {

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
