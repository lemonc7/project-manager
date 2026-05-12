-- name: ListTodosForCreator :many
SELECT
  t.id,
  t.name,
  t.description,
  t.remark,
  s.name AS status,
  t.priority,
  t.assignee_id,
  e.name AS assignee_name,
  t.started_at,
  t.finished_at
FROM todos t
INNER JOIN employees e ON t.assignee_id = e.id
INNER JOIN todo_status s ON t.status_id = s.id
WHERE
  t.creator_id = sqlc.arg('creator_id')
ORDER BY
  t.priority ASC,
  s.sort_order ASC,
  t.updated_at DESC;

-- name: ListTodosForAssignee :many
SELECT
  t.id,
  t.name,
  t.description,
  t.remark,
  s.name AS status,
  t.priority,
  t.creator_id,
  e.name AS creator_name,
  t.started_at,
  t.finished_at
FROM todos t
INNER JOIN employees e ON t.creator_id = e.id
INNER JOIN todo_status s ON t.status_id = s.id
WHERE
  t.assignee_id = sqlc.arg('assignee_id')
ORDER BY
  t.priority ASC,
  s.sort_order ASC,
  t.updated_at DESC;


-- name: CreateTodo :exec
INSERT INTO todos (
  name,
  description,
  remark,
  priority,
  status_id,
  creator_id,
  assignee_id
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7
);

-- name: UpdateTodo :execrows
UPDATE todos
SET
  name = COALESCE(sqlc.narg('name'), name),
  description = COALESCE(sqlc.narg('description'), description),
  remark = COALESCE(sqlc.narg('remark'), remark),
  priority = COALESCE(sqlc.narg('priority'), priority),
  status_id = COALESCE(sqlc.narg('status_id'), status_id),
  creator_id = COALESCE(sqlc.narg('creator_id'), creator_id),
  assignee_id = COALESCE(sqlc.narg('assignee_id'), assignee_id),
  started_at = COALESCE(sqlc.narg('started_at'), started_at),
  finished_at = COALESCE(sqlc.narg('finished_at'), finished_at)
WHERE id = sqlc.arg('id');

-- name: DeleteTodo :execrows
DELETE FROM todos WHERE id = $1;
