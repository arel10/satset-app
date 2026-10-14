-- name: CheckDatabaseHealth :one
SELECT 1 AS alive;

-- name: GetSystemSetting :one
SELECT setting_key, setting_value, created_at, updated_at
FROM system_settings
WHERE setting_key = ? LIMIT 1;
