package fulfilment

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/TukTuk/common"
	"github.com/TukTuk/maps"
	"github.com/TukTuk/model"
)

func (ff *FFClient) RequestRide(ctx context.Context, customerID int64, sLat, sLong, dLat, dLong float64, vehicleType string) (interface{}, error) {

	// if origin == "" {
	// 	log.Println("[RequestRide][Error]Empty Origin details")
	// 	return nil, errors.New("Empty Origin details")
	// }

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

	resp, err = ff.findDriver(ctx, ride, rideId, vehicleType)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (ff *FFClient) findDriver(ctx context.Context, ride model.RideDetailModel, rideId int64, vehicleType string) (*RideBookResponse, error) {
	var (
		//	defaultResp RideBookResponse
		err error
	)

	if ride.SourceLat == 0 || ride.SourceLong == 0 {
		log.Println("[FindDriver][Error] Err Lat long")
		return nil, errors.New("Empty Lat long")
	}

	drivers, err := ff.getAvailableDriverVehicle(ctx, ride.SourceLat, ride.SourceLong, vehicleType)
	if err != nil {
		return nil, err
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
	driverAloted, err := ff.sendPushNotification(ctx, driversList, rideId)
	if err != nil {
		log.Println("[FindDriver][Error] Err in sending push notification", err)
		return nil, err
	}

	return ff.rideResponse(ctx, driverAloted, rideId)
}

func (ff *FFClient) rideResponse(ctx context.Context, driverId, rideId int64) (*RideBookResponse, error) {

	var (
		resp RideBookResponse
		err  error
	)

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
		RideId:      rideId,
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
	case common.RideStatus.COMPLETED.ID:
		if ride.Status == 2 {
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

func (ff *FFClient) sendPushNotification(ctx context.Context, drivers []DriverData, rideId int64) (int64, error) {

	var (
		res  int64
		ride model.RideDetailModel
		err  error
	)

	for _, driver := range drivers {
		ride, err = model.TukTuk.GetRideDetailsByRideId(ctx, rideId)
		if err != nil {
			log.Println("[sendPushNotification][Error] Error in fetching data", err)
			return res, err
		}

		if ride.Status == common.RideStatus.BOOKED.ID {
			break
		}

		log.Printf("Sending Push notification to driver id:%d", driver.Id)
		//Sending push notification for 20 sec.
		time.Sleep(20 * time.Second)
	}

	res = ride.DriverId
	log.Printf("Ride booked for ride id:%d , driver id:%d", rideId, ride.DriverId)

	return res, err
}
