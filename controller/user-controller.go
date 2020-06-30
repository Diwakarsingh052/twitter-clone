package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"twitter-clone/graph/model"
	"twitter-clone/rand"
)

// NewUsers is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewUsers(us *model.UserService) *Users {
	return &Users{
		us: us,
	}
}

type Users struct {
	us *model.UserService
}
type jsonData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *Users) SayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello  How are you")
}

// New is used to render the form where a user can
// create a new user account.
//
// GET /signup

// Create is used to process the signup form when a user
// submits it. This is used to create a new user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var u1 jsonData
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.Unmarshal(b, &u1)

	user := model.User{
		Email:    u1.Email,
		Password: u1.Password,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = u.signIn(w, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// Login is used to verify the provided email address and
// password and then log the user in if they are correct.
//
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var userDetail model.User
	json.Unmarshal(b, &userDetail)

	user, err := u.us.Authenticate(userDetail.Email, userDetail.Password)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address.")
		case model.ErrInvalidPassword:
			fmt.Fprintln(w, "Invalid password provided.")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// signIn is used to sign the given user in via cookies
func (u *Users) signIn(w http.ResponseWriter, user *model.User) error {
	if user.Remember == nil {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = &token
		err = u.us.Update(user)

		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:  "remember_token",
		Value: *user.Remember,
	}
	u.us.Token = *user.Remember
	http.SetCookie(w, &cookie)
	return nil
}

// CookieTest is used to display cookies set on the current user
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}
