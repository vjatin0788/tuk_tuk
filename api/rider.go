package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/TukTuk/authentication"
	"github.com/TukTuk/errs"
	"github.com/TukTuk/lib"

	"github.com/TukTuk/fulfilment"
)

func (api *APIMod) DriverAvailableHandler(rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()

	lat := r.FormValue("lat")
	latVal, _ := strconv.ParseFloat(lat, 64)

	long := r.FormValue("long")
	longVal, _ := strconv.ParseFloat(long, 64)

	log.Printf("[DriverAvailableHandler]Customer Lat:%s,Long:%s", lat, long)

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
		return nil, errs.Err("TT_AU_401")
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

func (api *APIMod) RequestRideHandler(rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error

	ctx := r.Context()

	source := r.FormValue("source")
	sourceVal, err := lib.StringToFloatArray(source)
	if err != nil {
		log.Println("[RequestRide][Error] empty source")
		return nil, errors.New("Empty Source")
	}

	destination := r.FormValue("destination")
	destVal, err := lib.StringToFloatArray(destination)
	if err != nil {
		log.Println("[RequestRide][Error] empty destination")
		return nil, errors.New("Empty Destination")
	}

	vehicleType := r.FormValue("vehicle_type")

	paymentMethod := r.FormValue("payment_method")

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

	data, err := fulfilment.FF.RequestRide(ctx, uid, sourceVal[0], sourceVal[1], destVal[0], destVal[1], vehicleType, paymentMethod)
	if err != nil {
		log.Println("[RequestRide][Error] Err in request ride", err)

		return nil, err
	}

	return data, err
}

func (api *APIMod) RiderCancelHandler(rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var (
		err     error
		reqBody fulfilment.RideCancelRequest
	)
	ctx := r.Context()

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("[RiderCancelHandler] error:", err)
		return nil, errors.New("Error Body Not found")
	}

	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		log.Println("[RiderCancelHandler] error:", err)
		return nil, errors.New("Error Unmarshal body")
	}

	userid := r.Header.Get("User-Id")
	if userid == "" {
		log.Println("[RiderCancelHandler][Error] empty user id")
		return nil, errors.New("Empty User ID")
	}

	uid, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		log.Println("[RiderCancelHandler][Error] Parsing int")
		return nil, errors.New("Error parsing int")
	}

	data, err := fulfilment.FF.CustomerRideCancel(ctx, uid, reqBody)
	if err != nil {
		log.Println("[RiderCancelHandler][Error] Err in request ride", err)
		return nil, err
	}

	return data, err
}

func (api *APIMod) RiderStatusHandler(rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error

	ctx := r.Context()

	userid := r.Header.Get("User-Id")
	if userid == "" {
		log.Println("[RequestStatusHandler][Error] empty user id")
		return nil, errors.New("Empty User ID")
	}

	uid, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		log.Println("[RequestStatusHandler][Error] Parsing int")
		return nil, errors.New("Error parsing int")
	}

	data, err := fulfilment.FF.GetCustomerRideStatus(ctx, uid)
	if err != nil {
		log.Println("[RequestStatusHandler][Error] Err in request ride", err)

		return nil, err
	}

	return data, err
}
