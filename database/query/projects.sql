-- name: ListProjectsFiltered :many
SELECT
  p.id,
  p.name,
  p.description,
  p.remark,
  p.status,
  COALESCE(e.name,'') AS manager_name,
  p.started_at,
  p.duration,
  p.finished_at
FROM projects p
LEFT JOIN employees e ON p.manager_id = e.id
WHERE 
  (p.status = sqlc.narg('status') OR sqlc.narg('status') IS NULL)
  AND (p.manager_id = sqlc.narg('manager_id') OR sqlc.narg('manager_id') IS NULL)
  AND (p.started_at >= sqlc.narg('start_at') OR sqlc.narg('start_at') IS NULL)
  AND (p.started_at <= sqlc.narg('end_at') OR sqlc.narg('end_at') IS NULL)
ORDER BY
  CASE p.status
    WHEN '进行' THEN 1
    WHEN '排单' THEN 2
    WHEN '暂停' THEN 3
    WHEN '完成' THEN 4
    WHEN '终止' THEN 5
    ELSE 7
  END,
  (p.started_at + (p.duration || 'days')::interval) NULLS LAST,
  p.id;

-- name: CreateProject :one
INSERT INTO projects (
  name,
  description,
  remark,
  manager_id,
  duration
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
) RETURNING *;

-- name: UpdateProject :one
UPDATE projects
SET
  name = COALESCE(sqlc.narg('name'),name),
  description = COALESCE(sqlc.narg('description'),description),
  remark = COALESCE(sqlc.narg('remark'),remark),
  status = COALESCE(sqlc.narg('status'),status),
  manager_id = COALESCE(sqlc.narg('manager_id'),manager_id),
  started_at = COALESCE(sqlc.narg('started_at'),started_at),
  duration = COALESCE(sqlc.narg('duration'),duration),
  finished_at = COALESCE(sqlc.narg('finished_at'),finished_at)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;