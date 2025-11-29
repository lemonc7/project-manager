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
  finished_at
FROM projects
WHERE 
  (status = sqlc.narg('status') OR sqlc.narg('status') IS NULL)
  AND (manager_id = sqlc.narg('manager_id') OR sqlc.narg('manager_id') IS NULL)
  AND (started_at >= sqlc.narg('start_at') OR sqlc.narg('start_at') IS NULL)
  AND (started_at <= sqlc.narg('end_at') OR sqlc.narg('end_at') IS NULL)
ORDER BY
  priority,
  CASE status
    WHEN '进行' THEN 1
    WHEN '排单' THEN 2
    WHEN '暂停' THEN 3
    ELSE 4
  END,
  started_at NULLS LAST,
  id;

-- name: ListProjectsTodo :many
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
  finished_at
FROM projects
WHERE 
  status IN ('暂停','排单')
  AND (manager_id = sqlc.narg('manager_id') OR sqlc.narg('manager_id') IS NULL)
ORDER BY
  priority,
  CASE status
    WHEN '暂停' THEN 1
    WHEN '排单' THEN 2
  END,
  id;

-- name: ListProjectsDoing :many
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
  finished_at
FROM projects
WHERE 
  status = '进行'
  AND (manager_id = sqlc.narg('manager_id') OR sqlc.narg('manager_id') IS NULL)
  AND (started_at >= sqlc.narg('start_at') OR sqlc.narg('start_at') IS NULL)
  AND (started_at <= sqlc.narg('end_at') OR sqlc.narg('end_at') IS NULL)
ORDER BY
  priority,
  started_at NULLS LAST,
  id;

-- name: ListProjectsDone :many
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
  finished_at
FROM projects
WHERE 
  status = '完成'
  AND (manager_id = sqlc.narg('manager_id') OR sqlc.narg('manager_id') IS NULL)
  AND (started_at >= sqlc.narg('start_at') OR sqlc.narg('start_at') IS NULL)
  AND (started_at <= sqlc.narg('end_at') OR sqlc.narg('end_at') IS NULL)
ORDER BY
  priority,
  started_at NULLS LAST,
  id;

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