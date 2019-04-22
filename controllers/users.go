package controllers

import (
	"net/http"

	"uapp.com/models"
	"uapp.com/views"

	"github.com/gorilla/schema"
)

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "/views/users/new.gohtml"),
		// LoginView: views.NewView("bootstrap", "/views/users/login.gohtml"),
		HomeView: views.NewView("bootstrap", "/views/users/home.gohtml"),
		us:       us,
	}
}

type Users struct {
	NewView *views.View
	// LoginView *views.View
	HomeView *views.View
	us       *models.UserService
}

// New is used to render the form from where a user can
// create a new user account
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// List us used to render the home page with the last users listed
func (u *Users) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	users, err := u.us.OrderRecent(5)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := u.HomeView.Render(w, users); err != nil {
		panic(err)
	}
}

// SignupForm stores data from the form where a user a user can
// register
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create is used to process the signup form when a user
// tries to create a new user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", 301)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(dst, r.PostForm); err != nil {
		panic(err)
	}
	return nil
}
