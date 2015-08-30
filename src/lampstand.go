package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	db, err := sql.Open("sqlite3", "./hcsb.db")
	checkErr(err)

	dat, err := ioutil.ReadFile("setup_tables.sql")
	setupSql := string(dat)
	checkErr(err)

	_, err = db.Exec(setupSql)
	checkErr(err)

	dirName := "./hcsb"
	fileList := []string{}
	filepath.Walk(dirName, func(path string, f os.FileInfo, _ error) error {
		if(!f.IsDir()) {
			fileList = append(fileList, path)
		}
		return nil
	})

	for _, f := range fileList {
		fmt.Println("Inserting " + f + " into database")
		insertBookIntoDb(f, db)
	}

	fmt.Println("finished!")

	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func insertBookIntoDb(path string, db *sql.DB) {
	dat, err := ioutil.ReadFile(path)
	bookSql := string(dat)
	checkErr(err)

	_, err = db.Exec(bookSql)
	checkErr(err)
}
