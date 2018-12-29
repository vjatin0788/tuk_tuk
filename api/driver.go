package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/TukTuk/authentication"
	"github.com/TukTuk/fulfilment"
)

func (api *APIMod) DriverWebhook(rw http.ResponseWriter, r *http.Request) (interface{}, error) {

	ctx := r.Context()

	lat := r.FormValue("lat")
	latVal, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		log.Println("[DriverWebhook][Error] error parsing lat")
		return nil, errors.New("Wrong Lat value")
	}

	long := r.FormValue("long")
	longVal, err := strconv.ParseFloat(long, 64)
	if err != nil {
		log.Println("[DriverWebhook][Error]  error parsing long")
		return nil, errors.New("Wrong Long value")
	}

	userid := r.Header.Get("User-Id")
	//verify user id. this will be driver id
	if userid == "" {
		log.Println("[DriverWebhook][Error] empty user id")
		return nil, errors.New("Empty Driver ID")
	}

	driverId, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		log.Println("[DriverWebhook][Error] Error in parsing integer", err)
		return nil, errors.New("Error parsing integer")
	}

	authToken := r.Header.Get("TUKTUK_TOKEN")
	if authToken == "" {
		log.Println("[DriverAvailableHandler][Error] empty token")
		return nil, errors.New("Empty Auth Token")
	}

	duser, err := authentication.Auth.Authentication(ctx, false, true, authToken)
	if err != nil {
		log.Println("[DriverWebhook][Error] Error in fetching authentication details", err)
		return nil, err
	}

	if duser.Driver.Id != driverId {
		log.Printf("[DriverWebhook][Error] User id mismatch required: %d, found: %d", driverId, duser.Driver.Id)
		return nil, errors.New("User Id mismatch")
	}

	data, err := fulfilment.FF.DriverTracking(ctx, latVal, longVal, driverId)
	if err != nil {
		log.Println("[DriverWebhook][Error] Error in updating details", err)
		return nil, err
	}

	return data, err
}
