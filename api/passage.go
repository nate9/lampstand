package api

type Verse struct {
	Book    string  `json:"book"`
	Chapter float64 `json:"chapter"`
	VerseNo float64 `json:"verseNo"`
	Text    string  `json:"text"`
}

type Passage struct {
	Reference string  `json:"reference"`
	Version   string  `json:"version"`
	Verses    []Verse `json:"verses"`
}
