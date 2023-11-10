package sitemap

type filterFn func(string) bool

func filterInPlace(arr *[]string, f filterFn) *[]string {
	size := 0
	for i := range *arr {
		if f((*arr)[i]) {
			(*arr)[size] = (*arr)[i]
			size++
		}
	}
	*arr = (*arr)[:size]
	return arr
}
