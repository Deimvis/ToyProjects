package cyoa

type Story struct {
	Meta     Meta               `json:"meta"`
	Chapters map[string]Chapter `json:"chapters"`
}

type Meta struct {
	FirstChapterTitle string `json:"first_chapter_title"`
}

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text         string `json:"text"`
	ChapterTitle string `json:"arc"`
}
