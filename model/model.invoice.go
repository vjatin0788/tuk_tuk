package model

import (
	"context"
	"database/sql"
	"log"
)

type InvoiceModel struct {
	InvoiceId      int64   `db:"invoice_id" json:"invoice_id"`
	CustomerId     int64   `db:"customer_id" json:"customer_id"`
	DriverId       int64   `db:"driver_id" json:"driver_id"`
	SourceLat      float64 `db:"source_lat" json:"source_lat"`
	SourceLng      float64 `db:"source_lng" json:"source_lng"`
	SourceAddress  string  `db:"source_address" json:"source_address"`
	DestinationLat float64 `db:"destination_lat" json:"destination_lat"`
	DestinationLng float64 `db:"destination_lng" json:"destination_lng"`
	SourceTime     string  `db:"source_time" json:"source_time"`
	TotalMinutes   int64   `db:"total_minutes" json:"total_minutes"`
	CostPerMinute  float64 `db:"cost_per_minute" json:"cost_per_minute"`
	TimeCost       float64 `db:"time_cost" json:"time_cost"`
	Distance       float64 `db:"distance" json:"distance"`
	CostPerKm      float64 `db:"cost_per_km" json:"cost_per_km"`
	DistanceCost   float64 `db:"distance_cost" json:"distance_cost"`
	BaseFare       float64 `db:"base_fare" json:"base_fare"`
	ExtraCharges   float64 `db:"extra_charges" json:"extra_charges"`
	Discount       float64 `db:"discount" json:"discount"`
	TotalCost      float64 `db:"total_cost" json:"total_cost"`
	GstPercentage  float64 `db:"gst_percentage" json:"gst_percentage"`
	Gst            float64 `db:"gst" json:"gst"`
	FinalCost      float64 `db:"final_cost" json:"final_cost"`
	CreatedOn      string  `db:"created_on" json:"created_on"`
}

type InvoiceTable struct {
	InvoiceId      int64           `db:"invoice_id" json:"invoice_id"`
	CustomerId     sql.NullInt64   `db:"customer_id" json:"customer_id"`
	DriverId       sql.NullInt64   `db:"driver_id" json:"driver_id"`
	SourceLat      sql.NullFloat64 `db:"source_lat" json:"source_lat"`
	SourceLng      sql.NullFloat64 `db:"source_lng" json:"source_lng"`
	SourceAddress  sql.NullString  `db:"source_address" json:"source_address"`
	DestinationLat sql.NullFloat64 `db:"destination_lat" json:"destination_lat"`
	DestinationLng sql.NullFloat64 `db:"destination_lng" json:"destination_lng"`
	SourceTime     sql.NullString  `db:"source_time" json:"source_time"`
	TotalMinutes   sql.NullInt64   `db:"total_minutes" json:"total_minutes"`
	CostPerMinute  sql.NullFloat64 `db:"cost_per_minute" json:"cost_per_minute"`
	TimeCost       sql.NullFloat64 `db:"time_cost" json:"time_cost"`
	Distance       sql.NullFloat64 `db:"distance" json:"distance"`
	CostPerKm      sql.NullFloat64 `db:"cost_per_km" json:"cost_per_km"`
	DistanceCost   sql.NullFloat64 `db:"distance_cost" json:"distance_cost"`
	BaseFare       sql.NullFloat64 `db:"base_fare" json:"base_fare"`
	ExtraCharges   sql.NullFloat64 `db:"extra_charges" json:"extra_charges"`
	Discount       sql.NullFloat64 `db:"discount" json:"discount"`
	TotalCost      sql.NullFloat64 `db:"total_cost" json:"total_cost"`
	GstPercentage  sql.NullFloat64 `db:"gst_percentage" json:"gst_percentage"`
	Gst            sql.NullFloat64 `db:"gst" json:"gst"`
	FinalCost      sql.NullFloat64 `db:"final_cost" json:"final_cost"`
	CreatedOn      sql.NullString  `db:"created_on" json:"created_on"`
}

func (table InvoiceTable) GetModel() InvoiceModel {
	return InvoiceModel{
		InvoiceId:      table.InvoiceId,
		CustomerId:     table.CustomerId.Int64,
		DriverId:       table.DriverId.Int64,
		SourceLat:      table.SourceLat.Float64,
		SourceLng:      table.SourceLng.Float64,
		SourceAddress:  table.SourceAddress.String,
		DestinationLat: table.DestinationLat.Float64,
		DestinationLng: table.DestinationLng.Float64,
		SourceTime:     table.SourceTime.String,
		TotalMinutes:   table.TotalMinutes.Int64,
		CostPerMinute:  table.CostPerMinute.Float64,
		TimeCost:       table.TimeCost.Float64,
		Distance:       table.Distance.Float64,
		CostPerKm:      table.CostPerKm.Float64,
		DistanceCost:   table.DistanceCost.Float64,
		BaseFare:       table.BaseFare.Float64,
		ExtraCharges:   table.ExtraCharges.Float64,
		Discount:       table.Discount.Float64,
		TotalCost:      table.TotalCost.Float64,
		GstPercentage:  table.GstPercentage.Float64,
		Gst:            table.Gst.Float64,
		FinalCost:      table.FinalCost.Float64,
		CreatedOn:      table.CreatedOn.String,
	}
}

func (model InvoiceModel) GetTable() InvoiceTable {
	return InvoiceTable{
		InvoiceId:      model.InvoiceId,
		CustomerId:     sql.NullInt64{model.CustomerId, false},
		DriverId:       sql.NullInt64{model.DriverId, false},
		SourceLat:      sql.NullFloat64{model.SourceLat, false},
		SourceLng:      sql.NullFloat64{model.SourceLng, false},
		SourceAddress:  sql.NullString{model.SourceAddress, false},
		DestinationLat: sql.NullFloat64{model.DestinationLat, false},
		DestinationLng: sql.NullFloat64{model.DestinationLng, false},
		SourceTime:     sql.NullString{model.SourceTime, false},
		TotalMinutes:   sql.NullInt64{model.TotalMinutes, false},
		CostPerMinute:  sql.NullFloat64{model.CostPerMinute, false},
		TimeCost:       sql.NullFloat64{model.TimeCost, false},
		Distance:       sql.NullFloat64{model.Distance, false},
		CostPerKm:      sql.NullFloat64{model.CostPerKm, false},
		DistanceCost:   sql.NullFloat64{model.DistanceCost, false},
		BaseFare:       sql.NullFloat64{model.BaseFare, false},
		ExtraCharges:   sql.NullFloat64{model.ExtraCharges, false},
		Discount:       sql.NullFloat64{model.Discount, false},
		TotalCost:      sql.NullFloat64{model.TotalCost, false},
		GstPercentage:  sql.NullFloat64{model.GstPercentage, false},
		Gst:            sql.NullFloat64{model.Gst, false},
		FinalCost:      sql.NullFloat64{model.FinalCost, false},
		CreatedOn:      sql.NullString{model.CreatedOn, false},
	}
}

func (db *DBTuktuk) CreateInvoice(ctx context.Context, invoiceModel InvoiceModel) error {
	//validations neeed to be inserted here
	err := invoiceModel.GetTable().InsertInvoiceByCustomerID(ctx)
	return err
}

func (table InvoiceTable) InsertInvoiceByCustomerID(ctx context.Context) error {

	var err error

	_, err = statement.InsertInvoiceByCustomerId.Exec(table.CustomerId.Int64, table.DriverId.Int64,
		table.SourceLat.Float64, table.SourceLng.Float64, table.SourceAddress.String, table.DestinationLat.Float64,
		table.DestinationLng.Float64, table.SourceTime.String, table.TotalMinutes.Int64, table.CostPerMinute.Float64,
		table.TimeCost.Float64, table.Distance.Float64, table.CostPerKm.Float64, table.DistanceCost.Float64, table.BaseFare.Float64,
		table.ExtraCharges.Float64, table.Discount.Float64, table.TotalCost.Float64, table.GstPercentage.Float64,
		table.Gst.Float64, table.FinalCost.Float64)
	if err != nil {
		log.Println("[InsertInvoiceByCustomerID][Error] Err in inserting", err)
		return err
	}

	return err
}

func (db *DBTuktuk) GetInvoiceByCustomerId(ctx context.Context, custId int64) (InvoiceModel, error) {
	var (
		invoiceTabel InvoiceTable
		invoiceModel InvoiceModel
		err          error
	)

	err = statement.GetInvoiceByCustomerId.Get(&invoiceTabel, custId)
	if err != nil && sql.ErrNoRows == nil {
		log.Println("[GetAvailableDriver][Error] Err in fetching data from db", err)
		return invoiceModel, err
	}

	return invoiceTabel.GetModel(), err
}
