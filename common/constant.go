package common

const (
	MASTER_DB_LOCAL = "root:12345@/tuktuk?parseTime=true&loc=Local&allowNativePasswords=true"
	MASTER_DB       = "tuktuk_user:tuktuk_123@tcp(122.160.30.50:3306)/tuktuk?parseTime=true&loc=Local&allowNativePasswords=true"
	API_KEY         = "AIzaSyC2zwzwJP1SFBRGVt80SroTm-7ga-z1lcA"
	METHOD_GET      = "GET"
)

//Endpoints
const (
	DISTANCE_MATRIX = "/maps/api/distancematrix/json"
)

type STATUS_DETAIL struct {
	Label string
	ID    int64
}

var RideStatus = struct {
	REQUESTED STATUS_DETAIL
	BOOKED    STATUS_DETAIL
	COMPLETED STATUS_DETAIL
	FAILED    STATUS_DETAIL
}{
	REQUESTED: STATUS_DETAIL{"Ride Requested", 1},
	BOOKED:    STATUS_DETAIL{"Ride Booked", 2},
	COMPLETED: STATUS_DETAIL{"Ride Completed", 3},
	FAILED:    STATUS_DETAIL{"Ride Failed", 4},
}
