package main

import (
	"fmt"
	"github.com/nate9/lampstand/client"
	. "gopkg.in/check.v1"
	"log"
	"os"
	"testing"
)

var _ = fmt.Print
var _ = log.Print

func Test(t *testing.T) {
	TestingT(t)
}

type TestSuite struct {
	s    *PassageService
	host string
}

var _ = Suite(&TestSuite{})

func (s *TestSuite) SetUpSuite(c *C) {

	host := fmt.Sprintf("localhost:%d", 8080)
	s.host = fmt.Sprintf("http://%s", host)

	//Mock an interface or setup an actual database
	//s.s = new(PassageService)
	SetupTestDatabase("./test.db")
	s.s = NewPassageService("./test.db")

	//service run
	go RunPassageService(s.s, host)
}

func (s *TestSuite) TearDownSuite(c *C) {
	//delete test.db
	os.Remove("./test.db")
}

//It should get a single chapter from a book
func (s *TestSuite) TestGetSingleChapterPassage(c *C) {
	//Given a Passage Client
	client := client.PassageClient{Host: s.host}
	passage, err := client.GetSingleChapterPassage("NIV", "Matthew", 1)
	c.Assert(err, Equals, nil)
	c.Assert(len(passage.Verses), Equals, 25)
}

//It should get multiple chapters from a book
func (s *TestSuite) TestGetCrossChapterPassage(c *C) {
	client := client.PassageClient{Host: s.host}
	passage, err := client.GetMultipleChapterPassage("NIV", "Matthew", 1, 2, 25, 1)
	c.Assert(err, Equals, nil)
	c.Assert(len(passage.Verses), Equals, 2)
}

//It should get a single verse from a chapter from a book
func (s *TestSuite) TestGetSingleVersePassage(c *C) {
	client := client.PassageClient{Host: s.host}
	passage, err := client.GetSingleVersePassage("NIV", "Matthew", 1, 1)
	c.Assert(err, Equals, nil)
	c.Assert(len(passage.Verses), Equals, 1)
}

//It should get multiple verses from a chapter from a book
func (s *TestSuite) TestGetMultiVersePassage(c *C) {
	client := client.PassageClient{Host: s.host}
	passage, err := client.GetMultipleVersesPassage("NIV", "Matthew", 1, 1, 10)
	c.Assert(err, Equals, nil)
	c.Assert(len(passage.Verses), Equals, 10)
}

//It should get a passage from a different version
func (s *TestSuite) TestGetSingleVerseDiffVersion(c *C) {
	client := client.PassageClient{Host: s.host}
	passage, err := client.GetSingleChapterPassage("ESV", "Matthew", 1)
	c.Assert(err, Equals, nil)
	c.Assert(len(passage.Verses), Equals, 25)
}
