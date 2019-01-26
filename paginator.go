package flow

import (
	"strconv"
)

var (
	// PaginatorPerPageDefault is the amount of results per page
	PaginatorPerPageDefault = 20

	// PaginatorPageKey is the query parameter holding amount of results per page
	PaginatorPageKey = "page"

	// PaginatorPerPageKey is the query parameter holding the amount of results per page
	// to override the default one
	PaginatorPerPageKey = "per_page"
)

// Paginator is a type used to represent the pagination
type Paginator struct {
	// Current page you're on
	Page int `json:"page"`
	// Number of results you want per page
	PerPage int `json:"per_page"`
	// Page * PerPage (ex: 2 * 20, Offset == 40)
	Offset int `json:"offset"`
	// Total potential records matching the query
	TotalEntriesSize int `json:"total_entries_size"`
	// Total records returns, will be <= PerPage
	CurrentEntriesSize int `json:"current_entries_size"`
	// Total pages
	TotalPages int `json:"total_pages"`
}

// PaginationParams is a parameters provider interface to get the pagination params from
type PaginationParams interface {
	Get(key string) string
}

// NewPaginator returns a new `Paginator` value with the appropriate
// defaults set.
func NewPaginator(page int, perPage int) *Paginator {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = PaginatorPerPageDefault
	}
	p := &Paginator{Page: page, PerPage: perPage}
	p.Offset = (page - 1) * p.PerPage
	return p
}

// NewPaginatorFromParams takes an interface of type `PaginationParams`,
// the `url.Values` type works great with this interface, and returns
// a new `Paginator` based on the params or `PaginatorPageKey` and
// `PaginatorPerPageKey`. Defaults are `1` for the page and
// PaginatorPerPageDefault for the per page value.
func NewPaginatorFromParams(params PaginationParams) *Paginator {
	page := "1"
	if p := params.Get("page"); p != "" {
		page = p
	}

	perPage := strconv.Itoa(PaginatorPerPageDefault)
	if pp := params.Get("per_page"); pp != "" {
		perPage = pp
	}

	p, err := strconv.Atoi(page)
	if err != nil {
		p = 1
	}

	pp, err := strconv.Atoi(perPage)
	if err != nil {
		pp = PaginatorPerPageDefault
	}
	return NewPaginator(p, pp)
}
