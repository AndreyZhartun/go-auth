package main

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

//дополнительный маршрут для полноты картины
//accessProtectedResource проверяет access токен и выдает доступ к защищенному ресурсу
func accessProtectedResource(w http.ResponseWriter, r *http.Request) {
	// токен из куки
	c, err := r.Cookie("at")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value

	claims := &claims{}

	// парс JWT в claims
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		// проверка, что токен подписан с SHA
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Неожиданный алгоритм подписи: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	/* проверка соотнесения токенов в sub
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		if int(claims["sub"].(float64)) == 1 {

			newTokenPair, err := generateTokenPair()
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, newTokenPair)
		}

		return echo.ErrUnauthorized
	}*/

	// если access токен правильный, то дается доступ
	w.Write([]byte(fmt.Sprintf("Приветствую, %s!", claims.Username)))
}