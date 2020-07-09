package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// jwt.StandardClaims для полей истекания и пр.
type сlaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func checkToken(c *http.Cookie, err error) (int, *сlaims) {
	if err != nil {
		if err == http.ErrNoCookie {
			return http.StatusUnauthorized, nil
		}
		fmt.Println("cookie err")
		return http.StatusBadRequest, nil
	}
	tknStr := c.Value

	claims := &сlaims{}

	// парс JWT в claims
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		// проверка, что токен подписан с SHA
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("sign method unexpected: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("Неожиданный алгоритм подписи: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Println("tk sign not valid")
			return http.StatusUnauthorized, nil
		}
		fmt.Printf("parse err: %v", err)
		return http.StatusBadRequest, nil
	}
	if !tkn.Valid {
		fmt.Println("tk not valid")
		return http.StatusUnauthorized, nil
	}

	return http.StatusOK, claims
}

func refresh(w http.ResponseWriter, r *http.Request) {
	code, claims := checkToken(r.Cookie("at"))
	if code != http.StatusOK {
		w.WriteHeader(code)
		return
	}

	//проверка рефреш токена
	//проверка связи токенов и тд

	atExpiration := time.Now().Add(5 * time.Minute)
	atClaims := &сlaims{
		Username: claims.Username,
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

	http.SetCookie(w, &http.Cookie{
		Name:     "at",
		Value:    atSigned,
		Expires:  atExpiration,
		HttpOnly: true,
	})
}
