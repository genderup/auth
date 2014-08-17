# In-App Cloud Auth

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

You need to create a file named `.env` with the variables `DATABASE_URL` and `PRIVATE_KEY` like this:

```
DATABASE_URL=postgres://localhost/auth PRIVATE_KEY=asdfghjkl123456zxcvbnwertyu45678
```

Then you can use `make server`.

It runs on port `8080` by default but you can add the variable `PORT` in `.env` to change the port.

# Test

You need to create a file named `.testenv` with the variable `DATABASE_URL` like this:

```
DATABASE_URL=postgres://localhost/auth_test
```

Then you can use `make`.

# FAQ

## Is it ready for production?

No. It need more testing, refinements and features before a real release.
