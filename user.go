package main

// User - модель данных в БД
// guid из параметров запроса 1 маршрута, rts - хеши refresh токенов
type User struct {
	GUID string     `json:"guid"`
	Rts  []([]byte) `json:"rts"`
}
