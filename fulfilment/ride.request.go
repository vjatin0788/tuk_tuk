package fulfilment

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/TukTuk/firebase"
	"github.com/TukTuk/lib"

	"github.com/TukTuk/common"
	"github.com/TukTuk/maps"
	"github.com/TukTuk/model"
)

func (ff *FFClient) RequestRide(ctx context.Context, customerID int64, sLat, sLong, dLat, dLong float64, vehicleType string) (interface{}, error) {

	log.Printf("[RequestRide] Ride Source Lat:%f,Long:%f and Destination Lat:%f,Long:%f", sLat, sLong, dLat, dLong)

	ride, err := model.TukTuk.GetRideDetailsByCustomerId(ctx, customerID)
	if err != nil {
		log.Println("[RequestRide][Error] Error in fetching data", err)
		return nil, err
	}

	log.Println("[RequestRide] Ride :", ride.Status)

	//only new ride or first ride is allowed
	if ride.Status == common.RideStatus.REQUESTED.ID || ride.Status == common.RideStatus.BOOKED.ID {
		return nil, errors.New("Ride Invalid state")
	}

	data, err := ff.prepareRide(ctx, customerID, sLat, sLong, dLat, dLong, vehicleType)
	if err != nil {
		log.Println("[RequestRide][Error] Error in preparing ride", err)
		return nil, err
	}

	return data, err
}

func (ff *FFClient) prepareRide(ctx context.Context, custId int64, sLat, sLong, dLat, dLong float64, vehicleType string) (*RideBookResponse, error) {
	var (
		ride   model.RideDetailModel
		resp   *RideBookResponse
		err    error
		rideId int64
	)

	if custId == 0 {
		log.Println("[PrepareRide][Error] Empty customer id")
		return nil, errors.New("Empty cust id")
	}

	ride = model.RideDetailModel{
		CustomerId:      custId,
		SourceLat:       sLat,
		SourceLong:      sLong,
		DestinationLat:  dLat,
		DestinationLong: dLong,
	}

	err = ff.RideStateTransition(ctx, &ride, common.RideStatus.REQUESTED.ID)
	if err != nil {
		log.Println("[PrepareRide][Error] Ride state Transitiion", err)
		return nil, err
	}

	rideId, err = model.TukTuk.CreateRide(ctx, ride)
	if err != nil {
		log.Println("[PrepareRide][Error] Err creating ride", err)
		return nil, err
	}

	log.Println("RIDE CREATED, RIDE ID#:", rideId)
	ride.Id = rideId

	resp, finderr := ff.findDriver(ctx, &ride, vehicleType)
	if finderr != nil {
		//any error beyond this mark ride as fail in db
		log.Println("[rideResponse][Error] Ride Failed .setting ride status failed")

		err := ff.RideStateTransition(ctx, &ride, common.RideStatus.FAILED.ID)
		if err != nil {
			log.Println("[rideResponse][Error] Err in state transition", err)
			return nil, err
		}

		ride.RideFailedTime = time.Now().UTC().String()

		err = model.TukTuk.UpdateRideFail(ctx, ride)
		if err != nil {
			log.Println("[rideResponse][Error] Err in updating db", err)
		}

		return nil, finderr
	}

	return resp, err
}

func (ff *FFClient) findDriver(ctx context.Context, ride *model.RideDetailModel, vehicleType string) (*RideBookResponse, error) {
	var (
		err error
	)

	defaultRes := &RideBookResponse{
		Message: "No Driver Found",
		RideId:  ride.Id,
	}

	if ride.SourceLat == 0 || ride.SourceLong == 0 {
		log.Println("[FindDriver][Error] Err Lat long")
		return nil, errors.New("Empty Lat long")
	}

	drivers, err := ff.getAvailableDriverVehicle(ctx, ride.SourceLat, ride.SourceLong, vehicleType)
	if err != nil {
		return nil, err
	}

	if len(drivers) == 0 {
		log.Println("[FindDriver] Not Driver available for given vehicle type")
		return defaultRes, err
	}

	log.Printf("vehicles: %+v", drivers)

	//preparing destination and source for gmaps
	destination := ff.prepareDestinationForGmaps(ctx, drivers)
	source := fmt.Sprintf("%s,%s", ride.SourceLat, ride.SourceLong)

	distance, err := maps.MapsClient.GetDistance(ctx, destination, source)
	if err != nil {
		log.Println("[FindDriver][Error] Err in fetching distance from gmaps", err)
		return nil, err
	}

	log.Printf("Drivers:%+v, Distance:%+v", drivers, distance)

	driversList := ff.getDriversData(ctx, distance, drivers)

	// send push notification.
	driverAloted, err := ff.sendPushNotification(ctx, driversList, ride)
	if err != nil {
		log.Println("[FindDriver][Error] Err in sending push notification", err)
		return nil, err
	}

	return ff.rideResponse(ctx, driverAloted, ride)
}

func (ff *FFClient) rideResponse(ctx context.Context, driverId int64, ride *model.RideDetailModel) (*RideBookResponse, error) {

	var (
		resp RideBookResponse
		err  error
	)

	//if no driver aloted
	if driverId == 0 {
		err := ff.RideStateTransition(ctx, ride, common.RideStatus.FAILED.ID)
		if err != nil {
			log.Println("[rideResponse][Error] Err in state transition", err)
			return nil, err
		}

		ride.RideFailedTime = time.Now().UTC().String()

		err = model.TukTuk.UpdateRideFail(ctx, *ride)
		if err != nil {
			return nil, err
		}

		log.Println("[rideResponse]No Driver Found.")
		return &RideBookResponse{
			Message: "No Driver Found",
			RideId:  ride.Id,
		}, err
	}

	ddata, err := model.TukTuk.GetDriverUserById(ctx, driverId)
	if err != nil {
		log.Println("[rideResponse][Error] DB error", err)
		return nil, err
	}

	if ddata.Userid == 0 {
		log.Println("[RiderAvailable][Error] Empty Usserid data")
		return nil, errors.New("Driver not found in records")
	}

	driverIds := make([]int64, 0)
	driverIds = append(driverIds, driverId)

	vehicles, err := model.TukTuk.GetVehicleByAssignedDriver(ctx, driverIds)
	if err != nil {
		log.Println("[RiderAvailable][Error] Error in fetching vehicle data", err)
		return nil, err
	}

	if len(vehicles) == 0 {
		log.Println("[RiderAvailable][Error] Empty vehilce data")
		return nil, errors.New("Vehicle not assigned")
	}

	driverModel, err := model.TukTuk.GetDriverById(ctx, driverId)
	if err != nil {
		log.Println("[RiderAvailable][Error] Error in fetching data", err)
		return nil, err
	}

	resp = RideBookResponse{
		DriverDetail: DriverDetailsResponse{
			DriverId:      ddata.Userid,
			Name:          ddata.Name,
			Model:         vehicles[0].Model,
			VehicleNumber: vehicles[0].VehicleNumber,
		},
		CurrentLat:  driverModel.CurrentLatitude,
		CurrentLong: driverModel.CurrentLongitude,
		RideId:      ride.Id,
	}

	return &resp, err
}

func (ff *FFClient) getDriversData(ctx context.Context, distances maps.DistanceMatrix, driverModel []model.DriverTrackingModel) []DriverData {
	var drivers []DriverData

	for _, row := range distances.Rows {
		for idx, element := range row.Elements {
			drivers = append(drivers, DriverData{
				Id:       driverModel[idx].DriverID,
				Distance: element.Distance.Value,
			})
		}
	}

	//only sorting on the basis of meters and new conditions can be added
	sort.Slice(drivers, func(i, j int) bool { return drivers[i].Distance < drivers[j].Distance })
	log.Printf("SOrted List of Drivers: %+v", drivers)

	return drivers
}

func (ff *FFClient) prepareDestinationForGmaps(ctx context.Context, drivers []model.DriverTrackingModel) string {
	var destination string

	for idx, driver := range drivers {
		latLong := fmt.Sprintf("%f,%f", driver.CurrentLatitude, driver.CurrentLongitude)
		if idx == 0 {
			destination = fmt.Sprintf("%s", latLong)
		} else {
			destination = fmt.Sprintf("%s|%s", destination, latLong)
		}
	}

	log.Println("[prepareDestinationForGmaps] Destination string:", destination)
	return destination
}

func (ff *FFClient) RideStateTransition(ctx context.Context, ride *model.RideDetailModel, changeToState int64) error {

	var (
		err error
	)

	switch changeToState {
	case common.RideStatus.REQUESTED.ID:
		if ride.Status <= 1 {
			ride.Status = common.RideStatus.REQUESTED.ID
		} else {
			err = errors.New("Invalid state for request ride")
		}
	case common.RideStatus.BOOKED.ID:
		if ride.Status == 1 {
			ride.Status = common.RideStatus.BOOKED.ID
		} else {
			err = errors.New("Invalid state for booking ride")
		}
	case common.RideStatus.PROCESSING.ID:
		if ride.Status == 2 {
			ride.Status = common.RideStatus.BOOKED.ID
		} else {
			err = errors.New("Invalid state for Processing ride")
		}
	case common.RideStatus.COMPLETED.ID:
		if ride.Status == 3 {
			ride.Status = common.RideStatus.COMPLETED.ID
		} else {
			err = errors.New("Invalid state for complete ride")
		}
	case common.RideStatus.FAILED.ID:
		if ride.Status > 0 && ride.Status <= 2 {
			ride.Status = common.RideStatus.FAILED.ID
		} else {
			err = errors.New("Invalid state for failing ride")
		}
	}

	return err
}

func (ff *FFClient) sendPushNotification(ctx context.Context, drivers []DriverData, ride *model.RideDetailModel) (int64, error) {

	var (
		res         int64
		rideUpdated model.RideDetailModel
		err         error
	)

	//configure time for sending notification. It should not be more than 90 seconds.
	startTime := time.Now()

	for idx, driver := range drivers {

		if !ff.checkIfDriverLocValid(ctx, ride, driver) {
			log.Printf("[sendPushNotification] Driver location not valid to send push notification, skipping id: %d", driver.Id)
			continue
		}

		log.Printf("Sending Push notification to driver id:%d", driver.Id)
		//Sending push notification.
		go ff.sendNotification(ctx, ride, driver)
		//Wait for 10 sec before sending another notification
		time.Sleep(common.TIME_SLEEP)

		rideUpdated, err = model.TukTuk.GetRideDetailsByRideId(ctx, ride.Id)
		if err != nil {
			log.Println("[sendPushNotification][Error] Error in fetching data", err)
			return res, err
		}

		if rideUpdated.Status == common.RideStatus.BOOKED.ID {

			if !ff.checkIfDriverBookedIsValid(ctx, rideUpdated.DriverId, drivers[:idx+1]) {
				log.Printf("Driver booked for ride id:%d , driver id:%d , is not valid does not lies in range. %+v", ride.Id, rideUpdated.DriverId, drivers[:idx+1])

				log.Printf("Sending Push notification to wrong driver and cancel it's ride. id:%d", driver.Id)
				go ff.sendInvalidDriverNotification(ctx, &rideUpdated)

				return res, errors.New("Invalid Driver Booked")
			}

			log.Printf("RIDE BOOKED  for ride id:%d , driver id:%d", ride.Id, rideUpdated.DriverId)

			break
		}

		if ff.Cfg.Server.RideRequestTime < time.Since(startTime) {
			log.Printf("[sendPushNotification] Request Time out.", time.Since(startTime).Seconds)
			break
		}

	}

	res = rideUpdated.DriverId

	return res, err
}

func (ff *FFClient) liesInCustomerArea(ctx context.Context, sLat, sLong, dLat, dLong float64) bool {

	//haversine formula
	currentPoint := math.Acos(math.Sin(lib.Rad(sLat))*math.Sin(lib.Rad(dLat)) + math.Cos(lib.Rad(sLat))*math.Cos(lib.Rad(dLat))*math.Cos(lib.Rad(dLong)-lib.Rad(sLong)))

	log.Printf("[liesInCustomerArea] Current Point:%f, dist/rad:%f", currentPoint, DISTANCE/RADIUS)
	if currentPoint <= (DISTANCE / RADIUS) {
		log.Printf("[liesInCustomerArea] location valid for driver lat:%f ,long:%f", dLat, dLong)
		return true
	}
	return false
}

func (ff *FFClient) checkIfDriverLocValid(ctx context.Context, ride *model.RideDetailModel, driver DriverData) bool {

	driverTrackModel, err := model.TukTuk.GetDriverById(ctx, driver.Id)
	if err != nil {
		log.Println("[checkIfDriverLocValid][Error] Error in fetching data", err)
	}

	log.Printf("[checkIfDriverLocValid] driver result:%+v", driverTrackModel)

	return ff.liesInCustomerArea(ctx, ride.SourceLat, ride.SourceLong, driverTrackModel.CurrentLatitude, driverTrackModel.CurrentLongitude)
}

func (ff *FFClient) checkIfDriverBookedIsValid(ctx context.Context, driverId int64, drivers []DriverData) bool {

	var isValid bool

	for _, driver := range drivers {
		if driverId == driver.Id {
			isValid = true
			break
		}
	}

	return isValid
}

func (ff *FFClient) sendNotification(ctx context.Context, ride *model.RideDetailModel, driver DriverData) {

	fbclient := firebase.FClient

	data := map[string]string{
		"ride_id":          fmt.Sprintf("%d", ride.Id),
		"source_lat":       fmt.Sprintf("%f", ride.SourceLat),
		"source_long":      fmt.Sprintf("%f", ride.SourceLong),
		"destination_lat":  fmt.Sprintf("%f", ride.DestinationLat),
		"destination_long": fmt.Sprintf("%f", ride.DestinationLong),
	}

	fbclient.AddId(ctx, driver.DeviceId).SendPushNotification(ctx, data)
}

func (ff *FFClient) sendInvalidDriverNotification(ctx context.Context, ride *model.RideDetailModel) {

	fbclient := firebase.FClient

	ddata, err := model.TukTuk.GetDriverUserById(ctx, ride.DriverId)
	if err != nil {
		log.Println("[rideResponse][Error] DB error", err)
	}

	data := map[string]string{
		"ride_id": fmt.Sprintf("%d", ride.Id),
		"message": "Ride Cancelled. Invalid Ride.",
	}

	fbclient.AddId(ctx, ddata.DeviceId).SendPushNotification(ctx, data)
}
