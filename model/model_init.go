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
	InsertInvoiceByCustomerId  *sqlx.Stmt
	GetInvoiceByCustomerId     *sqlx.Stmt
	InsertRideDetails          *sqlx.Stmt
	UpdateRideDetails          *sqlx.Stmt
	GetRideDetailsByRideId     *sqlx.Stmt
	GetRideDetailsByCustomerID *sqlx.Stmt
	GetDriverUserById          *sqlx.Stmt
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
	getDriverUserById          = `SELECT * FROM tbusers WHERE userid=?`
	insertInvoiceByCustomerId  = `INSERT INTO tbinvoice(customer_id,driver_id,source_lat,source_lng,source_address,destination_lat,destination_lng,source_time,total_minutes,cost_per_minute,time_cost,distance,cost_per_km,distance_cost,base_fare,extra_charges,discount,total_cost,gst_percentage,gst,final_cost) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	getInvoiceByCustomerId     = `SELECT * FROM tbinvoice WHERE customer_id=? ORDER BY invoice_id DESC LIMIT `
	insertRideDetails          = `INSERT INTO tb_ride_details(customer_id,source_lat,source_long,destination_lat,destination_long,status) VALUES(?,?,?,?,?,?)`
	updateRideDetailsByRideId  = `UPDATE tb_ride_details SET driver_id=?,status=?,driver_cancelled=?,rider_cancelled=?,ride_booked_time=?,ride_completed_time=?,ride_failed_time=?  where id=?`
	getRideDetailsByRideId     = `SELECT * FROM tb_ride_details WHERE id=?`
	getRideDetailsByCustomerID = `SELECT * FROM tb_ride_details WHERE customer_id=? ORDER BY id DESC LIMIT 1`
)

func InitModel(db *DBTuktuk) {
	statement.InsertDriverData, _ = db.DBConnection.Preparex(insertDriver)
	statement.UpdateDriverData, _ = db.DBConnection.Preparex(updateDriver)
	statement.GetDriverById, _ = db.DBConnection.Preparex(getDriverById)
	statement.GetCustomerByAuth, _ = db.DBConnection.Preparex(getCustomerByAuth)
	statement.GetDriverUserByAuth, _ = db.DBConnection.Preparex(getDriverByAuth)
	statement.InsertInvoiceByCustomerId, _ = db.DBConnection.Preparex(insertInvoiceByCustomerId)
	statement.GetInvoiceByCustomerId, _ = db.DBConnection.Preparex(getInvoiceByCustomerId)
	statement.InsertRideDetails, _ = db.DBConnection.Preparex(insertRideDetails)
	statement.UpdateRideDetails, _ = db.DBConnection.Preparex(updateRideDetailsByRideId)
	statement.GetRideDetailsByCustomerID, _ = db.DBConnection.Preparex(getRideDetailsByCustomerID)
	statement.GetRideDetailsByRideId, _ = db.DBConnection.Preparex(getRideDetailsByRideId)
	statement.GetDriverUserById, _ = db.DBConnection.Preparex(getDriverUserById)
}
