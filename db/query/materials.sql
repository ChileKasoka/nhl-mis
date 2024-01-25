-- name: CreateMaterial :one
INSERT INTO materials (
  material_name, quantity, unit, cost_per_unit, project_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetMaterial :one
SELECT * FROM materials
WHERE material_id = $1
LIMIT 1;

-- name: DeleteMaterial :exec
DELETE FROM materials
WHERE material_id = $1;

-- name: ListMaterials :many
SELECT * FROM materials
ORDER BY material_id
LIMIT $1
OFFSET $2;

-- name: UpdateMaterial :one
UPDATE materials
SET material_name = $2,
    quantity = $3,
    unit = $4,
    cost_per_unit = $5,
    project_id = $6
WHERE material_id = $1
RETURNING *;
