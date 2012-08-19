package models 

import (
	"database/sql"
)

type UserWrapper struct {
	User User `json: user`
}

type User struct {
	Id int64 `sql:"_id"`
	Avatar string `sql:"avatar" json:"avatar"`
	Name string `sql:"name" json:"displayName"`
	UserId string `sql:"user_id" json:"encodedId"`
	AccessToken string `sql:"access_token"`
	AccessSecret string `sql:"access_secret"`
}

func UserById(db *sql.DB, id string) (*User, error) {
	q := "SELECT * FROM users WHERE user_id = ?"

	rows, er := db.Query(q, id)
	if er != nil {
		return nil, er
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var b User 

	if er := scan(rows, &b) ; er != nil {
		return nil, er
	}

	return &b, nil
}

func GetUsers(db *sql.DB) ([]*User, error) {
	q := "SELECT * FROM users"

	rows, er := db.Query(q)
	if er != nil {
		return nil, er
	}
	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var b User

		if er := scan(rows, &b) ; er != nil {
			return nil, er
		}

		users = append(users, &b)
	}

	return users, nil
}

func installUserTable(db *sql.DB) error {
	_, er := db.Exec(`
		CREATE TABLE IF NOT EXISTS users 
			( _id INTEGER PRIMARY KEY AUTOINCREMENT
			, avatar TEXT NOT NULL
			, name TEXT NOT NULL
			, user_id TEXT NOT NULL UNIQUE
			, access_token TEXT
			, access_secret TEXT
			);
	`)

	return er
}

func (u *User) Save(db *sql.DB) error {
	if u.Id > 0 {
		return update(db, "users", "_id", u)
	}

	id, er := insert(db, "users", "_id", u)
	if er != nil {
		return er
	}

	u.Id = id
	return nil
}
