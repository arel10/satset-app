package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func Connect(dsn string) (*DB, error) {
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return &DB{DB: sqlDB}, nil
}

func (d *DB) PingContext(ctx context.Context) error {
	if d == nil || d.DB == nil {
		return sql.ErrConnDone
	}
	return d.DB.PingContext(ctx)
}
