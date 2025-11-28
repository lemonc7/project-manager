-- name: AllowInitAdmin :one
SELECT NOT EXISTS (
  SELECT 1
  FROM system_settings
  WHERE key = 'admin_initialized'
    AND value = 'true'
) AS allow_init;

-- name: CreateSetting :exec
INSERT INTO system_settings (
  key,
  value,
  description
) VALUES (
  $1,
  $2,
  $3
)
ON CONFLICT (key) DO UPDATE
SET
  value = EXCLUDED.value,
  description = EXCLUDED.description;

