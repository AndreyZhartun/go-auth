package main

// Credentials - для чтения с запроса
type Credentials struct {
	Password string   `json:"password"`
	Username string   `json:"username"`
	Rts      []string `json:"rts"`
}

// User xd
type User struct {
	GUID string     `json:"guid"`
	Rts  []([]byte) `json:"rts"`
}
