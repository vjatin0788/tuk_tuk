package payment

type PaymentResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       PaymentData `json:"data"`
}

type PaymentData struct {
	TotalCost    float64 `json:"totalCost"`
	Distance     float64 `json:"distance"`
	TimeTaken    int64   `json:"timeTaken"`
	DistanceCost float64 `json:"distance_cost"`
	CostPerKm    string  `json:"costPerKm"`
	CostPerMin   float64 `json:"costPerMinute"`
	TimeCost     float64 `json:"timeCost"`
	Gst          string  `json:"gst"`
	BaseFare     float64 `json:"baseFare"`
}

type PaymentReq struct {
	RideId int64 `json:"ride_id"`
}
