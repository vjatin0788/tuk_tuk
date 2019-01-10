package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/TukTuk/errs"

	"github.com/TukTuk/authentication"
	"github.com/TukTuk/fulfilment"
)

func (api *APIMod) DriverWebhook(rw http.ResponseWriter, r *http.Request) (interface{}, error) {

	ctx := r.Context()

	dateTime := r.FormValue("date_time")
	locType := r.FormValue("loc_type")

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

	log.Printf("[DriverWebhook]Lat:%s,Long:%s", lat, long)

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
		return nil, errs.Err("TT_AU_401")

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

	data, err := fulfilment.FF.DriverTracking(ctx, latVal, longVal, driverId, dateTime, locType)
	if err != nil {
		log.Println("[DriverWebhook][Error] Error in updating details", err)
		return nil, err
	}

	return data, err
}

func (api *APIMod) DriverBookHandler(rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error

	ctx := r.Context()

	isBookedStr := r.FormValue("is_booked")
	isBooked, err := strconv.ParseBool(isBookedStr)
	if err != nil {
		log.Println("[DriverBookHandler][Error] Err parsing is booked", err)
		return nil, errors.New("Error Parsing IS BOOKED")
	}

	rideIdStr := r.FormValue("ride_id")
	rideId, err := strconv.ParseInt(rideIdStr, 10, 64)
	if err != nil {
		log.Println("[DriverBookHandler][Error] Err parsing ride id", err)
		return nil, errors.New("Error Parsing Ride ID")
	}

	userid := r.Header.Get("User-Id")
	if userid == "" {
		log.Println("[DriverBookHandler][Error] empty user id")
		return nil, errors.New("Empty User ID")
	}

	uid, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		log.Println("[DriverBookHandler][Error] Parsing int")
		return nil, errors.New("Error parsing int")
	}

	data, err := fulfilment.FF.DriverBooked(ctx, uid, rideId, isBooked)
	if err != nil {
		log.Println("[DriverBookHandler][Error] Err in request ride", err)

		return nil, err
	}

	return data, err
}

func (api *APIMod) DriverLocationHandler(rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error

	ctx := r.Context()

	rideIdStr := r.FormValue("ride_id")
	rideId, err := strconv.ParseInt(rideIdStr, 10, 64)
	if err != nil {
		log.Println("[DriverBookHandler][Error] Err parsing ride id", err)
		return nil, errors.New("Error Parsing Ride ID")
	}

	userid := r.Header.Get("User-Id")
	if userid == "" {
		log.Println("[DriverBookHandler][Error] empty user id")
		return nil, errors.New("Empty User ID")
	}

	uid, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		log.Println("[DriverBookHandler][Error] Parsing int")
		return nil, errors.New("Error parsing int")
	}

	data, err := fulfilment.FF.GetDriverCurrentLocation(ctx, uid, rideId)
	if err != nil {
		log.Println("[DriverBookHandler][Error] Err in request ride", err)
		return nil, err
	}

	return data, err
}

func (api *APIMod) RideStartHandler(rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error

	ctx := r.Context()

	rideIdStr := r.FormValue("ride_id")
	rideId, err := strconv.ParseInt(rideIdStr, 10, 64)
	if err != nil {
		log.Println("[RequestRide][Error] Parsing int")
		return nil, errors.New("Error parsing int")
	}

	userid := r.Header.Get("User-Id")
	if userid == "" {
		log.Println("[RequestRide][Error] empty user id")
		return nil, errors.New("Empty User ID")
	}

	uid, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		log.Println("[RequestRide][Error] Parsing int")
		return nil, errors.New("Error parsing int")
	}

	data, err := fulfilment.FF.StartRide(ctx, uid, rideId)
	if err != nil {
		log.Println("[RequestRide][Error] Err in request ride", err)
		return nil, err
	}

	return data, err
}

func (api *APIMod) RideCompleteHandler(rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var (
		err     error
		reqBody fulfilment.RideCompleteRequest
	)
	ctx := r.Context()

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("[RideCompleteHandler] error:", err)
		return nil, errors.New("Error Body Not found")
	}

	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		log.Println("[RideCompleteHandler] error:", err)
		return nil, errors.New("Error Unmarshal body")
	}

	userid := r.Header.Get("User-Id")
	if userid == "" {
		log.Println("[RideCompleteHandler][Error] empty user id")
		return nil, errors.New("Empty User ID")
	}

	uid, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		log.Println("[RideCompleteHandler][Error] Parsing int")
		return nil, errors.New("Error parsing int")
	}

	data, err := fulfilment.FF.RideComplete(ctx, uid, reqBody)
	if err != nil {
		log.Println("[RideCompleteHandler][Error] Err in request ride", err)
		return nil, err
	}

	return data, err
}

func (api *APIMod) DriverCancelHandler(rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var (
		err     error
		reqBody fulfilment.RideCancelRequest
	)
	ctx := r.Context()

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("[DriverCancelHandler] error:", err)
		return nil, errors.New("Error Body Not found")
	}

	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		log.Println("[DriverCancelHandler] error:", err)
		return nil, errors.New("Error Unmarshal body")
	}

	userid := r.Header.Get("User-Id")
	if userid == "" {
		log.Println("[DriverCancelHandler][Error] empty user id")
		return nil, errors.New("Empty User ID")
	}

	uid, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		log.Println("[DriverCancelHandler][Error] Parsing int")
		return nil, errors.New("Error parsing int")
	}

	data, err := fulfilment.FF.DriverRideCancel(ctx, uid, reqBody)
	if err != nil {
		log.Println("[DriverCancelHandler][Error] Err in request ride", err)
		return nil, err
	}

	return data, err
}

func (api *APIMod) DriverStatusHandler(rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error

	ctx := r.Context()

	userid := r.Header.Get("User-Id")
	if userid == "" {
		log.Println("[DriverStatusHandler][Error] empty user id")
		return nil, errors.New("Empty User ID")
	}

	uid, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		log.Println("[DriverStatusHandler][Error] Parsing int")
		return nil, errors.New("Error parsing int")
	}

	data, err := fulfilment.FF.GetDriverRideStatus(ctx, uid)
	if err != nil {
		log.Println("[DriverStatusHandler][Error] Err in request ride", err)

		return nil, err
	}

	return data, err
}
