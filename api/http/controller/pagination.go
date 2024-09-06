package httpcontroller

import (
	"github.com/puny-activity/music/internal/errs"
	"github.com/puny-activity/music/pkg/pagination"
	"github.com/puny-activity/music/pkg/werr"
	"net/http"
	"strconv"
)

type Pagination struct {
	PrevCursor *string `json:"prev_cursor,omitempty"`
	NextCursor *string `json:"next_cursor,omitempty"`
}

func ExtractCursorPagination(r *http.Request) (pagination.CursorPagination, error) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return pagination.CursorPagination{}, errs.LimitNotSpecified
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return pagination.CursorPagination{}, errs.InvalidLimitParameter
	}

	parametersStr := r.URL.Query().Get("sort")
	cursorStr := r.URL.Query().Get("cursor")
	if parametersStr == "" && cursorStr == "" {
		return pagination.CursorPagination{}, errs.SortAndCursorNotSpecified
	} else if parametersStr != "" && cursorStr != "" {
		return pagination.CursorPagination{}, errs.SortAndCursorSpecified
	}

	if parametersStr != "" {
		parameters, err := pagination.ParseParameters(parametersStr)
		if err != nil {
			return pagination.CursorPagination{}, werr.WrapES(errs.InvalidSortParameter, err.Error())
		}
		return pagination.CursorPagination{
			Limit:      limit,
			Parameters: parameters,
			Cursor:     nil,
		}, nil
	} else {
		cursor, err := pagination.DecodeCursor(cursorStr)
		if err != nil {
			return pagination.CursorPagination{}, werr.WrapES(errs.InvalidCursorParameter, err.Error())
		}
		return pagination.CursorPagination{
			Limit:      limit,
			Parameters: nil,
			Cursor:     &cursor,
		}, nil
	}
}
