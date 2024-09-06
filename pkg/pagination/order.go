package pagination

import (
	"fmt"
	"strings"
)

type Order string

const (
	UNKNOWN Order = "UNKNOWN"
	ASC     Order = "ASC"
	DESC    Order = "DESC"
)

func NewSortOrder(sortOrder string) (Order, error) {
	switch strings.ToUpper(sortOrder) {
	case ASC.SQLString(), ASC.ComparisonOperator(), ASC.MathOperator():
		return ASC, nil
	case DESC.SQLString(), DESC.ComparisonOperator(), DESC.MathOperator():
		return DESC, nil
	default:
		return UNKNOWN, fmt.Errorf("invalid sort order: %s", sortOrder)
	}
}

func (s Order) Invert() Order {
	switch s {
	case ASC:
		return DESC
	case DESC:
		return ASC
	default:
		return UNKNOWN
	}
}

func (s Order) SQLString() string {
	return string(s)
}

func (s Order) MathOperator() string {
	switch s {
	case ASC:
		return "+"
	case DESC:
		return "-"
	default:
		return "UNKNOWN"
	}
}

func (s Order) ComparisonOperator() string {
	switch s {
	case ASC:
		return ">"
	case DESC:
		return "<"
	default:
		return "UNKNOWN"
	}
}
