-- name: ListTodosFiltered :many
SELECT
  t.id,
  t.name,
  t.description,
  t.remark,
  t.status,
  t.priority,
  COALESCE(r.name,'') AS reporter_name,
  COALESCE(a.name,'') AS assignee_name,
  t.started_at,
  t.finished_at
FROM todos t
LEFT JOIN employees r ON t.reporter_id = r.id
LEFT JOIN employees a ON t.assignee_id = a.id
WHERE 
  t.status = ANY(sqlc.arg('status')::TEXT[]::todo_status[])
  AND (t.reporter_id = sqlc.narg('reporter_id') OR sqlc.narg('reporter_id') IS NULL)
  AND (t.assignee_id = sqlc.narg('assignee_id') OR sqlc.narg('assignee_id') IS NULL)
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

