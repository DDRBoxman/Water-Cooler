package main

import (
	"code.google.com/p/gorilla/mux"
	"database/sql"
	"net/http"
)

func BindHandlers(db *sql.DB, mux *mux.Router) {
	mux.HandleFunc("/", WrapHandler(db, UiHome))
	mux.HandleFunc("/auth", WrapHandler(db, UiAuth))
	mux.HandleFunc("/completeAuth", WrapHandler(db, UiCompleteAuth))
	
	mux.HandleFunc("/water/{user}", WrapHandler(db, UiWater))
	mux.HandleFunc("/water/{user}/{size}", WrapHandler(db, UiAddWater))
}

func WrapHandler(db *sql.DB, f func(*sql.DB, http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		er := f(db, w, r)

		if er != nil {
			http.Error(w, er.Error(), 500)
		}
	}
}
