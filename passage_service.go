package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/nate9/lampstand/api"
	"net/http"
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

	log.Info("Got a query request for passage: " + passagequery)
	pq := parsePassage(passagequery)
	passage := *new(api.Passage)
	var verses []api.Verse
	var reference string
	var err error

	book, err := s.dao.FindBook(pq.book)

	if pq.chapterEnd != -1 {
		verses, err = s.dao.FindMultiChapterPassage(version, pq.book, pq.chapter, pq.chapterEnd, pq.begin, pq.end)
		reference = fmt.Sprintf("%s %d:%d-%d:%d", book, pq.chapter, pq.begin, pq.chapterEnd, pq.end)
	} else if pq.end != -1 {
		verses, err = s.dao.FindVerses(version, pq.book, pq.chapter, pq.begin, pq.end)
		reference = fmt.Sprintf("%s %d:%d-%d", book, pq.chapter, pq.begin, pq.end)
	} else if pq.begin != -1 {
		verses, err = s.dao.FindVerse(version, pq.book, pq.chapter, pq.begin)
		reference = fmt.Sprintf("%s %d:%d", book, pq.chapter, pq.begin)
	} else {
		verses, err = s.dao.FindChapter(version, pq.book, pq.chapter)
		reference = fmt.Sprintf("%s %d", book, pq.chapter)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if len(verses) == 0 {
		log.Error("Passage " + passagequery + " not found")
		http.Error(w, "No such passage exists", 404)
		return
	}

	passage.Version = version
	passage.Reference = reference
	passage.Verses = verses

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
