package model

import (
	"log"

	"github.com/TukTuk/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const DriverMysql = "mysql"

var TukTuk *DBTuktuk

type DBTuktuk struct {
	DBConnection  *sqlx.DB
	DBString      string
	RetryInterval int
	MaxConn       int
	doneChannel   chan bool
}

func InitDatabase() error {

	var err error
	masterDBStr := common.MASTER_DB_PROD

	masterDB := &DBTuktuk{
		DBString: masterDBStr,
	}

	err = masterDB.Connect(DriverMysql)
	if err != nil {
		log.Println("[InitDatabase][Error] Err initializing db", err)
		return err
	}

	InitModel(masterDB)
	TukTuk = masterDB

	return err
}

func (d *DBTuktuk) Connect(driver string) error {
	var (
		db  *sqlx.DB
		err error
	)

	db, err = sqlx.Open(driver, d.DBString)
	if err != nil {
		log.Println("[Error]: DB open connection error", err.Error())
		return err
	}

	d.DBConnection = db

	err = db.Ping()
	if err != nil {
		log.Println("[Error]: DB connection error", err.Error())
		return err
	}

	db.SetMaxOpenConns(d.MaxConn)

	return err
}
