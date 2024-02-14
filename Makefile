build:
	@go build -o bin/goposts cmd/app/main.go

run: build
	@./bin/goposts
