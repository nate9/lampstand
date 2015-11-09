package main

import (
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestToJson(t *testing.T) {
	p := Passage{
		"Genesis 1:1",
		"NIV",
		[]Verse{
			Verse{
				Book:    "Genesis",
				Chapter: 1,
				VerseNo: 1,
				Text:    "In the beginning",
			},
		},
	}
	e := "{" +
		"\"reference\":\"Genesis 1:1\"," +
		"\"version\":\"NIV\"," +
		"\"verses\":[{\"book\":\"Genesis\",\"chapter\":1,\"verseNo\":1,\"text\":\"In the beginning\"}]}"
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

func TestToVerses(t *testing.T) {
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

	result := ToVerses(rs)

	want := []Verse{
		Verse{
			Book:    "Genesis",
			Chapter: 1,
			VerseNo: 1,
			Text:    "In the beginning",
		},
	}

	for i := 0; i < len(want); i++ {
		if result[i] != want[i] {
			t.Errorf("ToVerses = %+v, want %+v", result[i], want[i])
		}
	}
}
