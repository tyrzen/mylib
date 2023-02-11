package models

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// DataFilter build on of OData query options:
// $filter. represents mandatory query option which supports operations: `and`, `or`, `eq`, `ne`, `gt`, `lt`, `ge`, `le`.
// $orderby. optional param, represents sorting column which supports `acs` and `desc` operators.
// $top. optional param, represents limit of items from the resource.
// $skip. optional param, represents offset of records in the resource.
// Names of fields MUST correspond to struct field names and should provide in case-sensitive format.
// Example: /books?$filter=Rate lt 100 and Rate gte 400 and Genre eq 'Thriller'&$orderby=Title desc&$OptionTop=100&$skip=10
// TODO: add following properties: `from`, `to` (in UTC format), `in` Sequences (ids of sequences).
type DataFilter struct {
	URL     *url.URL
	Filter  Filter
	OrderBy string
	Top     int
	Skip    int
}

// Filter represent linked lists of OData expressions.
type Filter struct {
	Raw  string
	Head *FilterNode
}

// FilterNode represents OData expression.
type FilterNode struct {
	Field       string
	Operator    string
	Conjunction string
	Value       string
	Next        *FilterNode
}

type fieldData map[string]string

const (
	OptionFilter  = "$filter"
	OptionOrderBy = "$orderby"
	OptionTop     = "$top"
	OptionSkip    = "$skip"
)

const defaultTagName = "sql"

// UpdateURL makes query on top of parent URL.
func (df *DataFilter) UpdateURL() {
	q := df.URL.Query()

	q.Set(OptionFilter, df.Filter.Raw)

	if df.Top != 0 {
		q.Set(OptionTop, fmt.Sprintf("%v", df.Top))
	}

	if df.Skip != 0 {
		q.Set(OptionSkip, fmt.Sprintf("%v", df.Skip))
	}

	df.URL.RawQuery = q.Encode()
}

// Insert adds new expression to Filter chain.
func (f *Filter) Insert(exp *FilterNode) {
	if f.Head == nil {
		f.Head = exp
		return
	}

	node := f.Head
	for node.Next != nil {
		node = node.Next
	}

	node.Next = exp
}

// NewDataFilter creates *DataFilter of struct type T
// that build on OData OptionFilter of given URL.
// The input of OData query options will be validated moving along the process.
func NewDataFilter[T any](URL *url.URL) (*DataFilter, error) {
	data, err := getStructFieldData(*new(T))
	if err != nil {
		return nil, err
	}

	filter, err := parseFilter(URL.String(), data)
	if err != nil {
		return nil, err
	}

	orderBy, err := parseOrderBy(URL.String(), data)
	if err != nil {
		return nil, err
	}

	top, err := parseTop(URL.String())
	if err != nil {
		return nil, err
	}

	skip, err := parseSkip(URL.String())
	if err != nil {
		return nil, err
	}

	return &DataFilter{
		URL:     URL,
		Filter:  *filter,
		OrderBy: orderBy,
		Top:     top,
		Skip:    skip,
	}, nil
}

// parseQueryOption parses value of given QueryOption from URL query parameters.
func parseQueryOption(query, opt string) string {
	pattern := fmt.Sprintf(`(?P<option>\%s=)(?P<value>[^&$]*)`, opt)
	if match := regexp.MustCompile(pattern).
		FindStringSubmatch(query); match != nil {
		return match[2]
	}

	return ""
}

func parseSkip(url string) (int, error) {
	query := parseQueryOption(url, OptionSkip)
	if query == "" {
		return 0, nil
	}

	val, err := strconv.Atoi(query)
	if err != nil {
		return 0, fmt.Errorf("error parsing OptionSkip query option: %w", err)
	}

	return val, nil
}

func parseTop(url string) (int, error) {
	query := parseQueryOption(url, OptionTop)
	if query == "" {
		return 0, nil
	}

	val, err := strconv.Atoi(query)
	if err != nil {
		return 0, fmt.Errorf("error parsing OptionTop query option: %w", err)
	}

	return val, nil
}

func parseOrderBy(url string, fieldMap fieldData) (string, error) {
	query := parseQueryOption(url, OptionOrderBy)
	if query == "" {
		return "", nil
	}

	sortMap := map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
		"ASC":  "ASC",
		"DESC": "DESC",
	}

	var fieldList, sortList []string
	for k, v := range fieldMap {
		fieldList = append(fieldList, k, v)
	}

	for k, v := range sortMap {
		sortList = append(sortList, v, k)
	}

	pattern := fmt.Sprintf(`(%s)(\s(%s))*,*`,
		strings.Join(fieldList, "|"),
		strings.Join(sortList, "|"),
	)

	re := regexp.MustCompile(pattern)

	if match := re.ReplaceAllLiteralString(query, ""); strings.TrimSpace(match) != "" {
		return "", fmt.Errorf("query does not correspond pattern: %s", pattern)
	}

	for k, v := range fieldMap {
		query = strings.ReplaceAll(query, k, v)
	}

	for k, v := range sortMap {
		query = strings.ReplaceAll(query, k, v)
	}

	return query, nil
}

func parseFilter(url string, fieldMap fieldData) (*Filter, error) {
	query := parseQueryOption(url, OptionFilter)
	if query == "" {
		return nil, nil
	}

	operMap := map[string]string{
		"eq": "=",
		"ne": "!=",
		"gt": ">",
		"lt": "<",
		"le": "<=",
		"ge": ">=",
	}

	conjMap := map[string]string{
		"and": "AND",
		"or":  "OR",
		"AND": "AND",
		"OR":  "OR",
	}

	var operList, conjList, fieldList []string
	for k, v := range fieldMap {
		fieldList = append(fieldList, k, v)
	}

	for k := range operMap {
		operList = append(operList, k)
	}

	for k := range conjMap {
		conjList = append(conjList, k)
	}

	pattern := fmt.Sprintf(`(?P<field>%s)\s+(?P<operator>%s)\s+(?P<value>\d+|'[^']+')\s*(?P<conjunction>%s)*\s*`,
		strings.Join(fieldList, "|"),
		strings.Join(operList, "|"),
		strings.Join(conjList, "|"),
	)

	re := regexp.MustCompile(pattern)
	if match := re.ReplaceAllLiteralString(query, ""); strings.TrimSpace(match) != "" {
		return nil, fmt.Errorf("query does not correspond pattern: %s", pattern)
	}

	matches := re.FindAllStringSubmatch(query, -1)
	groups := re.SubexpNames()

	var fil = new(Filter)
	fil.Raw = query

	for _, match := range matches {
		var node FilterNode

		for i := 1; i < len(groups); i++ {
			switch groups[i] {
			case "field":
				node.Field = fieldMap[match[i]]
			case "operator":
				node.Operator = operMap[match[i]]
			case "value":
				node.Value = match[i]
			case "conjunction":
				node.Conjunction = conjMap[match[i]]
			}
		}

		fil.Insert(&node)
	}

	return fil, nil
}

// getStructFieldData retrieves a map of struct field names
// and tags of defaultTagName corresponding to them.
func getStructFieldData(src any) (fieldData, error) {
	var res = make(map[string]string, 0)

	srcValue := reflect.Indirect(reflect.ValueOf(src))
	if srcType := srcValue.Kind(); srcType != reflect.Struct {
		return nil, fmt.Errorf("input value must be struct, got: %v", srcType)
	}

	// iterate struct fields.
	for i := 0; i < srcValue.NumField(); i++ {
		fieldValue := srcValue.Field(i)
		fieldName := srcValue.Type().Field(i).Name
		// add only exported fields.
		if r := string(fieldName[0]); r >= `a` && r <= `z` {
			continue
		}

		tag := srcValue.Type().Field(i).Tag
		tagValue := tag.Get(defaultTagName)
		// add only field with tags.
		if tagValue == "" {
			continue
		}

		// add FieldName and value of defaultTagName.
		res[fieldName] = tagValue

		// recursive call for nested structs.
		if fieldValue.Type().Kind() != reflect.Struct {
			continue
		}

		nested, err := getStructFieldData(fieldValue.Interface())
		if err != nil {
			return nil, fmt.Errorf("error validating nested struct: %w", err)
		}

		for k, v := range nested {
			res[k] = v
		}
	}

	return res, nil
}
