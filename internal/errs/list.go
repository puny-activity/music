package errs

import "fmt"

var (
	InvalidAPIVersion         = fmt.Errorf("invalid api version")
	LimitNotSpecified         = fmt.Errorf("limit not specified")
	InvalidLimitParameter     = fmt.Errorf("invalid limit parameter")
	SortAndCursorNotSpecified = fmt.Errorf("sort and cursor not specified")
	InvalidSortParameter      = fmt.Errorf("invalid sort parameter")
	InvalidCursorParameter    = fmt.Errorf("invalid cursor parameter")
	SortAndCursorSpecified    = fmt.Errorf("sort and cursor specified")
	UnknownSortParameter      = fmt.Errorf("unknown sort parameter")

	Unexpected = fmt.Errorf("unexpected error")
)

// Unexpected
var unexpectedError = internalError{
	error: Unexpected,
	code:  "U-1",
}

var errorList = []internalError{
	// Request
	{
		error: InvalidAPIVersion,
		code:  "R-1",
	},
	{
		error: LimitNotSpecified,
		code:  "R-2",
	},
	{
		error: InvalidLimitParameter,
		code:  "R-3",
	},
	{
		error: SortAndCursorNotSpecified,
		code:  "R-4",
	},
	{
		error: InvalidSortParameter,
		code:  "R-5",
	},
	{
		error: InvalidCursorParameter,
		code:  "R-6",
	},
	{
		error: SortAndCursorSpecified,
		code:  "R-7",
	},
	{
		error: UnknownSortParameter,
		code:  "R-8",
	},
}
