-- name: ListTodosByReporter :many
SELECT
  t.id,
  t.name,
  t.description,
  t.remark,
  t.status,
  t.priority,
  t.assignee_id AS employee_id,
  COALESCE(e.name,'') AS employee_name,
  t.started_at,
  t.finished_at
FROM todos t
LEFT JOIN employees e ON t.assignee_id = e.id
WHERE
  t.reporter_id = sqlc.arg('reporter_id')
  AND t.status = ANY(sqlc.arg('status')::TEXT[])
ORDER BY
  t.priority DESC,
  t.updated_at DESC;

-- name: ListTodosByAssignee :many
SELECT
  t.id,
  t.name,
  t.description,
  t.remark,
  t.status,
  t.priority,
  t.reporter_id AS employee_id,
  COALESCE(e.name,'') AS employee_name,
  t.started_at,
  t.finished_at
FROM todos t
LEFT JOIN employees e ON t.reporter_id = e.id
WHERE 
  t.assignee_id = sqlc.arg('assignee_id')
  AND t.status = ANY(sqlc.arg('status')::TEXT[])
ORDER BY
  t.priority DESC,
  t.updated_at DESC;


-- name: CreateTodo :exec
INSERT INTO todos (
  name,
  description,
  remark,
  priority,
  reporter_id,
  assignee_id
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
);

-- name: UpdateTodo :execrows
UPDATE todos
SET
  name = COALESCE(sqlc.narg('name'),name),
  description = COALESCE(sqlc.narg('description'),description),
  remark = COALESCE(sqlc.narg('remark'),remark),
  priority = COALESCE(sqlc.narg('priority'),priority),
  status = COALESCE(sqlc.narg('status'),status),
  reporter_id = COALESCE(sqlc.narg('reporter_id'),reporter_id),
  assignee_id = COALESCE(sqlc.narg('assignee_id'),assignee_id)
WHERE id = sqlc.arg('id');

-- name: DeleteTodo :execrows
DELETE FROM todos
WHERE id = $1;

