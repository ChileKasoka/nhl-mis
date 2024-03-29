// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: projects.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createProject = `-- name: CreateProject :one
INSERT INTO projects (
  project_name, start_date, end_date, budget
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, project_name, start_date, end_date, budget
`

type CreateProjectParams struct {
	ProjectName string         `json:"project_name"`
	StartDate   sql.NullTime   `json:"start_date"`
	EndDate     sql.NullTime   `json:"end_date"`
	Budget      sql.NullString `json:"budget"`
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, createProject,
		arg.ProjectName,
		arg.StartDate,
		arg.EndDate,
		arg.Budget,
	)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.ProjectName,
		&i.StartDate,
		&i.EndDate,
		&i.Budget,
	)
	return i, err
}

const deleteProject = `-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1
`

func (q *Queries) DeleteProject(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteProject, id)
	return err
}

const getProject = `-- name: GetProject :one
SELECT id, project_name, start_date, end_date, budget FROM projects
WHERE id = $1 
LIMIT 1
`

func (q *Queries) GetProject(ctx context.Context, id uuid.UUID) (Project, error) {
	row := q.db.QueryRowContext(ctx, getProject, id)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.ProjectName,
		&i.StartDate,
		&i.EndDate,
		&i.Budget,
	)
	return i, err
}

const listProjects = `-- name: ListProjects :many
SELECT id, project_name, start_date, end_date, budget FROM projects
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListProjectsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListProjects(ctx context.Context, arg ListProjectsParams) ([]Project, error) {
	rows, err := q.db.QueryContext(ctx, listProjects, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Project{}
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ID,
			&i.ProjectName,
			&i.StartDate,
			&i.EndDate,
			&i.Budget,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProject = `-- name: UpdateProject :one
UPDATE projects
SET project_name = $2,
    start_date = $3,
    end_date = $4,
    budget = $5
WHERE id = $1
RETURNING id, project_name, start_date, end_date, budget
`

type UpdateProjectParams struct {
	ID          uuid.UUID      `json:"id"`
	ProjectName string         `json:"project_name"`
	StartDate   sql.NullTime   `json:"start_date"`
	EndDate     sql.NullTime   `json:"end_date"`
	Budget      sql.NullString `json:"budget"`
}

func (q *Queries) UpdateProject(ctx context.Context, arg UpdateProjectParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, updateProject,
		arg.ID,
		arg.ProjectName,
		arg.StartDate,
		arg.EndDate,
		arg.Budget,
	)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.ProjectName,
		&i.StartDate,
		&i.EndDate,
		&i.Budget,
	)
	return i, err
}
