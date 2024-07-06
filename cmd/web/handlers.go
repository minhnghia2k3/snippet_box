package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/minhnghia2k3/snippet_box/internal/models"
	"github.com/minhnghia2k3/snippet_box/internal/validator"
	"net/http"
	"strconv"
)

// represent the form data and validation errors for the form field.
type snippetCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` // Embedded type
}

// `home()` handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Get templateData struct, and add the snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}

// GET: /snippet/view/123
// Validate `id` param
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	// Convert id string to an integer
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
}

// Handler to show snippet form
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = &snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

// Handler to create a new record of snippet to database
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// ParseForm() parse raw data to PostForm() map
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form snippetCreateForm

	// Instead parse string to int at expires field, using formDecoder()
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate user input.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be longer than 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// If there is any error, redisplay the template.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	// Insert snippet data to mysql db
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
