CREATE TABLE IF NOT EXISTS "materials" (
    "material_id" UUID PRIMARY KEY,
    "material_name" VARCHAR(100) NOT NULL,
    "quantity" INT,
    "unit" INT,
    "cost_per_unit" DECIMAL(10, 2),
    "project_id" UUID,
    "status" VARCHAR(20),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);