package auth_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"code.google.com/p/go.crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/inappcloud/auth"
	"github.com/inappcloud/jsonapi"
)

var privateKey = []byte{
	3, 35, 53, 75, 43, 15, 165, 188, 131, 126, 6, 101, 119, 123, 166,
	143, 90, 179, 40, 230, 240, 84, 201, 40, 169, 15, 132, 178, 210, 80,
	46, 191, 211, 251, 90, 146, 210, 6, 71, 239, 150, 138, 180, 195, 119,
	98, 61, 34, 61, 46, 33, 114, 5, 46, 79, 8, 192, 205, 154, 245, 103,
	208, 128, 163}

const (
	cost = 4
)

func init() {
	os.Setenv("PRIVATE_KEY", string(privateKey))
}

func createToken(email string) string {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["email"] = email
	tokenString, _ := token.SignedString(privateKey)
	return tokenString
}

func withToken(email string, res string) string {
	return fmt.Sprintf(res, createToken(email))
}

var createUserTestData = []struct {
	path    string
	body    string
	expCode int
	expBody string
}{
	{"/users", `{"data":[{"email":"","password":"password"}]}`, 422, err(jsonapi.ErrInvalidParams(auth.EmailBlank, auth.EmailInvalid))},
	{"/users", `{"data":[{"email":"foo@bar.com","password":""}]}`, 422, err(jsonapi.ErrInvalidParams(auth.PasswordBlank))},
	{"/users", `{"data":[{"email":"not_an_email","password":"password"}]}`, 422, err(jsonapi.ErrInvalidParams(auth.EmailInvalid))},
	{"/users", `{"data":[{"email":"first@bar.com","password":"password"}]}`, 422, err(jsonapi.ErrInvalidParams(auth.EmailExists))},
	{"/users", `{"data":[{"email":"FIRST@bar.com","password":"password"}]}`, 422, err(jsonapi.ErrInvalidParams(auth.EmailExists))},
	{"/users", `{"data":[]}`, 422, err(jsonapi.ErrNoData)},
	{"/users", `{"other_key":[{"email":"foo@bar.com","password":"password"}]}`, 422, err(jsonapi.ErrNoData)},
	{"/users", `{}`, 422, err(jsonapi.ErrNoData)},
	{"/users", ``, 400, err(jsonapi.ErrBadRequest)},
	{"/users", `{"data":[{"email":"second@bar.com","password":"password"}]}`, 201, withToken("second@bar.com", `{"data":[{"id":2,"email":"second@bar.com","token":"%s"}]}`+"\n")},
	{"/users", `{"data":[{"email":"third @bar.  com\t\n","password":"password"}]}`, 201, withToken("third@bar.com", `{"data":[{"id":3,"email":"third@bar.com","token":"%s"}]}`+"\n")},
}

func TestCreateUser(t *testing.T) {
	db := NewTestDB()
	defer db.Close()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), cost)
	db.Exec(`INSERT INTO users(email, password) VALUES('first@bar.com', $1)`, string(hash))

	for _, test := range createUserTestData {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", test.path, bytes.NewBufferString(test.body))

		auth.Mux(db).ServeHTTP(w, r)

		eq(t, test.expCode, w.Code)
		eq(t, test.expBody, w.Body.String())
	}
}

var createSessionTestData = []struct {
	path    string
	body    string
	expCode int
	expBody string
}{
	{"/sessions", `{"data":[{"email":"foo@bar.com","password":"password"}]}`, 201, withToken("foo@bar.com", `{"data":[{"id":1,"email":"foo@bar.com","token":"%s"}]}`+"\n")},
	{"/sessions", `{"data":[{"email":"FOO@bar.com","password":"password"}]}`, 201, withToken("foo@bar.com", `{"data":[{"id":1,"email":"foo@bar.com","token":"%s"}]}`+"\n")},
	{"/sessions", `{"data":[{"email":"","password":"password"}]}`, 422, err(jsonapi.ErrInvalidParams(auth.EmailPasswordInvalid))},
	{"/sessions", `{"data":[{"email":"foo@bar.com","password":""}]}`, 422, err(jsonapi.ErrInvalidParams(auth.EmailPasswordInvalid))},
	{"/sessions", `{"data":[{"email":"not_an_email","password":"password"}]}`, 422, err(jsonapi.ErrInvalidParams(auth.EmailPasswordInvalid))},
	{"/sessions", `{"data":[{"email":"notfoo@bar.com","password":"password"}]}`, 422, err(jsonapi.ErrInvalidParams(auth.EmailPasswordInvalid))},
	{"/sessions", `{"data":[{"email":"foo@bar.com","password":"password1"}]}`, 422, err(jsonapi.ErrInvalidParams(auth.EmailPasswordInvalid))},
	{"/sessions", `{"data":[]}`, 422, err(jsonapi.ErrNoData)},
	{"/sessions", `{"other_key":[{"email":"foo@bar.com","password":"password"}]}`, 422, err(jsonapi.ErrNoData)},
	{"/sessions", `{}`, 422, err(jsonapi.ErrNoData)},
	{"/sessions", ``, 400, err(jsonapi.ErrBadRequest)},
}

func TestCreateSession(t *testing.T) {
	db := NewTestDB()
	defer db.Close()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), cost)
	db.Exec(`INSERT INTO users (email, password) VALUES('foo@bar.com', $1)`, string(hash))

	for _, test := range createSessionTestData {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", test.path, bytes.NewBufferString(test.body))

		auth.Mux(db).ServeHTTP(w, r)

		eq(t, test.expCode, w.Code)
		eq(t, test.expBody, w.Body.String())
	}
}

var getCurrentUserTestData = []struct {
	headers map[string]string
	path    string
	expCode int
	expBody string
}{
	{map[string]string{"Authorization": fmt.Sprintf("Bearer %s", createToken("foo@bar.com"))}, "/users/me", 200, withToken("foo@bar.com", `{"data":[{"id":1,"email":"foo@bar.com","token":"%s"}]}`+"\n")},
	{map[string]string{"Authorization": "Bearer none"}, "/users/me", 401, err(jsonapi.ErrUnauthorized)},
	{map[string]string{"Authorization": fmt.Sprintf("Bearer %s", createToken("something@unknown.com"))}, "/users/me", 401, err(jsonapi.ErrUnauthorized)},
}

func TestGetCurrentUser(t *testing.T) {
	db := NewTestDB()
	defer db.Close()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), cost)
	db.Exec(`INSERT INTO users(email, password, token) VALUES('foo@bar.com', $1, $2)`, string(hash), createToken("foo@bar.com"))

	for _, test := range getCurrentUserTestData {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", test.path, nil)
		for k, v := range test.headers {
			r.Header.Set(k, v)
		}

		auth.Mux(db).ServeHTTP(w, r)

		eq(t, test.expCode, w.Code)
		eq(t, test.expBody, w.Body.String())
	}
}
