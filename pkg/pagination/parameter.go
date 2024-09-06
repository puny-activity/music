package pagination

import (
	"fmt"
	"strings"
)

type Parameter struct {
	FieldName string
	SortOrder Order
}

func NewParameter(fieldName string, sortOrder Order) *Parameter {
	return &Parameter{
		FieldName: fieldName,
		SortOrder: sortOrder,
	}
}

func ParseParameters(parametersStr string) ([]Parameter, error) {
	fields := strings.Split(parametersStr, ",")
	parameters := make([]Parameter, 0, len(fields))

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if len(field) == 0 {
			continue
		}

		sortOrder := ASC
		fieldName := field

		if strings.HasPrefix(field, "-") {
			sortOrder = DESC
			fieldName = field[1:]
		} else if strings.HasPrefix(field, "+") {
			sortOrder = ASC
			fieldName = field[1:]
		}

		fieldName = strings.TrimSpace(fieldName)

		if fieldName == "" {
			return nil, fmt.Errorf("invalid field name in query: %s", field)
		}

		param := Parameter{
			FieldName: fieldName,
			SortOrder: sortOrder,
		}

		parameters = append(parameters, param)
	}

	if len(parameters) == 0 {
		return nil, fmt.Errorf("no valid parameters in query")
	}

	return parameters, nil
}
