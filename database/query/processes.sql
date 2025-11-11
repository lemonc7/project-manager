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
  p.project_id = COALESCE($1, p.project_id)
  AND p.status = COALESCE(NULLIF($2::text, '')::process_status, p.status)
  AND p.employee_id = COALESCE($3, p.employee_id)
ORDER BY p.sequence;
