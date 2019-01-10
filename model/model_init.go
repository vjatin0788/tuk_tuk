package model

import (
	"github.com/jmoiron/sqlx"
)

type PreparedStatement struct {
	GetAvailableDrivers                 *sqlx.Stmt
	InsertDriverData                    *sqlx.Stmt
	UpdateDriverData                    *sqlx.Stmt
	GetDriverById                       *sqlx.Stmt
	GetVehicleByAssignedDriver          *sqlx.Stmt
	GetCustomerByAuth                   *sqlx.Stmt
	GetCustomerById                     *sqlx.Stmt
	GetDriverUserByAuth                 *sqlx.Stmt
	InsertInvoiceByCustomerId           *sqlx.Stmt
	GetInvoiceByCustomerId              *sqlx.Stmt
	InsertRideDetails                   *sqlx.Stmt
	UpdateRideDetails                   *sqlx.Stmt
	UpdateRideDetailsWithStatus         *sqlx.Stmt
	UpdateRideStatusFailed              *sqlx.Stmt
	UpdateRideStart                     *sqlx.Stmt
	UpdateRideStatus                    *sqlx.Stmt
	UpdateRideComplete                  *sqlx.Stmt
	GetRideDetailsByRideId              *sqlx.Stmt
	GetRideDetailsByCustomerID          *sqlx.Stmt
	GetDriverUserById                   *sqlx.Stmt
	GetRideDetailsByDriverID            *sqlx.Stmt
	GetRideDetailStatusByCustomerID     *sqlx.Stmt
	InsertTrackingData                  *sqlx.Stmt
	GetRideDetailsByCustomerIDAndStatus *sqlx.Stmt
}

var (
	statement PreparedStatement
)

var (
	insertDriver                        = `INSERT INTO driver_tracking(driver_id,current_lat,current_long,current_lat_rad,current_long_rad) VALUES(?,?,?,?,?)`
	updateDriver                        = `UPDATE driver_tracking SET current_lat=?,current_long=?,current_lat_rad=?,current_long_rad=?,last_lat=?,last_long=?,last_lat_rad=?,last_long_rad=? WHERE  driver_id = ?`
	getDriverById                       = `SELECT id,driver_id,current_lat,current_long,current_lat_rad,current_long_rad FROM driver_tracking WHERE driver_id=?`
	getVehicleByAssignedDriver          = `SELECT * FROM tbvehicle WHERE assigned_driver_id IN (?)`
	getCustomerByAuth                   = `SELECT * FROM tbcustomers WHERE token=?`
	getCustomerById                     = `SELECT * FROM tbcustomers WHERE customer_id=?`
	getDriverByAuth                     = `SELECT * FROM tbusers WHERE FIND_IN_SET(?, token)`
	getDriverUserById                   = `SELECT * FROM tbusers WHERE userid=?`
	insertInvoiceByCustomerId           = `INSERT INTO tbinvoice(customer_id,driver_id,source_lat,source_lng,source_address,destination_lat,destination_lng,source_time,total_minutes,cost_per_minute,time_cost,distance,cost_per_km,distance_cost,base_fare,extra_charges,discount,total_cost,gst_percentage,gst,final_cost) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	getInvoiceByCustomerId              = `SELECT * FROM tbinvoice WHERE customer_id=? ORDER BY invoice_id DESC LIMIT `
	insertRideDetails                   = `INSERT INTO tb_ride_details(customer_id,source_lat,source_long,destination_lat,destination_long,status,payment_method) VALUES(?,?,?,?,?,?,?)`
	updateRideDetailsByRideId           = `UPDATE tb_ride_details SET status=?,driver_cancelled=?,rider_cancelled=?,ride_failed_time=? WHERE status<=2 AND id=?`
	updateRideDetailsByRideIdAndStatus  = `UPDATE tb_ride_details SET driver_id=?,status=?,ride_booked_time=? WHERE id=? and status=1`
	updateRideStatusStartByRideId       = `UPDATE tb_ride_details SET status=?,ride_start_time=? WHERE id=?`
	updateRideStatusFailedByRideId      = `UPDATE tb_ride_details SET status=?,ride_failed_time=? WHERE id=?`
	updateRideStatusById                = `UPDATE tb_ride_details SET driver_id=?,status=? WHERE id=?`
	updateRideCompleteById              = `UPDATE tb_ride_details SET status=?,ride_completed_time=?,destination_lat=?,destination_long=? WHERE id=?`
	getRideDetailsByRideId              = `SELECT * FROM tb_ride_details WHERE id=?`
	getRideDetailsByCustomerID          = `SELECT * FROM tb_ride_details WHERE customer_id=? ORDER BY id DESC LIMIT 1`
	getRideDetailsByCustomerIDAndStatus = `SELECT * FROM tb_ride_details WHERE customer_id=? AND status=? ORDER BY id DESC LIMIT 1`
	getRideDetailsByDriverID            = `SELECT * FROM tb_ride_details WHERE driver_id=? ORDER BY id DESC LIMIT 1`
	getRideDetailsStatusByCustomerID    = `SELECT * FROM tb_ride_details WHERE customer_id=? AND status IN (?) ORDER BY id`
	getRideDetailsStatusByDriverID      = `SELECT * FROM tb_ride_details WHERE driver_id=? AND status IN (?) ORDER BY id`
	insertTrackingData                  = `INSERT INTO tbtrackingdata(emailid,lat,lng,datetime,tracking_type,user_id) VALUES(?,?,?,?,?,?)`
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
	statement.GetRideDetailsByDriverID, _ = db.DBConnection.Preparex(getRideDetailsByDriverID)
	statement.UpdateRideDetailsWithStatus, _ = db.DBConnection.Preparex(updateRideDetailsByRideIdAndStatus)
	statement.UpdateRideStart, _ = db.DBConnection.Preparex(updateRideStatusStartByRideId)
	statement.GetCustomerById, _ = db.DBConnection.Preparex(getCustomerById)
	statement.UpdateRideStatusFailed, _ = db.DBConnection.Preparex(updateRideStatusFailedByRideId)
	statement.UpdateRideStatus, _ = db.DBConnection.Preparex(updateRideStatusById)
	statement.UpdateRideComplete, _ = db.DBConnection.Preparex(updateRideCompleteById)
	statement.GetRideDetailStatusByCustomerID, _ = db.DBConnection.Preparex(getRideDetailsStatusByCustomerID)
	statement.InsertTrackingData, _ = db.DBConnection.Preparex(insertTrackingData)
	statement.GetRideDetailsByCustomerIDAndStatus, _ = db.DBConnection.Preparex(getRideDetailsByCustomerIDAndStatus)
}
