package client

import (
	"encoding/json"
	"net/url"
	"fmt"
	"github.com/nate9/lampstand/api"
	"io/ioutil"
	"log"
	"net/http"
)

var _ = log.Print

type PassageClient struct {
	Host string
}

func (c *PassageClient) GetSingleChapterPassage(version string, book string, chapter int) (api.Passage, error) {
	var passage api.Passage

	passageString := url.QueryEscape(fmt.Sprintf("%s %d", book, chapter))
	passageQuery := fmt.Sprintf("%s/verses?passage=%s", version, passageString)
	url := fmt.Sprintf("%s/api/%s", c.Host, passageQuery)
	log.Println(url)
	r, err := http.Get(url)
	if err != nil {
		return passage, err
	}

	passage, err = toPassage(r)
	return passage, err
}

func (c *PassageClient) GetSingleVersePassage(version string, book string, chapter int, verse int) (api.Passage, error) {
	var passage api.Passage

	passageString := url.QueryEscape(fmt.Sprintf("%s %d:%d", book, chapter, verse))
	passageQuery := fmt.Sprintf("%s/verses?passage=%s", version, passageString)
	url := fmt.Sprintf("%s/api/%s", c.Host, passageQuery)
	log.Println(url)
	r, err := http.Get(url)
	if err != nil {
		return passage, err
	}

	passage, err = toPassage(r)
	return passage, err
}

func (c *PassageClient) GetMultipleVersesPassage(version string, book string, chapter int, verseStart int, verseEnd int) (api.Passage, error) {
	var passage api.Passage

	passageString := url.QueryEscape(fmt.Sprintf("%s %d:%d-%d", book, chapter, verseStart, verseEnd))
	passageQuery := fmt.Sprintf("%s/verses?passage=%s", version, passageString)
	url := fmt.Sprintf("%s/api/%s", c.Host, passageQuery)
	r, err := http.Get(url)
	if err != nil {
		return passage, err
	}

	passage, err = toPassage(r)
	return passage, err
}

func (c *PassageClient) GetMultipleChapterPassage(version string, book string, chapterStart int, chapterEnd int, verseStart int, verseEnd int) (api.Passage, error) {
	var passage api.Passage

	passageString := url.QueryEscape(fmt.Sprintf("%s %d:%d-%d:%d", book, chapterStart, verseStart, chapterEnd, verseEnd))
	passageQuery := fmt.Sprintf("%s/verses?passage=%s", version, passageString)
	url := fmt.Sprintf("%s/api/%s", c.Host, passageQuery)
	r, err := http.Get(url)
	if err != nil {
		return passage, err
	}

	passage, err = toPassage(r)
	return passage, err
}

func toPassage(r *http.Response) (api.Passage, error) {
	var passage api.Passage
	defer r.Body.Close()
	bodyContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bodyContent, &passage)
	return passage, err
}
