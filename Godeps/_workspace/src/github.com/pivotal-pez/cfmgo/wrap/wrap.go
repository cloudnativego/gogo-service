//Wrap is a simple helper to wrap API response data and errors in a consistent structure.
package wrap

//ResponseWrapper provides a standard structure for API responses.
type ResponseWrapper struct {
	//Status indicates the result of a request as "success" or "error"
	Status string `json:"status"`
	//Data holds the payload of the response
	Data interface{} `json:"data,omitempty"`
	//Message contains the nature of an error
	Message string `json:"message,omitempty"`
	//Count contains the number of records in the result set
	Count int `json:"count,omitempty"`
}

//One is for successful requests that yield a single result value.
func One(data interface{}) (rsp *ResponseWrapper) {
	rsp = &ResponseWrapper{
		Status: successStatus,
		Data:   data,
	}
	return
}

//Collection is for successful reuqests that have the potential
//to yield multiple results.
func Many(data interface{}, count int) (rsp *ResponseWrapper) {
	rsp = &ResponseWrapper{
		Status: successStatus,
		Data:   data,
		Count:  count,
	}
	return
}

//Error is for requests that yield no results due to an error.
func Error(message string) (rsp *ResponseWrapper) {
	rsp = &ResponseWrapper{
		Status:  errorStatus,
		Message: message,
	}
	return
}
