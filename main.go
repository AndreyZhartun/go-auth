package main

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

/*
	Задача: Написать сервис аутентификации.
	Четыре REST маршрута:
	•Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
	•Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов
	•Третий маршрут удаляет конкретный Refresh токен из базы
	•Четвертый маршрут удаляет все Refresh токены из базы для конкретного пользователя

	Технологии: Язык программирования Go.
	База данных MongoDB, топология Replica Set, использование транзакций обязательно.
	Access токен тип JWT, алгоритм SHA512.
	Refresh токен тип произвольный, формат передачи base64,
	хранится в базе исключительно в виде bcrypt хеша,
	должен быть защищен от изменения на стороне клиента и попыток повторного использования.
	Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена
	можно выполнить только тем Refresh токеном который был выдан вместе с ним.

	Результат:
	Результат выполнения задания нужно предоставить в виде исходного кода на Github,
	а также работающего приложения на Heroku.
*/

var jwtKey = []byte("very_secret_key")

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
	http.HandleFunc("/receive", receive)
	http.HandleFunc("/welcome", welcome)
	http.HandleFunc("/refresh", refresh)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
