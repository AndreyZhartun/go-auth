package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//receive  - первый маршрут, пока без обращения к БД
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

	//`Access токен тип JWT, алгоритм SHA512.`
	atExpiration := time.Now().Add(5 * time.Minute)
	atClaims := &claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: atExpiration.Unix(),
			//jti для связи токенов в паре
			Id: "test",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	atSigned, err := accessToken.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//`Refresh токен тип произвольный (jwt), формат передачи base64`
	rtExpiration := time.Now().Add(5 * time.Hour)
	rtClaims := &claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: rtExpiration.Unix(),
			Id:        "test",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rtSigned, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO: решить где хранить
	http.SetCookie(w, &http.Cookie{
		Name:     "at",
		Value:    atSigned,
		Expires:  atExpiration,
		HttpOnly: true, //защ от xss
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "rt",
		Value:    rtSigned,
		Expires:  rtExpiration,
		HttpOnly: true,
	})
}
