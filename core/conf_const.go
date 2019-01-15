package core

import "time"

const (
	GMAPS_SERVICE_HOSTNAME = "https://maps.googleapis.com"
	PAYMENT_STAGING        = "http://localhost:5000"
	PAYMENT_LOCAL          = "http://122.160.30.50:5092"
	FIREBASE_KEY           = ""
	RIDE_REQUEST_TIME      = 30 * 1000000000 * time.Nanosecond
	FIREBASE_TIMEOUT       = 10
	DRIVER_ARRIVED         = 0.100
	DRIVER_ARRIVAL         = 0.500
	RIDE_REQUEST_LOC       = 3.000
)
