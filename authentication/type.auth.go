package authentication

type User struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
}

type AuthUser struct {
	Driver   User `json:"driver"`
	Customer User `json:"customer"`
}
