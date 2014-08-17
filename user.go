package auth

import (
	"database/sql"
	"errors"
	"os"
	"regexp"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/inappcloud/jsonapi"

	"code.google.com/p/go.crypto/bcrypt"
)

const (
	PasswordBlank        = "password is blank"
	EmailBlank           = "email is blank"
	EmailInvalid         = "email format is invalid"
	EmailExists          = "email already exists"
	EmailNotFound        = "email does not exist"
	EmailPasswordInvalid = "email and/or password is invalid"
	TokenNotFound        = "token does not exist"
)

const (
	emailFormat = `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`
	cost        = 12
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

type UserRepo struct {
	user *User
	db   *sql.DB
}

func (u *User) EncryptPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), cost)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	return nil
}

func (u *User) Authenticate(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return jsonapi.ErrInvalidParams(EmailPasswordInvalid)
	}

	return nil
}

func (u *User) GenerateToken() {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["email"] = u.Email
	tokenString, _ := token.SignedString([]byte(os.Getenv("PRIVATE_KEY")))
	u.Token = tokenString
}

func (r *UserRepo) Validate() error {
	r.user.Email = normalizeEmail(r.user.Email)

	errs := []string{}

	if r.user.Password == "" {
		errs = append(errs, PasswordBlank)
	}

	if r.user.Email == "" {
		errs = append(errs, EmailBlank)
	}

	if !regexp.MustCompile(emailFormat).MatchString(r.user.Email) {
		errs = append(errs, EmailInvalid)
	}

	var count int
	r.db.QueryRow(`SELECT COUNT(id) FROM users WHERE email = $1`, r.user.Email).Scan(&count)

	if count > 0 {
		errs = append(errs, EmailExists)
	}

	if len(errs) > 0 {
		return jsonapi.ErrInvalidParams(errs...)
	}

	return nil
}

func (r *UserRepo) Create() error {
	err := r.Validate()
	if err != nil {
		return err
	}

	err = r.user.EncryptPassword()
	if err != nil {
		return err
	}

	r.user.GenerateToken()

	return r.db.QueryRow(`INSERT INTO users (email, password, token) VALUES($1, $2, $3) RETURNING id`, r.user.Email, r.user.Password, r.user.Token).Scan(&r.user.Id)
}

func (r *UserRepo) Update() error {
	_, err := r.db.Exec(`UPDATE users SET token = $1 WHERE id = $2`, r.user.Token, r.user.Id)
	return err
}

func (r *UserRepo) FetchByEmail(email string) error {
	err := r.db.QueryRow(`SELECT id, email, password, token FROM users WHERE email = $1`, normalizeEmail(email)).Scan(&r.user.Id, &r.user.Email, &r.user.Password, &r.user.Token)

	if err != nil {
		return err
	}

	if r.user.Id == 0 {
		return errors.New(EmailNotFound)
	}

	return nil
}

func normalizeEmail(email string) string {
	return string(regexp.MustCompile(`\s+`).ReplaceAll([]byte(strings.ToLower(email)), []byte{}))
}
