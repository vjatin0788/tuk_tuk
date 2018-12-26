package authentication

var Auth *TukTukAuth

type TukTukAuth struct {
}

func InitAuth() {
	Auth = &TukTukAuth{}
}
