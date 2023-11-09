package cyoa

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

type CLIRunner interface {
	Run() error
}

func NewCLIRunner(s Story) CLIRunner {
	r := runner{s}
	return r
}

type runner struct {
	s Story
}

func (r runner) Run() error {
	curChapterTitle := r.s.Meta.FirstChapterTitle
	for {
		c := r.s.Chapters[curChapterTitle]
		boldWhite := color.New(color.FgCyan).Add(color.Bold)
		boldWhite.Printf("\"%s\"\n", c.Title)
		for _, par := range c.Paragraphs {
			fmt.Printf("\t%s\n", par)
		}
		opts := c.Options
		if len(c.Options) == 0 {
			opts = []Option{{Text: "Play story again", ChapterTitle: r.s.Meta.FirstChapterTitle}}
		}
		opt, err := handlePrompt(opts)
		if err != nil {
			return err
		}
		curChapterTitle = opt.ChapterTitle
	}
}

func handlePrompt(opts []Option) (Option, error) {
	items := make([]string, len(opts))
	for i := range opts {
		items[i] = opts[i].Text
	}
	prompt := promptui.Select{
		Label: "Choose what to do next",
		Items: items,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return Option{}, err
	}
	return opts[i], nil
}
