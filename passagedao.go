package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
)

type PassageDaoImpl struct {
	db *sql.DB
}

type PassageDao interface {
	Setup(setupDir string)
	FindChapter(book string, chapterNo int) (Passage, error)
	FindVerse(book string, chapterNo int, verseNo int) (Passage, error)
	FindVerses(book string, chapterNo int, verseBegin int, verseEnd int) (Passage, error)
	Close()
}

func NewPassageDao(database string) (PassageDao, error) {
	db, err := sql.Open("sqlite3", database)
	p := &PassageDaoImpl{db: db}
	return p, err
}

func TestPassageDao(db *sql.DB) PassageDao {
	p := &PassageDaoImpl{db: db}
	return p
}

func (p *PassageDaoImpl) Setup(setupDir string) {
	fmt.Println("Setting up database")
	fileList := []string{}
	filepath.Walk(setupDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	for _, f := range fileList {
		fmt.Println("Inserting " + f + " into database")
		insertBookIntoDb(f, p.db)
	}
	fmt.Println("finished!")
}

func insertBookIntoDb(path string, db *sql.DB) {
	dat, err := ioutil.ReadFile(path)
	bookSql := string(dat)
	checkErr(err)

	_, err = db.Exec(bookSql)
	checkErr(err)
}

func (p *PassageDaoImpl) FindChapter(book string, chapterNo int) (result Passage, err error) {
	rows, err := p.db.Query("SELECT * FROM BIBLE WHERE BOOK LIKE ? + AND CHAPTER = ?", book + "%", chapterNo)
	checkErr(err)
	result = ToPassage(rows)
	return result, err
}

func (p *PassageDaoImpl) FindVerse(book string, chapterNo int, verseNo int) (result Passage, err error) {
	rows, err := p.db.Query("SELECT * FROM BIBLE WHERE BOOK LIKE ? AND CHAPTER = ? AND VERSE = ?", book + "%", chapterNo, verseNo)
	checkErr(err)
	result = ToPassage(rows)
	return result, err
}

func (p *PassageDaoImpl) FindVerses(book string, chapterNo int, verseBegin int, verseEnd int) (result Passage, err error) {
	query := "SELECT * FROM BIBLE WHERE BOOK LIKE ? AND CHAPTER = ? AND VERSE BETWEEN ? and ?"
	rows, err := p.db.Query(query, book + "%", chapterNo, verseBegin, verseEnd)
	checkErr(err)
	result = ToPassage(rows)
	return result, err
}

func (p *PassageDaoImpl) Close() {
	p.db.Close()
}
