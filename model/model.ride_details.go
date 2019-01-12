package model

import (
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
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
	PaymentMethod     string  `db:"payment_method" json:"payment_method"`
	RideCancelReason  string  `db:"ride_cancel_msg" json:"ride_cancel_msg"`
	CustomerRating    int64   `db:"customer_rating" json:"customer_rating"`
	DriverRating      int64   `db:"driver_rating" json:"driver_rating"`
}

type RideDetailTabel struct {
	Id                int64          `db:"id" json:"id"`
	CustomerId        int64          `db:"customer_id" json:"customer_id"`
	DriverId          int64          `db:"driver_id" json:"driver_id"`
	SourceLat         float64        `db:"source_lat" json:"source_lat"`
	SourceLong        float64        `db:"source_long" json:"source_long"`
	DestinationLat    float64        `db:"destination_lat" json:"destination_lat"`
	DestinationLong   float64        `db:"destination_long" json:"destination_long"`
	Status            int64          `db:"status" json:"status"`
	CreatedAt         string         `db:"created_at" json:"created_at"`
	UpdatedAt         string         `db:"updated_at" json:"updated_at"`
	DriverCancelled   sql.NullInt64  `db:"driver_cancelled" json:"driver_cancelled"`
	RiderCancelled    sql.NullInt64  `db:"rider_cancelled" json:"rider_cancelled"`
	RideBookedTime    sql.NullString `db:"ride_booked_time" json:"ride_booked_time"`
	RideCompletedTime sql.NullString `db:"ride_completed_time" json:"ride_completed_time"`
	RideFailedTime    sql.NullString `db:"ride_failed_time" json:"ride_failed_time"`
	RideStartTime     sql.NullString `db:"ride_start_time" json:"ride_start_time"`
	PaymentMethod     sql.NullString `db:"payment_method" json:"payment_method"`
	RideCancelReason  sql.NullString `db:"ride_cancel_msg" json:"ride_cancel_msg"`
	CustomerRating    sql.NullInt64  `db:"customer_rating" json:"customer_rating"`
	DriverRating      sql.NullInt64  `db:"driver_rating" json:"driver_rating"`
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
		RideBookedTime:    table.RideBookedTime.String,
		RideCompletedTime: table.RideCompletedTime.String,
		RideFailedTime:    table.RideFailedTime.String,
		RideStartTime:     table.RideStartTime.String,
		PaymentMethod:     table.PaymentMethod.String,
		RideCancelReason:  table.RideCancelReason.String,
		CustomerRating:    table.CustomerRating.Int64,
		DriverRating:      table.DriverRating.Int64,
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
		RideBookedTime:    sql.NullString{model.RideBookedTime, false},
		RideCompletedTime: sql.NullString{model.RideCompletedTime, false},
		RideFailedTime:    sql.NullString{model.RideFailedTime, false},
		RideStartTime:     sql.NullString{model.RideStartTime, false},
		PaymentMethod:     sql.NullString{model.PaymentMethod, false},
		RideCancelReason:  sql.NullString{model.RideCancelReason, false},
		DriverRating:      sql.NullInt64{model.DriverRating, false},
		CustomerRating:    sql.NullInt64{model.CustomerRating, false},
	}
}

func (db *DBTuktuk) CreateRide(ctx context.Context, rideModel RideDetailModel) (int64, error) {
	//validations neeed to be inserted here
	return rideModel.GetTable().InsertRideDetails(ctx)
}

func (table RideDetailTabel) InsertRideDetails(ctx context.Context) (int64, error) {

	var err error

	res, err := statement.InsertRideDetails.Exec(table.CustomerId, table.SourceLat, table.SourceLong, table.DestinationLat, table.DestinationLong, table.Status, table.PaymentMethod.String)
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

//it should be used if you have all data that need to be updated.
func (db *DBTuktuk) UpdateRide(ctx context.Context, rideModel RideDetailModel) (int64, error) {
	//validations neeed to be inserted here
	return rideModel.GetTable().UpdateRideDetails(ctx)
}

func (table RideDetailTabel) UpdateRideDetails(ctx context.Context) (int64, error) {

	var (
		err      error
		rowCount int64
	)

	row, err := statement.UpdateRideDetails.Exec(table.Status,
		table.DriverCancelled.Int64, table.RiderCancelled.Int64,
		table.RideFailedTime.String, table.Id)
	if err != nil {
		log.Println("[UpdateRideDetails][Error] Err in inserting", err)
		return rowCount, err
	}

	rowsAffectedCount, err := row.RowsAffected()
	if err != nil {
		log.Println("[UpdateRideDetails][Error] Err in getting row affected count", err)
		return rowCount, err
	}

	rowCount = rowsAffectedCount

	return rowCount, err
}

func (db *DBTuktuk) UpdateRideWithStatus(ctx context.Context, rideModel RideDetailModel, status int64) (int64, error) {
	//validations neeed to be inserted here
	return rideModel.GetTable().updateRideDetailsAndStatus(ctx, status)
}

func (table RideDetailTabel) updateRideDetailsAndStatus(ctx context.Context, status int64) (int64, error) {

	var (
		err      error
		rowCount int64
	)

	row, err := statement.UpdateRideDetailsWithStatus.Exec(table.DriverId, table.Status,
		table.RideBookedTime.String, table.Id, status)
	if err != nil {
		log.Println("[UpdateRideDetails][Error] Err in inserting", err)
		return rowCount, err
	}

	rowsAffectedCount, err := row.RowsAffected()
	if err != nil {
		log.Println("[UpdateRideDetails][Error] Err in getting row affected count", err)
		return rowCount, err
	}

	rowCount = rowsAffectedCount

	return rowCount, err
}

func (db *DBTuktuk) GetRideDetailsByRideId(ctx context.Context, id int64) (RideDetailModel, error) {
	var (
		rideTable RideDetailTabel
		rideModel RideDetailModel
		err       error
	)

	//convert into slice
	err = statement.GetRideDetailsByRideId.Get(&rideTable, id)
	if err != nil && sql.ErrNoRows == nil {
		log.Println("[GetRideDetailsByRideId][Error] Err in fetching data from db", err)
		return rideModel, err
	}

	return rideTable.GetModel(), nil
}

func (db *DBTuktuk) GetRideDetailsByCustomerId(ctx context.Context, id int64) (RideDetailModel, error) {
	var (
		rideTable RideDetailTabel
		rideModel RideDetailModel
		err       error
	)

	//convert into slice
	err = statement.GetRideDetailsByCustomerID.Get(&rideTable, id)
	if err != nil && sql.ErrNoRows == nil {
		log.Println("[GetRideDetailsByCustomerId][Error] Err in fetching data from db", err)
		return rideModel, err
	}

	return rideTable.GetModel(), nil
}

func (db *DBTuktuk) GetRideDetailsByDriverId(ctx context.Context, driverId int64) (RideDetailModel, error) {
	var (
		rideTable RideDetailTabel
		rideModel RideDetailModel
		err       error
	)

	//convert into slice
	err = statement.GetRideDetailsByDriverID.Get(&rideTable, driverId)
	if err != nil && sql.ErrNoRows == nil {
		log.Println("[GetRideDetailsByDriverId][Error] Err in fetching data from db", err)
		return rideModel, err
	}

	return rideTable.GetModel(), nil
}

func (db *DBTuktuk) UpdateRideFail(ctx context.Context, rideModel RideDetailModel) error {
	//validations neeed to be inserted here
	err := rideModel.GetTable().updateRideDetailsFailedStatus(ctx)
	return err
}

func (table RideDetailTabel) updateRideDetailsFailedStatus(ctx context.Context) error {

	var err error

	_, err = statement.UpdateRideStatusFailed.Exec(table.Status,
		table.RideFailedTime.String, table.Id)
	if err != nil {
		log.Println("[UpdateRideDetails][Error] Err in inserting", err)
		return err
	}

	return err
}

func (db *DBTuktuk) UpdateRideStart(ctx context.Context, rideModel RideDetailModel) (int64, error) {
	//validations neeed to be inserted here
	return rideModel.GetTable().updateRideDetailsStart(ctx)
}

func (table RideDetailTabel) updateRideDetailsStart(ctx context.Context) (int64, error) {

	var (
		err      error
		rowCount int64
	)

	row, err := statement.UpdateRideStart.Exec(table.Status,
		table.RideStartTime.String, table.Id)
	if err != nil {
		log.Println("[UpdateRideDetailsStart][Error] Err in inserting", err)
		return rowCount, err
	}

	rowsAffectedCount, err := row.RowsAffected()
	if err != nil {
		log.Println("[UpdateRideDetailsStart][Error] Err in getting row affected count", err)
		return rowCount, err
	}

	rowCount = rowsAffectedCount

	return rowCount, err
}

func (db *DBTuktuk) UpdateRideStatus(ctx context.Context, rideModel RideDetailModel) (int64, error) {
	//validations neeed to be inserted here
	return rideModel.GetTable().updateRideDetailsStatus(ctx)
}

func (table RideDetailTabel) updateRideDetailsStatus(ctx context.Context) (int64, error) {

	var (
		err      error
		rowCount int64
	)

	row, err := statement.UpdateRideStatus.Exec(table.DriverId, table.Status, table.Id)
	if err != nil {
		log.Println("[UpdateRideDetailsStatus][Error] Err in inserting", err)
		return rowCount, err
	}

	rowsAffectedCount, err := row.RowsAffected()
	if err != nil {
		log.Println("[UpdateRideDetailsStatus][Error] Err in getting row affected count", err)
		return rowCount, err
	}

	rowCount = rowsAffectedCount

	return rowCount, err
}

func (db *DBTuktuk) UpdateRideComplete(ctx context.Context, rideModel RideDetailModel) (int64, error) {
	//validations neeed to be inserted here
	return rideModel.GetTable().updateRideDetailsComplete(ctx)
}

func (table RideDetailTabel) updateRideDetailsComplete(ctx context.Context) (int64, error) {

	var (
		err      error
		rowCount int64
	)

	row, err := statement.UpdateRideComplete.Exec(table.Status, table.RideCompletedTime.String, table.DestinationLat, table.DestinationLong, table.Id)
	if err != nil {
		log.Println("[updateRideDetailsComplete][Error] Err in inserting", err)
		return rowCount, err
	}

	rowsAffectedCount, err := row.RowsAffected()
	if err != nil {
		log.Println("[updateRideDetailsComplete][Error] Err in getting row affected count", err)
		return rowCount, err
	}

	rowCount = rowsAffectedCount

	return rowCount, err
}

func (db *DBTuktuk) GetRideDetailStatusByCustomerId(ctx context.Context, id int64, status []int64) ([]RideDetailModel, error) {
	var (
		rideTable []RideDetailTabel
		rideModel []RideDetailModel
		err       error
	)

	//convert into slice
	query, args, err := sqlx.In(getRideDetailsStatusByCustomerID, id, status)
	if err != nil {
		log.Println("[GetRideDetailStatusByCustomerId][Error] Err in sqlx IN", err)
		return rideModel, err
	}

	query = db.DBConnection.Rebind(query)
	err = db.DBConnection.Select(&rideTable, query, args...)
	if err != nil {
		log.Println("[GetRideDetailStatusByCustomerId][Error] Err in fetching data from db", err)
		return rideModel, err
	}

	for _, table := range rideTable {
		rideModel = append(rideModel, table.GetModel())
	}

	return rideModel, err
}

func (db *DBTuktuk) GetRideDetailStatusByDriverId(ctx context.Context, id int64, status []int64) ([]RideDetailModel, error) {
	var (
		rideTable []RideDetailTabel
		rideModel []RideDetailModel
		err       error
	)

	log.Printf("[GetRideDetailStatusByDriverId]Ride Details id:%d,status:%+v", id, status)
	//convert into slice
	query, args, err := sqlx.In(getRideDetailsStatusByDriverID, id, status)
	if err != nil {
		log.Println("[GetRideDetailStatusByDriverId][Error] Err in sqlx IN", err)
		return rideModel, err
	}

	query = db.DBConnection.Rebind(query)
	err = db.DBConnection.Select(&rideTable, query, args...)
	if err != nil {
		log.Println("[GetRideDetailStatusByDriverId][Error] Err in fetching data from db", err)
		return rideModel, err
	}

	for _, table := range rideTable {
		rideModel = append(rideModel, table.GetModel())
	}

	return rideModel, err
}

func (db *DBTuktuk) GetRideDetailsByCustomerIdAndStatus(ctx context.Context, id, status int64) (RideDetailModel, error) {
	var (
		rideTable RideDetailTabel
		rideModel RideDetailModel
		err       error
	)

	//convert into slice
	err = statement.GetRideDetailsByCustomerIDAndStatus.Get(&rideTable, id, status)
	if err != nil && sql.ErrNoRows == nil {
		log.Println("[GetRideDetailsByCustomerIdAndStatus][Error] Err in fetching data from db", err)
		return rideModel, err
	}

	return rideTable.GetModel(), nil
}
