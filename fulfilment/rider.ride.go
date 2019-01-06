package fulfilment

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/TukTuk/common"
	"github.com/TukTuk/firebase"
	"github.com/TukTuk/model"
)

func (ff *FFClient) CustomerRideCancel(ctx context.Context, userId int64, rideReq RideCancelRequest) (interface{}, error) {
	var (
		defaultResp *RideCancelResponse
		err         error
	)

	if rideReq.RideId == 0 {
		log.Println("[CustomerRideCancel][Error] Error Ride Id is 0.")
		return defaultResp, errors.New("Ride Id is 0")
	}

	rideDetail, err := model.TukTuk.GetRideDetailsByRideId(ctx, rideReq.RideId)
	if err != nil {
		log.Println("[CustomerRideCancel][Error] Error in fetching ride data", err)
		return defaultResp, err
	}

	//it's check in case there is no ride of requested ride id.
	if rideReq.RideId != rideDetail.Id {
		log.Println("[CustomerRideCancel][Error] Invalid Ride id", rideDetail)
		return defaultResp, errors.New("Invalid Ride ID.")
	}

	log.Printf("[CustomerRideCancel] Ride:%+v ", rideDetail)

	return ff.prepareRideCancelReq(ctx, rideDetail, rideReq)
}

func (ff *FFClient) prepareRideCancelReq(ctx context.Context, ride model.RideDetailModel, rideReq RideCancelRequest) (*RideCancelResponse, error) {
	var (
		defaultResp *RideCancelResponse
		err         error
	)

	//state transition
	err = ff.RideStateTransition(ctx, &ride, common.RideStatus.FAILED.ID)
	if err != nil {
		log.Println("[prepareRideCancelReq][Error] Ride state Transitiion", err)
		return defaultResp, err
	}

	//update ride details
	ride.RideFailedTime = time.Now().UTC().String()
	ride.RiderCancelled = 1
	ride.RideCancelReason = rideReq.Reason
	log.Printf("[prepareRideCancelReq] Ride Updated:%+v ", ride)

	rowAffectedCount, err := model.TukTuk.UpdateRide(ctx, ride)
	if err != nil {
		log.Println("[prepareRideCancelReq][Error] Err in updating db", err)
		return defaultResp, err
	}

	if rowAffectedCount == 0 {
		log.Println("[prepareRideCancelReq][Error] Ride is not in valid state db,Row Affected:", rowAffectedCount)
		return defaultResp, errors.New("Something Went Wrong.")
	}

	//hit payment method api
	log.Printf("[prepareRideCancelReq] check other validations")

	defaultResp = &RideCancelResponse{
		Success: true,
	}

	ff.sendPushNotificationRideCancel(ctx, ride)

	return defaultResp, err
}

func (ff *FFClient) sendPushNotificationRideCancel(ctx context.Context, ride model.RideDetailModel) {
	fbclient := firebase.FClient

	data, err := model.TukTuk.GetDriverUserById(ctx, ride.DriverId)
	if err != nil {
		log.Println("[sendPushNotificationRideCancel][Error] Error in fetching data", err)
	}

	log.Printf("[sendPushNotificationRideCancel] Driver Tracking data:%+v ", data)

	if data.Userid != ride.DriverId {
		log.Printf("[sendPushNotificationRideCancel][Error] Driver ID mismatch in ride details. Found id:%d, required id:%d", data.Userid, ride.DriverId)
	}

	payLoad := PushNotification{
		Type: "ride_cancel",
		Data: PushNotificationRideCancel{
			RideId:  ride.Id,
			Message: "RIDE CANCEL",
		},
	}

	go fbclient.SendPushNotification(ctx, payLoad, data.DeviceId)
}
