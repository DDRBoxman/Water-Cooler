package main

import (
	"code.google.com/p/gorilla/mux"
	"database/sql"
	"net/http"
	"html/template"
	"watercooler/models"
	"fmt"
	"time"
	"io/ioutil"
	"encoding/json"
	"math"
)

func UiWater(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userId := vars["user"]
	user, err := models.UserById(db, userId)
	if err != nil {
		return err
	}

	var water [3]map[string]interface{}
		water[0] = map[string]interface{}{
		"Name" : "8oz",
		"Avatar" : "/img/8oz.png",
		"Size" : 8,}
		
		water[1] = map[string]interface{}{
		"Name" : "12oz",
		"Avatar" : "/img/12oz.png",
		"Size" : 12,}
		
		water[2] = map[string]interface{}{
		"Name" : "16.9oz",
		"Avatar" : "/img/169oz.png",
		"Size" : 16.9}

	tmpls, _ := template.ParseGlob("./public/tpl/*.tmpl")

	tmpls.New("content").Parse(`{{template "begin"}}{{template "water" .}}{{template "end"}}`)

	args := map[string]interface{}{
		"User" : user,
		"Water" : water,
	}

        tmpls.ExecuteTemplate(w, "content", args)


	return nil
}


func UiAddWater(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
        userId := vars["user"]
        user, err := models.UserById(db, userId)
        if err != nil {
                return err
        }

	o := LoadOAuth(user.AccessToken, user.AccessSecret)

	date := time.Now().Format("2006-01-02")

	data := map[string]string {
		"amount" : vars["size"],	
		"date" : date,	
		"unit" : "fl oz",
	}

	response, err := o.Post("http://api.fitbit.com/1/user/-/foods/log/water.json", data)
	if err != nil {
		return err
	}

	waterUrl := fmt.Sprintf("http://api.fitbit.com/1/user/-/foods/log/water/date/%s.json", date) 

	response, err = o.Get(waterUrl, map[string]string{})
	if err != nil {
		return err
	}
	bodyBytes, _ := ioutil.ReadAll(response.Body)

	var f interface{}
	err = json.Unmarshal(bodyBytes, &f)
	if err != nil {
		return err
	}

	m := f.(map[string]interface{})
	m = m["summary"].(map[string]interface{})
	if total, ok := m["water"].(float64) ; ok {
		
		percent := int64(math.Min((total / 1419.53) * 100, 100))

		tmpls, _ := template.ParseGlob("./public/tpl/*.tmpl")

		tmpls.New("content").Parse(`{{template "begin"}}{{template "waterfill" .}}{{template "end"}}`)

		args := map[string]interface{}{
			"User" : user,
			"Percent" : percent ,
		}

		tmpls.ExecuteTemplate(w, "content", args)

		return nil
	}

	return nil
}
