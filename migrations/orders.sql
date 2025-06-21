CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_name VARCHAR(255),
    amount NUMERIC,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);