package main // import "github.com/nate9/lampstand"

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type PassageService struct {
	dao PassageDao
}

type PassageQuery struct {
	book       string
	chapter    int
	chapterEnd int
	begin      int
	end        int
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

func (s *PassageService) findVerses(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	version := ps.ByName("version")
	params := r.URL.Query()
	passagequery := params.Get("passage")
	if passagequery == "" {
		http.Error(w, "passage query parameter not defined", 400)
		return
	}

	fmt.Println(passagequery)
	pq := parsePassage(passagequery)
	var passage Passage
	var err error
	if pq.chapterEnd != -1 {
		passage, err = s.dao.FindMultiChapterPassage(version, pq.book, pq.chapter, pq.chapterEnd, pq.begin, pq.end)
	} else if pq.end != -1 {
		passage, err = s.dao.FindVerses(version, pq.book, pq.chapter, pq.begin, pq.end)
	} else if pq.begin != -1 {
		passage, err = s.dao.FindVerse(version, pq.book, pq.chapter, pq.begin)
	} else {
		passage, err = s.dao.FindChapter(version, pq.book, pq.chapter)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if len(passage.Verses) == 0 {
		log.Println("passage not found")
		http.Error(w, "No such passage exists", 404)
		return
	}

	passageJson := ToJson(passage)
	fmt.Fprint(w, string(passageJson))
}

func parsePassage(passagequery string) (pq PassageQuery) {
	pq = PassageQuery{book: "", chapter: -1, chapterEnd: -1, begin: -1, end: -1}
	passagequery = strings.TrimSpace(passagequery)
	end := len(passagequery)
	lastspace := strings.LastIndex(passagequery, " ")
	pq.book = passagequery[:lastspace]
	colon := strings.Index(passagequery, ":")
	secondColon := strings.LastIndex(passagequery, ":")

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

	if colon != secondColon {
		pq.chapterEnd, _ = strconv.Atoi(passagequery[hyphen+1 : secondColon])
		pq.end, _ = strconv.Atoi(passagequery[secondColon+1 : end])
	}

	return pq
}

func main() {
	service := NewPassageService("./bible.db")
	router := httprouter.New()
	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	if bind == ":" {
		bind = "localhost:8080"
	}
	fmt.Println("Serving lampstand on : " + bind)
	router.GET("/api/:version/verses", service.findVerses)
	router.NotFound = http.FileServer(http.Dir("static"))
	log.Fatal(http.ListenAndServe(bind, router))
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
}
