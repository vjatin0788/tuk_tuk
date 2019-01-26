package common

import (
	"time"
)

const (
	MASTER_DB_LOCAL    = "root:12345@/tuktuk?parseTime=true&loc=Local&allowNativePasswords=true"
	MASTER_DB          = "tuktuk_user:tuktuk_123@tcp(122.160.30.50:3306)/tuktuk?parseTime=true&loc=Local&allowNativePasswords=true"
	MASTER_DB_PROD     = "tuktuk:Tadmin123@tcp(tuktukdb.cx4mkb6ac5le.ap-south-1.rds.amazonaws.com)/tuktuk?parseTime=true&loc=Local&allowNativePasswords=true"
	API_KEY            = ""
	METHOD_GET         = "GET"
	METHOD_POST        = "POST"
	TIME_SLEEP         = 10 * time.Second
	DRIVER_CANCELLED   = 1
	RIDER_CANCELLED    = 1
	NOTIFY_RIDER       = true
	STATUS_ACTIVATE    = "Activate"
	DRIVER_DUTY_STATUS = "On"
	CASH               = "cash"
	DRIVING_MODE       = "Driving"
)

//Endpoints
const (
	DISTANCE_MATRIX = "/maps/api/distancematrix/json"
	GET_PAYMENT     = "/getInvoice"
)

type STATUS_DETAIL struct {
	Label string
	ID    int64
}

var RideStatus = struct {
	REQUESTED  STATUS_DETAIL
	BOOKED     STATUS_DETAIL
	PROCESSING STATUS_DETAIL
	COMPLETED  STATUS_DETAIL
	FAILED     STATUS_DETAIL
}{
	REQUESTED:  STATUS_DETAIL{"Ride Requested", 1},
	BOOKED:     STATUS_DETAIL{"Ride Booked", 2},
	PROCESSING: STATUS_DETAIL{"Ride Processing", 3},
	COMPLETED:  STATUS_DETAIL{"Ride Completed", 4},
	FAILED:     STATUS_DETAIL{"Ride Failed", 5},
}

var RideStatusMap = map[int64]STATUS_DETAIL{
	RideStatus.REQUESTED.ID:  RideStatus.REQUESTED,
	RideStatus.BOOKED.ID:     RideStatus.BOOKED,
	RideStatus.PROCESSING.ID: RideStatus.PROCESSING,
	RideStatus.COMPLETED.ID:  RideStatus.COMPLETED,
	RideStatus.FAILED.ID:     RideStatus.FAILED,
}
