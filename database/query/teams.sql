-- name: ListTeams :many
SELECT
  id,
  name,
  description
FROM teams
ORDER BY name;

-- name: CreateTeam :exec
INSERT INTO teams (
  name,
  description
) VALUES (
  $1,
  $2
);

-- name: UpdateTeam :exec
UPDATE teams
SET
  name = COALESCE(sqlc.narg('name'), name),
  description = COALESCE(sqlc.narg('description'), description)
WHERE id = sqlc.arg('id');

-- name: DeleteTeam :exec
DELETE FROM teams WHERE id = $1;
