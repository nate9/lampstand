package main

import (
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestToJson(t *testing.T) {
	p := Passage{
		[]Verse{
			Verse{
				Book:    "Genesis",
				Chapter: 1,
				VerseNo: 1,
				Text:    "In the beginning",
			},
		},
	}
	e := "{\"Verses\":[{\"Book\":\"Genesis\",\"Chapter\":1,\"VerseNo\":1,\"Text\":\"In the beginning\"}]}"
	cases := []struct {
		in   Passage
		want string
	}{
		{p, e},
	}

	for _, c := range cases {
		got := ToJson(c.in)
		if got != c.want {
			t.Errorf("ToJson() == %q, want %q", got, c.want)
		}
	}
}

func TestToPassage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("failed to open sqlmock database:", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"VERSION", "BOOK", "CHAPTER", "VERSE", "TEXT"}).
		AddRow("HCSB", "Genesis", 1, 1, "In the beginning")

	mock.ExpectQuery("^SELECT (.+) FROM BIBLE").WillReturnRows(rows)

	rs, err := db.Query("SELECT * FROM BIBLE")
	if err != nil {
		fmt.Println("failed to query mock database:", err)
	}

	defer rs.Close()

	result := ToPassage(rs)

	want := Passage{
		[]Verse{
			Verse{
				Book:    "Genesis",
				Chapter: 1,
				VerseNo: 1,
				Text:    "In the beginning",
			},
		},
	}

	gotVerse := result.Verses[0]
	wantVerse := want.Verses[0]

	if gotVerse != wantVerse {
		t.Errorf("ToPassage = %q, want %q", gotVerse, wantVerse)
	}

}
