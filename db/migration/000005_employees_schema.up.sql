CREATE TABLE IF NOT EXISTS employees (
    employee_id UUID PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    phone_number VARCHAR(15),
    hire_date DATE,
    position VARCHAR(50),
    salary DECIMAL(15, 2),
    user_id UUID,
    FOREIGN KEY (user_id) REFERENCES users(id)
);