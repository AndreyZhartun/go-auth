package main

import (
	"log"
	"net/http"
	"os"
)

//TODO +транзакции, replica set - есть в атласе
//+jit для пар и фикс бага с 3 маршрутом
//тесты
//$ mongod --dbpath=C:/Users/Пользователь/go/src/go-auth/db
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

	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
