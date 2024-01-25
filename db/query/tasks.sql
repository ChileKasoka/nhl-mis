-- name: CreateTask :one
INSERT INTO tasks (
  task_name, description, start_date, end_date, project_id, employee_id, status
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE task_id = $1
LIMIT 1;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE task_id = $1;

-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY task_id
LIMIT $1
OFFSET $2;

-- name: UpdateTask :one
UPDATE tasks
SET task_name = $2,
    description = $3,
    start_date = $4,
    end_date = $5,
    project_id = $6,
    employee_id = $7,
    status = $8
WHERE task_id = $1
RETURNING *;
