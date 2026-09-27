CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  full_name VARCHAR(42) NOT NULL,
  age VARCHAR(3) NOT NULL,
  phone_number VARCHAR(20),
  habits VARCHAR(200),
  alive BOOLEAN NOT NULL,
  created_at TIMESTAMP NOT NULL
);

-- migrate -path migrations -database "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" up 1