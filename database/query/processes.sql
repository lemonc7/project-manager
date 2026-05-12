-- name: ListProcesses :many
SELECT
  id,
  name,
  description,
  remark,
  status_id,
  employee_id,
  step_order,
  spend_hours,
  started_at,
  finished_at
FROM processes
WHERE
  project_id = $1
ORDER BY step_order;

-- name: ListTasks :many
SELECT
  p.step_order,
  p.name,
  p.description,
  p.project_id,
  pr.name AS project_name,
  pr.remaining_days AS remaining_days
FROM processes p
INNER JOIN projects_view pr ON p.project_id = pr.id
INNER JOIN process_status ps ON p.status_id = ps.id
WHERE
  p.employee_id = $1
  AND ps.category = $2
ORDER BY
  remaining_days NULLS LAST,
  p.id;

-- name: UpdateProcess :execrows
UPDATE processes
SET
  name = COALESCE(sqlc.narg('name'), name),
  description = COALESCE(sqlc.narg('description'), description),
  remark = COALESCE(sqlc.narg('remark'), remark),
  status_id = COALESCE(sqlc.narg('status_id'), status_id),
  employee_id = COALESCE(sqlc.narg('employee_id'), employee_id),
  spend_hours = COALESCE(sqlc.narg('spend_hours'), spend_hours)
WHERE id = sqlc.arg('id');

-- name: InsertProcessLast :exec
WITH project_lock AS (
  SELECT id FROM projects WHERE id = $1 FOR UPDATE
),
max_order AS (
  SELECT COALESCE(MAX(step_order), 0) AS val
  FROM processes
  WHERE project_id = $1
)
INSERT INTO processes (
  project_id,
  name,
  description,
  remark,
  status_id,
  employee_id,
  step_order
) SELECT
  $1,
  $2,
  $3,
  $4,
  1,
  $5,
  m.val + 1
FROM max_order m, project_lock pl;

-- 用于InsertProcessLast,获取最大step_order
-- name: GetMaxOrder :one
SELECT COALESCE(MAX(step_order), 0)::INT AS max_order
FROM processes
WHERE project_id = $1;

-- name: InsertProcessAt :exec
INSERT INTO processes (
  project_id,
  name,
  description,
  remark,
  status_id,
  employee_id,
  step_order
) VALUES (
  $1,
  $2,
  $3,
  $4,
  1,
  $5,
  $6
);

-- 用于InsertProcessAt时,调整step_order
-- name: ShiftProcessesUp :exec
WITH project_lock AS (
  SELECT id FROM projects WHERE id = $1 FOR UPDATE
)
UPDATE processes
SET step_order = step_order + 1
FROM project_lock
WHERE project_id = $1 AND step_order >= $2;

-- name: DeleteProcess :one
DELETE FROM processes
WHERE id = $1
RETURNING project_id;

-- name: DeleteProcesses :execrows
DELETE FROM processes
WHERE
  project_id = $1
  AND step_order = ANY($2::INT[]);

-- 用于DeleteProcess,调整step_order
-- name: ReorderProcesses :exec
WITH project_lock AS (
  SELECT id FROM projects WHERE id = $1 FOR UPDATE
)
UPDATE processes AS p
SET step_order = sub.new_order
FROM project_lock, (
  SELECT id, ROW_NUMBER() OVER (ORDER BY step_order) AS new_order
  FROM processes
  WHERE processes.project_id = $1
) AS sub
WHERE p.id = sub.id;

-- 用于同步更新process的status
-- name: GetNextTwoProcesses :many
SELECT
  id,
  status_id
FROM processes
WHERE
  project_id = $1
  AND step_order > $2
ORDER BY step_order
LIMIT 2;

-- name: GetProcessByID :one
SELECT
  project_id,
  name,
  description,
  employee_id,
  step_order,
  started_at
FROM processes
WHERE id = $1;

-- name: InsertReworkProcess :exec
INSERT INTO processes (
  project_id,
  name,
  description,
  remark,
  status_id,
  employee_id,
  step_order,
  spend_hours,
  started_at,
  finished_at
) VALUES (
  $1,
  $2,
  $3,
  $4,
  1,
  $5,
  $6,
  $7,
  $8,
  CURRENT_TIMESTAMP
);

-- 用于重置项目所有流程的状态
-- name: ResetProcessByProjectID :execrows
UPDATE processes
SET
  status_id = 1,
  spend_time = NULL,
  started_at = NULL,
  finished_at = NULL
WHERE project_id = $1;

-- 用于复制项目所有流程,并重置状态
-- name: CopyProcessesByProjectID :execrows
INSERT INTO processes (
  project_id,
  name,
  description,
  remark,
  status_id,
  employee_id,
  step_order
)
SELECT
  sqlc.arg('target_project'),
  p.name,
  p.description,
  p.remark,
  1,
  p.employee_id,
  p.step_order
FROM processes p
WHERE p.project_id = sqlc.arg('original_project');


-- name: UpdateProcessStatus :execrows
UPDATE processes
SET
  status_id = sqlc.arg('target_status')
WHERE
  id = sqlc.arg('id')
  AND status_id = sqlc.arg('current_status');
