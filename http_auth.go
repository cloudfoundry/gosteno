package steno

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"
)

var TOKEN_LIFE_TIME = 5 * time.Minute // A token is valid only for 5 minutes
var wsTokens = make(map[string]time.Time)

// If the client only requests the token while doesn't make any authentication
// request,memoey leak could occur, so this function will run periodically.
// The interval is set to twice the life time of token
func cleanExpiredTokens() {
	c := time.Tick(TOKEN_LIFE_TIME * 2)
	for now := range c {
		log.Println("Cleaning expired tokens...")
		for k, v := range wsTokens {
			if v.Sub(now) < 0 {
				delete(wsTokens, k)
			}
		}
	}
}

func generateToken() string {
	var randBytes [32]byte
	rand.Read(randBytes[:])
	encoder := base64.URLEncoding
	d := make([]byte, encoder.EncodedLen(len(randBytes)))
	encoder.Encode(d, randBytes[:])
	token := strings.Replace(string(d), "=", "", -1) // remove the padding character '='
	wsTokens[token] = time.Now().Add(TOKEN_LIFE_TIME)
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
		if wsTokens[token].Sub(time.Now()) > 0 { // if token not exists, wsToken[token] is January 1, year 1, 00:00:00.000000000 UTC
			delete(wsTokens, token)
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
