-- name: CreateProject :one
INSERT INTO projects (
  project_name, start_date, end_date, budget
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetProject :one
SELECT * FROM projects
WHERE id = $1 
LIMIT 1;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;

-- name: ListProjects :many
SELECT * FROM projects
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateProject :one
UPDATE projects
SET project_name = $2,
    start_date = $3,
    end_date = $4,
    budget = $5
WHERE id = $1
RETURNING *;
