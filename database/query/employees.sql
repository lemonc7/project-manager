-- name: ListEmployees :many
SELECT
  e.id,
  e.name,
  t.name AS team_name,
  array_remove(array_agg(r.name), NULL)::TEXT[] AS roles
FROM employees e
LEFT JOIN teams t ON e.team_id = t.id
LEFT JOIN employee_roles er ON e.id = er.employee_id
LEFT JOIN roles r ON er.role_id = r.id
WHERE e.is_active = true
GROUP BY e.id, t.name
ORDER BY e.id;

-- name: CreateEmployee :one
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
) RETURNING id;

-- name: UpdateEmployee :exec
UPDATE employees
SET
  name = COALESCE(sqlc.narg('name'), name),
  username = COALESCE(sqlc.narg('username'), username),
  team_id = COALESCE(sqlc.narg('team_id'), team_id),
  is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE id = sqlc.arg('id');

-- name: GetEmployeeByUsername :one
SELECT
  e.id,
  e.name,
  e.username,
  e.password_hash,
  array_remove(array_agg(r.name), NULL)::TEXT[] AS roles
FROM employees e
LEFT JOIN employee_roles er ON e.id = er.employee_id
LEFT JOIN roles r ON er.role_id = r.id
WHERE
  e.username = $1
  AND e.is_active = true
GROUP BY e.id, e.name, e.username, e.password_hash;

-- name: UpdatePassword :exec
UPDATE employees
SET password_hash = $2
WHERE id = $1;

-- name: AddEmployeeRoles :exec
INSERT INTO employee_roles (employee_id, role_id)
SELECT sqlc.arg('employee_id'), unnest(sqlc.arg('role_ids')::INT[])
ON CONFLICT DO NOTHING;

-- name: DeleteEmployeeRoles :exec
DELETE FROM employee_roles WHERE employee_id = $1;
