package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wisbery/oxyde/common"
	"github.com/wisbery/oxyde/doc"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	httpGET    = "GET"
	httpPOST   = "POST"
	httpPUT    = "PUT"
	httpDELETE = "DELETE"
)

// Interface for request context. Instances of this interface
// provide data required to successfully execute HTTP requests.
type Context interface {
	GetUrl() string                // Returns URL of the endpoint to be called.
	GetAuthorizationToken() string // Returns access token to be passed in 'Authorization' header.
	GetHeaders() map[string]string // Returns map of HTTP headers to be passed to endpoint call.
	GetVerbose() bool              // Returns flag indicating if executing process should be more verbose.
}

// Function HttpGETString executes HTTP GET request and returns simple text result (not JSON string!)
func HttpGETString(c Context, dc *doc.Context, path string, params interface{}, result interface{}, status int) {
	requestPath, err := prepareRequestPath(path, params)
	uri := prepareUri(c, requestPath)
	displayRequestDetails(c, httpGET, uri)
	req, err := http.NewRequest(httpGET, uri, nil)
	common.PanicOnError(err)
	setRequestHeaders(c, req)
	client := http.Client{}
	res, err := client.Do(req)
	common.PanicOnError(err)
	panicOnUnexpectedStatusCode(c, status, res)
	responseBody := readResponseBody(c, res)
	collectDocumentationData(c, dc, res, httpGET, path, requestPath, params, nil, result, nil, responseBody)
	resultFields := doc.ParseObject(result)
	if len(resultFields) == 1 && resultFields[0].JsonName == "-" && resultFields[0].JsonType == "string" {
		reflect.ValueOf(result).Elem().Field(0).SetString(string(responseBody))
	}
}

// Function HttpGET executes HTTP GET request and returns JSON result.
func HttpGET(c Context, dc *doc.Context, path string, params interface{}, result interface{}, status int) {
	var responseBody []byte
	requestPath, err := prepareRequestPath(path, params)
	uri := prepareUri(c, requestPath)
	displayRequestDetails(c, httpGET, uri)
	req, err := http.NewRequest(httpGET, uri, nil)
	common.PanicOnError(err)
	setRequestHeaders(c, req)
	client := http.Client{}
	res, err := client.Do(req)
	common.PanicOnError(err)
	panicOnUnexpectedStatusCode(c, status, res)
	if common.NilValue(result) {
		responseBody = nil
	} else {
		responseBody = readResponseBody(c, res)
		err = json.Unmarshal(responseBody, result)
		common.PanicOnError(err)
	}
	collectDocumentationData(c, dc, res, httpGET, path, requestPath, params, nil, result, nil, responseBody)
}

// Function HttpPOST executes HTTP POST request.
func HttpPOST(c Context, dc *doc.Context, path string, payload interface{}, result interface{}, status int) {
	httpCall(c, dc, httpPOST, path, nil, payload, result, status)
}

// Function HttpPUT executes HTTP PUT request.
func HttpPUT(c Context, dc *doc.Context, path string, payload interface{}, result interface{}, status int) {
	httpCall(c, dc, httpPUT, path, nil, payload, result, status)
}

// Function HttpDELETE executes HTTP DELETE request.
func HttpDELETE(c Context, dc *doc.Context, path string, params interface{}, payload interface{}, result interface{}, status int) {
	httpCall(c, dc, httpDELETE, path, params, payload, result, status)
}

// Function httpCall executes HTTP request with specified HTTP method and parameters.
func httpCall(c Context, dc *doc.Context, method string, path string, params interface{}, payload interface{}, result interface{}, status int) {
	var req *http.Request
	var requestBody []byte
	var responseBody []byte
	var err error
	requestPath, err := prepareRequestPath(path, params)
	common.PanicOnError(err)
	uri := prepareUri(c, requestPath)
	displayRequestDetails(c, method, uri)
	if common.NilValue(payload) {
		requestBody = nil
		displayRequestPayload(c, nil)
		req, err = http.NewRequest(method, uri, nil)
		common.PanicOnError(err)
	} else {
		requestBody, err = json.Marshal(payload)
		common.PanicOnError(err)
		displayRequestPayload(c, requestBody)
		req, err = http.NewRequest(method, uri, bytes.NewReader(requestBody))
		common.PanicOnError(err)
		req.Header.Add("Content-Type", "application/json")
	}
	setRequestHeaders(c, req)
	client := http.Client{}
	res, err := client.Do(req)
	common.PanicOnError(err)
	panicOnUnexpectedStatusCode(c, status, res)
	if common.NilValue(result) {
		responseBody = nil
	} else {
		responseBody = readResponseBody(c, res)
		err = json.Unmarshal(responseBody, result)
		common.PanicOnError(err)
	}
	collectDocumentationData(c, dc, res, method, path, requestPath, params, payload, result, requestBody, responseBody)
}

func collectDocumentationData(c Context, dc *doc.Context, res *http.Response, method string, path string, requestPath string, params interface{}, payload interface{}, result interface{}, requestBody []byte, responseBody []byte) {
	if endpoint := dc.GetEndpoint(); endpoint != nil && dc.CollectDescriptionMode() {
		endpoint.Method = method
		endpoint.UrlRoot = c.GetUrl()
		endpoint.UrlPath = path
		if common.NilValue(params) {
			endpoint.Parameters = nil
		} else {
			endpoint.Parameters = doc.ParseObject(params)
		}
		if common.NilValue(payload) {
			endpoint.RequestBody = nil
		} else {
			endpoint.RequestBody = doc.ParseObject(payload)
		}
		if common.NilValue(result) {
			endpoint.ResponseBody = nil
		} else {
			endpoint.ResponseBody = doc.ParseObject(result)
		}
	}
	if endpoint := dc.GetEndpoint(); endpoint != nil && dc.CollectExamplesMode() {
		examples := endpoint.Examples
		if examples == nil {
			endpoint.Examples = make([]doc.Example, 0)
		}
		example := doc.Example{
			Summary:      dc.GetExampleSummary(),
			Description:  dc.GetExampleDescription(),
			Method:       method,
			Uri:          c.GetUrl() + requestPath,
			StatusCode:   res.StatusCode,
			RequestBody:  common.PrettyPrint(requestBody),
			ResponseBody: common.PrettyPrint(responseBody)}
		endpoint.Examples = append(endpoint.Examples, example)
	}
	dc.SaveRole(method, path, res.StatusCode)
	dc.StopCollecting()
}

func prepareRequestPath(path string, params interface{}) (string, error) {
	if common.NilValue(params) {
		return path, nil
	}
	paramsType := common.TypeOfValue(params)
	if paramsType.Kind().String() != "struct" {
		return "", errors.New("only struct parameters are allowed")
	}
	firstParameter := true
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)
		fieldJsonName := field.Tag.Get(common.JsonTagName)
		placeholder := "{" + fieldJsonName + "}"
		value := common.ValueOfValue(params).Field(i).Interface()
		if !common.NilValue(value) {
			valueStr := url.PathEscape(fmt.Sprintf("%v", common.ValueOfValue(value)))
			if strings.Contains(path, placeholder) {
				path = strings.ReplaceAll(path, placeholder, valueStr)
			} else {
				if firstParameter {
					path = path + "?"
				} else {
					path = path + "&"
				}
				path = path + fieldJsonName + "=" + valueStr
				firstParameter = false
			}
		}
	}
	return path, nil
}

// Function setRequestHeaders adds to the request authorization header and user defined headers.
func setRequestHeaders(c Context, req *http.Request) {
	if len(c.GetAuthorizationToken()) > 0 {
		req.Header.Add("Authorization", c.GetAuthorizationToken())
	}
	if c.GetHeaders() != nil {
		for name, value := range c.GetHeaders() {
			req.Header.Add(name, value)
		}
	}
}

// Function readResponseBody reads and returns the body of HTTP response.
func readResponseBody(c Context, res *http.Response) []byte {
	body, err := ioutil.ReadAll(res.Body)
	common.PanicOnError(err)
	err = res.Body.Close()
	common.PanicOnError(err)
	displayResponseBody(c, body)
	return body
}

// Function prepareUri concatenates URL defined in context with
// request path and returns full URI of HTTP request.
func prepareUri(c Context, path string) string {
	return c.GetUrl() + path
}

// Function displayRequestDetails writes to standard output
// request method and URI.
func displayRequestDetails(c Context, method string, uri string) {
	if c.GetVerbose() {
		fmt.Printf("\n\n===> %s:\n%s\n", method, uri)
	}
}

// Function displayRequestPayload writes to standard output
// pretty-printed request payload.
func displayRequestPayload(c Context, payload []byte) {
	if c.GetVerbose() {
		if payload == nil {
			fmt.Printf("\n===> REQUEST PAYLOAD:\n(none)\n")
		} else {
			fmt.Printf("\n===> REQUEST PAYLOAD:\n%s\n", common.PrettyPrint(payload))
		}
	}
}

// Function displayResponseBody writes to standard output
// pretty-printed response body when verbose mode is on.
func displayResponseBody(c Context, body []byte) {
	if c.GetVerbose() {
		if body == nil {
			fmt.Printf("\n<=== RESPONSE BODY:\n(none)\n")
		} else {
			fmt.Printf("\n<=== RESPONSE BODY:\n%s\n", common.PrettyPrint(body))
		}
	}
}

// Function panicOnUnexpectedStatusCode displays error message and panics when
// actual HTTP response status code differs from the expected one.
func panicOnUnexpectedStatusCode(c Context, expectedCode int, res *http.Response) {
	// display the returned status code if the same as expected
	if c.GetVerbose() {
		fmt.Printf("\n<=== STATUS:\n%d\n", res.StatusCode)
	}
	// check if the expected status code is the same as returned by server
	if res.StatusCode != expectedCode {
		readResponseBody(c, res)
		separator := common.MakeString('-', 120)
		fmt.Printf("\n\n%s\n>     ERROR: unexpected status code\n>  Expected: %d\n>    Actual: %d\n%s\n\n",
			separator,
			expectedCode,
			res.StatusCode,
			separator)
		common.BrExit()
	}
}
