package main

import (
	"testing"
)

func TestParseVerse(t *testing.T) {
	cases := []struct {
		in string
		pq PassageQuery
	}{
		{"Genesis 1:1-10", PassageQuery{"Genesis", 1, -1, 1, 10}},
		{"Genesis 1:1", PassageQuery{"Genesis", 1, -1, 1, -1}},
		{"Genesis 1", PassageQuery{"Genesis", 1, -1, -1, -1}},
		{"1 Chronicles 2:3-4", PassageQuery{"1 Chronicles", 2, -1, 3, 4}},
		{"2 Corinthians 2:1-3:5", PassageQuery{"2 Corinthians", 2, 3, 1, 5}},
	}

	for _, c := range cases {
		pq := parsePassage(c.in)
		if pq != c.pq {
			t.Errorf("parsePassage(%q) == %+v\n, want %+v\n", c.in, pq, c.pq)
		}
	}
}
