package models

import (
	"database/sql"
	"reflect"
	"bytes"
	"fmt"
	"errors"
	"strings"
)

func fixQueryParameters(q string) string {
	i := 1
	var buf bytes.Buffer

	for _, b := range []byte(q) {
		if b == '?' {
			placeholder := fmt.Sprintf("$%d", i)
			buf.WriteString(placeholder)
			i += 1
		} else {
			buf.WriteByte(b)
		}
	}

	return buf.String()
}

func scan(rs *sql.Rows, args ...interface{}) error {
	assignmentMap := make(map[string]interface{})

	prefix := ""

	/* build a mapping from
	 *  colname -> interface{}, type, fieldName
	 */

	for _, arg := range args {
		val := reflect.ValueOf(arg)
		ty := val.Type()
		kind := val.Kind()

		if kind == reflect.String {
			prefix = arg.(string)
			continue

		} else if kind == reflect.Ptr {
			val = val.Elem()
			ty = val.Type()
			kind = val.Kind()

			if kind == reflect.Struct {
				/* Need to enumerate the fields or some shit */
				for i := 0 ; i < val.NumField() ; i += 1 {
					field := ty.Field(i)

					if sqlName := field.Tag.Get("sql") ; sqlName != "" {
						assignmentMap[prefix + sqlName] = val.Field(i).Addr().Interface()
					}
				}

			} else {
				panic("pointer to non-struct")
			}

		} else {
			panic("Invalid arg")
		}
	}

	cols, er := rs.Columns()
	if er != nil {
		return er
	}

	ifaces := []interface{}{}
	for _, col := range cols {
		if target, ok := assignmentMap[col] ; ok {
			ifaces = append(ifaces, target)
		} else {
			ifaces = append(ifaces, new(interface{}))
		}
	}

	if er := rs.Scan(ifaces...) ; er != nil {
		fmt.Printf("Error encountered, columns: %#v\n", cols)
		return er
	}

	return nil
}

func update(dbish interface{}, table, idField string, arg interface{}) error {
	val := reflect.ValueOf(arg)

	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.New("must pass a struct (or pointer to one)")
	}

	ty := val.Type()

	sqlFields := []string{}
	newValues := []interface{}{}
	var id int64 = 0

	for i := 0 ; i < ty.NumField() ; i += 1 {
		field := ty.Field(i)
		sqlField := field.Tag.Get("sql")

		if sqlField == "" {
			continue
		}

		if sqlField == idField {
			k := val.Field(i).Kind()
			if k != reflect.Int && k != reflect.Int8 && k != reflect.Int16 && k != reflect.Int32 && k != reflect.Int64 {
				return errors.New("idField is non-integer, cannot coerse")
			}

			id = val.Field(i).Int()
		}

		sqlFields = append(sqlFields, fmt.Sprintf("%s = ?", sqlField))
		newValues = append(newValues, val.Field(i).Interface())
	}

	if id == 0 {
		return errors.New("idField is 0 or not set, cannot update")
	}

	newValues = append(newValues, id)

	q := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", table, strings.Join(sqlFields, ", "), idField)

	var er error = nil

	if db, ok := dbish.(*sql.DB) ; ok {
		_, er = db.Exec(q, newValues...)

	} else if tx, ok := dbish.(*sql.Tx) ; ok {
		_, er = tx.Exec(q, newValues...)

	} else {
		er = errors.New("dbish is neither a *sql.DB nor a *sql.Tx")
	}

	if er != nil {
		return er
	}

	return nil
}

func insert(dbish interface{}, table, idField string, arg interface{}) (int64, error) {
	val := reflect.ValueOf(arg)

	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return 0, errors.New("must pass a struct (or pointer to one)")
	}

	ty := val.Type()

	sqlFields := []string{}
	placeholders := []string{}
	newValues := []interface{}{}

	for i := 0 ; i < ty.NumField() ; i += 1 {
		field := ty.Field(i)
		sqlField := field.Tag.Get("sql")

		if sqlField == "" {
			continue
		}

		if sqlField == idField {
			continue
		}

		sqlFields = append(sqlFields, fmt.Sprintf("%s", sqlField))
		placeholders = append(placeholders, "?")
		newValues = append(newValues, val.Field(i).Interface())
	}

	q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(sqlFields, ", "), strings.Join(placeholders, ","))

	var er error = nil
	var res sql.Result

	if db, ok := dbish.(*sql.DB) ; ok {
		res, er = db.Exec(q, newValues...)

	} else if tx, ok := dbish.(*sql.Tx) ; ok {
		res, er = tx.Exec(q, newValues...)

	} else {
		er = errors.New("dbish is neither a *sql.DB nor a *sql.Tx")
	}

	if er != nil {
		return 0, er
	}

	lastId, er := res.LastInsertId()

	if er != nil {
		return 0, er
	}

	return lastId, nil
}
