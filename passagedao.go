package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"path/filepath"
)

type PassageDaoImpl struct {
	db *sql.DB
}

type PassageDao interface {
	Setup(setupDir string)
	GetVersions() ([]string, error)
	FindChapter(version string, book string, chapterNo int) (Passage, error)
	FindVerse(version string, book string, chapterNo int, verseNo int) (Passage, error)
	FindVerses(version string, book string, chapterNo int, verseBegin int, verseEnd int) (Passage, error)
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

func (p *PassageDaoImpl) GetVersions() (result []string, err error) {
	query := "SELECT * FROM VERSIONS"
	rows, err := p.db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	result = []string{}
	for rows.Next() {
		var version string
		var info string
		var copyright string
		rows.Scan(&version, &info, &copyright)
		result = append(result, version)
	}
	return result, err
}

func (p *PassageDaoImpl) FindChapter(version string, book string, chapterNo int) (result Passage, err error) {
	query := "SELECT * FROM BIBLE WHERE VERSION = ? AND BOOK LIKE ? + AND CHAPTER = ?"
	rows, err := p.db.Query(query, version, book+"%", chapterNo)
	checkErr(err)
	result = ToPassage(rows)
	return result, err
}

func (p *PassageDaoImpl) FindVerse(version string, book string, chapterNo int, verseNo int) (result Passage, err error) {
	query := "SELECT * FROM BIBLE WHERE VERSION = ? AND BOOK LIKE ? AND CHAPTER = ? AND VERSE = ?"
	rows, err := p.db.Query(query, version, book+"%", chapterNo, verseNo)
	checkErr(err)
	result = ToPassage(rows)
	return result, err
}

func (p *PassageDaoImpl) FindVerses(version string, book string, chapterNo int, verseBegin int, verseEnd int) (result Passage, err error) {
	query := "SELECT * FROM BIBLE WHERE VERSION = ? AND BOOK LIKE ? AND CHAPTER = ? AND VERSE BETWEEN ? and ?"
	rows, err := p.db.Query(query, version, book+"%", chapterNo, verseBegin, verseEnd)
	checkErr(err)
	result = ToPassage(rows)
	return result, err
}

func (p *PassageDaoImpl) Close() {
	p.db.Close()
}
