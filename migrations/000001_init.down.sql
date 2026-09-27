DROP TABLE users;

-- migrate -path migrations -database "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" down 1