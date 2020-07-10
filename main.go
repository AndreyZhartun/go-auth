package main

import (
	"log"
	"net/http"
)

//TODO транзакции, replica set
//$ mongod --dbpath=C:/Users/Пользователь/go/src/go-auth/db
var jwtKey = []byte("very_secret_key")

func main() {
	//1. выдача пары токенов
	http.HandleFunc("/receive", receive)
	//2. обновление access токена на основе refresh токена
	http.HandleFunc("/refresh", refresh)
	//3. удаление заданного токена из БД
	http.HandleFunc("/remove", remove)
	//4.удаление всех токенов из БД
	http.HandleFunc("/removeall", removeAll)
	//дополнительные маршруты для тестирования
	//access - типичный запрос к защищенному ресурсу, не изменяет токены
	http.HandleFunc("/access", accessProtectedResource)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
