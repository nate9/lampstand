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
	Verses []Verse `json:"verses"`
}

func ToPassage(rs *sql.Rows) Passage {
	p := Passage{Verses: []Verse{}}
	for rs.Next() {
		var version string
		v := new(Verse)
		rs.Scan(&version, &v.Book, &v.Chapter, &v.VerseNo, &v.Text)
		p.Verses = append(p.Verses, *v)
	}
	return p
}

func ToJson(p Passage) string {
	verseJson, err := json.Marshal(p)
	if err != nil {
		fmt.Println("json err:", err)
	}
	return string(verseJson)
}
