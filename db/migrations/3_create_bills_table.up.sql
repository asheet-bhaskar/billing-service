CREATE TABLE bills (
    id SERIAL PRIMARY KEY,
    description VARCHAR(100) NOT NULL,
    customer_id SERIAL NOT NULL,
    currency_id VARCHAR(3) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('open', 'closed')),
    total_amount DECIMAL(18, 2) DEFAULT 0.00,
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT timezone('UTC', NOW()),
    updated_at TIMESTAMP DEFAULT timezone('UTC', NOW()),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
    FOREIGN KEY (currency_id) REFERENCES currencies(id) ON DELETE CASCADE
);