package main

import (
	"database/sql"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/nate9/lampstand/api"
)

func ToVerses(rs *sql.Rows) []api.Verse {
	verses := []api.Verse{}
	for rs.Next() {
		v := new(api.Verse)
		var version string
		rs.Scan(&version, &v.Book, &v.Chapter, &v.VerseNo, &v.Text)
		verses = append(verses, *v)
	}
	return verses
}

func ToJson(p api.Passage) string {
	verseJson, err := json.Marshal(p)
	if err != nil {
		log.Error("json err:", err)
	}
	return string(verseJson)
}
