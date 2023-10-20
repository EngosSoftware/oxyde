package oxyde

// JsonApiError is a JSON API implementation of an error.
type JsonApiError struct {
	Status *string `json:"status" api:"The HTTP status code applicable to reported problem."`
	Code   *string `json:"code"   api:"An application-specific error code."`
	Title  *string `json:"title"  api:"A short, human-readable summary of the problem that never changed from occurrence to occurrence of the problem."`
	Detail *string `json:"detail" api:"A human-readable explanation specific to the occurrence of the problem."`
}

// JsonApiErrors is an array of JSON API errors.
type JsonApiErrors = []JsonApiError
