package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/TukTuk/authentication"

	"github.com/TukTuk/fulfilment"
)

func (api *APIMod) DriverAvailableHandler(rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()

	lat := r.FormValue("lat")
	latVal, _ := strconv.ParseFloat(lat, 64)

	long := r.FormValue("long")
	longVal, _ := strconv.ParseFloat(long, 64)

	vehicleType := r.FormValue("vehicle_type")

	userid := r.Header.Get("User-Id")
	if userid == "" {
		log.Println("[DriverAvailableHandler][Error] empty user id")
		return nil, errors.New("Empty User ID")
	}

	uid, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		log.Println("[DriverAvailableHandler][Error] Parsing int")
		return nil, errors.New("Err parsing int")
	}

	authToken := r.Header.Get("TUKTUK_TOKEN")
	if authToken == "" {
		log.Println("[DriverAvailableHandler][Error] empty token")
		return nil, errors.New("Empty Auth Token")
	}

	user, err := authentication.Auth.Authentication(ctx, true, false, authToken)
	if err != nil {
		log.Println("[DriverAvailableHandler][Error] Error in fetching authentication details", err)
		return nil, err
	}

	if user.Customer.Id != uid {
		log.Printf("[DriverAvailableHandler][Error] User id mismatch required: %d, found: %d", uid, user.Driver.Id)
		return nil, errors.New("User Id mismatch")
	}

	data, err := fulfilment.FF.DriverAvailable(ctx, latVal, longVal, vehicleType)
	if err != nil {
		log.Println("[DriverAvailableHandler][Error] Error in fetching details", err)
		return nil, err
	}

	return data, err
}

func (api *APIMod) DriverWebhook(rw http.ResponseWriter, r *http.Request) (interface{}, error) {

	ctx := r.Context()

	lat := r.FormValue("lat")
	latVal, _ := strconv.ParseFloat(lat, 64)
	long := r.FormValue("long")
	longVal, _ := strconv.ParseFloat(long, 64)

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
