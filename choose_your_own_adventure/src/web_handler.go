package cyoa

import (
	"log"
	"net/http"
	"strings"
	"text/template"
)

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithTitleFn(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.titleFn = fn
	}
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, defaultTpl, defaultTitleFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s       Story
	t       *template.Template
	titleFn func(r *http.Request) string
}

func defaultTitleFn(r *http.Request) string {
	return strings.Trim(strings.TrimSpace(r.URL.Path), "/")
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/static/") {
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
		return
	}
	chapterTitle := h.titleFn(r)
	if chapterTitle == "" {
		chapterTitle = h.s.Meta.FirstChapterTitle
	}
	if chapter, ok := h.s.Chapters[chapterTitle]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}
