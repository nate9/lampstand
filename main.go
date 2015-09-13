package main // import "github.com/nate9/lampstand"

import (
	"fmt"
	"github.com/nate9/lampstand/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type PassageService struct {
	dao PassageDao
}

type PassageQuery struct {
	book    string
	chapter int
	begin   int
	end     int
}

func NewPassageService(db string) *PassageService {
	service := new(PassageService)
	dao, err := NewPassageDao(db)
	service.dao = dao
	if err != nil {
		fmt.Println("couldn't open passage service: ", err)
	}
	return service
}

func (s *PassageService) findVerses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := r.URL.Query()
	passagequery := params.Get("passage")
	fmt.Println(passagequery)
	pq := parsePassage(passagequery)
	var passage Passage
	var err error
	if pq.end != -1 {
		passage, err = s.dao.FindVerses(pq.book, pq.chapter, pq.begin, pq.end)
	} else if pq.begin != -1 {
		passage, err = s.dao.FindVerse(pq.book, pq.chapter, pq.begin)
	} else {
		passage, err = s.dao.FindChapter(pq.book, pq.chapter)
	}
	checkErr(err)
	passageJson := ToJson(passage)
	fmt.Fprint(w, string(passageJson))
}

func parsePassage(passagequery string) (pq PassageQuery) {
	pq = PassageQuery{book: "", chapter: -1, begin: -1, end: -1}
	passagequery = strings.TrimSpace(passagequery)
	end := len(passagequery)
	lastspace := strings.LastIndex(passagequery, " ")
	pq.book = passagequery[:lastspace]
	colon := strings.Index(passagequery, ":")

	if colon != -1 {
		pq.chapter, _ = strconv.Atoi(passagequery[lastspace+1 : colon])
	} else {
		pq.chapter, _ = strconv.Atoi(passagequery[lastspace+1 : end])
	}

	hyphen := strings.Index(passagequery, "-")
	if hyphen != -1 && colon != -1 {
		pq.begin, _ = strconv.Atoi(passagequery[colon+1 : hyphen])
		pq.end, _ = strconv.Atoi(passagequery[hyphen+1 : end])
	} else if hyphen == -1 && colon != -1 {
		pq.begin, _ = strconv.Atoi(passagequery[colon+1 : end])
	}

	return pq
}

func main() {
	service := NewPassageService("./hcsb.db")
	router := httprouter.New()
	router.GET("/api/verses", service.findVerses)
	router.ServeFiles("/lampstand/*filepath", http.Dir("static"))
	log.Fatal(http.ListenAndServe(":8080", router))
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}
