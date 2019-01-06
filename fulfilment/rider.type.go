package fulfilment

var DriverBookedNotifiedMap map[int64]chan bool

type RiderAvailableResponse struct {
	DriverDetails    DriverDetailsResponse `json:"driver_details"`
	CurrentLatitude  float64               `json:"current_lat"`
	CurrentLongitude float64               `json:"current_long"`
}

type DriverDetailsResponse struct {
	DriverId      int64   `json:"driver_id,omitempty"`
	Name          string  `json:"name,omitempty"`
	Model         string  `json:"vehicle_model,omitempty"`
	VehicleNumber string  `json:"vehicle_number,omitempty"`
	PhoneNumber   string  `json:"phone_number,omitempty"`
	DriverImage   string  `json:"driver_image,omitempty"`
	Rating        string  `json:"rating,omitempty"`
	CurrentLat    float64 `json:"current_lat,omitempty"`
	CurrentLong   float64 `json:"current_long,omitempty"`
	VehicleType   string  `json:"vehicle_type,omitempty"`
}

type DriverTrackingResponse struct {
	Success bool `json:"success"`
}

type RideBookResponse struct {
	DriverDetail    *DriverDetailsResponse `json:"driver_details,omitempty"`
	SourceLat       float64                `json:"source_lat,omitempty"`
	SourceLong      float64                `json:"source_long,omitempty"`
	DestinationLat  float64                `json:"destination_lat,omitempty"`
	DestinationLong float64                `json:"destination_long,omitempty"`
	RideId          int64                  `json:"ride_id,omitempty"`
	Message         string                 `json:"message,omitempty"`
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

type PushNotificationRideCancel struct {
	RideId  int64  `json:"ride_id"`
	Message string `json:"message"`
}

type RideCancelRequest struct {
	RideId int64  `json:"ride_id"`
	Reason string `json:"ride_cancel_msg"`
}

type RideCancelResponse struct {
	Success bool `json:"success"`
}
