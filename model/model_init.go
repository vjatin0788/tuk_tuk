package model

import (
	"github.com/jmoiron/sqlx"
)

type PreparedStatement struct {
	GetAvailableDrivers        *sqlx.Stmt
	InsertDriverData           *sqlx.Stmt
	UpdateDriverData           *sqlx.Stmt
	GetDriverById              *sqlx.Stmt
	GetVehicleByAssignedDriver *sqlx.Stmt
	GetCustomerByAuth          *sqlx.Stmt
	GetDriverUserByAuth        *sqlx.Stmt
}

var (
	statement PreparedStatement
)

var (
	insertDriver               = `INSERT INTO driver_tracking(driver_id,current_lat,current_long,current_lat_rad,current_long_rad) VALUES(?,?,?,?,?)`
	updateDriver               = `UPDATE driver_tracking SET current_lat=?,current_long=?,current_lat_rad=?,current_long_rad=?,last_lat=?,last_long=?,last_lat_rad=?,last_long_rad=? WHERE  driver_id = ?`
	getDriverById              = `SELECT id,driver_id,current_lat,current_long,current_lat_rad,current_long_rad FROM driver_tracking WHERE driver_id=?`
	getVehicleByAssignedDriver = `SELECT * FROM tbvehicle WHERE assigned_driver_id IN (?)`
	getCustomerByAuth          = `SELECT * FROM tbcustomers WHERE token=?`
	getDriverByAuth            = `SELECT * FROM tbusers WHERE FIND_IN_SET(?, token)`
)

func InitModel(db *DBTuktuk) {
	statement.InsertDriverData, _ = db.DBConnection.Preparex(insertDriver)
	statement.UpdateDriverData, _ = db.DBConnection.Preparex(updateDriver)
	statement.GetDriverById, _ = db.DBConnection.Preparex(getDriverById)
	statement.GetCustomerByAuth, _ = db.DBConnection.Preparex(getCustomerByAuth)
	statement.GetDriverUserByAuth, _ = db.DBConnection.Preparex(getDriverByAuth)
}
