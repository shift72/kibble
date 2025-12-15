package models

import (
	"strconv"
	"strings"
)

type Paginator struct {
	Route       *Route
	RoutePrefix string
	TotalItems  int
}

func (ctx *RenderContext) Paginate(totalItems int) []Pagination {
	return Paginator{
		RoutePrefix: ctx.RoutePrefix,
		Route:       ctx.Route,
		TotalItems:  totalItems,
	}.GetAll()
}

// Pagination describes a single page of results
type Pagination struct {
	Index          int    // current page, 1-based index
	Size           int    // nominal page size. Actual number of results may differ
	Total          int    // Total pages
	PreviousURL    string // Prev page href, or blank
	CurrentURL     string // Current page href
	NextURL        string // Next page href, or blank
	ItemSliceStart int    // First item index (inclusive) for use with [:] slice operator
	ItemSliceEnd   int    // Last item index (exclusive) for use with [:] slice operator
	ItemCount      int    // Number of items on the page
}

func (p Paginator) GetAll() []Pagination {
	firstPage := p.GetPagination(1)

	results := make([]Pagination, firstPage.Total)
	results[0] = firstPage
	for i := 1; i < firstPage.Total; i++ {
		results[i] = p.GetPagination(i + 1)
	}

	return results
}

func (p Paginator) GetPagination(currentPage int) Pagination {
	totalItems := p.TotalItems
	pageSize := p.Route.PageSize
	if pageSize <= 0 {
		// slightly weird edge case here: if there are no items the page size is 0
		pageSize = totalItems
	}

	var totalPages int
	if pageSize <= 0 || totalItems <= 0 {
		// ensure that empty collections always render at least one page
		totalPages = 1
	} else {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}

	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPages {
		currentPage = totalPages
	}

	pagination := Pagination{
		Index:          currentPage,
		Total:          totalPages,
		Size:           pageSize,
		CurrentURL:     p.GetPaginationURLPath(currentPage),
		PreviousURL:    "",
		NextURL:        "",
		ItemSliceStart: 0,
		ItemSliceEnd:   0,
		ItemCount:      0,
	}

	// NOTE: follows Go slice semantics [start:end)
	pagination.ItemSliceStart = (pagination.Index - 1) * pageSize
	pagination.ItemSliceEnd = pagination.ItemSliceStart + pageSize
	if pagination.ItemSliceEnd > totalItems {
		pagination.ItemSliceEnd = totalItems
	}

	pagination.ItemCount = pagination.ItemSliceEnd - pagination.ItemSliceStart

	if pagination.Index > 1 {
		pagination.PreviousURL = p.GetPaginationURLPath(pagination.Index - 1)
	}

	if pagination.Index < pagination.Total {
		pagination.NextURL = p.GetPaginationURLPath(pagination.Index + 1)
	}

	return pagination
}

func (p Paginator) GetPaginationURLPath(pageIndex int) string {
	var urlPath string
	if len(p.Route.FirstPageURLPath) > 0 && pageIndex <= 1 {
		urlPath = p.Route.FirstPageURLPath
	} else {
		urlPath = p.Route.URLPath
	}

	return p.RoutePrefix + strings.Replace(urlPath, ":index", strconv.Itoa(pageIndex), 1)
}
