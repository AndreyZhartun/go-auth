package main

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("very_secret_key")

//TODO: mongodb
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// claims для jwt
// jwt.StandardClaims для полей истекания и пр.
type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	//1. выдача пары токенов
	http.HandleFunc("/receive", receive)
	//2. обновление access токена на основе refresh токена
	http.HandleFunc("/refresh", refresh)
	//3.
	//4.
	//дополнительные маршруты для тестирования
	//access - типичный запрос к защищенному ресурсу, не изменяет токены
	http.HandleFunc("/access", accessProtectedResource)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
