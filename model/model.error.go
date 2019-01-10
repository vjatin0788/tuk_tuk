package model

import (
	"database/sql"
	"log"
	"time"
)

type ErrorModel struct {
	ID        int64     `json:"id" db:"id"`
	Status    int       `json:"status" db:"status"`
	ErrorCode string    `json:"error_code" db:"error_code"`
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type ErrorTable struct {
	ID        int64     `json:"id" db:"id"`
	Status    int       `json:"status" db:"status"`
	ErrorCode string    `json:"error_code" db:"error_code"`
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (table ErrorTable) GetModel() ErrorModel {
	return ErrorModel{
		ID:        table.ID,
		Status:    table.Status,
		ErrorCode: table.ErrorCode,
		Message:   table.Message,
		CreatedAt: table.CreatedAt,
		UpdatedAt: table.UpdatedAt,
	}
}

func (db *DBTuktuk) GetErrors() ([]ErrorModel, error) {
	var (
		table []ErrorTable
		model []ErrorModel
	)

	err := statement.GetErrors.Select(&table)
	if err != nil && sql.ErrNoRows == nil {
		log.Println("[GetErrors][Error] Err in fetching data from db", err)
		return model, err
	}

	for _, et := range table {
		model = append(model, et.GetModel())
	}

	return model, nil
}
