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

func Connect(dsn string, maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) (*DB, error) {
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if maxOpenConns <= 0 {
		maxOpenConns = 25
	}
	if maxIdleConns <= 0 {
		maxIdleConns = 10
	}
	if connMaxLifetime <= 0 {
		connMaxLifetime = 30 * time.Minute
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	return &DB{DB: sqlDB}, nil
}

func (d *DB) PingContext(ctx context.Context) error {
	if d == nil || d.DB == nil {
		return sql.ErrConnDone
	}
	return d.DB.PingContext(ctx)
}
