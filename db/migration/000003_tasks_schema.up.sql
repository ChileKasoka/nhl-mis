CREATE TABLE IF NOT EXISTS "tasks" (
    "task_id" UUID PRIMARY KEY,
    "task_name" VARCHAR(100) NOT NULL,
    "description" TEXT,
    "start_date" DATE,
    "end_date" DATE,
    "project_id" UUID,
    "employee_id" UUID,
    "status" VARCHAR(20),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);