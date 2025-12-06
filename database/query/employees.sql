-- name: ListEmployees :many
SELECT
  id,
  name,
  roles,
  team_id
FROM employees
WHERE is_active = true
ORDER BY id;

-- name: CreateEmployee :exec
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
);

-- name: UpdateEmployee :exec
UPDATE employees
SET
  name = COALESCE(sqlc.narg('name'),name),
  roles = COALESCE(sqlc.narg('roles'),roles),
  team_id = COALESCE(sqlc.narg('team_id'),team_id),
  is_active = COALESCE(sqlc.narg('is_active'),is_active)
WHERE id = sqlc.arg('id');

-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE id = $1;

-- name: GetEmployeeByUsername :one
SELECT
  id,
  name,
  username,
  password_hash,
  roles,
  team_id
FROM employees
WHERE 
  username = $1
  AND is_active = true;

-- name: UpdatePassword :exec
UPDATE employees
SET password_hash = $2
WHERE id = $1;
