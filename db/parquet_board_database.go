package db

import (
	"errors"
	"os"
	"soloboard/model"

	"github.com/parquet-go/parquet-go"
)

type parquetBoardDatabase struct {
	dbfilename string
}

func NewBoardDatabase(dbfilename string) parquetBoardDatabase {
	if _, err := os.Stat(dbfilename); errors.Is(err, os.ErrNotExist) {
		parquet.WriteFile(dbfilename, []model.Board{})
	}
	return parquetBoardDatabase{dbfilename}
}

func (p parquetBoardDatabase) Read() ([]model.Board, error) {
	return parquet.ReadFile[model.Board](p.dbfilename)
}

func (p parquetBoardDatabase) Write(rows []model.Board) error {
	return parquet.WriteFile(p.dbfilename, rows)
}
