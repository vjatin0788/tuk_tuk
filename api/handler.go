package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TukTuk/errs"

	"github.com/gorilla/mux"
)

type Base struct {
	StatusCode   int    `json:"statusCode"`
	ErrorMessage string `json:"message"`
}

type Response struct {
	Base
	Data interface{} `json:"data"`
}

// each handler can return the data and error, and serveHTTP can chose how to convert this
type HandlerFunc func(rw http.ResponseWriter, r *http.Request) (interface{}, error)

func (fn HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	response := Response{}
	response.Base.StatusCode = 200

	var data interface{}
	var err error

	//errStatus := http.StatusInternalServerError

	data, err = fn(rw, r)
	var buf []byte

	rw.Header().Set("Content-Type", "application/json")

	if data != nil && err == nil {
		response.Data = data
		log.Println(data)
	}

	if err != nil {
		switch t := err.(type) {
		case errs.APIError:
			response.ErrorMessage = t.Message
			response.StatusCode = t.Statuscode
			rw.WriteHeader(t.Statuscode)
		default:
			response.ErrorMessage = err.Error()
			response.StatusCode = 400
			rw.WriteHeader(400)
		}

	}

	if buf, err = json.Marshal(response); err != nil {
		rw.WriteHeader(400)
	}

	rw.Write(buf)

}

func (api *APIMod) InitHandler() {
	r := mux.NewRouter()

	r.Handle("/v1/tuktuk/driver/available", HandlerFunc(api.DriverAvailableHandler)).Methods("GET")

	r.Handle("/v1/tuktuk/driver/hotspot", HandlerFunc(api.DriverWebhook)).Methods("GET")

	r.Handle("/v1/tuktuk/driver/book", HandlerFunc(api.DriverBookHandler)).Methods("GET")

	r.Handle("/v1/tuktuk/driver/start", HandlerFunc(api.RideStartHandler)).Methods("GET")

	r.Handle("/v1/tuktuk/driver/ride/cancel", HandlerFunc(api.DriverCancelHandler)).Methods("POST")

	r.Handle("/v1/tuktuk/driver/complete", HandlerFunc(api.RideCompleteHandler)).Methods("POST")

	r.Handle("/v1/tuktuk/driver/status", HandlerFunc(api.DriverStatusHandler)).Methods("GET")

	r.Handle("/v1/tuktuk/rider/driver/location", HandlerFunc(api.DriverLocationHandler)).Methods("GET")

	r.Handle("/v1/tuktuk/rider/ride/cancel", HandlerFunc(api.RiderCancelHandler)).Methods("POST")

	r.Handle("/v1/tuktuk/rider/request", HandlerFunc(api.RequestRideHandler)).Methods("GET")

	r.Handle("/v1/tuktuk/rider/status", HandlerFunc(api.RiderStatusHandler)).Methods("GET")

	http.Handle("/", r)
	log.Println("Handler initialized")
}
