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
		rideDetail  model.RideDetailModel
	)

	defaultResp = &RideCancelResponse{}

	if !rideReq.RideRequestCancel {
		if rideReq.RideId == 0 {
			log.Println("[CustomerRideCancel][Error] Error Ride Id is 0.")
			return defaultResp, errors.New("Ride Id is 0")
		}

		rideDetail, err = model.TukTuk.GetRideDetailsByRideId(ctx, rideReq.RideId)
		if err != nil {
			log.Println("[CustomerRideCancel][Error] Error in fetching ride data", err)
			return defaultResp, err
		}

		//it's check in case there is no ride of requested ride id.
		if rideReq.RideId != rideDetail.Id {
			log.Println("[CustomerRideCancel][Error] Invalid Ride id", rideDetail)
			return defaultResp, errors.New("Invalid Ride ID.")
		}

	} else {
		rideDetail, err = model.TukTuk.GetRideDetailsByCustomerIdAndStatus(ctx, userId, common.RideStatus.REQUESTED.ID)
		if err != nil {
			log.Println("[CustomerRideCancel][Error] Error in fetching ride data", err)
			return defaultResp, err
		}

		if userId != rideDetail.CustomerId {
			log.Println("[CustomerRideCancel][Error] Invalid Customer id", rideDetail)
			return defaultResp, errors.New("Invalid Customer ID.")
		}

	}

	//notify if ride is in requested state.
	if rideDetail.Status <= common.RideStatus.REQUESTED.ID {
		log.Printf("NOTIFYING RIDER. Request ride cancel map:%+v", RequestRideCancel)
		if val, ok := RequestRideCancel[rideDetail.Id]; ok {
			val <- common.NOTIFY_RIDER
		} else {
			//Register in NSQ
			log.Println("[CustomerRideCancel][Error] Error in getting value from map.Unable to notify.")
			return nil, errors.New("Unable to notify.")
		}
		defaultResp.Success = true
		return defaultResp, err
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

func (ff *FFClient) GetCustomerRideStatus(ctx context.Context, custId int64) (interface{}, error) {
	var (
		err        error
		defaultRes []RideBookResponse
	)

	userData, err := model.TukTuk.GetCustomerById(ctx, custId)
	if err != nil {
		log.Println("[GetCustomerRideStatus]DB err:", err)
		return defaultRes, err
	}

	if userData.CustomerId != custId {
		log.Printf("[GetCustomerRideStatus][Error] Invalid customer id. found:%d, req: %d", userData.CustomerId, custId)
		return defaultRes, err
	}

	//status
	status := []int64{common.RideStatus.REQUESTED.ID, common.RideStatus.BOOKED.ID, common.RideStatus.PROCESSING.ID}

	rideData, err := model.TukTuk.GetRideDetailStatusByCustomerId(ctx, custId, status)
	if err != nil {
		log.Printf("[GetCustomerRideStatus][Error] err", err)
		return defaultRes, err
	}

	for idx, ride := range rideData {
		ddata, err := model.TukTuk.GetDriverUserById(ctx, ride.DriverId)
		if err != nil {
			log.Println("[GetCustomerRideStatus][Error] DB error", err)
		}

		idForVehicle := []int64{ride.DriverId}
		vehicles, err := model.TukTuk.GetVehicleByAssignedDriver(ctx, idForVehicle)
		if err != nil {
			log.Println("[GetCustomerRideStatus][Error] Error in fetching vehicle data", err)
			return nil, err
		}

		driver, err := model.TukTuk.GetDriverById(ctx, ride.DriverId)
		if err != nil {
			log.Println("[GetCustomerRideStatus][Error] Error in fetching data", err)
			return nil, errors.New("Empty fetching data")
		}

		defaultRes = append(defaultRes, RideBookResponse{
			DriverDetail: &DriverDetailsResponse{
				DriverId:    ddata.Userid,
				Name:        ddata.Name,
				CurrentLat:  driver.CurrentLatitude,
				CurrentLong: driver.CurrentLongitude,
				DriverImage: ddata.Driverpic,
				PhoneNumber: ddata.Mobileno,
			},
			SourceLat:       ride.SourceLat,
			SourceLong:      ride.SourceLong,
			DestinationLat:  ride.DestinationLat,
			DestinationLong: ride.DestinationLong,
			Status:          common.RideStatusMap[ride.Status].Label,
			RideId:          ride.Id,
		})

		if len(vehicles) != 0 {
			defaultRes[idx].DriverDetail.Model = vehicles[0].Model
			defaultRes[idx].DriverDetail.VehicleNumber = vehicles[0].VehicleNumber
			defaultRes[idx].DriverDetail.VehicleType = vehicles[0].VehicleType
		}
	}

	return defaultRes, err
}
