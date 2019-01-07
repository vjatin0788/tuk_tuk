package model

import (
	"context"
	"database/sql"
	"log"
)

type TrackingModel struct {
	TrackingId   int64   `db:"tracking_id" json:"tracking_id"`
	EmailId      string  `db:"emailid" json:"emailid"`
	Latitude     float64 `db:"lat" json:"lat"`
	Longitutde   float64 `db:"lng" json:"lng"`
	Date         string  `db:"datetime" json:"datetime"`
	CreatedOn    string  `db:"created_on" json:"created_on"`
	TrackingType string  `db:"tracking_type" json:"tracking_type"`
	UserId       int64   `db:"user_id" json:"user_id"`
}

type TrackingTable struct {
	TrackingId   int64           `db:"tracking_id" json:"tracking_id"`
	EmailId      sql.NullString  `db:"emailid" json:"emailid"`
	Latitude     sql.NullFloat64 `db:"lat" json:"lat"`
	Longitutde   sql.NullFloat64 `db:"lng" json:"lng"`
	Date         sql.NullString  `db:"datetime" json:"datetime"`
	CreatedOn    string          `db:"created_on" json:"created_on"`
	TrackingType sql.NullString  `db:"tracking_type" json:"tracking_type"`
	UserId       sql.NullInt64   `db:"user_id" json:"user_id"`
}

func (table TrackingTable) GetModel() TrackingModel {
	return TrackingModel{
		TrackingId:   table.TrackingId,
		EmailId:      table.EmailId.String,
		Latitude:     table.Latitude.Float64,
		Longitutde:   table.Longitutde.Float64,
		Date:         table.Date.String,
		TrackingType: table.TrackingType.String,
		UserId:       table.UserId.Int64,
	}
}

func (model TrackingModel) GetTable() TrackingTable {
	return TrackingTable{
		TrackingId:   model.TrackingId,
		EmailId:      sql.NullString{model.EmailId, false},
		Latitude:     sql.NullFloat64{model.Latitude, false},
		Longitutde:   sql.NullFloat64{model.Longitutde, false},
		Date:         sql.NullString{model.Date, false},
		TrackingType: sql.NullString{model.TrackingType, false},
		UserId:       sql.NullInt64{model.UserId, false},
	}
}

func (db *DBTuktuk) CreateTracking(ctx context.Context, model TrackingModel) (int64, error) {
	//validations neeed to be inserted here
	return model.GetTable().InsertTrackingData(ctx)
}

func (table TrackingTable) InsertTrackingData(ctx context.Context) (int64, error) {

	var (
		err       error
		defaultId int64
	)

	res, err := statement.InsertTrackingData.ExecContext(ctx, table.EmailId.String, table.Latitude.Float64, table.Longitutde.Float64, table.Date.String, table.TrackingType.String, table.UserId.Int64)
	if err != nil {
		log.Println("[InsertTrackingData][Error] Err in inserting", err)
		return defaultId, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("[InsertTrackingData][Error] Err in getting last id", err)
		return defaultId, err
	}

	defaultId = id

	return defaultId, err
}
