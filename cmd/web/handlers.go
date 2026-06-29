package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.fayokunmiosho.com/internal/models"
	"snippetbox.fayokunmiosho.com/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Server", "Go")

	snippets, err := app.snippets.Latest()
	data := app.newTemplateData(r)
	data.Snippets = snippets

	if err != nil {
		app.serverError(w,r, err)
		return
	}

	app.render(w,r, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		http.NotFound(w,r)
		
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w,r,err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet


	app.render(w,r, http.StatusOK, "view.tmpl", data)

}
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

		data := app.newTemplateData(r)

		data.Form = snippetCreateForm{
			Expires: 365,
		}

		app.render(w,r,http.StatusOK, "create.tmpl", data)
}

type snippetCreateForm struct {
	Title string `form:"title"`
	Content string `form:"content"`
	Expires int `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	var form snippetCreateForm

	err := app.decodePostForm(r, &form)

	if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
	}

		form.CheckField(validator.NotBlank(form.Title), "title", "Title cannot be empty")
		form.CheckField(validator.NotBlank(form.Content), "content", "Content cannot be empty")
		form.CheckField(validator.MaxChars(form.Content, 100), "content","Content cannot be more than 100 characters long")
		form.CheckField(validator.PermittedValue(form.Expires, 1,7,365), "expires", "This field must equal 1, 7 or 365")
    

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(w,r,err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetDeletePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		app.logger.Error(err.Error())
		app.serverError(w,r,err)
	}
	id, err := strconv.Atoi(r.PostForm.Get("id"))

	if err != nil {
		app.logger.Error(err.Error())
		app.serverError(w,r,err)
	}


	_, err = app.snippets.Delete(id)

	if err != nil {
		app.logger.Error(err.Error())
		app.serverError(w,r,err)
	}

	app.logger.Info("snippet deleted", "id", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}