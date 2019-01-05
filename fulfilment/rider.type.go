package fulfilment

var DriverBookedNotifiedMap map[int64]chan bool

type RiderAvailableResponse struct {
	DriverDetails    DriverDetailsResponse `json:"driver_details"`
	CurrentLatitude  float64               `json:"current_lat"`
	CurrentLongitude float64               `json:"current_long"`
}

type DriverDetailsResponse struct {
	DriverId      int64   `json:"driver_id"`
	Name          string  `json:"name"`
	Model         string  `json:"vehicle_model"`
	VehicleNumber string  `json:"vehicle_number"`
	PhoneNumber   string  `json:"phone_number"`
	DriverImage   string  `json:"driver_image"`
	Rating        string  `json:"rating"`
	CurrentLat    float64 `json:"current_lat"`
	CurrentLong   float64 `json:"current_long"`
	VehicleType   string  `json:"vehicle_type"`
}

type DriverTrackingResponse struct {
	Success bool `json:"success"`
}

type RideBookResponse struct {
	DriverDetail    DriverDetailsResponse `json:"driver_details"`
	SourceLat       float64               `json:"source_lat"`
	SourceLong      float64               `json:"source_long"`
	DestinationLat  float64               `json:"destination_lat"`
	DestinationLong float64               `json:"destination_long"`
	RideId          int64                 `json:"ride_id"`
	Message         string                `json:"message"`
}

type DriverData struct {
	Id       int64
	Distance int64
	DeviceId string
}

type DriverBookedResponse struct {
	RiderDetail CustomerDetailsResponse `json:"ride_details"`
	CurrentLat  float64                 `json:"current_lat"`
	CurrentLong float64                 `json:"current_long"`
	RideId      int64                   `json:"ride_id"`
	Message     string                  `json:"message"`
}

type CustomerDetailsResponse struct {
	CustomerId int64  `json:"customer_id"`
	Name       string `json:name`
}

type RideStartResponse struct {
	Success     bool    `json:"success"`
	CurrentLat  float64 `json:"destination_lat"`
	CurrentLong float64 `json:"destination_long"`
	RideId      int64   `json:"ride_id"`
}

type DriverLocationResponse struct {
	CurrentLat  float64 `json:"current_lat"`
	CurrentLong float64 `json:"current_long"`
	RideId      int64   `json:"ride_id"`
	Message     string  `json:"message"`
}

type PushNotificationRideRequest struct {
	CurrentLat  float64 `json:"current_lat"`
	CurrentLong float64 `json:"current_long"`
	RideId      int64   `json:"ride_id"`
	Name        string  `json:"name"`
}

type PushNotificationInvalidRide struct {
	RideId  int64  `json:"ride_id"`
	Message string `json:"message"`
}

type PushNotification struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type RideCompleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Amount  int64  `json:"amount"`
}

type RideCompleteRequest struct {
	RideId          int64   `json:"ride_id"`
	DestinationLat  float64 `json:"destination_lat"`
	DestinationLong float64 `json:"destination_long"`
}

type PushNotificationRideComplete struct {
	RideId  int64  `json:"ride_id"`
	Message string `json:"message"`
}
