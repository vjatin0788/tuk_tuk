package errs

import (
	"encoding/json"
	"log"

	"github.com/TukTuk/model"
)

var ErrorMap map[string]APIError

type APIError struct {
	Statuscode int
	Message    string
}

func (err APIError) Error() string {
	byteStr, er := json.Marshal(&err)
	if er != nil {
		log.Println("[Error] Error in marshaling")
	}

	return string(byteStr)
}

func InitError() {
	ErrorMap = make(map[string]APIError)

	modelErr, err := model.TukTuk.GetErrors()
	if err != nil {
		log.Fatal("[InitError][Error] Failed to init api errors", err)
	}

	for _, me := range modelErr {
		ErrorMap[me.ErrorCode] = APIError{
			Statuscode: me.Status,
			Message:    me.Message,
		}
	}
}

func Err(code string) APIError {
	return ErrorMap[code]
}
