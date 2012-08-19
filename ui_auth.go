package main

import (
	"net/http"
	"fmt"
	"github.com/alloy-d/goauth"
	"code.google.com/p/gorilla/sessions"
	"io/ioutil"
	"encoding/json"
	"watercooler/models"
	"database/sql"
	"html/template"
)

var store = sessions.NewFilesystemStore("", []byte("something-very-secret"))

func UiAuth(db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	o := SetupOAuth()

	err := o.GetRequestToken()
	if err != nil {
		return err
	}

	url, err := o.AuthorizationURL()
	if err != nil {
		return err
	}

	session, _ := store.Get(r, "oauth")
	session.Values["oauth"] = o 
	store.Save(r, w, session)

	http.Redirect(w, r, url, http.StatusFound)	

	return nil
}

func UiCompleteAuth(db *sql.DB, w http.ResponseWriter, r *http.Request) error {


	values := r.URL.Query()
	if verifier, ok := values["oauth_verifier"] ; ok {
	if session, er := store.Get(r, "oauth") ; er == nil {
		if val, ok := session.Values["oauth"] ; ok {
			if o, ok := val.(oauth.OAuth) ; ok {
				err := o.GetAccessToken(verifier[0])
				if err != nil {
					return err
				}
	
				response, err := o.Post("http://api.fitbit.com/1/user/-/profile.json",
        map[string]string{})
				if err != nil {
					return err
				} else {
					defer response.Body.Close()
				}
				bodyBytes, _ := ioutil.ReadAll(response.Body) 
				var userWrapper models.UserWrapper
				err = json.Unmarshal(bodyBytes, &userWrapper)
				if (err != nil) {
					return err
				}

				user := userWrapper.User
				user.AccessToken = o.AccessToken
				user.AccessSecret = o.AccessSecret
				
				er := user.Save(db)
				if er != nil {
					return er
				}

				tmpls, _ := template.ParseGlob("./public/tpl/*.tmpl")
				tmpls.New("content").Parse(`{{template "begin"}}{{template "success" .}}{{template "end"}}`)

				tmpls.ExecuteTemplate(w, "content", user)

			} else {
				fmt.Printf("%v", ok)
			}
		} else {
			fmt.Printf("Token ! %v", ok)
		}
	} else {
		fmt.Printf(" session : %v", er)
	}
	}

	return nil
}

