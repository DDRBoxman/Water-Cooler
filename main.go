package main

import (
	"code.google.com/p/gorilla/mux"
	"net/http"
	"fmt"
	"github.com/alloy-d/goauth"
	"encoding/gob"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"watercooler/models"
)

func doMain() error {

	gob.Register(oauth.OAuth{})

	db, er := sql.Open("sqlite3", "db.sqlite")
	if er != nil {
		return er
	}

	if er := models.InstallTables(db) ; er != nil {
		return er
	}

	mux := mux.NewRouter()

	BindHandlers(db, mux)

	mux.NotFoundHandler = http.FileServer(http.Dir("public/"))
	return http.ListenAndServe("localhost:8000", mux)
}

func SetupOAuth() *oauth.OAuth {
	return LoadOAuth("", "")
}

func LoadOAuth(accessToken, accessSecret string) *oauth.OAuth {
	o := new (oauth.OAuth)

	o.SignatureMethod = oauth.HMAC_SHA1

	o.ConsumerKey = "yourkeyhere"
	o.ConsumerSecret = "yourkeyhere"
	o.Callback = "http://127.0.0.1:8000/completeAuth"

	o.RequestTokenURL = "http://api.fitbit.com/oauth/request_token"
	o.OwnerAuthURL = "http://www.fitbit.com/oauth/authorize"
	o.AccessTokenURL = "http://api.fitbit.com/oauth/access_token"

	if accessToken != "" {
		o.AccessToken = accessToken
	}

	if accessSecret != "" {
		o.AccessSecret = accessSecret
	}

	return o
}

func main() {
	if er := doMain() ; er != nil {
		fmt.Printf("Failed to start: %#v\n", er)
	}
}
