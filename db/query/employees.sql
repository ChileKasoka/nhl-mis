-- name: CreateEmployee :one
INSERT INTO employees (
  first_name, last_name, email, phone_number, hire_date, position, salary
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetEmployee :one
SELECT * FROM employees
WHERE employee_id = $1
LIMIT 1;

-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE employee_id = $1;

-- name: ListEmployees :many
SELECT * FROM employees
ORDER BY employee_id
LIMIT $1
OFFSET $2;

-- name: UpdateEmployee :one
UPDATE employees
SET first_name = $2,
    last_name = $3,
    email = $4,
    phone_number = $5,
    hire_date = $6,
    position = $7,
    salary = $8
WHERE employee_id = $1
RETURNING *;

-- name: GetUserEmployeeInfoByID :many
SELECT
    users.id AS user_id,
    users.first_name,
    users.last_name,
    users.email,
    employees.employee_id,
    employees.phone_number,
    employees.hire_date,
    employees.position,
    employees.salary
FROM
    users
INNER JOIN
    employees ON users.id = employees.user_id
WHERE
    users.id = $1;
