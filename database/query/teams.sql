-- name: ListTeams :many
SELECT
  id,
  name,
  description,
  created_at,
  updated_at
FROM teams
ORDER BY name;

-- name: GetTeamByID :one
SELECT
  id,
  name,
  description,
  created_at,
  updated_at
FROM teams
WHERE id = $1;

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
  name = $2,
  description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM teams
WHERE id = $1;
