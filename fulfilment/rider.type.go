package fulfilment

type RiderAvailableResponse struct {
	DriverDetails    DriverDetailsResponse `json:"driver_details"`
	CurrentLatitude  float64               `json:"current_lat"`
	CurrentLongitude float64               `json:"current_long"`
}

type DriverDetailsResponse struct {
	DriverId      int64  `json:"driver_id"`
	Name          string `json:name`
	Model         string `json:"model"`
	VehicleNumber string `json:"vehicle_number"`
}

type DriverTrackingResponse struct {
	Success bool `json:"success"`
}

type RideBookResponse struct {
	DriverDetail DriverDetailsResponse `json:"driver_details"`
	CurrentLat   float64               `json:"current_lat"`
	CurrentLong  float64               `json:"current_long"`
	RideId       int64                 `json:"ride_id"`
}

type DriverData struct {
	Id       int64
	Distance int64
	DeviceId string
}
