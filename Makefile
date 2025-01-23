.PHONY: run
run:
	@echo "Loading environment variables from .env file"
	@set -o allexport; source ./load_env.sh; set +o allexport; \
	echo "Running Go application"; \
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
	@echo "Loading environment variables from .env file"
	@set -o allexport; source ./load_env.sh; set +o allexport; \
	echo "Running air"; \
	air -c .air.toml
swagger:
	swag init -g cmd/main.go