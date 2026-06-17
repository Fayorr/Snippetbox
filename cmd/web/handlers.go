package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"snippetbox.fayokunmiosho.com/internal/models"
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
	Title string
	Content string
	Expires int
	FieldError map[string]string
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &snippetCreateForm{
		Title: r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
		FieldError: map[string]string{},
	}

	if strings.TrimSpace(form.Title) == ""  {
        form.FieldError["title"] = "Title cannot be empty"
    } 
    if strings.TrimSpace(form.Content) == "" {
        form.FieldError["content"] = "Content cannot be empty"
    } else if utf8.RuneCountInString(form.Content) > 100 {
        // This 'else if' is fine because it applies to the same field
        form.FieldError["content"] = "Content cannot be more than 100 characters long"
    } 
    if expires != 1 && expires != 7 && expires != 365 {
        form.FieldError["expires"] = "This field must equal 1, 7 or 365"
    }

	if len(form.FieldError) > 0 {
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
		app.serverError(w,r,err)
	}
	id := r.PostForm.Get("ID")
	num, err := app.snippets.Delete(11)

	if err != nil {
		app.serverError(w,r,err)
	}

	fmt.Fprintf(w, "snippet ID %d deleted", num)
}