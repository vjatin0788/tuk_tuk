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
	)

	if userLat == 0 || userLong == 0 {
		return nil, errors.New("Empty lat or long")
	}

	driverData, err := ff.getAvailableDriverVehicle(ctx, userLat, userLong, vehicleType)
	if err != nil {
		return driverData, err
	}

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
			log.Printf("Vehicle type found:%s, required:%s for vehicle id:%d", vehicle.VehicleType, vehicleType, val.ID)
			if vehicleType == vehicle.VehicleType {
				driverModel = append(driverModel, val)
			}
		}
	}

	log.Printf("[getGetAvailableVehicles] Available Driver Data. :%+v , vehicle type:%s", driverModel, vehicleType)

	return driverModel
}

func (ff *FFClient) getAvailableDriverVehicle(ctx context.Context, userLat, userLong float64, vehicleType string) ([]model.DriverTrackingModel, error) {
	var (
		driverData []model.DriverTrackingModel
		err        error
	)

	data, err := model.TukTuk.GetAvailableDriver(ctx, lib.Rad(userLat), lib.Rad(userLong), DISTANCE, RADIUS)
	if err != nil {
		log.Println("[getAvailableDriverVehicle][Error] Error in fetching data", err)
		return driverData, err
	}

	dataMap, driverIds := getGetAvailableDriverMap(ctx, data)
	if len(driverIds) == 0 {
		log.Println("[getAvailableDriverVehicle][Error] No Driver available")
		return driverData, errors.New("No Driver Available")
	}

	log.Printf("[getAvailableDriverVehicle]Available Drivers:%+v", driverIds)

	vehicles, err := model.TukTuk.GetVehicleByAssignedDriver(ctx, driverIds)
	if err != nil {
		log.Println("[getAvailableDriverVehicle][Error] Error in fetching vehicle data", err)
		return driverData, err
	}

	log.Printf("[getAvailableDriverVehicle]Available Drivers after vehicles filter:%+v", vehicles)

	driverData = getGetAvailableVehicles(ctx, vehicles, dataMap, vehicleType)

	return driverData, err
}
