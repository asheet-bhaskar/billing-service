CREATE TABLE line_items (
    id VARCHAR(36) PRIMARY KEY,
    bill_id VARCHAR(36) NOT NULL,
    description TEXT NOT NULL,
    amount DECIMAL(18, 2) NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMP DEFAULT timezone('UTC', NOW()),
    removed BOOLEAN DEFAULT false,
    FOREIGN KEY (bill_id) REFERENCES bills(id) ON DELETE CASCADE
);