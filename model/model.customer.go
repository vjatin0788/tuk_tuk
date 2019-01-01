package model

import (
	"context"
	"database/sql"
	"log"
)

type CustomerModel struct {
	CustomerId     int64  `db:"customer_id" json:"customer_id"`
	UserId         string `db:"user_id" json:"user_id"`
	EmailType      string `db:"email_type" json:"email_type"`
	EmailId        string `db:"email_id" json:"email_id"`
	LoginVia       string `db:"login_via" json:"login_via"`
	Name           string `db:"name" json:"name"`
	Gender         string `db:"gender" json:"gender"`
	Password       string `db:"password" json:"password"`
	Dob            string `db:"dob" json:"dob"`
	Mobile         string `db:"mobile_no" json:"mobile_no"`
	PhotoUrl       string `db:"photo_url" json:"photo_url"`
	Token          string `db:"token" json:"token"`
	LastLogin      string `db:"last_login" json:"last_login"`
	CreatedBy      string `db:"created_by" json:"created_by"`
	CreatedOn      string `db:"created_on" json:"created_on"`
	UpdatedBy      string `db:"updated_by" json:"updated_by"`
	UpdatedOn      string `db:"updated_on" json:"updated_on"`
	MobileVerified string `db:"mobile_verified" json:"mobile_verified"`
	EmailVerified  string `db:"email_verified" json:"email_verified"`
	MobileOtp      string `db:"mobile_otp" json:"mobile_otp"`
	EmailOtp       string `db:"email_otp" json:"email_otp"`
	DeviceId       string `db:"device_id" json:"device_id"`
	Refferal       int64  `db:"refferal" json:"refferal"`
}

type CustomerTable struct {
	CustomerId     int64          `db:"customer_id" json:"customer_id"`
	UserId         sql.NullString `db:"user_id" json:"user_id"`
	EmailType      sql.NullString `db:"email_type" json:"email_type"`
	EmailId        sql.NullString `db:"email_id" json:"email_id"`
	LoginVia       sql.NullString `db:"login_via" json:"login_via"`
	Name           sql.NullString `db:"name" json:"name"`
	Gender         sql.NullString `db:"gender" json:"gender"`
	Password       sql.NullString `db:"password" json:"password"`
	Dob            sql.NullString `db:"dob" json:"dob"`
	Mobile         sql.NullString `db:"mobile_no" json:"mobile_no"`
	PhotoUrl       sql.NullString `db:"photo_url" json:"photo_url"`
	Token          sql.NullString `db:"token" json:"token"`
	LastLogin      sql.NullString `db:"last_login" json:"last_login"`
	CreatedBy      sql.NullString `db:"created_by" json:"created_by"`
	CreatedOn      sql.NullString `db:"created_on" json:"created_on"`
	UpdatedBy      sql.NullString `db:"updated_by" json:"updated_by"`
	UpdatedOn      sql.NullString `db:"updated_on" json:"updated_on"`
	MobileVerified sql.NullString `db:"mobile_verified" json:"mobile_verified"`
	EmailVerified  sql.NullString `db:"email_verified" json:"email_verified"`
	MobileOtp      sql.NullString `db:"mobile_otp" json:"mobile_otp"`
	EmailOtp       sql.NullString `db:"email_otp" json:"email_otp"`
	DeviceId       sql.NullString `db:"device_id" json:"device_id"`
	Refferal       sql.NullInt64  `db:"refferal" json:"refferal"`
}

func (table CustomerTable) GetModel() CustomerModel {
	return CustomerModel{
		CustomerId:     table.CustomerId,
		UserId:         table.UserId.String,
		EmailId:        table.EmailId.String,
		EmailType:      table.EmailType.String,
		LoginVia:       table.LoginVia.String,
		Name:           table.Name.String,
		Gender:         table.Gender.String,
		Password:       table.Password.String,
		Dob:            table.Dob.String,
		Mobile:         table.Mobile.String,
		PhotoUrl:       table.PhotoUrl.String,
		Token:          table.Token.String,
		LastLogin:      table.LastLogin.String,
		CreatedBy:      table.CreatedBy.String,
		CreatedOn:      table.CreatedOn.String,
		UpdatedBy:      table.UpdatedBy.String,
		UpdatedOn:      table.UpdatedOn.String,
		EmailVerified:  table.EmailVerified.String,
		MobileVerified: table.MobileVerified.String,
		EmailOtp:       table.EmailOtp.String,
		MobileOtp:      table.MobileOtp.String,
		DeviceId:       table.DeviceId.String,
		Refferal:       table.Refferal.Int64,
	}
}

func (model CustomerModel) GetTable() CustomerTable {
	return CustomerTable{
		CustomerId:     model.CustomerId,
		UserId:         sql.NullString{model.UserId, false},
		EmailId:        sql.NullString{model.EmailId, false},
		EmailType:      sql.NullString{model.EmailType, false},
		LoginVia:       sql.NullString{model.LoginVia, false},
		Name:           sql.NullString{model.Name, false},
		Gender:         sql.NullString{model.Gender, false},
		Password:       sql.NullString{model.Password, false},
		Dob:            sql.NullString{model.Dob, false},
		Mobile:         sql.NullString{model.Mobile, false},
		PhotoUrl:       sql.NullString{model.PhotoUrl, false},
		Token:          sql.NullString{model.Token, false},
		LastLogin:      sql.NullString{model.LastLogin, false},
		CreatedBy:      sql.NullString{model.CreatedBy, false},
		CreatedOn:      sql.NullString{model.CreatedOn, false},
		UpdatedBy:      sql.NullString{model.UpdatedBy, false},
		UpdatedOn:      sql.NullString{model.UpdatedOn, false},
		EmailVerified:  sql.NullString{model.EmailVerified, false},
		MobileVerified: sql.NullString{model.MobileVerified, false},
		EmailOtp:       sql.NullString{model.EmailOtp, false},
		MobileOtp:      sql.NullString{model.MobileOtp, false},
		DeviceId:       sql.NullString{model.DeviceId, false},
		Refferal:       sql.NullInt64{model.Refferal, false},
	}
}

func (db *DBTuktuk) GetCustomerByToken(ctx context.Context, authToken string) (CustomerModel, error) {
	var (
		custModel CustomerModel
		custTable []CustomerTable
		err       error
	)

	err = statement.GetCustomerByAuth.SelectContext(ctx, &custTable, authToken)
	if err != nil {
		log.Println("[GetCustomerByToken][Error] Err in fetching data from db", err)
		return custModel, err
	}

	for _, cust := range custTable {
		custModel = cust.GetModel()
	}

	return custModel, nil
}

func (db *DBTuktuk) GetCustomerById(ctx context.Context, id int64) (CustomerModel, error) {
	var (
		custModel CustomerModel
		custTable []CustomerTable
		err       error
	)

	err = statement.GetCustomerById.SelectContext(ctx, &custTable, id)
	if err != nil {
		log.Println("[GetCustomerById][Error] Err in fetching data from db", err)
		return custModel, err
	}

	for _, cust := range custTable {
		custModel = cust.GetModel()
	}

	return custModel, nil
}
