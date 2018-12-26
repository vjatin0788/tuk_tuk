package model

import (
	"context"
	"database/sql"
	"log"
)

type DriverUserModel struct {
	Userid                          int64  `db:"userid" json:"userid"`
	Emailid                         string `db:"emailid" json:"emailid"`
	Gender                          string `db:"gender" json:"gender"`
	Dob                             string `db:"dob" json:"dob"`
	Mobileno                        string `db:"mobile_no" json:"mobile_no"`
	Photourl                        string `db:"photo_url" json:"photo_url"`
	Token                           string `db:"token" json:"token"`
	Lastlogin                       string `db:"last_login" json:"last_login"`
	Usertype                        string `db:"user_type" json:"user_type"`
	Createdby                       string `db:"created_by" json:"created_by"`
	Createdon                       string `db:"created_on" json:"created_on"`
	Password                        string `db:"password" json:"password"`
	Address                         string `db:"address" json:"address"`
	Logoutat                        string `db:"logout_at" json:"logout_at"`
	City                            string `db:"city" json:"city"`
	Drivinglicencefront             string `db:"driving_licence_front" json:"driving_licence_front"`
	Drivinglicencenumber            string `db:"driving_licence_number" json:"driving_licence_number"`
	Driverpic                       string `db:"driver_pic" json:"driver_pic"`
	Pancard                         string `db:"pancard" json:"pancard"`
	Pancardnumber                   string `db:"pan_card_number" json:"pan_card_number"`
	Registrationcertificate         string `db:"registration_certificate" json:"registration_certificate"`
	Certificateofregistrationnumber string `db:"certificate_of_registration_number" json:"certificate_of_registration_number"`
	Motorinsurence                  string `db:"motor_insurence" json:"motor_insurence"`
	Motorinsurencenumber            string `db:"motor_insurence_number" json:"motor_insurence_number"`
	Policeverification              string `db:"police_verification" json:"police_verification"`
	Policeverificationnumber        string `db:"police_verification_number" json:"police_verification_number"`
	Adharcard                       string `db:"adhar_card" json:"adhar_card"`
	Aadharcardnumber                string `db:"aadhar_card_number" json:"aadhar_card_number"`
	Intime                          string `db:"in_time" json:"in_time"`
	Outtime                         string `db:"out_time" json:"out_time"`
	Loginstatus                     string `db:"login_status" json:"login_status"`
	Updatedon                       string `db:"updated_on" json:"updated_on"`
	Vehicletype                     string `db:"vehicle_type" json:"vehicle_type"`
	Status                          string `db:"status" json:"status"`
	Driverassigned                  string `db:"driver_assigned" json:"driver_assigned"`
	Driverdutystatus                string `db:"driver_duty_status" json:"driver_duty_status"`
	Name                            string `db:"name" json:"name"`
}

type DriverUserTable struct {
	Userid                          int64          `db:"userid" json:"userid"`
	Emailid                         sql.NullString `db:"emailid" json:"emailid"`
	Gender                          sql.NullString `db:"gender" json:"gender"`
	Dob                             sql.NullString `db:"dob" json:"dob"`
	Mobileno                        sql.NullString `db:"mobile_no" json:"mobile_no"`
	Photourl                        sql.NullString `db:"photo_url" json:"photo_url"`
	Token                           sql.NullString `db:"token" json:"token"`
	Lastlogin                       sql.NullString `db:"last_login" json:"last_login"`
	Usertype                        sql.NullString `db:"user_type" json:"user_type"`
	Createdby                       sql.NullString `db:"created_by" json:"created_by"`
	Createdon                       sql.NullString `db:"created_on" json:"created_on"`
	Password                        sql.NullString `db:"password" json:"password"`
	Address                         sql.NullString `db:"address" json:"address"`
	Logoutat                        sql.NullString `db:"logout_at" json:"logout_at"`
	City                            sql.NullString `db:"city" json:"city"`
	Drivinglicencefront             sql.NullString `db:"driving_licence_front" json:"driving_licence_front"`
	Drivinglicencenumber            sql.NullString `db:"driving_licence_number" json:"driving_licence_number"`
	Driverpic                       sql.NullString `db:"driver_pic" json:"driver_pic"`
	Pancard                         sql.NullString `db:"pancard" json:"pancard"`
	Pancardnumber                   sql.NullString `db:"pan_card_number" json:"pan_card_number"`
	Registrationcertificate         sql.NullString `db:"registration_certificate" json:"registration_certificate"`
	Certificateofregistrationnumber sql.NullString `db:"certificate_of_registration_number" json:"certificate_of_registration_number"`
	Motorinsurence                  sql.NullString `db:"motor_insurence" json:"motor_insurence"`
	Motorinsurencenumber            sql.NullString `db:"motor_insurence_number" json:"motor_insurence_number"`
	Policeverification              sql.NullString `db:"police_verification" json:"police_verification"`
	Policeverificationnumber        sql.NullString `db:"police_verification_number" json:"police_verification_number"`
	Adharcard                       sql.NullString `db:"adhar_card" json:"adhar_card"`
	Aadharcardnumber                sql.NullString `db:"aadhar_card_number" json:"aadhar_card_number"`
	Intime                          sql.NullString `db:"in_time" json:"in_time"`
	Outtime                         sql.NullString `db:"out_time" json:"out_time"`
	Loginstatus                     sql.NullString `db:"login_status" json:"login_status"`
	Updatedon                       sql.NullString `db:"updated_on" json:"updated_on"`
	Vehicletype                     sql.NullString `db:"vehicle_type" json:"vehicle_type"`
	Status                          sql.NullString `db:"status" json:"status"`
	Driverassigned                  sql.NullString `db:"driver_assigned" json:"driver_assigned"`
	Driverdutystatus                sql.NullString `db:"driver_duty_status" json:"driver_duty_status"`
	Name                            sql.NullString `db:"name" json:"name"`
}

func (table DriverUserTable) GetModel() DriverUserModel {
	return DriverUserModel{
		Userid:                          table.Userid,
		Emailid:                         table.Emailid.String,
		Gender:                          table.Gender.String,
		Dob:                             table.Dob.String,
		Mobileno:                        table.Mobileno.String,
		Photourl:                        table.Photourl.String,
		Token:                           table.Token.String,
		Lastlogin:                       table.Lastlogin.String,
		Usertype:                        table.Usertype.String,
		Createdby:                       table.Createdby.String,
		Createdon:                       table.Createdon.String,
		Password:                        table.Password.String,
		Address:                         table.Address.String,
		Logoutat:                        table.Logoutat.String,
		City:                            table.City.String,
		Drivinglicencefront:             table.Drivinglicencefront.String,
		Drivinglicencenumber:            table.Drivinglicencenumber.String,
		Driverpic:                       table.Driverpic.String,
		Pancard:                         table.Pancard.String,
		Pancardnumber:                   table.Pancardnumber.String,
		Registrationcertificate:         table.Registrationcertificate.String,
		Certificateofregistrationnumber: table.Certificateofregistrationnumber.String,
		Motorinsurence:                  table.Motorinsurence.String,
		Motorinsurencenumber:            table.Motorinsurencenumber.String,
		Policeverification:              table.Policeverification.String,
		Policeverificationnumber:        table.Policeverificationnumber.String,
		Adharcard:                       table.Adharcard.String,
		Aadharcardnumber:                table.Aadharcardnumber.String,
		Intime:                          table.Intime.String,
		Outtime:                         table.Outtime.String,
		Loginstatus:                     table.Loginstatus.String,
		Updatedon:                       table.Updatedon.String,
		Vehicletype:                     table.Vehicletype.String,
		Status:                          table.Status.String,
		Driverassigned:                  table.Driverassigned.String,
		Driverdutystatus:                table.Driverdutystatus.String,
		Name:                            table.Name.String,
	}
}

func (model DriverUserModel) GetTable() DriverUserTable {
	return DriverUserTable{
		Userid:                          model.Userid,
		Emailid:                         sql.NullString{model.Emailid, false},
		Gender:                          sql.NullString{model.Gender, false},
		Dob:                             sql.NullString{model.Dob, false},
		Mobileno:                        sql.NullString{model.Mobileno, false},
		Photourl:                        sql.NullString{model.Photourl, false},
		Token:                           sql.NullString{model.Token, false},
		Lastlogin:                       sql.NullString{model.Lastlogin, false},
		Usertype:                        sql.NullString{model.Usertype, false},
		Createdby:                       sql.NullString{model.Createdby, false},
		Createdon:                       sql.NullString{model.Createdon, false},
		Password:                        sql.NullString{model.Password, false},
		Address:                         sql.NullString{model.Address, false},
		Logoutat:                        sql.NullString{model.Logoutat, false},
		City:                            sql.NullString{model.City, false},
		Drivinglicencefront:             sql.NullString{model.Drivinglicencefront, false},
		Drivinglicencenumber:            sql.NullString{model.Drivinglicencenumber, false},
		Driverpic:                       sql.NullString{model.Driverpic, false},
		Pancard:                         sql.NullString{model.Pancard, false},
		Pancardnumber:                   sql.NullString{model.Pancardnumber, false},
		Registrationcertificate:         sql.NullString{model.Registrationcertificate, false},
		Certificateofregistrationnumber: sql.NullString{model.Certificateofregistrationnumber, false},
		Motorinsurence:                  sql.NullString{model.Motorinsurence, false},
		Motorinsurencenumber:            sql.NullString{model.Motorinsurencenumber, false},
		Policeverification:              sql.NullString{model.Policeverification, false},
		Policeverificationnumber:        sql.NullString{model.Policeverificationnumber, false},
		Adharcard:                       sql.NullString{model.Adharcard, false},
		Aadharcardnumber:                sql.NullString{model.Aadharcardnumber, false},
		Intime:                          sql.NullString{model.Intime, false},
		Outtime:                         sql.NullString{model.Outtime, false},
		Loginstatus:                     sql.NullString{model.Loginstatus, false},
		Updatedon:                       sql.NullString{model.Updatedon, false},
		Vehicletype:                     sql.NullString{model.Vehicletype, false},
		Status:                          sql.NullString{model.Status, false},
		Driverassigned:                  sql.NullString{model.Driverassigned, false},
		Driverdutystatus:                sql.NullString{model.Driverdutystatus, false},
		Name:                            sql.NullString{model.Name, false},
	}
}

func (db *DBTuktuk) GetDriverByToken(ctx context.Context, authToken string) (DriverUserModel, error) {
	var (
		driverModel DriverUserModel
		driverTable []DriverUserTable
		err         error
	)

	err = statement.GetDriverUserByAuth.SelectContext(ctx, &driverTable, authToken)
	if err != nil {
		log.Println("[GetDriverByToken][Error] Err in fetching data from db", err)
		return driverModel, err
	}

	for _, driver := range driverTable {
		driverModel = driver.GetModel()
	}

	return driverModel, nil
}
