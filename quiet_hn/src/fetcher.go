package quiet_hn

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"

	"github.com/Deimvis/toyprojects/quiet_hn/hn"
)

// HNAPIFetcher is thread-safe
type HNAPIFetcher struct {
	client hn.Client                    // thread-safe
	cache  Cache[int, storyFetchResult] // thread-safe
}

func NewHNAPIFetcher() HNAPIFetcher {
	return HNAPIFetcher{client: hn.NewClient(), cache: NewCacheWithTTL[int, storyFetchResult]()}
}

func (f HNAPIFetcher) FetchTopStories(n int) ([]Story, error) {
	ids, err := f.client.TopItems()
	if err != nil {
		return nil, fmt.Errorf("failed to load top stories %s", err.Error())
	}
	stories := f.fetchStoriesWithLimit(ids, n)
	f.orderStoriesByIds(&stories, ids)
	return stories, nil
}

func (f HNAPIFetcher) fetchStoriesWithLimit(ids []int, limit int) []Story {
	var results []Story
	start := 0
	for len(results) < limit && start < len(ids) {
		need := limit - len(results)
		left := len(ids) - start
		batchSize := max(int(float64(need)*1.1), need+2)
		batchSize = min(batchSize, left)
		batchResults := f.fetchMany(ids[start : start+batchSize])
		for _, r := range batchResults {
			if r.err != nil || !isStoryLink(r.story) {
				continue
			}
			results = append(results, r.story)
		}
		start += batchSize
	}
	return results[:limit]
}

func (f HNAPIFetcher) fetchMany(ids []int) []storyFetchResult {
	log.Printf("fetch %d stories\n", len(ids))
	resultCh := make(chan storyFetchResult)
	for _, id := range ids {
		go func(id int) {
			resultCh <- f.fetch(id)
		}(id)
	}
	var results []storyFetchResult
	for range ids {
		results = append(results, <-resultCh)
	}
	return results
}

func (f HNAPIFetcher) fetch(id int) storyFetchResult {
	if f.cache != nil {
		if v, ok := f.cache.Get(id); ok {
			return v
		}
	}
	hnStory, err := f.client.GetItem(id)
	if err != nil {
		return storyFetchResult{err: err}
	}
	ret := storyFetchResult{story: parseHNItem(hnStory)}
	if f.cache != nil {
		f.cache.Put(id, ret)
	}
	return ret
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
