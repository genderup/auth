TEST?=./

default: test

server:
	env $$(cat .env) go run cmd/auth-server/main.go

test:
	env $$(cat .testenv) go test $(TEST)

updatedeps:
	go get -u -v ./...
