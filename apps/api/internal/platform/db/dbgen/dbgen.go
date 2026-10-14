package dbgen

import (
	"context"
	"database/sql"
	"time"
)

type SystemSetting struct {
	ID           uint64    `json:"id"`
	SettingKey   string    `json:"setting_key"`
	SettingValue string    `json:"setting_value"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Queries struct {
	db DBTX
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{db: tx}
}

const checkDatabaseHealth = `-- name: CheckDatabaseHealth :one
SELECT 1 AS alive
`

func (q *Queries) CheckDatabaseHealth(ctx context.Context) (int32, error) {
	row := q.db.QueryRowContext(ctx, checkDatabaseHealth)
	var alive int32
	err := row.Scan(&alive)
	return alive, err
}

const getSystemSetting = `-- name: GetSystemSetting :one
SELECT setting_key, setting_value, created_at, updated_at
FROM system_settings
WHERE setting_key = ? LIMIT 1
`

type GetSystemSettingRow struct {
	SettingKey   string    `json:"setting_key"`
	SettingValue string    `json:"setting_value"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (q *Queries) GetSystemSetting(ctx context.Context, settingKey string) (GetSystemSettingRow, error) {
	row := q.db.QueryRowContext(ctx, getSystemSetting, settingKey)
	var i GetSystemSettingRow
	err := row.Scan(
		&i.SettingKey,
		&i.SettingValue,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
