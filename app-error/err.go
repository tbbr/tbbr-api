package appError

import "net/http"

// Err obj used to create application specific errors that conform
// to the JSONApi spec
type Err struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// Error function satisfies the error interface
func (e Err) Error() string {
	return e.Code
}

// NewErr creates a new error and returns it
func NewErr(status int, code string, detail string) Err {
	return Err{
		Status: status,
		Code:   code,
		Title:  "title",
		Detail: detail,
	}
}

// Application specific error codes
var (
	RecordNotFound = Err{
		http.StatusNotFound,
		"1000",
		"RecordNotFound",
		"The requested record was not found in the database",
	}
	RecordValidationFailure = Err{
		http.StatusBadRequest,
		"1001",
		"RecordValidationFailure",
		"The validations on the record failed",
	}
	InvalidParams = Err{
		http.StatusBadRequest,
		"2001",
		"IncorrectParams",
		"The request sent had invalid params",
	}
	JSONParseFailure = Err{
		http.StatusInternalServerError,
		"3000",
		"JSONParseFailure",
		"The server encountered an error while parsing JSON",
	}
)
