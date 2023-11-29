package quiet_hn

import (
	"errors"
	"net/url"
	"sort"
	"strings"

	"github.com/Deimvis/toyprojects/quiet_hn/hn"
)

type HNAPIFetcher struct {
	client hn.Client
}

func NewHNAPIFetcher() HNAPIFetcher {
	return HNAPIFetcher{client: hn.Client{}}
}

func (f HNAPIFetcher) FetchTopStories(n int) ([]Story, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("Failed to load top stories")
	}
	stories := f.fetchStoriesWithLimit(ids, n)
	f.orderStoriesByIds(&stories, ids)
	return stories, nil
}

func (f HNAPIFetcher) fetchStoriesWithLimit(ids []int, limit int) []Story {
	var results []Story
	start := 0
	for len(results) < limit && start < len(ids) {
		batchSize := max(int(float64(limit)*1.1), limit+2)
		batchSize = min(batchSize, limit-start)
		batchResults := f.fetch(ids[start : start+batchSize])
		for _, r := range batchResults {
			if r.err != nil || !isStoryLink(r.story) {
				continue
			}
			results = append(results, r.story)
		}
		start += batchSize
	}
	return results
}

func (f HNAPIFetcher) fetch(ids []int) []storyFetchResult {
	resultCh := make(chan storyFetchResult)
	for _, id := range ids {
		go func(id int) {
			hnStory, err := f.client.GetItem(id)
			if err != nil {
				resultCh <- storyFetchResult{err: err}
			} else {
				resultCh <- storyFetchResult{story: parseHNItem(hnStory)}
			}
		}(id)
	}
	var results []storyFetchResult
	for range ids {
		results = append(results, <-resultCh)
	}
	return results
}

func (f HNAPIFetcher) orderStoriesByIds(stories *[]Story, orderedIds []int) {
	id2ind := make(map[int]int)
	for ind, id := range orderedIds {
		id2ind[id] = ind
	}
	sort.Slice(*stories, func(i, j int) bool {
		return id2ind[(*stories)[i].ID] < id2ind[(*stories)[j].ID]
	})
}

func isStoryLink(item Story) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) Story {
	ret := Story{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

type storyFetchResult struct {
	story Story
	err   error
}

type Story struct {
	hn.Item
	Host string
}
