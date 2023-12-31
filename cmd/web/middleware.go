package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		)
		w.Header().Set(
			"Referrer-Policy",
			"origin-when-cross-origin",
		)
		w.Header().Set(
			"X-Content-Type-Options",
			"nosniff",
		)
		w.Header().Set(
			"X-Frame-Options",
			"deny",
		)
		w.Header().Set(
			"X-XSS-Protection",
			"0",
		)
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func ()  {
			// use built in recover function to check if panic
			if err := recover(); err!=nil {
				// set Connection: close header
				w.Header().Set("Connection", "close")
				// call serverError to return response with internal server error
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			// add path user is trying to access to session data before redirecting
			app.sessionManager.Put(r.Context(), "redirectPathAfterLogin", r.URL.Path)
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// set Cache-Control: no-store header so that pages that require authentication are not cached
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (app *application) noSurf(next http.Handler) http.Handler {
	var csrfHandler = nosurf.New(next)
	csrfHandler.SetBaseCookie(
		http.Cookie{
			HttpOnly: true,
			Path: "/",
			Secure: true,
		},
	)
	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get authenticatedUserId from session
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserId")
		if id==0 {
			// call next handler if no authenticatedUserId is present
			next.ServeHTTP(w, r)
			return
		}
		// check if user exists with id
		exists, err := app.users.Exists(id)
		if err!=nil {
			app.serverError(w, err)
			return
		}
		// if users exists
		// create copy of request with context containing isAuthenticatedContextKey set to true
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, exists)
			r = r.WithContext(ctx)
		}
		// call next handler
		next.ServeHTTP(w, r)
	})
}