-- name: ListProjectsFiltered :many
SELECT
  p.id,
  p.name,
  COALESCE(e.name,'') AS manager_name,
  p.description,
  p.remark,
  p.status,
  p.start_time,
  p.duration,
  p.created_at,
  p.updated_at
FROM projects p
LEFT JOIN employees e ON p.manager_id = e.id
WHERE 
  p.status = COALESCE(sqlc.narg('status'), p.status)
  AND p.manager_id = COALESCE(sqlc.narg('manager_id'), p.manager_id)
  AND p.created_at >= COALESCE(sqlc.narg('start_time'), p.created_at)
  AND p.created_at <= COALESCE(sqlc.narg('end_time'), p.created_at)
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
  manager_id,
  remark
) VALUES (
  $1,
  $2,
  $3,
  $4
) RETURNING *;

-- name: UpdateProject :one
UPDATE projects
SET
  name = COALESCE(sqlc.narg('name'),name),
  description = COALESCE(sqlc.narg('description'),description),
  status = COALESCE(sqlc.narg('status'),status),
  manager_id = COALESCE(sqlc.narg('manager_id'),manager_id),
  start_time = COALESCE(sqlc.narg('start_time'),start_time),
  duration = COALESCE(sqlc.narg('duration'),duration),
  remark = COALESCE(sqlc.narg('remark'),remark)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;