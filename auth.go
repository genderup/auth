package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/inappcloud/jsonapi"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

type body struct {
	Data []User `json:"data"`
}

func Mux(db *sql.DB) http.Handler {
	m := web.New()
	m.NotFound(jsonapi.NotFoundHandler)
	m.Use(jsonapi.ContentTypeHandler)
	m.Use(middleware.EnvInit)
	m.Use(func(c *web.C, next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Env["db"] = db
			next.ServeHTTP(w, r)
		})
	})

	m.Use(bodyParserHandler)
	m.Use(dataCheckerHandler)

	m.Post("/users", userCreationHandler)
	m.Post("/sessions", sessionCreationHandler)
	m.Get("/users/me", currentUserHandler)

	return m
}

func bodyParserHandler(c *web.C, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			c.Env["body"] = new(body)
			jsonapi.BodyParserHandler(c.Env["body"], next).ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func dataCheckerHandler(c *web.C, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if len(c.Env["body"].(*body).Data) == 0 {
				jsonapi.Error(w, jsonapi.ErrNoData)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func newBody(users ...User) *body {
	return &body{users}
}

func userCreationHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	u := &c.Env["body"].(*body).Data[0]
	ur := &UserRepo{u, c.Env["db"].(*sql.DB)}

	err := ur.Create()
	if err != nil {
		jsonapi.Error(w, err)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(newBody(User{Id: u.Id, Email: u.Email, Token: u.Token}))
}

func sessionCreationHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	reqUser := c.Env["body"].(*body).Data[0]
	u := new(User)
	ur := &UserRepo{u, c.Env["db"].(*sql.DB)}

	ur.FetchByEmail(reqUser.Email)

	err := u.Authenticate(reqUser.Password)
	if err != nil {
		jsonapi.Error(w, err)
		return
	}

	u.GenerateToken()

	err = ur.Update()
	if err != nil {
		jsonapi.Error(w, err)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(newBody(User{Id: u.Id, Email: u.Email, Token: u.Token}))
}

func currentUserHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	u := new(User)
	ur := &UserRepo{u, c.Env["db"].(*sql.DB)}

	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) ([]byte, error) {
		return []byte(os.Getenv("PRIVATE_KEY")), nil
	})

	if err != nil {
		jsonapi.Error(w, jsonapi.ErrUnauthorized)
		return
	}

	err = ur.FetchByEmail(token.Claims["email"].(string))

	if err == nil && token.Valid {
		json.NewEncoder(w).Encode(newBody(User{Id: u.Id, Email: u.Email, Token: u.Token}))
	} else {
		jsonapi.Error(w, jsonapi.ErrUnauthorized)
	}
}
