# In-App Cloud Auth

[![Build Status](https://travis-ci.org/inappcloud/auth.svg?branch=master)](https://travis-ci.org/inappcloud/auth)

An authentication microservice.
You can use it directly as a REST API or use the package as part of your Go web application.
It is built with Goji and requires PostgreSQL.
`cmd/auth-server/main.go` will show you how to use the package in your own code.

# API

You can call the API directly with curl to test it.

## POST /users

name | type | description
---- | ---- | -----------
email | string | **Required**. A valid email
password | string | **Required**. A password

``` json
{
  "data": [{
    "email": "john@doe.com",
    "password": "p@ssw0rd"
  }]
}
```

## POST /sessions

``` json
{
  "data": [{
    "email": "john@doe.com",
    "password": "p@ssw0rd"
  }]
}
```

## GET /users/me

You must set the `Authorization` header with value `Bearer token` where `token` is the value you will get from the other endpoints.

## Response for all endpoints

name | type | description
---- | ---- | -----------
email | string | User's email
password | string | User's password
token | string | A JWT token

``` json
{
  "data": [{
    "id": 1,
    "email": "john@doe.com",
    "token": "jnshionefv2r3rnvdfoi239"
  }]
}
```

## Errors

**TODO**: Add documentation

# Run

## Locally

```
make setup
make server
```

If you have a custom environment for PostgreSQL, please update variable `DATABASE_URL` in file `.env`.
It runs on port `8080` by default but you can add `PORT` in `.env` to change the port.

# Test

```
make
```

If you have a custom environment for PostgreSQL, please update `TEST_DB` in file `Makefile`.

# FAQ

## Is it ready for production?

No. It need more testing, refinements and features before a real release.
