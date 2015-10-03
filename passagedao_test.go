package main

import (
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestShouldReturnVersions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"VERSION", "INFO", "COPYRIGHT"})
	mockRows.AddRow("HCSB", "published 2003", "Copyright etc.")
	mockRows.AddRow("NIV", "Never Incorrect Version", "Copyright 1999")
	mock.ExpectQuery("^SELECT (.+) FROM VERSIONS").WillReturnRows(mockRows)
	p := TestPassageDao(db)
	versions, err := p.GetVersions()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	expected := []string{"HCSB", "NIV"}
	for i := 0; i < len(expected); i++ {
		a := versions[i]
		b := expected[i]
		if a != b {
			t.Errorf("versions: %q, expected: %q", versions, expected)
		}
	}
}

func TestShouldLookForAVerse(t *testing.T) {
	fmt.Println("Starting....")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	mockRows := sqlmock.NewRows([]string{"VERSION", "BOOK", "CHAPTER", "VERSE", "TEXT"})
	mockRows.AddRow("HCSB", "Genesis", 1, 1, "In the beginning")
	mock.ExpectQuery("^SELECT (.+) FROM BIBLE WHERE VERSION = (.+) AND BOOK LIKE (.+) AND CHAPTER = (.+) AND VERSE = (.+)").
		WithArgs("HCSB", "Genesis%", 1, 1).WillReturnRows(mockRows)
	p := TestPassageDao(db)
	_, err = p.FindVerse("HCSB", "Genesis", 1, 1)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldLookForAPassage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	mockRows := sqlmock.NewRows([]string{"VERSION", "BOOK", "CHAPTER", "VERSE", "TEXT"})
	mockRows.AddRow("HCSB", "Genesis", 1, 1, "In the beginning")
	mockRows.AddRow("HCSB", "Genesis", 1, 2, "God created the heavens")
	mock.ExpectQuery("^SELECT (.+) FROM BIBLE WHERE VERSION = (.+) AND BOOK LIKE (.+) AND CHAPTER = (.+) AND VERSE BETWEEN (.+) AND (.+)").
		WithArgs("HCSB", "Genesis%", 1, 1, 2).WillReturnRows(mockRows)
	p := TestPassageDao(db)
	_, err = p.FindVerses("HCSB", "Genesis", 1, 1, 2)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldLookForAChapter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	mockRows := sqlmock.NewRows([]string{"VERSION", "BOOK", "CHAPTER", "VERSE", "TEXT"})
	mockRows.AddRow("HCSB", "Genesis", 1, 1, "In the beginning")
	mockRows.AddRow("HCSB", "Genesis", 1, 2, "God created the heavens")
	mock.ExpectQuery("^SELECT (.+) FROM BIBLE WHERE VERSION = (.+) AND BOOK LIKE (.+) AND CHAPTER = (.+)").WithArgs("HCSB", "Genesis%", 1).WillReturnRows(mockRows)
	p := TestPassageDao(db)
	_, err = p.FindChapter("HCSB", "Genesis", 1)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldLookForMultiChapterPassage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	mockRows := sqlmock.NewRows([]string{"VERSION", "BOOK", "CHAPTER", "VERSE", "TEXT"})
	mockRows.AddRow("HCSB", "Genesis", 1, 1, "In the beginning")
	mockRows.AddRow("HCSB", "Genesis", 2, 10, "God created the heavens")
	mock.ExpectQuery("^SELECT (.+) FROM BIBLE WHERE VERSION = (.+) AND BOOK LIKE (.+) AND \\(\\(CHAPTER = (.+) AND VERSE >= (.+)\\) OR \\(CHAPTER = (.+) AND VERSE <= (.+)[\\)\\)]").
		WithArgs("HCSB", "Genesis%", 1, 1, 2, 10).WillReturnRows(mockRows)
	p := TestPassageDao(db)
	_, err = p.FindMultiChapterPassage("HCSB", "Genesis", 1, 2, 1, 10)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
