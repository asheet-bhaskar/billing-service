CREATE TABLE customers (
    id VARCHAR(36) PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email     VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT timezone('UTC', NOW()),
    updated_at TIMESTAMP NOT NULL DEFAULT timezone('UTC', NOW())
);