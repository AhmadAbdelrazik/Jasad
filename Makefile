build:
	@go build -o bin/jasad

run: build
	@./bin/jasad

test:
	@go test -v ./...

play:
	@go run playground.go
