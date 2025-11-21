-- name: ListProcessesFiltered :many
SELECT
  p.id,
  p.name,
  p.sequence,
  p.status,
  COALESCE(pr.name,'') AS project_name,
  COALESCE(e.name,'') AS employee_name,
  p.description,
  p.remark,
  p.start_time,
  p.end_time,
  p.created_at,
  p.updated_at
FROM processes p
LEFT JOIN employees e ON p.employee_id = e.id
LEFT JOIN projects pr ON p.project_id = pr.id
WHERE 
  p.project_id = COALESCE(sqlc.narg('project_id'), p.project_id)
  AND p.status = COALESCE(sqlc.narg('status'), p.status)
  AND p.employee_id = COALESCE(sqlc.narg('employee_id'), p.employee_id)
ORDER BY p.sequence;

-- name: ListTasksTodo :many
SELECT
  p.id,
  p.name,
  p.project_id,
  COALESCE(pr.name,'') AS project_name
FROM processes p
LEFT JOIN projects pr ON p.project_id = pr.id
WHERE
  p.employee_id = $1
  AND p.status NOT IN ('取消', '终止')
ORDER BY p.id;

-- name: ListTasksReady :many
SELECT
  p.id,
  p.name,
  p.project_id,
  COALESCE(pr.name,'') AS project_name
FROM processes p
LEFT JOIN projects pr ON p.project_id = pr.id
WHERE
  p.employee_id = $1
  AND p.status IN ('准备', '暂停')
ORDER BY p.id;

-- name: ListTasksDoing :many
SELECT
  p.id,
  p.name,
  p.start_time,
  pr.duration,
  p.project_id,
  COALESCE(pr.name,'') AS project_name
FROM processes p
LEFT JOIN projects pr ON p.project_id = pr.id
WHERE
  p.employee_id = $1
  AND p.status IN ('进行', '返工')
ORDER BY p.id;

-- name: UpdateProcess :one
UPDATE processes
SET
  name = COALESCE(sqlc.narg('name'),name),
  description = COALESCE(sqlc.narg('description'),description),
  status = COALESCE(sqlc.narg('status'),status),
  employee_id = COALESCE(sqlc.narg('employee_id'),employee_id),
  remark = COALESCE(sqlc.narg('remark'),remark)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: InsertProcessLast :one
WITH project_lock AS (
  SELECT id FROM projects WHERE id = $1 FOR UPDATE
),
max_seq AS (
  SELECT COALESCE(MAX(sequence), 0) AS val
  FROM processes
  WHERE project_id = $1
)
INSERT INTO processes (
  name,
  project_id,
  employee_id,
  sequence,
  description,
  remark
) SELECT
  $2,
  $1,
  $3,
  m.val + 1,
  $4,
  $5
FROM max_seq m, project_lock pl
RETURNING *;

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
  name,
  project_id,
  employee_id,
  sequence,
  description,
  remark
) VALUES (
  $3,
  $1,
  $4,
  $2,
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

