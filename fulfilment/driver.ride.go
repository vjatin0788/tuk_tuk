package fulfilment

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/TukTuk/common"
	"github.com/TukTuk/firebase"

	"github.com/TukTuk/model"
)

func (ff *FFClient) StartRide(ctx context.Context, userId, rideId int64) (interface{}, error) {
	var (
		err         error
		defaultResp *RideStartResponse
	)

	defaultResp = &RideStartResponse{
		Success: false,
	}

	if rideId == 0 {
		log.Println("[StartRide][Error] Error Ride Id is 0.")
		return defaultResp, errors.New("Ride Id is 0")
	}

	rideDetail, err := model.TukTuk.GetRideDetailsByRideId(ctx, rideId)
	if err != nil {
		log.Println("[StartRide][Error] Error in fetching ride data", err)
		return defaultResp, err
	}

	//it's check in case there is no ride of requested ride id.
	if rideId != rideDetail.Id {
		log.Println("[StartRide][Error] Invalid Ride id", rideId)
		return defaultResp, errors.New("Invalid Ride ID.")
	}

	log.Printf("[StartRide] Ride:%+v ", rideDetail)

	ddata, err := model.TukTuk.GetDriverUserById(ctx, userId)
	if err != nil {
		log.Println("[StartRide][Error] Error in fetching ride data", err)
		return defaultResp, err
	}

	//verify ride details
	err = ff.verifyRideDetails(ctx, &rideDetail, ddata)
	if err != nil {
		return defaultResp, err
	}

	//state transition
	err = ff.RideStateTransition(ctx, &rideDetail, common.RideStatus.PROCESSING.ID)
	if err != nil {
		log.Println("[StartRide][Error] Ride state Transitiion", err)
		return defaultResp, err
	}

	//update ride details
	rideDetail.RideStartTime = time.Now().UTC().String()
	log.Printf("[StartRide] Ride Updated:%+v ", rideDetail)
	rowAffectedCount, err := model.TukTuk.UpdateRideStart(ctx, rideDetail)
	if err != nil {
		log.Println("[StartRide][Error] Err in updating db", err)
		return defaultResp, err
	}

	if rowAffectedCount == 0 {
		log.Println("[StartRide][Error] Ride is not in valid state db,Row Affected:", rowAffectedCount)
		return defaultResp, errors.New("Something Went Wrong.")
	}

	ff.sendPushNotificationToCustomer(ctx, &rideDetail, ddata)

	defaultResp.Success = true

	return defaultResp, err
}

func (ff *FFClient) verifyRideDetails(ctx context.Context, ride *model.RideDetailModel, ddata model.DriverUserModel) error {
	var err error

	//check ride status
	if ride.Status != common.RideStatus.BOOKED.ID {
		log.Println("[verifyRideDetails][Error] Error in fetching ride data", err)
		return errors.New("Invalid Ride Status")
	}

	if ddata.Userid != ride.DriverId {
		log.Printf("[verifyRideDetails][Error] Driver ID mismatch in ride details. Found id:%d, required id:%d", ddata.Userid, ride.DriverId)
		return errors.New("Driver Id mismatch")
	}

	if !strings.EqualFold(ddata.Status, common.STATUS_ACTIVATE) {
		log.Printf("[verifyRideDetails][Error] DDriver status not valid", ddata.Status)
		return errors.New("Invalid Driver Status")
	}

	//add payment verifications.

	return err
}

func (ff *FFClient) sendPushNotificationToCustomer(ctx context.Context, ride *model.RideDetailModel, driver model.DriverUserModel) {
	fbclient := firebase.FClient

	data := map[string]string{
		"ride_id": fmt.Sprintf("%d", ride.Id),
		"message": fmt.Sprint("Driver Arrived"),
	}

	go fbclient.AddId(ctx, driver.DeviceId).SendPushNotification(ctx, data)
}

func (ff *FFClient) GetDriverCurrentLocation(ctx context.Context, userId, rideId int64) (interface{}, error) {
	var (
		defaultResp DriverLocationResponse
		err         error
	)

	if rideId == 0 {
		log.Println("[GetDriverCurrentLocation][Error] Error Ride Id is 0.")
		return defaultResp, errors.New("Ride Id is 0")
	}

	rideDetail, err := model.TukTuk.GetRideDetailsByRideId(ctx, rideId)
	if err != nil {
		log.Println("[GetDriverCurrentLocation][Error] Error in fetching ride data", err)
		return defaultResp, err
	}

	//it's check in case there is no ride of requested ride id.
	if rideId != rideDetail.Id {
		log.Println("[GetDriverCurrentLocation][Error] Invalid Ride id", rideId)
		return defaultResp, errors.New("Invalid Ride ID.")
	}

	log.Printf("[GetDriverCurrentLocation] Ride:%+v ", rideDetail)

	ddata, err := model.TukTuk.GetCustomerById(ctx, userId)
	if err != nil {
		log.Println("[GetDriverCurrentLocation][Error] Error in fetching ride data", err)
		return defaultResp, err
	}

	if ddata.CustomerId != rideDetail.CustomerId {
		log.Printf("[GetDriverCurrentLocation][Error] Customer ID mismatch in ride details. Found id:%d, required id:%d", ddata.CustomerId, rideDetail.CustomerId)
		return defaultResp, errors.New("Driver Id mismatch")
	}

	log.Printf("[GetDriverCurrentLocation] Customer data:%+v ", ddata)

	return ff.prepareDriverLocationResponse(ctx, rideDetail)
}

func (ff *FFClient) prepareDriverLocationResponse(ctx context.Context, ride model.RideDetailModel) (*DriverLocationResponse, error) {
	var (
		defaultResp *DriverLocationResponse
		err         error
	)

	data, err := model.TukTuk.GetDriverById(ctx, ride.DriverId)
	if err != nil {
		log.Println("[prepareDriverLocationResponse][Error] Error in fetching data", err)
		return defaultResp, errors.New("DB Error")
	}

	log.Printf("[GetDriverCurrentLocation] Driver Tracking data:%+v ", data)

	if data.DriverID != ride.DriverId {
		log.Printf("[prepareDriverLocationResponse][Error] Driver ID mismatch in ride details. Found id:%d, required id:%d", data.DriverID, ride.DriverId)
		return defaultResp, errors.New("Driver Id mismatch")
	}

	defaultResp = &DriverLocationResponse{
		CurrentLat:  data.CurrentLatitude,
		CurrentLong: data.CurrentLongitude,
		RideId:      ride.Id,
	}

	return defaultResp, err
}
