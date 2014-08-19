# In-App Cloud Auth [![Build Status](https://travis-ci.org/inappcloud/auth.svg?branch=master)](https://travis-ci.org/inappcloud/auth)

An authentication microservice.
You can use it directly as a REST API or use the package as part of your Go web application.
It is built with Goji and requires PostgreSQL.
`cmd/auth-server/main.go` will show you how to use the package in your own code.

# API

It uses [JSON API](http://jsonapi.org) for request and response format.

## Creating a user

To create a user, you need the following information:

name | type | description
---- | ---- | -----------
email | string | **Required**. A valid email
password | string | **Required**. A password

Then you can make a request like this:

```
POST /users
Content-Type: application/json
Accept: application/json

{
  "data": [{
    "email": "john@doe.com",
    "password": "p@ssw0rd"
  }]
}
```

Endpoint will respond with `400 Bad Request` if the body is not JSON and with `422 Unprocessable Entity` if the body has invalid formatting or params.

## Log in

It's almost the same as user creation, the only change is the URI.

```
POST /sessions
Content-Type: application/json
Accept: application/json

{
  "data": [{
    "email": "john@doe.com",
    "password": "p@ssw0rd"
  }]
}
```

## Current user

You must set the `Authorization` header with value `Bearer token` where `token` is the value you will get from the other endpoints.

```
GET /users/me
Accept: application/json
Authorization: Bearer jnshionefv2r3rnvdfoi239
```

## Response for all endpoints

name  | type   | description
----- | ------ | -----------
id    | int    | User ID
email | string | User Email
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

Endpoint will respond with `401 Unauthorized` if the token is invalid.

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
