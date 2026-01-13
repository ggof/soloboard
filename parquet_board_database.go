package main

import (
	"github.com/parquet-go/parquet-go"
)

type parquetBoardDatabase struct {
	dbfilename string
}

func NewBoardDatabase(dbfilename string) BoardDatabase {
	return parquetBoardDatabase{dbfilename}
}

func (p parquetBoardDatabase) Read() ([]Board, error) {
	return parquet.ReadFile[Board](p.dbfilename)
}

func (p parquetBoardDatabase) Write(rows []Board) error {
	return parquet.WriteFile(p.dbfilename, rows)
}
