package fulfilment

import (
	"context"
	"errors"
	"log"

	"github.com/TukTuk/lib"

	"github.com/TukTuk/model"
)

const DISTANCE float64 = 10
const RADIUS float64 = 6371.01

func (ff *FFClient) DriverAvailable(ctx context.Context, userLat, userLong float64, vehicleType string) (interface{}, error) {
	//logic comes here
	var (
		response []RiderAvailableResponse
		err      error
		data     []model.DriverTrackingModel
	)

	if userLat == 0 || userLong == 0 {
		return nil, errors.New("Empty lat or long")
	}

	data, err = model.TukTuk.GetAvailableDriver(ctx, lib.Rad(userLat), lib.Rad(userLong), DISTANCE, RADIUS)
	if err != nil {
		log.Println("[RiderAvailable][Error] Error in fetching data", err)
		return nil, err
	}

	dataMap, driverIds := getGetAvailableDriverMap(ctx, data)

	vehicles, err := model.TukTuk.GetVehicleByAssignedDriver(ctx, driverIds)
	if err != nil {
		log.Println("[RiderAvailable][Error] Error in fetching vehicle data", err)
		return nil, err
	}

	driverData := getGetAvailableVehicles(ctx, vehicles, dataMap, vehicleType)

	for _, driver := range driverData {
		response = append(response, RiderAvailableResponse{
			DriverDetails: DriverDetailsResponse{
				DriverId: driver.DriverID,
			},
			CurrentLatitude:  driver.CurrentLatitude,
			CurrentLongitude: driver.CurrentLongitude,
		})
	}

	return response, err
}

func verifyLocation(userLat, userLong float64) bool {
	//some validation needs to be added here.
	return false
}

func (ff *FFClient) DriverTracking(ctx context.Context, userLat, userLong float64, driverId int64) (interface{}, error) {
	//logic comes here
	var (
		err    error
		driver model.DriverTrackingModel
	)

	defaultRes := DriverTrackingResponse{}

	if userLat == 0 || userLong == 0 {
		return nil, errors.New("Empty lat or long")
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

	log.Println("driver id:", driver.DriverID)
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

	defaultRes.Success = true

	return defaultRes, err
}

func getGetAvailableDriverMap(ctx context.Context, driverModel []model.DriverTrackingModel) (map[int64]model.DriverTrackingModel, []int64) {

	driverMap := make(map[int64]model.DriverTrackingModel)
	driverIds := make([]int64, 0)

	for _, driver := range driverModel {
		driverMap[driver.DriverID] = driver
		driverIds = append(driverIds, driver.DriverID)
	}

	return driverMap, driverIds
}

func getGetAvailableVehicles(ctx context.Context, vehicles []model.VehicleModel, driverMap map[int64]model.DriverTrackingModel, vehicleType string) []model.DriverTrackingModel {

	driverModel := make([]model.DriverTrackingModel, 0)

	for _, vehicle := range vehicles {
		if val, ok := driverMap[vehicle.AssignedDriverId]; ok {
			if vehicleType == vehicle.VehicleType {
				driverModel = append(driverModel, val)
			}
		}
	}

	return driverModel
}
