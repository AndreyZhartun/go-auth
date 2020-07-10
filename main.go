package main

import (
	"log"
	"net/http"
)

//TODO транзакции
//$ mongod --dbpath=C:/Users/Пользователь/go/src/go-auth/db
var jwtKey = []byte("very_secret_key")

func main() {
	//1. выдача пары токенов
	http.HandleFunc("/receive", receive)
	//2. обновление access токена на основе refresh токена
	http.HandleFunc("/refresh", refresh)
	//3.
	http.HandleFunc("/remove", remove)
	//4.
	//дополнительные маршруты для тестирования
	//access - типичный запрос к защищенному ресурсу, не изменяет токены
	http.HandleFunc("/access", accessProtectedResource)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
