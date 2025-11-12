-- name: ListProcessesFiltered :many
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
WHERE 
  p.project_id = COALESCE(sqlc.narg('project_id'), p.project_id)
  AND p.status = COALESCE(sqlc.narg('status'), p.status)
  AND p.employee_id = COALESCE(sqlc.narg('employee_id'), p.employee_id)
ORDER BY p.sequence;

-- name: CreateProcess :exec
INSERT INTO processes (
  name,
  project_id,
  employee_id,
  description,
  sequence,
  remark
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
);

-- name: UpdateProcess :exec
UPDATE processes
SET
  name = COALESCE(sqlc.narg('name'),name),
  description = COALESCE(sqlc.narg('description'),description),
  status = COALESCE(sqlc.narg('status'),status),
  employee_id = COALESCE(sqlc.narg('employee_id'),employee_id),
  sequence = COALESCE(sqlc.narg('sequence'),sequence),
  remark = COALESCE(sqlc.narg('remark'),remark)
WHERE id = sqlc.arg('id');

-- name: DeleteProcess :exec
DELETE FROM processes
WHERE id = $1;

