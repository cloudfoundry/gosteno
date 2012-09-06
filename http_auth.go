package steno

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"
)

const TOKEN_LIFE_TIME = 5 * time.Minute

var wsTokens = make(map[string]bool)

func deleteExpiredToken(token string, c <-chan time.Time) {
	<-c
	delete(wsTokens, token)
	log.Printf("Expired token deleted : %s\n", token)
}

func generateToken() string {
	var randBytes [32]byte
	rand.Read(randBytes[:])

	encoder := base64.URLEncoding
	d := make([]byte, encoder.EncodedLen(len(randBytes)))
	encoder.Encode(d, randBytes[:])

	// remove the padding character '='
	token := strings.Replace(string(d), "=", "", -1)
	wsTokens[token] = true
	go deleteExpiredToken(token, time.After(TOKEN_LIFE_TIME))

	return token
}

func wsTokenHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		token := generateToken()
		_, err := w.Write([]byte(token))
		if err != nil {
			log.Println(err)
		}

	default:
		http.NotFound(w, r)
	}
}

func checkAuth(req *http.Request, user string, password string) bool {
	log.Printf("Authenticating for request: %s ...", req.RequestURI)

	if req.Header.Get("Upgrade") == "websocket" {
		token := req.FormValue("token")
		if wsTokens[token] {
			log.Printf("WebSocket request authorized: %s", req.RequestURI)
			return true
		}
		return false
	}

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
