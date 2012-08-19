package models

import (
	"database/sql"
)

func InstallTables(db *sql.DB) error {
	funcs := []func(*sql.DB)error{
		installUserTable,
	}

	for _, f := range funcs {
		if er := f(db) ; er != nil {
			return er
		}
	}

	return nil
}
