package daum

import (
	"encoding/json"
	"errors"
	"fmt"
	"internal/common"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// BookResult represents a document of a Daum Book search result.
type BookResult struct {
	WebResult
	ISBN        string   `json:"isbn"`
	Authors     []string `json:"authors"`
	Publisher   string   `json:"publisher"`
	Translators []string `json:"translators"`
	Price       int      `json:"price"`
	SalePrice   int      `json:"sale_price"`
	Thumbnail   string   `json:"thumbnail"`
	Status      string   `json:"status"`
}

// BookSearchResult represents a Daum Book search result.
type BookSearchResult struct {
	Meta      common.PageableMeta `json:"meta"`
	Documents []BookResult        `json:"documents"`
}

// String implements fmt.Stringer.
func (br BookSearchResult) String() string { return common.String(br) }

type BookSearchResults []BookSearchResult

// SaveAs saves brs to @filename.
func (brs BookSearchResults) SaveAs(filename string) error { return common.SaveAsJSON(brs, filename) }

// BookSearchIterator is a lazy book search iterator.
type BookSearchIterator struct {
	Query   string
	AuthKey string
	Sort    string
	Page    int
	Size    int
	Target  string
	end     bool
}

// BookSearch allows to search books by @query in the Daum Book service.
//
// See https://developers.kakao.com/docs/latest/ko/daum-search/dev-guide#search-book for more details.
func BookSearch(query string) *BookSearchIterator {
	return &BookSearchIterator{
		Query:   url.QueryEscape(strings.TrimSpace(query)),
		AuthKey: common.KeyPrefix,
		Sort:    "accuracy",
		Page:    1,
		Size:    10,
		Target:  "",
		end:     false,
	}
}

// AuthorizeWith sets the authorization key to @key.
func (bi *BookSearchIterator) AuthorizeWith(key string) *BookSearchIterator {
	bi.AuthKey = common.FormatKey(key)
	return bi
}

// SortBy sets the sorting order of the document results to @order.
//
// @order can be accuracy or latest. (default is accuracy)
func (bi *BookSearchIterator) SortBy(order string) *BookSearchIterator {
	switch order {
	case "accuracy", "latest":
		bi.Sort = order
	default:
		panic(common.ErrUnsupportedSortingOrder)
	}
	if r := recover(); r != nil {
		log.Panicln(r)
	}
	return bi
}

// Result sets the result page number (a value between 1 and 50).
func (bi *BookSearchIterator) Result(page int) *BookSearchIterator {
	if 1 <= page && page <= 50 {
		bi.Page = page
	} else {
		panic(common.ErrPageOutOfBound)
	}
	if r := recover(); r != nil {
		log.Panicln(r)
	}
	return bi
}

// Display sets the number of documents displayed on a single page (a value between 1 and 50).
func (bi *BookSearchIterator) Display(size int) *BookSearchIterator {
	if 1 <= size && size <= 50 {
		bi.Size = size
	} else {
		panic(common.ErrSizeOutOfBound)
	}
	if r := recover(); r != nil {
		log.Panicln(r)
	}
	return bi
}

// Filter limits the search field.
//
// @target can be one of the following options:
//
// title, isbn, publisher, person
func (bi *BookSearchIterator) Filter(target string) *BookSearchIterator {
	switch target {
	case "title", "isbn", "publisher", "person", "":
		bi.Target = target
	default:
		panic(errors.New(
			`target must be one of the following options:
			title, isbn, publisher, person`))
	}
	if r := recover(); r != nil {
		log.Panicln(r)
	}
	return bi
}

// Next returns the book search result and proceeds the iterator to the next page.
func (bi *BookSearchIterator) Next() (res BookSearchResult, err error) {
	if bi.end {
		return res, ErrEndPage
	}

	client := new(http.Client)
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("https://dapi.kakao.com/v3/search/book?query=%s&sort=%s&page=%d&size=%d&target=%s",
			bi.Query, bi.Sort, bi.Page, bi.Size, bi.Target), nil)

	if err != nil {
		return
	}

	req.Close = true

	req.Header.Set(common.Authorization, bi.AuthKey)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return
	}

	bi.Page++

	bi.end = res.Meta.IsEnd || 50 < bi.Page

	return
}
