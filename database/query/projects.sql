-- name: ListProjects :many
SELECT
  id,
  name,
  description,
  remark,
  status,
  manager_id,
  priority,
  started_at,
  duration_days,
  finished_at,
  remaining_days,
  category
FROM projects_view
WHERE
  status = ANY(sqlc.arg('status')::TEXT[])
  AND (manager_id = sqlc.narg('manager_id') OR sqlc.narg('manager_id') IS NULL)
  AND (started_at >= sqlc.narg('start_at') OR sqlc.narg('start_at') IS NULL)
  AND (started_at <= sqlc.narg('end_at') OR sqlc.narg('end_at') IS NULL)
ORDER BY
  priority ASC,
  remaining_days NULLS LAST,
  updated_at DESC;

-- name: CreateProject :exec
INSERT INTO projects (
  name,
  description,
  remark,
  status_id,
  manager_id,
  priority,
  duration_days
) VALUES (
  $1,
  $2,
  $3,
  1,
  $4,
  $5,
  $6,
);

-- name: UpdateProject :execrows
UPDATE projects
SET
  name = COALESCE(sqlc.narg('name'), name),
  description = COALESCE(sqlc.narg('description'), description),
  remark = COALESCE(sqlc.narg('remark'), remark),
  status_id = COALESCE(sqlc.narg('status'), status_id),
  manager_id = COALESCE(sqlc.narg('manager_id'), manager_id),
  priority = COALESCE(sqlc.narg('priority'), priority),
  started_at = COALESCE(sqlc.narg('started_at'), started_at),
  duration_days = COALESCE(sqlc.narg('duration'), duration_days),
  finished_at = COALESCE(sqlc.narg('finished_at'), finished_at)
WHERE id = sqlc.arg('id');

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;

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
  status_id = 1,
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
  duration_days
)
SELECT
  $2,
  p.description,
  p.remark,
  $3,
  p.priority,
  p.duration_days
FROM projects p
WHERE p.id = $1
RETURNING projects.id; -- 返回新的id,方便复制processes

-- name: IsProjectExists :one
SELECT EXISTS (
  SELECT 1
  FROM projects
  WHERE name = $1
) AS is_project_exists;
