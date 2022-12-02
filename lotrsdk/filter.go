package lotrsdk

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// filters can be of several forms
// var operator value
// ?name=Gandalf
// var exists
// ?name
// var does not exist
// ?!name

//  match/negate
//include/exclude
// exists, doesn't exist
// regex
// inequality

// we can not chain inequalites (ie, budgetInMillions<300,250 does not work)

// FilterCompareType represents all the comparisons available to via a query param
type FilterCompareType int

const (
	FilterCompareEqual              FilterCompareType = iota
	FilterCompareNotEqual           FilterCompareType = iota
	FilterCompareLessThan           FilterCompareType = iota
	FilterCompareGreaterThan        FilterCompareType = iota
	FilterCompareLessThanOrEqual    FilterCompareType = iota
	FilterCompareGreaterThanOrEqual FilterCompareType = iota
)

func (ft FilterCompareType) ToString() (string, error) {
	switch ft {
	case FilterCompareEqual:
		return "=", nil
	case FilterCompareNotEqual:
		return "!=", nil
	case FilterCompareLessThan:
		return "<", nil
	case FilterCompareGreaterThan:
		return ">", nil
	case FilterCompareLessThanOrEqual:
		return "<=", nil
	case FilterCompareGreaterThanOrEqual:
		return ">=", nil
	}
	return "", fmt.Errorf("invalid compare operator")
}

// SortOrder is the different was to sort results (asc or desc)
type SortOrder int

const (
	SortOrderAscending  SortOrder = iota
	SortOrderDescending SortOrder = iota
)

func (so SortOrder) ToString() (string, error) {
	switch so {
	case SortOrderAscending:
		return "asc", nil
	case SortOrderDescending:
		return "desc", nil
	}
	return "", fmt.Errorf("invalid sort order")
}

// Filter represents a way to modify the request, be it for filtering for specific values
// limiting to a number of results, or sorting the values
type Filter interface {
	// GenerateRawQuery creates the raw string that will be added to the query params
	GenerateRawQuery() (string, error)
}

// BinaryFilter represents a binary option for selecting the data (ie, budget<100)
//  key - the field to filter by
//  operator - the operator for comparision (=, <, etc)
//  value - the value to search for
//  values - more that one value can be provided
func BinaryFilter(key string, operator FilterCompareType, value string, values ...string) Filter {
	bf := binaryFilter{
		key:      key,
		operator: operator,
	}
	bf.values = append([]string{value}, values...)
	return bf
}

type binaryFilter struct {
	key      string
	values   []string
	operator FilterCompareType
}

func (bf binaryFilter) GenerateRawQuery() (string, error) {
	// key-op-value
	operatorStr, err := bf.operator.ToString()
	if err != nil {
		return "", fmt.Errorf("cannot filter with invalid operator")
	}

	// it is invalid to chain together inequalities
	if bf.operator != FilterCompareEqual && bf.operator != FilterCompareNotEqual && len(bf.values) > 1 {
		return "", fmt.Errorf("cannot filter with operator %s on more than one value", operatorStr)
	}

	var sb strings.Builder
	sb.WriteString(bf.key)

	if err != nil {
		return "", fmt.Errorf("failed to generate operator string: %w", err)
	}
	sb.WriteString(operatorStr)

	if len(bf.values) == 0 {
		return "", fmt.Errorf("trying to generate query without value")
	}
	sb.WriteString(url.QueryEscape(bf.values[0]))
	for _, val := range bf.values[1:] {
		sb.WriteByte(',')
		sb.WriteString(url.QueryEscape(val))
	}

	return sb.String(), nil
}

// ExistFilter only selects the data if it contains the provided field
//   key - the field to search for
func ExistFilter(key string) Filter {
	return existFilter{
		key: key,
	}
}

type existFilter struct {
	key string
}

func (ef existFilter) GenerateRawQuery() (string, error) {
	return ef.key, nil
}

// NotExistFilter only selects the data if it does not contain the provided field
//   key - the field to search for
func NotExistFilter(key string) Filter {
	return notExistFilter{
		key: key,
	}
}

type notExistFilter struct {
	key string
}

func (nf notExistFilter) GenerateRawQuery() (string, error) {
	return fmt.Sprintf("!%s", nf.key), nil
}

// Sort and pagination are technically not filters, but they are applied the same way

// Sort sorts the output
//   value - the value to sort by
//   order - the sorting order (asc/desc)
func Sort(value string, order SortOrder) Filter {
	return sortFilter{
		value: value,
		order: order,
	}
}

type sortFilter struct {
	value string
	order SortOrder
}

func (sf sortFilter) GenerateRawQuery() (string, error) {
	orderStr, err := sf.order.ToString()
	if err != nil {
		return "", fmt.Errorf("failed to generate asc/desc string for sort: %w", err)
	}
	return fmt.Sprintf("sort=%s:%s", sf.value, orderStr), nil
}

// struct used to back Limit, Page, and Offset
type paginationFilter struct {
	key   string
	value string
}

func (pf paginationFilter) GenerateRawQuery() (string, error) {
	return fmt.Sprintf("%s=%s", pf.key, pf.value), nil
}

// Limit adds limit=%d to the queryparams
//   value - the amount to limit
func Limit(value int) Filter {
	return paginationFilter{
		key:   "limit",
		value: strconv.Itoa(value),
	}
}

// Page adds page=%d to the queryparams
//   value - the amount to page
func Page(value int) Filter {
	return paginationFilter{
		key:   "page",
		value: strconv.Itoa(value),
	}
}

// Offset adds offset=%d to the queryparams
//   value - the amount to offset
func Offset(value int) Filter {
	return paginationFilter{
		key:   "offset",
		value: strconv.Itoa(value),
	}
}

// Filters is a slice of Filter(s) that also implements the Filter interface
type Filters []Filter

func (fs Filters) GenerateRawQuery() (string, error) {
	if len(fs) == 0 {
		return "", nil
	}

	sb := strings.Builder{}
	str, err := fs[0].GenerateRawQuery()
	if err != nil {
		return "", fmt.Errorf("failed to generate query params: %w", err)
	}
	sb.WriteString(str)

	for _, elt := range fs[1:] {
		sb.WriteByte('&')
		str, err := elt.GenerateRawQuery()
		if err != nil {
			return "", fmt.Errorf("failed to generate query params: %w", err)
		}
		sb.WriteString(str)
	}

	return sb.String(), nil
}

// MergeFilters combines multiple filters together as a single filter
//   filters - the filters to merge
func MergeFilters(filters ...Filter) Filter {
	result := Filters{}
	for _, f := range filters {
		result = append(result, f)
	}
	return result
}
