package fulfilment

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/TukTuk/errs"

	"github.com/TukTuk/payment"

	"github.com/TukTuk/common"
	"github.com/TukTuk/firebase"
	"github.com/TukTuk/lib"

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

	dataPush := PushNotification{
		Type: "ride_start",
		Data: PushNotificationRideStart{
			RideId:      rideDetail.Id,
			Message:     "Ride Started",
			PhoneNumber: ddata.Mobileno,
		},
	}

	ff.sendPushNotificationToCustomer(ctx, rideDetail, dataPush)

	defaultResp.Success = true
	defaultResp.CurrentLat = rideDetail.DestinationLat
	defaultResp.CurrentLong = rideDetail.DestinationLong
	defaultResp.RideId = rideDetail.Id

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

	return err
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

	return ff.prepareDriverLocationResponse(ctx, rideDetail, ddata)
}

func (ff *FFClient) prepareDriverLocationResponse(ctx context.Context, ride model.RideDetailModel, customerData model.CustomerModel) (*DriverLocationResponse, error) {
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

	ddata, err := model.TukTuk.GetDriverUserById(ctx, data.DriverID)
	if err != nil {
		log.Println("[prepareDriverLocationResponse][Error] DB error", data.DriverID)
	}

	if ddata.Userid != data.DriverID {
		log.Printf("[prepareDriverLocationResponse][Error] Invalid driver id. found:%d, req: %d", ddata.Userid, data.DriverID)
		return defaultResp, err
	}

	defaultResp = &DriverLocationResponse{
		CurrentLat:  data.CurrentLatitude,
		CurrentLong: data.CurrentLongitude,
		RideId:      ride.Id,
	}

	go ff.sendNotificationDriverArrived(ctx, &ride, data, customerData, ddata)

	return defaultResp, err
}

func (ff *FFClient) RideComplete(ctx context.Context, userId int64, rideReq RideCompleteRequest) (interface{}, error) {
	var (
		defaultResp *RideCompleteResponse
		err         error
	)

	if rideReq.RideId == 0 {
		log.Println("[RideComplete][Error] Error Ride Id is 0.")
		return defaultResp, errors.New("Ride Id is 0")
	}

	rideDetail, err := model.TukTuk.GetRideDetailsByRideId(ctx, rideReq.RideId)
	if err != nil {
		log.Println("[RideComplete][Error] Error in fetching ride data", err)
		return defaultResp, err
	}

	//it's check in case there is no ride of requested ride id.
	if rideReq.RideId != rideDetail.Id {
		log.Println("[RideComplete][Error] Invalid Ride id", rideDetail)
		return defaultResp, errors.New("Invalid Ride ID.")
	}

	log.Printf("[RideComplete] Ride:%+v ", rideDetail)

	return ff.prepareRideComplete(ctx, rideDetail, rideReq)
}

func (ff *FFClient) prepareRideComplete(ctx context.Context, ride model.RideDetailModel, rideReq RideCompleteRequest) (*RideCompleteResponse, error) {
	var (
		defaultResp *RideCompleteResponse
		err         error
	)

	data, err := model.TukTuk.GetDriverById(ctx, ride.DriverId)
	if err != nil {
		log.Println("[prepareRideComplete][Error] Error in fetching data", err)
		return defaultResp, errors.New("DB Error")
	}

	log.Printf("[prepareRideComplete] Driver Tracking data:%+v ", data)

	if data.DriverID != ride.DriverId {
		log.Printf("[prepareRideComplete][Error] Driver ID mismatch in ride details. Found id:%d, required id:%d", data.DriverID, ride.DriverId)
		return defaultResp, errors.New("Driver Id mismatch")
	}

	//state transition
	err = ff.RideStateTransition(ctx, &ride, common.RideStatus.COMPLETED.ID)
	if err != nil {
		log.Println("[prepareRideComplete][Error] Ride state Transitiion", err)
		return defaultResp, err
	}

	//update ride details
	ride.RideCompletedTime = time.Now().UTC().String()
	ride.DestinationLat = rideReq.DestinationLat
	ride.DestinationLong = rideReq.DestinationLong
	log.Printf("[prepareRideComplete] Ride Updated:%+v ", ride)
	rowAffectedCount, err := model.TukTuk.UpdateRideComplete(ctx, ride)
	if err != nil {
		log.Println("[prepareRideComplete][Error] Err in updating db", err)
		return defaultResp, err
	}

	if rowAffectedCount == 0 {
		log.Println("[prepareRideComplete][Error] Ride is not in valid state db,Row Affected:", rowAffectedCount)
		return defaultResp, errors.New("Something Went Wrong.")
	}

	//hit payment method api
	log.Printf("[RideComplete] Hitting payment api")
	var message string
	if strings.EqualFold(ride.PaymentMethod, common.CASH) {
		message = "COLLECT CASH"
	}

	amount, err := ff.initiatePayment(ctx, ride.Id)
	if err != nil {
		log.Println("[RideComplete][Error] Error in Payments:", err)
		return defaultResp, errs.Err("PA_RI_400")
	}

	dataPush := PushNotification{
		Type: "ride_complete",
		Data: PushNotificationRideComplete{
			RideId:  ride.Id,
			Message: "RIDE COMPLETE",
			Amount:  amount,
		},
	}

	ff.sendPushNotificationToCustomer(ctx, ride, dataPush)

	//Delete ride from map after completion
	delete(DriverBookedNotifiedMap, ride.Id)
	delete(RequestRideCancel, ride.Id)

	defaultResp = &RideCompleteResponse{
		Success: true,
		Message: message,
		Amount:  amount,
	}

	return defaultResp, err
}

func (ff *FFClient) sendPushNotificationToCustomer(ctx context.Context, ride model.RideDetailModel, data PushNotification) {
	fbclient := firebase.FClient

	ddata, err := model.TukTuk.GetCustomerById(ctx, ride.CustomerId)
	if err != nil {
		log.Println("[sendPushNotificationRideComplete][Error] Error in fetching ride data", err)
	}

	if ddata.CustomerId != ride.CustomerId {
		log.Printf("[sendPushNotificationRideComplete][Error] Customer ID mismatch in ride details. Found id:%d, required id:%d", ddata.CustomerId, ride.Id)
	}

	go fbclient.SendPushNotification(ctx, data, ddata.DeviceId)
}

func (ff *FFClient) DriverRideCancel(ctx context.Context, userId int64, rideReq RideCancelRequest) (interface{}, error) {
	var (
		defaultResp *RideCancelResponse
		err         error
	)

	if rideReq.RideId == 0 {
		log.Println("[DriverRideCancel][Error] Error Ride Id is 0.")
		return defaultResp, errors.New("Ride Id is 0")
	}

	rideDetail, err := model.TukTuk.GetRideDetailsByRideId(ctx, rideReq.RideId)
	if err != nil {
		log.Println("[DriverRideCancel][Error] Error in fetching ride data", err)
		return defaultResp, err
	}

	//it's check in case there is no ride of requested ride id.
	if rideReq.RideId != rideDetail.Id {
		log.Println("[DriverRideCancel][Error] Invalid Ride id", rideDetail)
		return defaultResp, errors.New("Invalid Ride ID.")
	}

	log.Printf("[DriverRideCancel] Ride:%+v ", rideDetail)

	return ff.prepareDriverRideCancelReq(ctx, rideDetail, rideReq)
}

func (ff *FFClient) prepareDriverRideCancelReq(ctx context.Context, ride model.RideDetailModel, rideReq RideCancelRequest) (*RideCancelResponse, error) {
	var (
		defaultResp *RideCancelResponse
		err         error
	)

	//state transition
	err = ff.RideStateTransition(ctx, &ride, common.RideStatus.FAILED.ID)
	if err != nil {
		log.Println("[prepareDriverRideCancelReq][Error] Ride state Transitiion", err)
		return defaultResp, err
	}

	//update ride details
	ride.RideFailedTime = time.Now().UTC().String()
	ride.DriverCancelled = 1
	ride.RideCancelReason = rideReq.Reason
	log.Printf("[prepareDriverRideCancelReq] Ride Updated:%+v ", ride)

	rowAffectedCount, err := model.TukTuk.UpdateRide(ctx, ride)
	if err != nil {
		log.Println("[prepareDriverRideCancelReq][Error] Err in updating db", err)
		return defaultResp, err
	}

	if rowAffectedCount == 0 {
		log.Println("[prepareDriverRideCancelReq][Error] Ride is not in valid state db,Row Affected:", rowAffectedCount)
		return defaultResp, errors.New("Something Went Wrong.")
	}

	//check other validations
	log.Printf("[prepareDriverRideCancelReq] check other validations")

	defaultResp = &RideCancelResponse{
		Success: true,
	}

	//add error in all sending notifications
	ff.sendPushNotificationDriverRideCancel(ctx, ride)

	return defaultResp, err
}

func (ff *FFClient) sendPushNotificationDriverRideCancel(ctx context.Context, ride model.RideDetailModel) {
	fbclient := firebase.FClient

	cdata, err := model.TukTuk.GetCustomerById(ctx, ride.CustomerId)
	if err != nil {
		log.Println("[sendPushNotificationDriverRideCancel][Error] Error in fetching ride data", err)
	}

	if cdata.CustomerId != ride.CustomerId {
		log.Printf("[sendPushNotificationDriverRideCancel][Error] Customer ID mismatch in ride details. Found id:%d, required id:%d", cdata.CustomerId, ride.CustomerId)
	}

	log.Printf("[sendPushNotificationRideCancel] customer data:%+v ", cdata)

	payLoad := PushNotification{
		Type: "ride_cancel",
		Data: PushNotificationRideCancel{
			RideId:  ride.Id,
			Message: "RIDE CANCEL",
		},
	}

	go fbclient.SendPushNotification(ctx, payLoad, cdata.DeviceId)
}

func (ff *FFClient) GetDriverRideStatus(ctx context.Context, id int64) (interface{}, error) {
	var (
		err        error
		defaultRes []RideBookResponse
	)

	ddata, err := model.TukTuk.GetDriverUserById(ctx, id)
	if err != nil {
		log.Println("[GetDriverRideStatus][Error] DB error", err)
	}

	if ddata.Userid != id {
		log.Printf("[GetDriverRideStatus][Error] Invalid driver id. found:%d, req: %d", ddata.Userid, id)
		return defaultRes, err
	}

	//status
	status := []int64{common.RideStatus.BOOKED.ID, common.RideStatus.PROCESSING.ID}

	rideData, err := model.TukTuk.GetRideDetailStatusByDriverId(ctx, id, status)
	if err != nil {
		log.Printf("[GetDriverRideStatus][Error] err", err)
		return defaultRes, err
	}

	log.Printf("[GetDriverRideStatus]Ride Details:%+v", rideData)

	for _, ride := range rideData {

		userData, err := model.TukTuk.GetCustomerById(ctx, ride.CustomerId)
		if err != nil {
			log.Println("[GetDriverRideStatus]DB err:", err)
			return defaultRes, err
		}

		defaultRes = append(defaultRes, RideBookResponse{
			CustomerDetail: &CustomerDetailsResponse{
				Name:        userData.Name,
				PhoneNumber: userData.Mobile,
			},
			SourceLat:       ride.SourceLat,
			SourceLong:      ride.SourceLong,
			DestinationLat:  ride.DestinationLat,
			DestinationLong: ride.DestinationLong,
			Status:          common.RideStatusMap[ride.Status].Label,
			RideId:          ride.Id,
		})

	}

	log.Printf("[GetDriverRideStatus]Ride resp:%+v", defaultRes)
	return defaultRes, err
}

func (ff *FFClient) DriverTracking(ctx context.Context, userLat, userLong float64, driverId int64, dateTime, locType string) (interface{}, error) {
	//logic comes here
	var (
		err    error
		driver model.DriverTrackingModel
	)

	defaultRes := DriverTrackingResponse{}

	if userLat == 0 || userLong == 0 {
		return nil, errors.New("Empty lat or long")
	}

	driverData, err := model.TukTuk.GetDriverUserById(ctx, driverId)
	if err != nil {
		log.Println("[updateTrackingDetails][Error] Error ", err)
		return nil, err
	}

	if driverData.Userid != driverId {
		log.Println("[updateTrackingDetails][Error] Driver ID mismatch ")
		return nil, errors.New("Driver ID Mismatch")
	}

	driver, err = model.TukTuk.GetDriverById(ctx, driverId)
	if err != nil {
		log.Println("[DriverTracking][Error] Error in fetching data", err)
		return nil, errors.New("Empty fetching data")
	}

	driverModel := model.DriverTrackingModel{
		DriverID:               driverId,
		CurrentLatitude:        userLat,
		CurrentLongitude:       userLong,
		CurrentLatitudeRadian:  lib.Rad(userLat),
		CurrentLongitudeRadian: lib.Rad(userLong),
	}

	if driver.DriverID == 0 {
		err = model.TukTuk.Create(ctx, driverModel)
		if err != nil {
			log.Println("[DriverTracking][Error] Error in inserting data ", err)
			return nil, errors.New("Empty inserting driver details")
		}
	} else {
		driverModel.LastLatitude = driver.CurrentLatitude
		driverModel.LastLongitude = driver.CurrentLongitude
		driverModel.LastLatitudeRadian = driver.CurrentLatitudeRadian
		driverModel.LastLongitudeRadian = driver.CurrentLongitudeRadian

		//we can add check if last and current location same than no need to update db

		err = model.TukTuk.Update(ctx, driverModel)
		if err != nil {
			log.Println("[DriverTracking][Error] Error in updating data ", err)
			return nil, errors.New("Empty updating driver details")
		}
	}

	err = ff.updateTrackingDetails(ctx, driverModel, dateTime, locType, driverData)
	if err != nil {
		log.Println("[DriverTracking][Error] Error in updating tracking data ", err)
		return nil, err
	}

	log.Printf("[DriverTracking]Driver id:%d ,Lat:%f, Long:%f", driver.DriverID, driver.CurrentLatitude, driver.CurrentLongitude)
	defaultRes.Success = true

	return defaultRes, err
}

func (ff *FFClient) updateTrackingDetails(ctx context.Context, drModel model.DriverTrackingModel, dateTime, locType string, driver model.DriverUserModel) error {
	var err error

	trackingModel := model.TrackingModel{
		Latitude:     drModel.CurrentLatitude,
		Longitutde:   drModel.CurrentLongitude,
		UserId:       driver.Userid,
		EmailId:      driver.Emailid,
		Date:         dateTime,
		TrackingType: locType,
	}

	_, err = model.TukTuk.CreateTracking(ctx, trackingModel)
	if err != nil {
		log.Println("[updateTrackingDetails][Error] Error inserting data ", err)
		return err
	}

	return err
}

func (ff *FFClient) initiatePayment(ctx context.Context, rideId int64) (float64, error) {
	var (
		defaultVal float64
		err        error
	)
	//add verifications
	resp, err := payment.PayClient.InitiatePaymentRequest(ctx, rideId)
	if err != nil {
		return defaultVal, err
	}

	defaultVal = resp.Data.TotalCost

	return defaultVal, err
}

func (ff *FFClient) sendNotificationDriverArrived(ctx context.Context, ride *model.RideDetailModel, driverTrackData model.DriverTrackingModel, customerData model.CustomerModel, driverData model.DriverUserModel) {

	fbclient := firebase.FClient

	var payLoad PushNotification
	if ff.liesInCustomerArea(ctx, ride.SourceLat, ride.SourceLong, ride.DestinationLat, ride.DestinationLong, ff.Cfg.Ride.DriverArrived) {
		payLoad = PushNotification{
			Type: "ride_driver_arrived",
			Data: PushNotificationDriverArrived{
				RideId:      ride.Id,
				Message:     "DRIVER ARRIVED",
				PhoneNumber: driverData.Mobileno,
			},
		}

		log.Println("[sendNotificationDriverArrived] Driver arrived")
	} else if ff.liesInCustomerArea(ctx, ride.SourceLat, ride.SourceLong, ride.DestinationLat, ride.DestinationLong, ff.Cfg.Ride.DriverArrival) {
		payLoad = PushNotification{
			Type: "ride_driver_arriving",
			Data: PushNotificationDriverArrived{
				RideId:      ride.Id,
				Message:     "DRIVER ARRIVING",
				PhoneNumber: driverData.Mobileno,
			},
		}

		log.Println("[sendNotificationDriverArrived] Driver arrival")
	} else {
		log.Println("[sendNotificationDriverArrived] Driver Not arrived yet")
		return
	}

	fbclient.SendPushNotification(ctx, payLoad, customerData.DeviceId)
}
