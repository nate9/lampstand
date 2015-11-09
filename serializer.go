package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Verse struct {
	Book    string  `json:"book"`
	Chapter float64 `json:"chapter"`
	VerseNo float64 `json:"verseNo"`
	Text    string  `json:"text"`
}

type Passage struct {
	Reference string  `json:"reference"`
	Version   string  `json:"version"`
	Verses    []Verse `json:"verses"`
}

func ToVerses(rs *sql.Rows) []Verse {
	verses := []Verse{}
	for rs.Next() {
		v := new(Verse)
		var version string
		rs.Scan(&version, &v.Book, &v.Chapter, &v.VerseNo, &v.Text)
		verses = append(verses, *v)
	}
	return verses
}

func ToJson(p Passage) string {
	verseJson, err := json.Marshal(p)
	if err != nil {
		fmt.Println("json err:", err)
	}
	return string(verseJson)
}
