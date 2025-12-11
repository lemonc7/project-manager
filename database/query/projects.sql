-- name: ListProjectsFiltered :many
SELECT
  id,
  name,
  description,
  remark,
  status,
  manager_id,
  priority,
  started_at,
  duration,
  finished_at,
  remaining_days
FROM projects_view
WHERE 
  (started_at >= sqlc.narg('start_at') OR sqlc.narg('start_at') IS NULL)
  AND (started_at <= sqlc.narg('end_at') OR sqlc.narg('end_at') IS NULL)
ORDER BY
  priority DESC,
  remaining_days NULLS LAST,
  updated_at DESC;

-- name: CreateProject :exec
INSERT INTO projects (
  name,
  description,
  remark,
  manager_id,
  priority,
  duration
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
);

-- name: UpdateProject :exec
UPDATE projects
SET
  name = COALESCE(sqlc.narg('name'),name),
  description = COALESCE(sqlc.narg('description'),description),
  remark = COALESCE(sqlc.narg('remark'),remark),
  status = COALESCE(sqlc.narg('status'),status),
  manager_id = COALESCE(sqlc.narg('manager_id'),manager_id),
  priority = COALESCE(sqlc.narg('priority'),priority),
  started_at = COALESCE(sqlc.narg('started_at'),started_at),
  duration = COALESCE(sqlc.narg('duration'),duration),
  finished_at = COALESCE(sqlc.narg('finished_at'),finished_at)
WHERE id = sqlc.arg('id');

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;

-- name: StartProject :one
UPDATE projects
SET
  status = '进行',
  started_at = CURRENT_TIMESTAMP
WHERE
  id = $1
  AND status = '排单'
RETURNING id;

-- name: FinishProject :exec
UPDATE projects
SET
  status = '完成',
  finished_at = CURRENT_TIMESTAMP
WHERE
  id = $1;

-- name: GetProjectByID :one
SELECT
  p.name,
  p.description,
  p.remark,
  p.status,
  COALESCE(e.name,'') AS manager_name,
  p.priority,
  p.remaining_days
FROM projects_view p
LEFT JOIN employees e ON p.manager_id = e.id
WHERE p.id = $1;
