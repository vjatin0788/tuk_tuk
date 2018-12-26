package model

import (
	"context"
	"fmt"
	"log"
	"time"
)

type DriverTrackingTable struct {
	ID                     int64     `db:"id" json:"id"`
	DriverID               int64     `db:"driver_id" json:"driver_id"`
	CurrentLatitude        float64   `db:"current_lat" json:"current_lat"`
	CurrentLongitude       float64   `db:"current_long" json:"current_long"`
	LastLatitude           float64   `db:"last_lat" json:"last_lat"`
	LastLongitude          float64   `db:"last_long" json:"last_long"`
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`
	CurrentLatitudeRadian  float64   `db:"current_lat_rad" json:"current_lat_rad"`
	CurrentLongitudeRadian float64   `db:"current_long_rad" json:"current_long_rad"`
	LastLatitudeRadian     float64   `db:"last_lat_rad" json:"last_lat_rad"`
	LastLongitudeRadian    float64   `db:"last_long_rad" json:"last_long_rad"`
}

type DriverTrackingModel struct {
	ID                     int64     `db:"id" json:"id"`
	DriverID               int64     `db:"driver_id" json:"driver_id"`
	CurrentLatitude        float64   `db:"current_lat" json:"current_lat"`
	CurrentLongitude       float64   `db:"current_long" json:"current_long"`
	LastLatitude           float64   `db:"last_lat" json:"last_lat"`
	LastLongitude          float64   `db:"last_long" json:"last_long"`
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`
	CurrentLatitudeRadian  float64   `db:"current_lat_rad" json:"current_lat_rad"`
	CurrentLongitudeRadian float64   `db:"current_long_rad" json:"current_long_rad"`
	LastLatitudeRadian     float64   `db:"last_lat_rad" json:"last_lat_rad"`
	LastLongitudeRadian    float64   `db:"last_long_rad" json:"last_long_rad"`
}

func (table DriverTrackingTable) GetModel() DriverTrackingModel {
	return DriverTrackingModel{
		ID:                     table.ID,
		DriverID:               table.DriverID,
		CurrentLatitude:        table.CurrentLatitude,
		CurrentLongitude:       table.CurrentLongitude,
		LastLatitude:           table.LastLatitude,
		LastLongitude:          table.LastLongitude,
		CreatedAt:              table.CreatedAt,
		UpdatedAt:              table.UpdatedAt,
		CurrentLatitudeRadian:  table.CurrentLatitudeRadian,
		CurrentLongitudeRadian: table.CurrentLongitudeRadian,
		LastLatitudeRadian:     table.LastLatitudeRadian,
		LastLongitudeRadian:    table.LastLongitudeRadian,
	}
}

func (model DriverTrackingModel) GetTable() DriverTrackingTable {
	return DriverTrackingTable{
		ID:                     model.ID,
		DriverID:               model.DriverID,
		CurrentLatitude:        model.CurrentLatitude,
		CurrentLongitude:       model.CurrentLongitude,
		LastLatitude:           model.LastLatitude,
		LastLongitude:          model.LastLongitude,
		CreatedAt:              model.CreatedAt,
		UpdatedAt:              model.UpdatedAt,
		CurrentLatitudeRadian:  model.CurrentLatitudeRadian,
		CurrentLongitudeRadian: model.CurrentLongitudeRadian,
		LastLatitudeRadian:     model.LastLatitudeRadian,
		LastLongitudeRadian:    model.LastLongitudeRadian,
	}
}

func (db *DBTuktuk) PrepareQueryForAvailableDriver(ctx context.Context, userLat, userLong float64, distance, radius float64) {
	getAvailableDrivers := "SELECT * FROM driver_tracking WHERE acos(sin(%f) * sin(current_lat_rad) + cos(%f) * cos(current_lat_rad) * cos(current_long_rad - (%f)))  <= %f"

	query := fmt.Sprintf(getAvailableDrivers, userLat, userLat, userLong, distance/radius)
	log.Println("[PrepareQueryForAvailableDriver] Query", query)

	statement.GetAvailableDrivers, _ = db.DBConnection.Preparex(query)
}

func (db *DBTuktuk) GetAvailableDriver(ctx context.Context, userLat, userLong float64, distance, radius float64) ([]DriverTrackingModel, error) {
	var (
		driverTable []DriverTrackingTable
		driverModel []DriverTrackingModel
		err         error
	)

	db.PrepareQueryForAvailableDriver(ctx, userLat, userLong, distance, radius)

	err = statement.GetAvailableDrivers.SelectContext(ctx, &driverTable)
	if err != nil {
		log.Println("[GetAvailableDriver][Error] Err in fetching data from db", err)
		return driverModel, err
	}

	for _, table := range driverTable {
		driverModel = append(driverModel, table.GetModel())
	}

	return driverModel, err
}

func (db *DBTuktuk) Create(ctx context.Context, driverModel DriverTrackingModel) error {
	//validations neeed to be inserted here
	err := driverModel.GetTable().InsertDriver(ctx)
	return err
}

func (table DriverTrackingTable) InsertDriver(ctx context.Context) error {

	var err error

	_, err = statement.InsertDriverData.ExecContext(ctx, table.DriverID, table.CurrentLatitude, table.CurrentLongitude, table.CurrentLatitudeRadian, table.CurrentLongitudeRadian)
	if err != nil {
		log.Println("[InsertDriver][Error] Err in inserting", err)
		return err
	}

	return err
}

func (db *DBTuktuk) Update(ctx context.Context, driverModel DriverTrackingModel) error {
	//validations neeed to be inserted here
	err := driverModel.GetTable().UpdateDriver(ctx)
	return err
}

func (table DriverTrackingTable) UpdateDriver(ctx context.Context) error {

	var err error

	_, err = statement.UpdateDriverData.ExecContext(ctx, table.CurrentLatitude, table.CurrentLongitude, table.CurrentLatitudeRadian, table.CurrentLongitudeRadian, table.LastLatitude, table.LastLongitude, table.LastLatitudeRadian, table.LastLongitudeRadian, table.DriverID)
	if err != nil {
		log.Println("[UpdateDriver][Error] Err in inserting", err)
		return err
	}

	return err
}

func (db *DBTuktuk) GetDriverById(ctx context.Context, driverId int64) (DriverTrackingModel, error) {
	var (
		driverTable DriverTrackingTable
		driverModel DriverTrackingModel
		err         error
	)

	//convert into slice
	row, err := statement.GetDriverById.QueryxContext(ctx, driverId)
	if err != nil {
		log.Println("[GetDriverById][Error] Err in fetching data from db", err)
		return driverModel, err
	}

	for row.Next() {
		err := row.StructScan(&driverTable)
		if err != nil {
			log.Println("[GetDriverById][Error] Err in scanning row", err)
			return driverModel, err
		}
	}

	return driverTable.GetModel(), err
}
