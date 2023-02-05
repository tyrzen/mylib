package psql

import (
	"fmt"

	"github.com/delveper/mylib/app/models"
)

func evaluateQuery(filter models.DataFilter) string {
	var query = "WHERE "

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

	if filter.OrderBy != nil {
		query = fmt.Sprintf("%v\nORDER BY %v", query, *filter.OrderBy)
	}

	if filter.Skip != nil {
		query = fmt.Sprintf("%v\nOFFSET %v", query, *filter.Skip)
	}

	if filter.Top != nil {
		query = fmt.Sprintf("%v\nLIMIT %v", query, *filter.Top)
	}

	return query
}
