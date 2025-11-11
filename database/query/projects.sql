-- name: ListProjectsFiltered :many
SELECT
  p.id,
  p.name AS project_name,
  e.name AS manager_name,
  p.description,
  p.status,
  p.start_time,
  p.duration,
  p.created_at,
  p.updated_at
FROM projects p
LEFT JOIN employees e ON p.manager_id = e.id
WHERE 
  p.status = COALESCE(NULLIF($1::text, '')::project_status, p.status)
  AND p.manager_id = COALESCE($2, p.manager_id)
  AND p.start_time >= COALESCE($3, p.start_time)
  AND p.start_time <= COALESCE($4, p.start_time)
ORDER BY
  CASE p.status
    WHEN '进行' THEN 1
    WHEN '排单' THEN 2
    WHEN '暂停' THEN 3
    WHEN '完成' THEN 4
    WHEN '终止' THEN 5
    WHEN '归档' THEN 6
    ELSE 7
  END,
  (p.start_time + (p.duration || 'days')::interval) NULLS LAST,
  p.id;

-- name: CreateProject :one
INSERT INTO projects (
  name,
  description,
  manager_id
) VALUES (
  $1,
  $2,
  $3
) RETURNING *;

-- name: UpdateProject :one
UPDATE projects
SET
  name = $2,
  description = $3,
  status = $4,
  manager_id = $5,
  start_time = $6,
  duration = $7
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;