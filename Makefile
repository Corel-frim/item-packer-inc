.PHONY: swagger test linter run install run-docker

swagger:
	swag init -g cmd/main.go -o docs/swagger

test:
	go test -v ./...

linter:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

run:
	go run cmd/main.go

run-docker:
	docker build -t item-packer-inc .
	docker run --name item-packer-inc -p 8080:8080 -d --rm item-packer-inc

install:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
