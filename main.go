package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

//jwt token
var mySigningKey = []byte("secret")
var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["admin"] = true
	claims["name"] = "dafdgdsds"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	/* Finally, write the token to the browser window */
	w.Write([]byte(tokenString))
})

//password
type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Hash     string `json:"-"`
	Password string `json:"password"`
}

func Get(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

var ErrInvalidPassword = errors.New("Invalid Password")
var ErrPasswordMismatch = errors.New("Password cannot be blank")
var ErrEmptyPassword = errors.New("No password provided")

func ChangePassword(r *http.Request) error {
	u := Get(r, "user").(User)
	currentPw := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_new_password")
	// Check the current password
	err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(currentPw))
	if err != nil {
		return ErrInvalidPassword
	}
	// Check that the new password isn't blank
	if newPassword == "" {
		return ErrEmptyPassword
	}
	// Check that new passwords match
	if newPassword != confirmPassword {
		return ErrPasswordMismatch
	}
	// Generate the new hash
	h, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Hash = string(h)
	return nil
}

func main() {
	//Init Router
	r := mux.NewRouter()

	r.Handle("/get-token", GetTokenHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8001", r))
}
