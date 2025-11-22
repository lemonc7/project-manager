-- name: ListEmployees :many
SELECT
  e.id,
  e.name AS employee_name,
  t.name AS team_name,
  e.roles
FROM employees e
LEFT JOIN teams t ON e.team_id = t.id
WHERE e.is_active = true
ORDER BY e.name;

-- name: GetEmployeeByUsername :one
SELECT
  id,
  username,
  password_hash,
  roles
FROM employees
WHERE 
  username = $1
  AND is_active = true;

-- name: CreateEmployee :one
INSERT INTO employees (
  name,
  username,
  password_hash,
  roles,
  team_id
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
) RETURNING *;

-- name: UpdateEmployee :one
UPDATE employees
SET
  name = COALESCE(sqlc.narg('name'),name),
  roles = COALESCE(sqlc.narg('roles'),roles),
  team_id = COALESCE(sqlc.narg('team_id'),team_id),
  is_active = COALESCE(sqlc.narg('is_active'),is_active)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE id = $1;