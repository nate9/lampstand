package lampstand

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Verse struct {
	Book    string
	Chapter float64
	VerseNo float64
	Text    string
}

type Passage struct {
	Verses []Verse
}

func ToPassage(rs *sql.Rows) Passage {
	p := Passage{Verses: []Verse{}}
	for rs.Next() {
		v := new(Verse)
		rs.Scan(&v.Book, &v.Chapter, &v.VerseNo, &v.Text)
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
