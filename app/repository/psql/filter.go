package psql

import (
	"fmt"

	"github.com/delveper/mylib/app/models"
)

func evalQuery(filter models.DataFilter) string {
	var query string

	if filter.Filter != nil {
		query = "WHERE "

		eval := func(node *models.FilterNode) string {
			return fmt.Sprintf("%v%v%v %v ",
				node.Field,
				node.Operator,
				node.Value,
				node.Conjunction,
			)
		}

		node := filter.Filter.Head
		for node.Next != nil {
			query += eval(node)
			node = node.Next
		}

		query += eval(node)
	}

	if filter.OrderBy != "" {
		query = fmt.Sprintf("%v\nORDER BY %v", query, filter.OrderBy)
	}

	if filter.Skip != 0 {
		query = fmt.Sprintf("%v\nOFFSET %v", query, filter.Skip)
	}

	if filter.Top != 0 {
		query = fmt.Sprintf("%v\nLIMIT %v", query, filter.Top)
	}

	return query
}
