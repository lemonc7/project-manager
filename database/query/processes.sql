-- name: ListProcesses :many
SELECT
  id,
  project_id,
  name,
  description,
  remark,
  status,
  employee_id,
  sequence,
  spend_time,
  started_at,
  finished_at
FROM processes p
WHERE 
  project_id = $1
ORDER BY sequence;

-- name: ListTasksTodo :many
SELECT
  p.id,
  p.name,
  p.description,
  p.project_id,
  COALESCE(pr.name,'') AS project_name
FROM processes p
LEFT JOIN projects pr ON p.project_id = pr.id
WHERE
  p.employee_id = $1
  AND p.status NOT IN ('终止','完成')
  AND pr.status NOT IN ('终止','完成')
ORDER BY p.id;

-- name: ListTasksReady :many
SELECT
  p.id,
  p.name,
  p.description,
  p.project_id,
  COALESCE(pr.name,'') AS project_name
FROM processes p
LEFT JOIN projects pr ON p.project_id = pr.id
WHERE
  p.employee_id = $1
  AND p.status IN ('准备', '暂停')
  AND pr.status NOT IN ('终止','完成')
ORDER BY p.id;

-- name: ListTasksDoing :many
SELECT
  p.id,
  p.name,
  p.description,
  p.project_id,
  COALESCE(pr.name,'') AS project_name,
  pr.remaining_days AS remaining_days
FROM processes p
LEFT JOIN projects_view pr ON p.project_id = pr.id
WHERE
  p.employee_id = $1
  AND p.status = '进行'
  AND pr.status NOT IN ('终止','完成')
ORDER BY 
  remaining_days NULLS LAST,
  p.id;

-- name: UpdateProcess :exec
UPDATE processes
SET
  name = COALESCE(sqlc.narg('name'),name),
  description = COALESCE(sqlc.narg('description'),description),
  remark = COALESCE(sqlc.narg('remark'),remark),
  status = COALESCE(sqlc.narg('status'),status),
  employee_id = COALESCE(sqlc.narg('employee_id'),employee_id),
  spend_time = COALESCE(sqlc.narg('spend_time'),spend_time),
  started_at = COALESCE(sqlc.narg('started_at'),started_at),
  finished_at = COALESCE(sqlc.narg('finished_at'),finished_at)
WHERE id = sqlc.arg('id');

-- name: InsertProcessLast :exec
WITH project_lock AS (
  SELECT id FROM projects WHERE id = $1 FOR UPDATE
),
max_seq AS (
  SELECT COALESCE(MAX(sequence), 0) AS val
  FROM processes
  WHERE project_id = $1
)
INSERT INTO processes (
  project_id,
  name,
  description,
  remark,
  employee_id,
  sequence
) SELECT
  $1,
  $2,
  $3,
  $4,
  $5,
  m.val + 1
FROM max_seq m, project_lock pl;

-- name: ShiftProcessesUp :exec
WITH project_lock AS (
  SELECT id FROM projects WHERE id = $1 FOR UPDATE
),
lock_rows AS (
  SELECT id
  FROM processes
  WHERE project_id = $1 AND sequence >= $2
  FOR UPDATE
)
UPDATE processes
SET sequence = sequence + 1
WHERE processes.project_id = $1 AND processes.sequence >= $2;

-- name: GetMaxSequence :one
SELECT COALESCE(MAX(sequence), 0)::INT AS max_seq
FROM processes
WHERE project_id = $1;

-- name: InsertProcessAt :exec
INSERT INTO processes (
  project_id,
  name,
  description,
  remark,
  employee_id,
  sequence
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
);

-- name: DeleteProcess :one
DELETE FROM processes
WHERE id = $1
RETURNING project_id;

-- name: ReorderProcesses :exec
WITH project_lock AS (
  SELECT id FROM projects WHERE id = $1 FOR UPDATE
)
UPDATE processes AS p
SET sequence = sub.new_seq
FROM (
  SELECT id,ROW_NUMBER() OVER (ORDER BY sequence) AS new_seq
  FROM processes
  WHERE processes.project_id = $1
) AS sub
WHERE p.id = sub.id;

-- name: FinishProcess :one
UPDATE processes
SET 
  status = '完成',
  finished_at = CURRENT_TIMESTAMP,
  spend_time = $2
WHERE
  id = $1
  AND status = '进行'
RETURNING id;

