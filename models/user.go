package models 

import (
	"database/sql"
	"github.com/lye/crud"
)

type UserWrapper struct {
	User User `json: user`
}

type User struct {
	Id int64 `crud:"_id"`
	Avatar string `crud:"avatar" json:"avatar"`
	Name string `crud:"name" json:"displayName"`
	UserId string `crud:"user_id" json:"encodedId"`
	AccessToken string `crud:"access_token"`
	AccessSecret string `crud:"access_secret"`
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

	var u User
	if er := crud.Scan(rows, &u) ; er != nil {
		return nil, er
	}

	return &u, nil
}

func GetUsers(db *sql.DB) ([]User, error) {
	q := "SELECT * FROM users"

	rows, er := db.Query(q)
	if er != nil {
		return nil, er
	}

	users := []User{}

	if er := crud.ScanAll(rows, &users) ; er != nil {
		return nil, er
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
		return crud.Update(db, "users", "_id", u)
	}

	id, er := crud.Insert(db, "users", "_id", u)
	if er != nil {
		return er
	}

	u.Id = id
	return nil
}
