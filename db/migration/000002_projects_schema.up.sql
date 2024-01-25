CREATE TABLE IF NOT EXISTS "projects" (
    "id" UUID PRIMARY KEY,
    "project_name" VARCHAR(100) NOT NULL,
    "start_date" DATE,
    "end_date" DATE,
    "budget" DECIMAL(15, 2)
);