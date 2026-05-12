-- name: ListEmployees :many
SELECT
  e.id,
  e.name,
  t.name AS team_name,
  COALESCE(array_agg(r.name) FILTER (WHERE r.name IS NOT NULL), '{}') AS roles
FROM employees e
LEFT JOIN teams t ON e.team_id = t.id
LEFT JOIN employee_roles er ON e.id = er.employee_id
LEFT JOIN roles r ON er.role_id = r.id
WHERE is_active = true
ORDER BY e.id, t.name;

-- name: CreateEmployee :exec
INSERT INTO employees (
  name,
  username,
  password_hash,
  team_id
) VALUES (
  $1,
  $2,
  $3,
  $4
);

-- name: UpdateEmployee :exec
UPDATE employees
SET
  name = COALESCE(sqlc.narg('name'), name),
  username = COALESCE(sqlc.narg('username'), username),
  team_id = COALESCE(sqlc.narg('team_id'), team_id),
  is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE id = sqlc.arg('id');

-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE id = $1;

-- name: GetEmployeeByUsername :one
SELECT
  e.id,
  e.name,
  e.username,
  e.password_hash,
  COALESCE(array_agg(r.name) FILTER (WHERE r.name IS NOT NULL), '{}') AS roles
FROM employees e
LEFT JOIN employee_roles er ON e.id = er.employee_id
LEFT JOIN roles r ON er.role_id = r.id
WHERE 
  username = $1
  AND is_active = true;

-- name: UpdatePassword :exec
UPDATE employees
SET password_hash = $2
WHERE id = $1;

-- name: AddEmployeeRole :exec
INSERT INTO employee_roles (employee_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteEmployeeRoles :exec
DELETE FROM employee_roles WHERE employee_id = $1;