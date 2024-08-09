build:
	@go build -o bin/jasad cmd/jasad/main.go

run: build
	@./bin/jasad

test:
	@go test -v ./...
