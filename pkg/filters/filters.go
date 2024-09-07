package filters

import (
	"fmt"
	"strings"
)

type Filter struct {
	AndFilters []AndFilter
}

func NewFilter() Filter {
	return Filter{
		AndFilters: make([]AndFilter, 0),
	}
}

func (f Filter) Add(andFilter AndFilter) Filter {
	return Filter{
		AndFilters: append(f.AndFilters, andFilter),
	}
}

func (f Filter) IsEmpty() bool {
	return len(f.AndFilters) == 0
}

type AndFilter struct {
	OrFilter []OrFilter
}

func NewAndFilter() AndFilter {
	return AndFilter{
		OrFilter: make([]OrFilter, 0),
	}
}

func (f AndFilter) Add(orFilter OrFilter) AndFilter {
	return AndFilter{
		OrFilter: append(f.OrFilter, orFilter),
	}
}

type OrFilter struct {
	Key   string
	Value any
}

func NewOrFilter(key string, value any) OrFilter {
	return OrFilter{
		Key:   key,
		Value: value,
	}
}

func (f Filter) BuildWhere() string {
	if len(f.AndFilters) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(" ( ")
	for i, andFilter := range f.AndFilters {
		if i > 0 {
			sb.WriteString(" AND ")
		}

		for j, orFilter := range andFilter.OrFilter {
			if len(andFilter.OrFilter) > 1 && j == 0 {
				sb.WriteString(" ( ")
			}
			if j > 0 {
				sb.WriteString(" OR ")
			}

			sb.WriteString(fmt.Sprintf(" %s = %v ", orFilter.Key, formatValue(orFilter.Value)))

			if len(andFilter.OrFilter) > 1 && j == len(andFilter.OrFilter)-1 {
				sb.WriteString(" ) ")
			}
		}
	}
	sb.WriteString(" ) ")

	return sb.String()
}

func formatValue(value any) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
	case int, int64, float64:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("'%v'", v)
	}
}
