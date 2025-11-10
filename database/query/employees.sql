-- name: ListEmployees :many
SELECT
  e.id,
  e.name AS employee_name,
  t.name AS team_name,
  e.is_active,
  e.created_at,
  e.updated_at
FROM employees e
LEFT JOIN teams t ON e.team_id = t.id;

-- name: GetEmployeeByID :one
SELECT
  e.id,
  e.name AS employee_name,
  t.name AS team_name,
  e.is_active,
  e.created_at,
  e.updated_at
FROM employees e
LEFT JOIN teams t ON e.team_id = t.id
WHERE e.id = $1;

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
  name = $2,
  team_id = $3,
  is_active = $4
WHERE id = $1
RETURNING *;

-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE id = $1;