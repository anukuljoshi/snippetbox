package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.anukuljoshi/internals/models"
	"snippetbox.anukuljoshi/internals/validator"
)

// struct to hold form data and embedded validator
// added struct tags for decoding form field names to struct fields
type snippetCreateForm struct {
	Title string `form:"title"`
	Content string `form:"content"`
	Expires int `form:"expires"`
	validator.Validator `form:"-"`
}

// struct to hold form data and embedded validator
// added struct tags for decoding form field names to struct fields
type userSignupForm struct {
	Name string `form:"name"`
	Email string `form:"email"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
}

// struct to hold form data and embedded validator
// added struct tags for decoding form field names to struct fields
type userLoginForm struct {
	Email string `form:"email"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
}

// handler for home
func (app *application) home(w http.ResponseWriter,  r *http.Request){
	snippets, err := app.snippets.Latest()
	if err!=nil {
		app.serverError(w, err)
		return
	}
	// call newTemplateData to create templateData with CurrentYear
	data := app.newTemplateData(r)
	data.Snippets = snippets
	// use the render helper method
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

// about section for app
func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "about.tmpl.html", data)
}

// handler for viewing a snippet
func (app *application) viewSnippet(w http.ResponseWriter,  r *http.Request){
	params := httprouter.ParamsFromContext(r.Context())
	// get id sent in query params
	id, err := strconv.Atoi(params.ByName("id"))	
	if err != nil || id < 1{
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err!=nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
	}
	// call newTemplateData to create templateData with CurrentYear
	data := app.newTemplateData(r)
	data.Snippet = snippet
	// use render helper method
	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	// initialize snippetCreateForm struct to pass to template
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

// handler for creating a snippet
func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {
	// initialize empty snippetCreateForm
	var form snippetCreateForm
	// call decodePostForm helper method to decode data into snippetCreateForm struct
	var err = app.decodePostForm(r, &form)
	if err!=nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// use our custom validator to check for validations
	// validations check for title
	// 1. title is not empty
	form.CheckField(
		validator.NotBlank(form.Title),
		"title",
		"This field cannot be blank",
	)
	// 2. title is less than 100 characters
	form.CheckField(
		validator.MaxLen(form.Title, 100),
		"title",
		"This field cannot be more than 100 characters long",
	)
	// validations check for content
	// 1. content is not empty
	form.CheckField(
		validator.NotBlank(form.Content),
		"content",
		"This field cannot be blank",
	)
	// validation checks for expires
	// expires should be either 1, 7 or 365
	form.CheckField(
		validator.PermittedValue(form.Expires, 1, 7, 365),
		"expires",
		"This field must be equal to 1, 7 or 365",
	)
	// return bad request if form.FieldErrors are present
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusBadRequest, "create.tmpl.html", data)
		return
	}
	// call insert for snippet model with data
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err!=nil {
		app.serverError(w, err)
		return
	}
	// use Put method of sessionManager to add a flash message to session
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created")
	// redirect to snippet view for the created snippet id
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// user handlers
func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignUpPost(w http.ResponseWriter, r *http.Request) {
	// initialize empty snippetCreateForm
	var form userSignupForm
	// call decodePostForm helper method to decode data into snippetCreateForm struct
	err := app.decodePostForm(r, &form)
	if err!=nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// use our custom validator to check for validations
	// validations check for title
	// 1. name is not empty
	form.CheckField(
		validator.NotBlank(form.Name),
		"name",
		"This field cannot be blank",
	)
	// email is not empty
	form.CheckField(
		validator.NotBlank(form.Email),
		"email",
		"This field cannot be blank",
	)
	// valid email
	form.CheckField(
		validator.Matches(form.Email, validator.EmailRX),
		"email",
		"This field must be a valid email address",
	)
	// password it not empty
	form.CheckField(
		validator.NotBlank(form.Password),
		"password",
		"This field cannot be blank",
	)
	// password len is at least 8
	form.CheckField(
		validator.MinLen(form.Password, 8),
		"password",
		"This field must be at least 8 characters long",
	)
	// re render form with data if invalid
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusBadRequest, "signup.tmpl.html", data)
		return
	}
	_, err = app.users.Insert(form.Name, form.Email, form.Password)
	if err!=nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusBadRequest, "signup.tmpl.html", data)
			return
		}
		app.serverError(w, err)
		return
	}
	// add confirmation flash message
	app.sessionManager.Put(r.Context(), "flash", "Your sign up was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK,  "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err!=nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// form validation
	// email is not empty
	form.CheckField(
		validator.NotBlank(form.Email),
		"email",
		"This field cannot be empty",
	)
	// email is valid
	form.CheckField(
		validator.Matches(form.Email, validator.EmailRX),
		"email",
		"This field must be a valid email address",
	)
	// password is not blank
	form.CheckField(
		validator.NotBlank(form.Password),
		"password",
		"This field cannot be empty",
	)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusBadRequest, "login.tmpl.html", data)
		return
	}
	// check if credentials are correct
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err!=nil {
		// add non field error to form if invalid credentials
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or Password is incorrect")
			data :=  app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusBadRequest, "login.tmpl.html", data)
			return
		}
		app.serverError(w, err)
		return
	}
	// use RenewToken() method on the current session
	err = app.sessionManager.RenewToken(r.Context())
	if err!=nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserId", id)
	redirectURL := app.sessionManager.PopString(r.Context(), "redirectPathAfterLogin")
	if redirectURL=="" {
		redirectURL = "/"
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// user RenewToken() to change current session id
	err := app.sessionManager.RenewToken(r.Context())
	if err!=nil {
		app.serverError(w, err)
		return
	}
	// remove authenticatedUserId from session data
	app.sessionManager.Remove(r.Context(), "authenticatedUserId")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userAccount(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserId")
	user, err := app.users.Get(id)
	if err!=nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData((r))
	data.User = user
	app.render(w, http.StatusOK, "account.tmpl.html", data)
}

type updatePasswordForm struct {
	CurrentPassword string `form:"current_password"`
	NewPassword string `form:"new_password"`
	ConfirmPassword string `form:"confirm_password"`
	validator.Validator `form:"-"`
}

func (app *application) updatePassword(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = updatePasswordForm{}
	app.render(w, http.StatusOK, "update_password.tmpl.html", data)
}

func (app *application) updatePasswordPost(w http.ResponseWriter, r *http.Request) {
	var form updatePasswordForm
	var err = app.decodePostForm(r, &form)
	if err!=nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(
		validator.NotBlank(form.CurrentPassword),
		"current_password",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.NotBlank(form.NewPassword),
		"new_password",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.MinLen(form.NewPassword, 8),
		"new_password",
		"This field must be at least 8 characters long",
	)
	form.CheckField(
		validator.NotBlank(form.ConfirmPassword),
		"confirm_password",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.PermittedValue(form.ConfirmPassword, form.NewPassword),
		"confirm_password",
		"This field must be same as New Password",
	)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusBadRequest, "update_password.tmpl.html", data)
		return
	}

	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserId")
	err = app.users.UpdatePassword(userId, form.CurrentPassword, form.NewPassword)
	if err!=nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("current_password", "Current password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusBadRequest, "update_password.tmpl.html", data)
			return
		}
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Password updated successfully")
	http.Redirect(w, r, "/user/account", http.StatusSeeOther)
}
