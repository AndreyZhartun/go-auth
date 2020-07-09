package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//receive  - первый маршрут, пока только с access токеном
func receive(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO: фикс сравнение чистых паролей
	expectedPassword, ok := users[creds.Username]

	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	atExpiration := time.Now().Add(5 * time.Minute)
	claims := &claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: atExpiration.Unix(),
		},
	}

	//`Access токен тип JWT, алгоритм SHA512.`
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	atSigned, err := accessToken.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "at",
		Value:   atSigned,
		Expires: atExpiration,
	})
}
