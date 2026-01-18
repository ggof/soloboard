package db

import (
	"soloboard/model"

	"github.com/parquet-go/parquet-go"
)

type parquetBoardDatabase struct {
	dbfilename string
}

func NewBoardDatabase(dbfilename string) parquetBoardDatabase {
	return parquetBoardDatabase{dbfilename}
}

func (p parquetBoardDatabase) Read() ([]model.Board, error) {
	return parquet.ReadFile[model.Board](p.dbfilename)
}

func (p parquetBoardDatabase) Write(rows []model.Board) error {
	return parquet.WriteFile(p.dbfilename, rows)
}
