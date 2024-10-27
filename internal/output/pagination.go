package output

func NewPaginatedPage[T any](data []T, perPage int, pageNum int) *PaginatedPage[T] {
	var results []T
	var from int
	var to int
	total := len(data)
	if total == 0 {
		results = []T{}
	} else {
		start := perPage * (pageNum - 1)
		end := perPage * pageNum
		// Make sure we don't go out of the slice bounds
		if start > len(data) {
			start = len(data)
		}
		if end > len(data) {
			end = len(data)
		}
		results = data[start:end]
	}

	numberOfResults := len(results)

	if numberOfResults > 0 {
		from = perPage*(pageNum-1) + 1
		to = perPage * pageNum
		if to > total {
			to = total
		}
	}

	nextPage := 0
	if to < total {
		nextPage = pageNum + 1
	}
	prevPage := 0
	if pageNum > 1 {
		prevPage = pageNum - 1
	}
	return &PaginatedPage[T]{
		Results:     results,
		PerPage:     perPage,
		PageNum:     pageNum,
		Total:       total,
		From:        from,
		To:          to,
		NextPage:    nextPage,
		PrevPage:    prevPage,
		HasNextPage: nextPage != 0,
		HasPrevPage: prevPage != 0,
	}
}

type PaginatedPage[T any] struct {
	PerPage     int
	PageNum     int
	Total       int
	Results     []T
	From        int
	To          int
	NextPage    int
	PrevPage    int
	HasNextPage bool
	HasPrevPage bool
}
