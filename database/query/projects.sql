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

-- name: UpdateProject :execrows
UPDATE projects
SET
  name = COALESCE(sqlc.narg('name'),name),
  description = COALESCE(sqlc.narg('description'),description),
  remark = COALESCE(sqlc.narg('remark'),remark),
  status = COALESCE(sqlc.narg('status'),status),
  manager_id = COALESCE(sqlc.narg('manager_id'),manager_id),
  priority = COALESCE(sqlc.narg('priority'),priority),
  duration = COALESCE(sqlc.narg('duration'),duration)
WHERE id = sqlc.arg('id');

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;

-- name: UpdateProjectStatus :execrows
UPDATE projects
SET
  status = sqlc.arg('target_status')
WHERE
  id = $1
  AND status = sqlc.arg('current_status');

-- name: GetProjectByID :one
SELECT
  name,
  description,
  remaining_days
FROM projects_view
WHERE id = $1;

-- name: ResetProject :execrows
UPDATE projects
SET
  status = '排单',
  started_at = NULL,
  finished_at = NULL
WHERE id = $1;

-- name: CopyProject :one
INSERT INTO projects (
  name,
  description,
  remark,
  manager_id,
  priority,
  duration
)
SELECT
  $2,
  p.description,
  p.remark,
  $3,
  p.priority,
  p.duration
FROM projects p
WHERE p.id = $1
RETURNING projects.id; -- 返回新的id,方便复制processes

-- name: IsProjectExists :one
SELECT EXISTS (
  SELECT 1
  FROM projects
  WHERE name = $1
) AS is_project_exists;
