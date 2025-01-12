.PHONY: run
run:
	go run ./cmd/main.go -port ${port}
.PHONY: lint
lint:
	golangci-lint run
test:
	go test  ./... 
coverage:
	go test  -coverprofile=coverage.out ./... ;
	go tool cover -func=coverage.out
.PHONY: build
build:
	go build -o ./bin ./cmd
.PHONY: air
air:
	air -c .air.toml