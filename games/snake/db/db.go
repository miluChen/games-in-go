package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "./game.db"

var gameDB *sql.DB

func Open() error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	// create table
	sqlStmt := "create table if not exists snake (id integer not null primary key, name text);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	gameDB = db
	return nil
}

func Insert(name string) error {
	stmt := "insert into snake(name) values(?)"
	_, err := gameDB.Exec(stmt, name)
	return err
}

func Read() ([]string, error) {
	rows, err := gameDB.Query("select name from snake")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return names, nil
}

func Close() error {
	if gameDB != nil {
		return gameDB.Close()
	}
	return nil
}
