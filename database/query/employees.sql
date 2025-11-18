-- name: ListEmployees :many
SELECT
  e.id,
  e.name AS employee_name,
  COALESCE(t.name,'') AS team_name,
  e.is_active,
  e.created_at,
  e.updated_at
FROM employees e
LEFT JOIN teams t ON e.team_id = t.id;

-- name: CreateEmployee :one
INSERT INTO employees (
  name,
  team_id
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: UpdateEmployee :one
UPDATE employees
SET
  name = COALESCE(sqlc.narg('name'),name),
  team_id = COALESCE(sqlc.narg('team_id'),team_id),
  is_active = COALESCE(sqlc.narg('is_active'),is_active)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE id = $1;