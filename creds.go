package main

// Credentials - для чтения с запроса
type Credentials struct {
	Password string   `json:"password"`
	Username string   `json:"username"`
	Rts      []string `json:"rts"`
}
