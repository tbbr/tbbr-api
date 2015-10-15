package appError

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
	RecordNotFound = &Err{404, "1000", "RecordNotFound", "The requested record was not found in the database"}
	// RecordValidationFail = "record_validation_fail"

	JSONParseFailure = &Err{404, "3000", "JSONParseFailure", "The server encountered an error while parsing JSON"}
)
