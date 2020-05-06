package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//DB - structure that is holding the Database
type DB struct {
	*sql.DB
}

type key string

const (
	//SQLKEY - key to store and retrieve sql reference in context.
	SQLKEY key = "SQL"
)

//InitDB - creates new =connection instance
func InitDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
