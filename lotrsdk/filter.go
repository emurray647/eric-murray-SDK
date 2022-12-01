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

type FilterType int

const (
	FilterCompareEqual              = iota
	FilterCompareNotEqual           = iota
	FilterCompareLessThan           = iota
	FilterCompareGreaterThan        = iota
	FilterCompareLessThanOrEqual    = iota
	FilterCompareGreaterThanOrEqual = iota
)

func (ft FilterType) ToString() (string, error) {
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

// A Filter needs to generate query params
type Filter interface {
	// GenerateQueryParam() string
	// AddQueryParam(url.Values)
	GenerateRawQuery() (string, error)
}

func BinaryFilter(variable string, operator FilterType, value string, values ...string) Filter {
	bf := binaryFilter{
		variable: variable,
		operator: operator,
	}
	bf.values = append([]string{value}, values...)
	return bf
}

type binaryFilter struct {
	variable string
	values   []string
	operator FilterType
}

func (bf binaryFilter) GenerateRawQuery() (string, error) {
	// key-op-value
	var sb strings.Builder
	sb.WriteString(bf.variable)
	operatorStr, err := bf.operator.ToString()
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

// not technically Filters, but they are applied the same way
type sortFilter struct {
	value string
	order SortOrder
}

func Sort(value string, order SortOrder) Filter {
	return sortFilter{
		value: value,
		order: order,
	}
}

func (sf sortFilter) GenerateRawQuery() (string, error) {
	orderStr, err := sf.order.ToString()
	if err != nil {
		return "", fmt.Errorf("failed to generate asc/desc string for sort: %w", err)
	}
	return fmt.Sprintf("sort=%s:%s", sf.value, orderStr), nil
}

type paginationFilter struct {
	key   string
	value string
}

func (pf paginationFilter) GenerateRawQuery() (string, error) {
	return fmt.Sprintf("%s=%s", pf.key, pf.value), nil
}

func Limit(value int) Filter {
	return paginationFilter{
		key:   "limit",
		value: strconv.Itoa(value),
	}
}

func Page(value int) Filter {
	return paginationFilter{
		key:   "page",
		value: strconv.Itoa(value),
	}
}

func Offset(value int) Filter {
	return paginationFilter{
		key:   "offset",
		value: strconv.Itoa(value),
	}
}

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

func MergeFilters(filters ...Filter) Filter {
	result := Filters{}
	for _, f := range filters {
		result = append(result, f)
	}
	return result
}
