-- name: ListTeams :many
SELECT
  id,
  name,
  description,
  created_at,
  updated_at
FROM teams
ORDER BY name;

-- name: CreateTeam :one
INSERT INTO teams (
  name,
  description
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: UpdateTeam :one
UPDATE teams
SET
  name = COALESCE(sqlc.narg('name'),name),
  description = COALESCE(sqlc.narg('description'),description)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM teams
WHERE id = $1;
