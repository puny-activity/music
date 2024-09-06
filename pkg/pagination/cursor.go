package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type CursorType string

const CursorPrev CursorType = "prev"
const CursorNext CursorType = "next"

type CursorParameter struct {
	Parameter
	Value any
}

func NewCursorParameter(param Parameter, value any) CursorParameter {
	return CursorParameter{
		Parameter: param,
		Value:     value,
	}
}

type Cursor struct {
	cursorType CursorType
	parameters []CursorParameter
}

func NewCursor(cursorType CursorType) Cursor {
	return Cursor{
		cursorType: cursorType,
		parameters: make([]CursorParameter, 0),
	}
}

func DecodeCursor(encodedCursor string) (Cursor, error) {
	data, err := base64.URLEncoding.DecodeString(encodedCursor)
	if err != nil {
		return Cursor{}, fmt.Errorf("failed to decode base64 cursor: %v", err)
	}

	var decoded struct {
		Type       CursorType        `json:"type"`
		Parameters []CursorParameter `json:"parameters"`
	}

	err = json.Unmarshal(data, &decoded)
	if err != nil {
		return Cursor{}, fmt.Errorf("failed to unmarshal cursor: %v", err)
	}

	return Cursor{
		cursorType: decoded.Type,
		parameters: decoded.Parameters,
	}, nil
}

func (c Cursor) Type() CursorType {
	return c.cursorType
}

func (c Cursor) Add(parameter CursorParameter) Cursor {
	c.parameters = append(c.parameters, parameter)
	return c
}

func (c Cursor) GetAll() []CursorParameter {
	return c.parameters
}

func (c Cursor) Encode() (string, error) {
	data := struct {
		Type       CursorType        `json:"type"`
		Parameters []CursorParameter `json:"parameters"`
	}{
		Type:       c.cursorType,
		Parameters: c.parameters,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cursor: %v", err)
	}

	return base64.URLEncoding.EncodeToString(jsonData), nil
}

func (c Cursor) BuildWhere() (string, error) {
	if len(c.parameters) == 0 {
		return "", nil
	}

	var sb strings.Builder
	for i, param := range c.parameters {
		if i > 0 {
			sb.WriteString(" OR ")
		}

		if i == 0 {
			if c.cursorType == CursorPrev {
				sb.WriteString(fmt.Sprintf("%s %s %v", param.FieldName, param.SortOrder.Invert().ComparisonOperator(), formatValue(param.Value)))
			} else if c.cursorType == CursorNext {
				sb.WriteString(fmt.Sprintf("%s %s %v", param.FieldName, param.SortOrder.ComparisonOperator(), formatValue(param.Value)))
			}
		} else {
			prevParam := c.parameters[i-1]
			if c.cursorType == CursorPrev {
				sb.WriteString(fmt.Sprintf("(%s = %v AND %s %s %v)", prevParam.FieldName, formatValue(prevParam.Value), param.FieldName, param.SortOrder.Invert().ComparisonOperator(), formatValue(param.Value)))
			} else if c.cursorType == CursorNext {
				sb.WriteString(fmt.Sprintf("(%s = %v AND %s %s %v)", prevParam.FieldName, formatValue(prevParam.Value), param.FieldName, param.SortOrder.ComparisonOperator(), formatValue(param.Value)))
			}
		}
	}

	return sb.String(), nil
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

type CursorPair struct {
	PrevCursor *Cursor
	NextCursor *Cursor
}

type CursorPagination struct {
	Limit      int
	Parameters []Parameter
	Cursor     *Cursor
}
