package api

var Api *APIMod

type APIMod struct {
}

func InitApiMod() *APIMod {
	return &APIMod{}
}
