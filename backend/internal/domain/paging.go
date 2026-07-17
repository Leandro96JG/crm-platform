package domain

type PagingFilter struct {
	Limit     int
	Offset    int
	SortBy    string
	SortOrder string
}

type Paging struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PagingResult[T any] struct {
	Result []T
	Paging Paging
}
