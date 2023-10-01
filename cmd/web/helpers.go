package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
)

// serverError writes error message and stack trace to error log
// return 500 Internal Server Error response to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s",  err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// serverError
// return specific error message and status code to the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound: helper method which send a 404 not found response using clientError
func (app *application) notFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// get template set from cache with key as page
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exists", page)
		app.serverError(w, err)
		return
	}
	// initialize a new buffer
	buf := new(bytes.Buffer)
	// write template to buffer instead of response writer 
	err := ts.ExecuteTemplate(buf, "base", data)
	if err!=nil {
		app.serverError(w, err)
		return
	}
	// continue if template if safely written to buffer
	// set status code
	w.WriteHeader(status)
	// write contents of buffer to response writer
	buf.WriteTo(w)
}

// helper method to decode PostForm data
// dst is target destination we want the data to be decoded into
func (app *application) decodePostForm(r *http.Request, dst any) error {
	var err = r.ParseForm()
	if err!=nil {
		return err
	}
	// call Decode on our formDecoder instance with destination dst
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err!=nil {
		// check type of error
		// panic if error is InvalidDecodeError
		var invalidDecodeError *form.InvalidDecoderError
		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}
		// return if other error
		return err
	}
	return nil
}
