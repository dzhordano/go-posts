build:
	@go build -o bin/goposts cmd/app/main.go

run: build
	@./bin/goposts

# CREATE TABLES
up:
	@migrate -path ./migrations -database 'postgres://postgres:qwerty@host.docker.internal:6000/postgres?sslmode=disable' up
# @migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:6000/postgres?sslmode=disable' up
# DROP TABLES
down:
	@migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5433/postgres?sslmode=disable' down
