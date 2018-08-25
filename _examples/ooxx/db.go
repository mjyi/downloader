package ooxx

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func InitDB() {
	DB = openSelfDB()
	handleError(createOOXXTable())
}

func CloseDB() {
	DB.Close()
}

func openSelfDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./ooxx.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func createOOXXTable() error {
	sqlStmt := `create table if not exists ooxx(
		id text not null primary key,
		post_id text,
		comment_date text, 
		pics text,
		ext text);`

	_, err := DB.Exec(sqlStmt)
	return err
}

// InsertModels 插入Model
func InsertModels(oxs []OOXXModel) {
	tx, err := DB.Begin()
	handleError(err)

	stmt, err := tx.Prepare("INSERT INTO ooxx(id, post_id, comment_date, pics) VALUES(?, ?, ?, ?)")
	handleError(err)

	defer stmt.Close()
	for _, ox := range oxs {
		stmt.Exec(ox.CommentID, ox.CommentPostID, ox.CommentDate, ox.PicsStr)
	}
	tx.Commit()
}


func handleError(err error)  {
	if err != nil {
		 log.Println(err)
	}
}
