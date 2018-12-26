package fulfilment

type RiderAvailableResponse struct {
	DriverDetails    DriverDetailsResponse `json:"driver_details"`
	CurrentLatitude  float64               `json:"current_lat"`
	CurrentLongitude float64               `json:"current_long"`
}

type DriverDetailsResponse struct {
	DriverId int64 `json:"driver_id"`
}

type DriverTrackingResponse struct {
	Success bool `json:"success"`
}
