TEST?=./
TEST_DB=postgres://localhost/inappcloud_auth_test?sslmode=disable

default: test

setup: updatedeps
	psql -c 'create database inappcloud_auth;'
	psql -c 'create table users (id serial, email text, password text, token text);' -d inappcloud_auth
	echo "DATABASE_URL=postgres://localhost/inappcloud_auth?sslmode=disable PRIVATE_KEY=$$(openssl rand -hex 64)" > .env

server:
	env $$(cat .env) go run cmd/auth-server/main.go

updatedeps:
	go get -u -v ./...
	go get -u -v github.com/lib/pq

test:
	psql -c 'drop database if exists inappcloud_auth_test;'
	psql -c 'create database inappcloud_auth_test;'
	psql -c 'create table users (id serial, email text, password text, token text);' -d inappcloud_auth_test
	env DATABASE_URL=$(TEST_DB) go test $(TEST)
