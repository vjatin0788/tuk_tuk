package fulfilment

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/TukTuk/firebase"
	"github.com/TukTuk/lib"

	"github.com/TukTuk/common"
	"github.com/TukTuk/maps"
	"github.com/TukTuk/model"
)

func (ff *FFClient) RequestRide(ctx context.Context, customerID int64, sLat, sLong, dLat, dLong float64, vehicleType, paymentMethod string) (interface{}, error) {

	log.Printf("[RequestRide] Ride Source Lat:%f,Long:%f and Destination Lat:%f,Long:%f", sLat, sLong, dLat, dLong)

	ride, err := model.TukTuk.GetRideDetailsByCustomerId(ctx, customerID)
	if err != nil {
		log.Println("[RequestRide][Error] Error in fetching data", err)
		return nil, err
	}

	log.Println("[RequestRide] Ride :", ride.Status)

	//add go routine to automatically make status=1 to status=5
	//only new ride or first ride is allowed
	if ride.Status == common.RideStatus.REQUESTED.ID || ride.Status == common.RideStatus.BOOKED.ID {
		return nil, errors.New("Ride Invalid state")
	}

	data, err := ff.prepareRide(ctx, customerID, sLat, sLong, dLat, dLong, vehicleType, paymentMethod)
	if err != nil {
		log.Println("[RequestRide][Error] Error in preparing ride", err)
		return nil, err
	}

	return data, err
}

func (ff *FFClient) prepareRide(ctx context.Context, custId int64, sLat, sLong, dLat, dLong float64, vehicleType, paymentMethod string) (*RideBookResponse, error) {
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
		PaymentMethod:   paymentMethod,
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

	rideCancelChan := make(chan bool)
	RequestRideCancel[ride.Id] = rideCancelChan
	log.Printf("[sendPushNotification] RequestRideCancelMap map:%+v", RequestRideCancel)

	resp, finderr := ff.findDriver(ctx, &ride, vehicleType)
	if finderr != nil {
		//any error beyond this mark ride as fail in db
		log.Println("[PrepareRide][Error] Ride Failed .setting ride status failed")

		err := ff.RideStateTransition(ctx, &ride, common.RideStatus.FAILED.ID)
		if err != nil {
			log.Println("[PrepareRide][Error] Err in state transition", err)
			return nil, err
		}

		ride.RideFailedTime = time.Now().UTC().String()

		err = model.TukTuk.UpdateRideFail(ctx, ride)
		if err != nil {
			log.Println("[PrepareRide][Error] Err in updating db", err)
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

	//preparing destination and source for gmaps
	destination := ff.prepareDestinationForGmaps(ctx, drivers)
	source := fmt.Sprintf("%f,%f", ride.SourceLat, ride.SourceLong)

	distance, err := maps.MapsClient.GetDistance(ctx, destination, source)
	if err != nil {
		log.Println("[FindDriver][Error] Err in fetching distance from gmaps", err)
		return nil, err
	}

	log.Printf("Distance:%+v", distance)

	driversList, err := ff.getDriversData(ctx, distance, drivers)
	if err != nil {
		return nil, err
	}

	log.Printf("Final Driver List:%+v", driversList)

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
		DriverDetail: &DriverDetailsResponse{
			DriverId:      ddata.Userid,
			Name:          ddata.Name,
			Model:         vehicles[0].Model,
			VehicleNumber: vehicles[0].VehicleNumber,
			PhoneNumber:   ddata.Mobileno,
			VehicleType:   vehicles[0].VehicleType,
			CurrentLat:    driverModel.CurrentLatitude,
			CurrentLong:   driverModel.CurrentLongitude,
			DriverImage:   ddata.Driverpic,
		},
		SourceLat:       ride.SourceLat,
		SourceLong:      ride.SourceLong,
		DestinationLat:  ride.DestinationLat,
		DestinationLong: ride.DestinationLong,
		RideId:          ride.Id,
	}

	return &resp, err
}

func (ff *FFClient) getDriversData(ctx context.Context, distances maps.DistanceMatrix, driverModel []model.DriverTrackingModel) ([]DriverData, error) {
	var (
		drivers []DriverData
		err     error
	)

	for _, row := range distances.Rows {
		for idx, element := range row.Elements {
			ddata, err := model.TukTuk.GetDriverUserById(ctx, driverModel[idx].DriverID)
			if err != nil {
				log.Println("[rideResponse][Error] DB error", err)
				return drivers, err
			}

			log.Printf("[getDriversData] Driver status:%s, id:%d, dutyStatus:%s,driverModel:%+v", ddata.Status, ddata.Userid, ddata.Driverdutystatus, driverModel[idx])

			if strings.EqualFold(ddata.Status, common.STATUS_ACTIVATE) && strings.EqualFold(ddata.Driverdutystatus, common.DRIVER_DUTY_STATUS) {
				drivers = append(drivers, DriverData{
					Id:       ddata.Userid,
					Distance: element.Distance.Value,
					DeviceId: ddata.DeviceId,
				})
			}
		}
	}

	//will be removed in future.
	if len(drivers) == 0 {
		for idx := range driverModel {
			ddata, err := model.TukTuk.GetDriverUserById(ctx, driverModel[idx].DriverID)
			if err != nil {
				log.Println("[rideResponse][Error] DB error", err)
				return drivers, err
			}

			if strings.EqualFold(ddata.Status, common.STATUS_ACTIVATE) && strings.EqualFold(ddata.Driverdutystatus, common.DRIVER_DUTY_STATUS) {
				drivers = append(drivers, DriverData{
					Id:       ddata.Userid,
					DeviceId: ddata.DeviceId,
				})
			}
		}
	}

	//only sorting on the basis of meters and new conditions can be added
	sort.Slice(drivers, func(i, j int) bool { return drivers[i].Distance < drivers[j].Distance })

	return drivers, err
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
			ride.Status = common.RideStatus.PROCESSING.ID
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

	//adding ride to driver map.
	riderChan := make(chan bool)
	DriverBookedNotifiedMap[ride.Id] = riderChan
	log.Printf("[sendPushNotification] DriverBookedNotifiedMap map:%+v", DriverBookedNotifiedMap)

	//sending notifications to all drivers
	var driverCount = make(map[int64]DriverData)
	for _, driver := range drivers {

		if ff.checkForRideRequestCancellation(ctx, ride.Id) {
			log.Println("[sendPushNotification] Cancel signal recieved ride Id:", ride.Id)
			go ff.sendNotificationIfRideRequestCancelled(ctx, driverCount, ride)
			return res, err
		}

		if !ff.checkIfDriverLocValid(ctx, ride, driver) {
			log.Printf("[sendPushNotification] Driver location not valid to send push notification, skipping id: %d", driver.Id)
			continue
		}

		//Sending push notification.
		log.Printf("Sending Push notification to driver id:%d", driver.Id)
		go ff.sendNotification(ctx, ride, driver)

		driverCount[driver.Id] = driver
	}

	//wait for 30 seconds for driver booking or cancellation.
	//configure time for sending notification. It should not be more than 90 seconds.
	startTime := time.Now()

rideLoop:
	for ff.Cfg.Server.RideRequestTime > time.Duration(time.Since(startTime).Nanoseconds()) {
		select {
		case <-riderChan:
			log.Printf("Booking recieved for id:%d , driver id:%d", ride.Id, rideUpdated.DriverId)

			if ff.checkForRideRequestCancellation(ctx, ride.Id) {
				log.Println("[sendPushNotification] Cancel signal recieved ride Id:", ride.Id)
				go ff.sendNotificationIfRideRequestCancelled(ctx, driverCount, ride)
				break rideLoop
			}

			rideUpdated, err = model.TukTuk.GetRideDetailsByRideId(ctx, ride.Id)
			if err != nil {
				log.Println("[sendPushNotification][Error] Error in fetching data", err)
				return res, err
			}

			updatedDriverId := rideUpdated.DriverId

			if rideUpdated.Status == common.RideStatus.BOOKED.ID {

				if !ff.checkIfDriverBookedIsValid(ctx, &rideUpdated, drivers) {
					log.Printf("Driver booked for ride id:%d , driver id:%d , is not valid does not lies in range. %+v", ride.Id, rideUpdated.DriverId, drivers)

					log.Printf("Sending Push notification to wrong driver and cancel it's ride. id:%d", updatedDriverId)
					ff.sendInvalidDriverNotification(ctx, updatedDriverId, &rideUpdated)
					break
				}
				log.Printf("RIDE BOOKED  for ride id:%d , driver id:%d", ride.Id, rideUpdated.DriverId)
			}
			break rideLoop
		case <-RequestRideCancel[ride.Id]:
			log.Println("[sendPushNotification] Cancel signal recieved ride Id:", ride.Id)
			go ff.sendNotificationIfRideRequestCancelled(ctx, driverCount, ride)

			break rideLoop
		case <-time.After(10 * time.Second):
			log.Printf("No Driver booked yet for id:%d , driver id:%d", ride.Id, rideUpdated.DriverId)
			break
		case <-time.After(ff.Cfg.Server.RideRequestTime - time.Duration(time.Since(startTime).Nanoseconds())):
			log.Printf("Request Time out for ride id:%d, time served:%v", ride.Id, time.Duration(time.Since(startTime).Nanoseconds()))
			break rideLoop
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

func (ff *FFClient) checkIfDriverBookedIsValid(ctx context.Context, ride *model.RideDetailModel, drivers []DriverData) bool {

	var isValid bool

	for _, driver := range drivers {
		if ride.DriverId == driver.Id {
			isValid = true
			break
		}
	}

	if !isValid {
		log.Println("[checkIfDriverBookedIsValid]Invalid Driver.Setting ride status to requsted and driver id to 0.")
		ride.Status = common.RideStatus.REQUESTED.ID
		ride.DriverId = 0

		log.Printf("[checkIfDriverBookedIsValid]ride details:%+v", ride)

		rowCount, err := model.TukTuk.UpdateRideStatus(ctx, *ride)
		if err != nil {
			log.Println("[checkIfDriverBookedIsValid][Error] Err in updating db", err)
		}

		if rowCount == 0 {
			log.Println("[checkIfDriverBookedIsValid] DB not Updated")
		}

	}

	return isValid
}

func (ff *FFClient) sendNotification(ctx context.Context, ride *model.RideDetailModel, driver DriverData) {

	fbclient := firebase.FClient

	userData, err := model.TukTuk.GetCustomerById(ctx, ride.CustomerId)
	if err != nil {
		log.Println("[sendNotification]DB err:", err)
	}

	data := PushNotification{
		Type: "ride_accept",
		Data: PushNotificationRideRequest{
			RideId:      ride.Id,
			CurrentLat:  ride.SourceLat,
			CurrentLong: ride.SourceLong,
			Name:        userData.Name,
		},
	}

	log.Printf("[sendNotification]Push notification:%+v", data)

	go fbclient.SendPushNotification(ctx, data, driver.DeviceId)
}

func (ff *FFClient) sendInvalidDriverNotification(ctx context.Context, driverId int64, ride *model.RideDetailModel) {

	fbclient := firebase.FClient

	ddata, err := model.TukTuk.GetDriverUserById(ctx, driverId)
	if err != nil {
		log.Println("[rideResponse][Error] DB error", err)
	}

	data := PushNotification{
		Type: "invalid_ride",
		Data: PushNotificationInvalidRide{
			RideId:  ride.Id,
			Message: "Invalid Ride",
		},
	}

	log.Printf("[sendNotification]Push notification:%+v", data)

	go fbclient.SendPushNotification(ctx, data, ddata.DeviceId)
}

func (ff *FFClient) checkForRideRequestCancellation(ctx context.Context, rideId int64) bool {
	var isCancelled bool

	log.Printf("[checkForRideRequestCancellation] RequestRideCancel map:%+v", RequestRideCancel)
	if val, ok := RequestRideCancel[rideId]; ok {
		select {
		case <-val:
			isCancelled = true
		default:
			log.Println("[checkForRideRequestCancellation] No cancelation recieved")
		}
	} else {
		//Register in NSQ
		log.Println("[checkForRideRequestCancellation][Error] Error in getting value from map.")
	}

	log.Println("[checkForRideRequestCancellation] Cancellation value:", isCancelled)

	return isCancelled
}

func (ff *FFClient) sendNotificationIfRideRequestCancelled(ctx context.Context, drivers map[int64]DriverData, ride *model.RideDetailModel) {
	fbclient := firebase.FClient

	for key := range drivers {
		ddata, err := model.TukTuk.GetDriverUserById(ctx, key)
		if err != nil {
			log.Println("[sendNotificationIfRideRequestCancelled][Error] DB error", err)
		}

		data := PushNotification{
			Type: "ride_cancel",
			Data: PushNotificationInvalidRide{
				RideId:  ride.Id,
				Message: "RIDE CANCELLED",
			},
		}

		log.Printf("[sendNotificationIfRideRequestCancelled]Push notification:%+v", data)

		fbclient.SendPushNotification(ctx, data, ddata.DeviceId)
	}
}
