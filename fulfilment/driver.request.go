package fulfilment

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/TukTuk/common"

	"github.com/TukTuk/model"
)

func (ff *FFClient) DriverBooked(ctx context.Context, userId, rideId int64, isBooked bool) (interface{}, error) {

	var (
		err error
	)

	defaultRes := &DriverBookedResponse{
		Message: "No Ride Booked",
		RideId:  rideId,
	}

	if rideId == 0 {
		log.Println("[DriverBooked][Error] Error Ride Id is 0.")
		return nil, errors.New("Ride Id is 0")
	}

	rideDetail, err := model.TukTuk.GetRideDetailsByRideId(ctx, rideId)
	if err != nil {
		log.Println("[DriverBooked][Error] Error in fetching ride data", err)
		return nil, err
	}

	log.Printf("[DriverBooked] Ride Details:%+v", rideDetail)

	//it's check in case there is no ride of requested ride id.
	if rideId != rideDetail.Id {
		log.Println("[DriverBooked][Error] Invalid Ride id", rideId)
		return nil, errors.New("Invalid Ride ID.")
	}

	//check whether request is accepted or declined.
	if !isBooked {
		// marking driver Cancelled.
		//we can add more analysis about how many rides cancelled by driver.
		log.Println("[DriverBooked] Driver Declined the request for ride id:", rideId)
		return defaultRes, err
	}

	return ff.prepareRideForDriver(ctx, userId, &rideDetail)
}

func (ff *FFClient) prepareRideForDriver(ctx context.Context, userId int64, ride *model.RideDetailModel) (*DriverBookedResponse, error) {

	var err error

	if ride == nil {
		log.Println("[prepareRideForDriver][Error] Empty ride struct")
		return nil, errors.New("Something Went Wrong")
	}

	//check if driver already aloted to current ride or not.
	if ride.DriverId > 0 {
		if ride.Status < common.RideStatus.BOOKED.ID {
			log.Println("[prepareRideForDriver] Driver Found but status mismatch", ride.Status)
			return nil, errors.New("Something Went Wrong")
		}

		log.Println("[prepareRideForDriver] Driver Already aloted")
		return nil, errors.New("Something Went Wrong")
	}

	//check if driver is allocated to another ride
	rdData, err := model.TukTuk.GetRideDetailsByDriverId(ctx, userId)
	if err != nil {
		return nil, err
	}

	log.Printf("[prepareRideForDriver] Ride Details of driver id:%d , details:%+v", ride.DriverId, rdData)

	if rdData.Status >= common.RideStatus.BOOKED.ID && rdData.Status < common.RideStatus.COMPLETED.ID {
		log.Println("[prepareRideForDriver] Driver already booked with some other ride", rdData.Status)
		return nil, errors.New("Driver already on ride")
	}

	//ride state transition
	err = ff.RideStateTransition(ctx, ride, common.RideStatus.BOOKED.ID)
	if err != nil {
		log.Printf("[prepareRideForDriver][Error] Ride state Transitiion:%s, ride state:%d", err, ride.Status)
		return nil, err
	}

	//Alot driver for ride ID.
	rideUpdated, err := ff.alotDriverForRide(ctx, userId, ride)
	if err != nil {
		log.Println("[prepareRideForDriver][Error] Alot driver error", err)
		return nil, err
	}

	return ff.prepareDriverBookedResponse(ctx, rideUpdated)
}

func (ff *FFClient) prepareDriverBookedResponse(ctx context.Context, ride *model.RideDetailModel) (*DriverBookedResponse, error) {

	var (
		err      error
		bookResp *DriverBookedResponse
	)

	if ride == nil {
		log.Println("[prepareDriverBookedResponse][Error] Empty ride struct")
		return nil, errors.New("Something Went Wrong")
	}

	userData, err := model.TukTuk.GetCustomerById(ctx, ride.CustomerId)
	if err != nil {
		return nil, err
	}

	log.Printf("[prepareDriverBookedResponse]Customer Data found:%+v, id: %d", userData, ride.CustomerId)

	if userData.CustomerId != ride.CustomerId {
		log.Printf("[prepareDriverBookedResponse][Error] Invalid customer id. found:%d, req: %d", userData.CustomerId, ride.CustomerId)
		return nil, err
	}

	bookResp = &DriverBookedResponse{
		RiderDetail: CustomerDetailsResponse{
			Name:       userData.Name,
			CustomerId: userData.CustomerId,
		},
		CurrentLat:  ride.SourceLat,
		CurrentLong: ride.SourceLong,
		RideId:      ride.Id,
	}

	return bookResp, err
}

func (ff *FFClient) alotDriverForRide(ctx context.Context, userId int64, ride *model.RideDetailModel) (*model.RideDetailModel, error) {
	var (
		err error
		//rideModel model.RideDetailModel
	)

	if userId == 0 {
		log.Println("[alotDriverForRide][Error] Empty User ID")
		return nil, errors.New("Empty User ID")
	}

	ddata, err := model.TukTuk.GetDriverUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	log.Printf("[validateAndUpdateRideStatus]Driver data :%+v, driver id:%d", ddata, userId)

	if ddata.Userid != userId {
		log.Println("[validateAndUpdateRideStatus][Error] User ID Mismatch")
		return nil, errors.New("User ID Mismatch")
	}

	ride.DriverId = ddata.Userid
	ride.RideBookedTime = time.Now().UTC().String()

	rowAffectedCount, err := model.TukTuk.UpdateRideWithStatus(ctx, *ride)
	if err != nil {
		log.Println("[validateAndUpdateRideStatus][Error] Err in updating db", err)
		return nil, err
	}

	if rowAffectedCount == 0 {
		log.Println("[validateAndUpdateRideStatus][Error] Ride is not in valid state db,Row Affected:", rowAffectedCount)
		return nil, errors.New("Driver Already Booked or Something went wrong")
	}

	log.Printf("RIDE BOOKED FOR DRIVER:%d, RIDE ID:%d", userId, ride.Id)

	log.Printf("NOTIFYING RIDER. DriverBookedNotifiedMap map:%+v", DriverBookedNotifiedMap)
	if val, ok := DriverBookedNotifiedMap[ride.Id]; ok {
		val <- common.NOTIFY_RIDER
	} else {
		//Register in NSQ
		log.Println("[validateAndUpdateRideStatus][Error] Error in getting value from map.Unable to notify.")
	}

	return ride, err
}
