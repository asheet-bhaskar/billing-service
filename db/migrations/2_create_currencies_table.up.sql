CREATE TABLE currencies (
    id VARCHAR(36) PRIMARY KEY,
    code CHAR(3) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT timezone('UTC', NOW()),
    updated_at TIMESTAMP DEFAULT timezone('UTC', NOW())
);