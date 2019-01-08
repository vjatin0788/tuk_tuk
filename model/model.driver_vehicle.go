package model

import (
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type VehicleTable struct {
	VehicleId        int64          `db:"vehicle_id" json:"vehicle_id"`
	VehicleType      sql.NullString `db:"vehicle_type" json:"vehicle_type"`
	Make             sql.NullString `db:"make" json:"make"`
	Model            sql.NullString `db:"model" json:"model"`
	VehicleNumber    sql.NullString `db:"vehicle_number" json:"vehicle_number"`
	RcNumber         sql.NullString `db:"rc_no" json:"rc_no"`
	RcImage          sql.NullString `db:"rc_image" json:"rc_image"`
	PermitNumber     sql.NullString `db:"permit_no" json:"permit_no"`
	PermitPath       sql.NullString `db:"permit_path" json:"permit_path"`
	InsuranceNumber  sql.NullString `db:"insurance_no" json:"insurance_no"`
	InsurancePath    sql.NullString `db:"insurance_path" json:"insurance_path"`
	CreatedOn        sql.NullString `db:"created_on" json:"created_on"`
	CreatedBy        sql.NullString `db:"created_by" json:"created_by"`
	UpdatedOn        sql.NullString `db:"updated_on" json:"updated_on"`
	UpdatedBy        sql.NullString `db:"updated_by" json:"updated_by"`
	Status           sql.NullString `db:"status" json:"status"`
	AssignedDriverId sql.NullInt64  `db:"assigned_driver_id" json:"assigned_driver_id"`
}

type VehicleModel struct {
	VehicleId        int64  `db:"vehicle_id" json:"vehicle_id"`
	VehicleType      string `db:"vehicle_type" json:"vehicle_type"`
	Make             string `db:"make" json:"make"`
	Model            string `db:"model" json:"model"`
	VehicleNumber    string `db:"vehicle_number" json:"vehicle_number"`
	RcNumber         string `db:"rc_no" json:"rc_no"`
	RcImage          string `db:"rc_image" json:"rc_image"`
	PermitNumber     string `db:"permit_no" json:"permit_no"`
	PermitPath       string `db:"permit_path" json:"permit_path"`
	InsuranceNumber  string `db:"insurance_no" json:"insurance_no"`
	InsurancePath    string `db:"insurance_path" json:"insurance_path"`
	CreatedOn        string `db:"created_on" json:"created_on"`
	CreatedBy        string `db:"created_by" json:"created_by"`
	UpdatedOn        string `db:"updated_on" json:"updated_on"`
	UpdatedBy        string `db:"updated_by" json:"updated_by"`
	Status           string `db:"status" json:"status"`
	AssignedDriverId int64  `db:"assigned_driver_id" json:"assigned_driver_id"`
}

func (table VehicleTable) GetModel() VehicleModel {
	return VehicleModel{
		VehicleId:        table.VehicleId,
		VehicleType:      table.VehicleType.String,
		Make:             table.Make.String,
		Model:            table.Model.String,
		VehicleNumber:    table.VehicleNumber.String,
		RcNumber:         table.RcNumber.String,
		RcImage:          table.RcImage.String,
		PermitNumber:     table.PermitNumber.String,
		PermitPath:       table.PermitPath.String,
		InsuranceNumber:  table.InsuranceNumber.String,
		InsurancePath:    table.InsurancePath.String,
		CreatedOn:        table.CreatedOn.String,
		CreatedBy:        table.CreatedBy.String,
		UpdatedBy:        table.UpdatedBy.String,
		UpdatedOn:        table.UpdatedOn.String,
		Status:           table.Status.String,
		AssignedDriverId: table.AssignedDriverId.Int64,
	}
}

func (model VehicleModel) GetModel() VehicleTable {
	return VehicleTable{
		VehicleId:        model.VehicleId,
		VehicleType:      sql.NullString{model.VehicleType, false},
		Make:             sql.NullString{model.Make, false},
		Model:            sql.NullString{model.Model, false},
		VehicleNumber:    sql.NullString{model.VehicleNumber, false},
		RcNumber:         sql.NullString{model.RcNumber, false},
		RcImage:          sql.NullString{model.RcImage, false},
		PermitNumber:     sql.NullString{model.PermitNumber, false},
		PermitPath:       sql.NullString{model.PermitPath, false},
		InsuranceNumber:  sql.NullString{model.InsuranceNumber, false},
		InsurancePath:    sql.NullString{model.InsurancePath, false},
		CreatedOn:        sql.NullString{model.CreatedOn, false},
		CreatedBy:        sql.NullString{model.CreatedBy, false},
		UpdatedBy:        sql.NullString{model.UpdatedBy, false},
		UpdatedOn:        sql.NullString{model.UpdatedOn, false},
		Status:           sql.NullString{model.Status, false},
		AssignedDriverId: sql.NullInt64{model.AssignedDriverId, false},
	}
}

func (db *DBTuktuk) GetVehicleByAssignedDriver(ctx context.Context, driverIds []int64) ([]VehicleModel, error) {

	var (
		model []VehicleModel
		table []VehicleTable
		err   error
	)

	query, args, err := sqlx.In(getVehicleByAssignedDriver, driverIds)
	if err != nil {
		log.Println("[GetVehicleByAssignedDriver][Error] Err in sqlx IN", err)
		return model, err
	}

	query = db.DBConnection.Rebind(query)
	err = db.DBConnection.Select(&table, query, args...)
	if err != nil {
		log.Println("[GetVehicleByAssignedDriver][Error] Err in fetching data from db", err)
		return model, err
	}

	for _, t1 := range table {
		model = append(model, t1.GetModel())
	}

	return model, err

}
