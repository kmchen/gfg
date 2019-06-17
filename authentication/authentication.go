package authentication

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Authentication struct {
	clientIds map[string]string
}

func (auth *Authentication) Populate() {
	auth.clientIds = make(map[string]string)
	auth.clientIds["clientId0"] = "clientId0"
}

func genHmacSecret() string {
	secret := "somesecret"
	data := "data"
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func parseTokenFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected error")
	}
	return []byte(genHmacSecret()), nil
}

func (auth *Authentication) validateToken(clientId interface{}) bool {
	if _, ok := auth.clientIds[fmt.Sprintf("%v", clientId)]; ok {
		return ok
	}
	return false
}

func (auth *Authentication) IsValidToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, parseTokenFunc)
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		isValidToken := auth.validateToken(claims["clientId"])
		return isValidToken, nil
	}
	return false, nil
}

func (auth *Authentication) GenToken(clientId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"clientId": clientId,
		"exp":      time.Now().Add(time.Minute * 15),
	})
	return token.SignedString([]byte(genHmacSecret()))
}

func (auth *Authentication) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")
		if isValid, _ := auth.IsValidToken(token); isValid {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	}
}

type Client struct {
	Id string `json:"id"`
}

func (auth *Authentication) TokenRequestHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var client Client
	err := decoder.Decode(&client)
	if err != nil || client.Id == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	token, err := auth.GenToken(client.Id)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	w.Write([]byte(token))
}
