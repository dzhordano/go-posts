build:
	@go build -o bin/goposts app/main.go

run: build
	@./bin/goposts
