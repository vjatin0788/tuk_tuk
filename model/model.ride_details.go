package model

import (
	"context"
	"database/sql"
	"log"
)

type RideDetailModel struct {
	Id                int64   `db:"id" json:"id"`
	CustomerId        int64   `db:"customer_id" json:"customer_id"`
	DriverId          int64   `db:"driver_id" json:"driver_id"`
	SourceLat         float64 `db:"source_lat" json:"source_lat"`
	SourceLong        float64 `db:"source_long" json:"source_long"`
	DestinationLat    float64 `db:"destination_lat" json:"destination_lat"`
	DestinationLong   float64 `db:"destination_long" json:"destination_long"`
	Status            int64   `db:"status" json:"status"`
	CreatedAt         string  `db:"created_at" json:"created_at"`
	UpdatedAt         string  `db:"updated_at" json:"updated_at"`
	DriverCancelled   int64   `db:"driver_cancelled" json:"driver_cancelled"`
	RiderCancelled    int64   `db:"rider_cancelled" json:"rider_cancelled"`
	RideBookedTime    string  `db:"ride_booked_time" json:"ride_booked_time"`
	RideCompletedTime string  `db:"ride_completed_time" json:"ride_completed_time"`
	RideFailedTime    string  `db:"ride_failed_time" json:"ride_failed_time"`
	RideStartTime     string  `db:"ride_start_time" json:"ride_start_time"`
}

type RideDetailTabel struct {
	Id                int64         `db:"id" json:"id"`
	CustomerId        int64         `db:"customer_id" json:"customer_id"`
	DriverId          int64         `db:"driver_id" json:"driver_id"`
	SourceLat         float64       `db:"source_lat" json:"source_lat"`
	SourceLong        float64       `db:"source_long" json:"source_long"`
	DestinationLat    float64       `db:"destination_lat" json:"destination_lat"`
	DestinationLong   float64       `db:"destination_long" json:"destination_long"`
	Status            int64         `db:"status" json:"status"`
	CreatedAt         string        `db:"created_at" json:"created_at"`
	UpdatedAt         string        `db:"updated_at" json:"updated_at"`
	DriverCancelled   sql.NullInt64 `db:"driver_cancelled" json:"driver_cancelled"`
	RiderCancelled    sql.NullInt64 `db:"rider_cancelled" json:"rider_cancelled"`
	RideBookedTime    string        `db:"ride_booked_time" json:"ride_booked_time"`
	RideCompletedTime string        `db:"ride_completed_time" json:"ride_completed_time"`
	RideFailedTime    string        `db:"ride_failed_time" json:"ride_failed_time"`
	RideStartTime     string        `db:"ride_start_time" json:"ride_start_time"`
}

func (table RideDetailTabel) GetModel() RideDetailModel {
	return RideDetailModel{
		Id:                table.Id,
		CustomerId:        table.CustomerId,
		DriverId:          table.DriverId,
		SourceLat:         table.SourceLat,
		SourceLong:        table.SourceLong,
		DestinationLat:    table.DestinationLat,
		DestinationLong:   table.DestinationLong,
		Status:            table.Status,
		CreatedAt:         table.CreatedAt,
		UpdatedAt:         table.UpdatedAt,
		DriverCancelled:   table.DriverCancelled.Int64,
		RiderCancelled:    table.RiderCancelled.Int64,
		RideBookedTime:    table.RideBookedTime,
		RideCompletedTime: table.RideCompletedTime,
		RideFailedTime:    table.RideFailedTime,
	}
}

func (model RideDetailModel) GetTable() RideDetailTabel {
	return RideDetailTabel{
		Id:                model.Id,
		CustomerId:        model.CustomerId,
		DriverId:          model.DriverId,
		SourceLat:         model.SourceLat,
		SourceLong:        model.SourceLong,
		DestinationLat:    model.DestinationLat,
		DestinationLong:   model.DestinationLong,
		Status:            model.Status,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
		DriverCancelled:   sql.NullInt64{model.DriverCancelled, false},
		RiderCancelled:    sql.NullInt64{model.RiderCancelled, false},
		RideBookedTime:    model.RideBookedTime,
		RideCompletedTime: model.RideCompletedTime,
		RideFailedTime:    model.RideFailedTime,
	}
}

func (db *DBTuktuk) CreateRide(ctx context.Context, rideModel RideDetailModel) (int64, error) {
	//validations neeed to be inserted here
	return rideModel.GetTable().InsertRideDetails(ctx)
}

func (table RideDetailTabel) InsertRideDetails(ctx context.Context) (int64, error) {

	var err error

	res, err := statement.InsertRideDetails.ExecContext(ctx, table.CustomerId, table.SourceLat, table.SourceLong, table.DestinationLat, table.DestinationLong, table.Status)
	if err != nil {
		log.Println("[InsertRideDetails][Error] Err in inserting", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("[InsertRideDetails][Error] Err in getting last id", err)
		return 0, err
	}

	return id, err
}

func (db *DBTuktuk) UpdateRide(ctx context.Context, rideModel RideDetailModel) error {
	//validations neeed to be inserted here
	err := rideModel.GetTable().UpdateRideDetails(ctx)
	return err
}

func (table RideDetailTabel) UpdateRideDetails(ctx context.Context) error {

	var err error

	_, err = statement.UpdateRideDetails.ExecContext(ctx, table.DriverId, table.Status,
		table.DriverCancelled, table.RiderCancelled, table.RideBookedTime, table.RideCompletedTime,
		table.RideFailedTime, table.RideStartTime, table.Id)
	if err != nil {
		log.Println("[UpdateRideDetails][Error] Err in inserting", err)
		return err
	}

	return err
}

func (db *DBTuktuk) GetRideDetailsByRideId(ctx context.Context, id int64) (RideDetailModel, error) {
	var (
		rideTable RideDetailTabel
		rideModel RideDetailModel
		err       error
	)

	//convert into slice
	row, err := statement.GetRideDetailsByRideId.QueryxContext(ctx, id)
	if err != nil {
		log.Println("[GetRideDetailsByRideId][Error] Err in fetching data from db", err)
		return rideModel, err
	}

	for row.Next() {
		err := row.StructScan(&rideTable)
		if err != nil {
			log.Println("[GetRideDetailsByRideId][Error] Err in scanning row", err)
			return rideModel, err
		}
	}

	return rideTable.GetModel(), err
}

func (db *DBTuktuk) GetRideDetailsByCustomerId(ctx context.Context, id int64) (RideDetailModel, error) {
	var (
		rideTable RideDetailTabel
		rideModel RideDetailModel
		err       error
	)

	//convert into slice
	row, err := statement.GetRideDetailsByCustomerID.QueryxContext(ctx, id)
	if err != nil {
		log.Println("[GetRideDetailsByCustomerId][Error] Err in fetching data from db", err)
		return rideModel, err
	}

	for row.Next() {
		err := row.StructScan(&rideTable)
		if err != nil {
			log.Println("[GetRideDetailsByCustomerId][Error] Err in scanning row", err)
			return rideModel, err
		}
	}

	return rideTable.GetModel(), err
}
