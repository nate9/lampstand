package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
)

type Verse struct {
	Book string
	Chapter float64
	VerseNo float64
	Text string
}

type VerseSlice struct {
	Verses []Verse
}

func findVerses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := r.URL.Query()
	passage := params.Get("passage")
	fmt.Println(passage)

	v := Verse{Book: "Genesis", Chapter: 1, VerseNo: 1, Text: "In the beginning"}
	verseJson, err := json.Marshal(v)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Fprint(w, string(verseJson))
}

func main() {
	router := httprouter.New()
	router.GET("/verses", findVerses)
	http.ListenAndServe(":8080", router)
}