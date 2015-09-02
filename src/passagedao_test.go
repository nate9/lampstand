package lampstand

import (
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"fmt"
)

func TestShouldLookForAVerse(t *testing.T) {
	fmt.Println("Starting....")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	mockRows := sqlmock.NewRows([]string{"BOOK", "CHAPTER", "VERSE", "TEXT"})
	mockRows.AddRow("Genesis", 1, 1, "In the beginning")
	mock.ExpectQuery("^SELECT (.+) FROM BIBLE WHERE BOOK = (.+) AND CHAPTER = (.+) AND VERSE = (.+)").WithArgs("Genesis", 1, 1).WillReturnRows(mockRows)
	p := TestPassageDao(db)
	_, err = p.FindVerse("Genesis", 1, 1)

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
	mockRows := sqlmock.NewRows([]string{"BOOK", "CHAPTER", "VERSE", "TEXT"})
	mockRows.AddRow("Genesis", 1, 1, "In the beginning")
	mockRows.AddRow("Genesis", 1, 2, "God created the heavens")
	mock.ExpectQuery("^SELECT (.+) FROM BIBLE WHERE BOOK = (.+) AND CHAPTER = (.+) AND VERSE BETWEEN (.+) and (.+)").WithArgs("Genesis", 1, 1, 2).WillReturnRows(mockRows)
	p := TestPassageDao(db)
	_, err = p.FindVerses("Genesis", 1, 1, 2)

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
	mockRows := sqlmock.NewRows([]string{"BOOK", "CHAPTER", "VERSE", "TEXT"})
	mockRows.AddRow("Genesis", 1, 1, "In the beginning")
	mockRows.AddRow("Genesis", 1, 2, "God created the heavens")
	mock.ExpectQuery("^SELECT (.+) FROM BIBLE WHERE BOOK = (.+) AND CHAPTER = (.+)").WithArgs("Genesis", 1).WillReturnRows(mockRows)
	p := TestPassageDao(db)
	_, err = p.FindChapter("Genesis", 1)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}