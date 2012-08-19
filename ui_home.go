package main

import (
	"net/http"
	"database/sql"
	"html/template"
	"watercooler/models"
)

func UiHome(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	
	users, _ := models.GetUsers(db)
	
	tmpls, _ := template.ParseGlob("./public/tpl/*.tmpl")
	tmpls.New("content").Parse(`{{template "begin"}}{{template "userlist" .}}{{template "end"}}`)

	tmpls.ExecuteTemplate(w, "content", users)

	return nil;
}
