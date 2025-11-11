-- name: ListProcessesByProject :many
SELECT
  p.id,
  p.name AS process_name,
  pr.name AS project_name,
  p.description,
  p.status,
  e.name AS employee_name,
  p.sequence,
  p.start_time,
  p.end_time,
  p.created_at,
  p.updated_at
FROM processes p
LEFT JOIN employees e ON p.employee_id = e.id
LEFT JOIN projects pr ON p.project_id = pr.id
WHERE p.project_id = $1
ORDER BY p.sequence;

-- name: GetProcessByID :one
SELECT
  p.id,
  p.name AS process_name,
  pr.name AS project_name,
  p.description,
  p.status,
  e.name AS employee_name,
  p.sequence,
  p.start_time,
  p.end_time,
  p.created_at,
  p.updated_at
FROM processes p
LEFT JOIN employees e ON p.employee_id = e.id
LEFT JOIN projects pr ON p.project_id = pr.id
WHERE p.id = $1;

-- name: ListProcessesByEmployee :many
