package oxyde

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

// Request context. Provides additional data required to successfully execute HTTP requests.
type Context struct {
	Url      string // The URL of tested endpoint.
	Token    string // Authorization token.
	UserName string // Name of the authorized user.
	RoleName string // Name of the current role of the authorized user.
	Verbose  bool   // Flag indicating if executing process should be more verbose.
	Version  string // API version to be used in endpoint URL.
}

func CreateContext() *Context {
	return &Context{
		Url:      "",
		Token:    "",
		UserName: "",
		RoleName: "",
		Verbose:  false,
		Version:  "v1",
	}
}

// Function HttpGETString executes HTTP GET request and returns simple text result (not JSON string!)
func HttpGETString(
	ctx *Context,
	dtx *DocContext,
	path string,
	params interface{},
	result interface{},
	status int) {
	requestPath, err := prepareRequestPath(path, ctx.Version, params)
	uri := prepareUri(ctx, requestPath)
	displayRequestDetails(ctx, httpGET, uri)
	req, err := http.NewRequest(httpGET, uri, nil)
	panicOnError(err)
	setRequestHeaders(ctx, req, nil)
	client := http.Client{}
	res, err := client.Do(req)
	panicOnError(err)
	panicOnUnexpectedStatusCode(ctx, status, res)
	responseBody := readResponseBody(ctx, res)
	collectDocumentationData(ctx, dtx, res, httpGET, path, requestPath, nil, params, nil, result, nil, responseBody)
	resultFields := ParseType(result)
	if len(resultFields) == 1 && resultFields[0].JsonName == "-" && resultFields[0].JsonType == "string" {
		reflect.ValueOf(result).Elem().Field(0).SetString(string(responseBody))
	}
}

// Function HttpGET executes HTTP GET request and returns JSON result.
func HttpGET(
	ctx *Context, /* Request context. */
	dtx *DocContext, /* Documentation context. */
	path string, /* Request path. */
	headers interface{}, /* Request headers. */
	params interface{}, /* Request parameters. */
	result interface{}, /* Response payload (response body). */
	status int /* Expected HTTP status code. */) {
	var responseBody []byte
	requestPath, err := prepareRequestPath(path, ctx.Version, params)
	uri := prepareUri(ctx, requestPath)
	displayRequestDetails(ctx, httpGET, uri)
	req, err := http.NewRequest(httpGET, uri, nil)
	panicOnError(err)
	setRequestHeaders(ctx, req, headers)
	client := http.Client{}
	res, err := client.Do(req)
	panicOnError(err)
	panicOnUnexpectedStatusCode(ctx, status, res)
	if nilValue(result) {
		responseBody = nil
	} else {
		responseBody = readResponseBody(ctx, res)
		err = json.Unmarshal(responseBody, result)
		panicOnError(err)
	}
	collectDocumentationData(ctx, dtx, res, httpGET, path, requestPath, headers, params, nil, result, nil, responseBody)
}

// Function HttpPOST executes HTTP POST request.
func HttpPOST(
	ctx *Context, /* Request context. */
	dtx *DocContext, /* Documentation context. */
	path string, /* Request path. */
	headers interface{}, /* Request headers. */
	params interface{}, /* Request parameters. */
	body interface{}, /* Request payload (request body). */
	result interface{}, /* Response payload (response body). */
	status int /* Expected HTTP status code. */) {
	httpCall(ctx, dtx, httpPOST, path, headers, params, body, result, status)
}

// Function HttpPUT executes HTTP PUT request.
func HttpPUT(
	ctx *Context, /* Request context. */
	dtx *DocContext, /* Documentation context. */
	path string, /* Request path. */
	headers interface{}, /* Request headers. */
	params interface{}, /* Request parameters. */
	body interface{}, /* Request payload (request body). */
	result interface{}, /* Response payload (response body). */
	status int /* Expected HTTP status code. */) {
	httpCall(ctx, dtx, httpPUT, path, headers, params, body, result, status)
}

// Function HttpDELETE executes HTTP DELETE request.
func HttpDELETE(
	ctx *Context, /* Request context. */
	dtx *DocContext, /* Documentation context. */
	path string, /* Request path. */
	headers interface{}, /* Request headers. */
	params interface{}, /* Request parameters. */
	body interface{}, /* Request payload (request body). */
	result interface{}, /* Response payload (response body). */
	status int /* Expected HTTP status code. */) {
	httpCall(ctx, dtx, httpDELETE, path, headers, params, body, result, status)
}

// Function httpCall executes HTTP request with specified HTTP method and parameters.
func httpCall(
	ctx *Context,
	dtx *DocContext,
	method string,
	path string,
	headers interface{},
	params interface{},
	body interface{},
	result interface{},
	status int) {
	var req *http.Request
	var requestBody []byte
	var responseBody []byte
	var err error
	requestPath, err := prepareRequestPath(path, ctx.Version, params)
	panicOnError(err)
	uri := prepareUri(ctx, requestPath)
	displayRequestDetails(ctx, method, uri)
	if nilValue(body) {
		requestBody = nil
		displayRequestPayload(ctx, nil)
		req, err = http.NewRequest(method, uri, nil)
		panicOnError(err)
	} else {
		bodyFields := ParseType(body)
		if len(bodyFields) == 1 && bodyFields[0].JsonName == "-" && bodyFields[0].JsonType == "string" {
			field := reflect.ValueOf(body).Elem().Field(0)
			if field.Kind() == reflect.Ptr {
				field = reflect.Indirect(field)
			}
			requestBody = []byte(field.String())
		} else {
			requestBody, err = json.Marshal(body)
			panicOnError(err)
		}
		displayRequestPayload(ctx, requestBody)
		req, err = http.NewRequest(method, uri, bytes.NewReader(requestBody))
		panicOnError(err)
		req.Header.Add("Content-Type", "application/json")
	}
	setRequestHeaders(ctx, req, headers)
	client := http.Client{}
	res, err := client.Do(req)
	panicOnError(err)
	panicOnUnexpectedStatusCode(ctx, status, res)
	if nilValue(result) {
		responseBody = nil
	} else {
		responseBody = readResponseBody(ctx, res)
		err = json.Unmarshal(responseBody, result)
		panicOnError(err)
	}
	collectDocumentationData(ctx, dtx, res, method, path, requestPath, headers, params, body, result, requestBody, responseBody)
}

func collectDocumentationData(
	ctx *Context,
	dc *DocContext,
	res *http.Response,
	method string,
	path string,
	requestPath string,
	headers interface{},
	params interface{},
	payload interface{},
	result interface{},
	requestBody []byte,
	responseBody []byte) {
	if endpoint := dc.GetEndpoint(); endpoint != nil && dc.CollectDescriptionMode() {
		endpoint.Method = method
		endpoint.RootPath = ctx.Url
		endpoint.RequestPath = preparePath(path, ctx.Version)
		if nilValue(headers) {
			endpoint.Headers = nil
		} else {
			endpoint.Headers = ParseType(headers)
		}
		if nilValue(params) {
			endpoint.Parameters = nil
		} else {
			endpoint.Parameters = ParseType(params)
		}
		if nilValue(payload) {
			endpoint.RequestBody = nil
		} else {
			endpoint.RequestBody = ParseType(payload)
		}
		if nilValue(result) {
			endpoint.ResponseBody = nil
		} else {
			endpoint.ResponseBody = ParseType(result)
		}
	}
	if endpoint := dc.GetEndpoint(); endpoint != nil && dc.CollectExamplesMode() {
		usages := endpoint.Usages
		if usages == nil {
			endpoint.Usages = make([]Usage, 0)
		}
		usage := Usage{
			Summary:      dc.GetExampleSummary(),
			Description:  dc.GetExampleDescription(),
			Method:       method,
			Headers:      parseHeaders(headers),
			Url:          ctx.Url + requestPath,
			RequestBody:  prettyPrint(requestBody),
			ResponseBody: prettyPrint(responseBody),
			StatusCode:   res.StatusCode}
		endpoint.Usages = append(endpoint.Usages, usage)
	}
	dc.SaveRole(method, preparePath(path, ctx.Version), res.StatusCode)
	dc.StopCollecting()
}

func preparePath(path string, version string) string {
	if version != "" {
		if strings.Contains(path, VersionPlaceholder) {
			path = strings.ReplaceAll(path, VersionPlaceholder, version)
		}
	}
	return path
}

func prepareRequestPath(path string, version string, params interface{}) (string, error) {
	if version != "" {
		if strings.Contains(path, VersionPlaceholder) {
			path = strings.ReplaceAll(path, VersionPlaceholder, version)
		}
	}
	if nilValue(params) {
		return path, nil
	}
	paramsType := TypeOfValue(params)
	if paramsType.Kind().String() != "struct" {
		return "", errors.New("only struct parameters are allowed")
	}
	firstParameter := true
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)
		fieldJsonName := field.Tag.Get(JsonTagName)
		placeholder := "{" + fieldJsonName + "}"
		value := ValueOfValue(params).Field(i).Interface()
		if !nilValue(value) {
			valueStr := url.PathEscape(fmt.Sprintf("%v", ValueOfValue(value)))
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

// Function setRequestHeaders adds headers to the request.
func setRequestHeaders(ctx *Context, req *http.Request, headers interface{}) {
	for name, value := range parseHeaders(headers) {
		req.Header.Add(name, value)
	}
}

// Function readResponseBody reads and returns the body of HTTP response.
func readResponseBody(ctx *Context, res *http.Response) []byte {
	body, err := ioutil.ReadAll(res.Body)
	panicOnError(err)
	err = res.Body.Close()
	panicOnError(err)
	displayResponseBody(ctx, body)
	return body
}

// Function prepareUri concatenates URL defined in context with
// request path and returns full URI of HTTP request.
func prepareUri(ctx *Context, path string) string {
	return ctx.Url + path
}

// Function displayRequestDetails writes to standard output
// request method and URI.
func displayRequestDetails(ctx *Context, method string, uri string) {
	if ctx.Verbose {
		fmt.Printf("\n\n===> %s:\n%s\n", method, uri)
	}
}

// Function displayRequestPayload writes to standard output
// pretty-printed request payload.
func displayRequestPayload(ctx *Context, payload []byte) {
	if ctx.Verbose {
		if payload == nil {
			fmt.Printf("\n===> REQUEST PAYLOAD:\n(none)\n")
		} else {
			fmt.Printf("\n===> REQUEST PAYLOAD:\n%s\n", prettyPrint(payload))
		}
	}
}

// Function displayResponseBody writes to standard output
// pretty-printed response body when verbose mode is on.
func displayResponseBody(ctx *Context, body []byte) {
	if ctx.Verbose {
		if body == nil {
			fmt.Printf("\n<=== RESPONSE BODY:\n(none)\n")
		} else {
			fmt.Printf("\n<=== RESPONSE BODY:\n%s\n", prettyPrint(body))
		}
	}
}

// Function panicOnUnexpectedStatusCode displays error message and panics when
// actual HTTP response status code differs from the expected one.
func panicOnUnexpectedStatusCode(ctx *Context, expectedCode int, res *http.Response) {
	// display the returned status code if the same as expected
	if ctx.Verbose {
		fmt.Printf("\n<=== STATUS:\n%d\n", res.StatusCode)
	}
	// check if the expected status code is the same as returned by server
	if res.StatusCode != expectedCode {
		readResponseBody(ctx, res)
		separator := makeText("-", 120)
		fmt.Printf("\n\n%s\n>     ERROR: unexpected status code\n>  Expected: %d\n>    Actual: %d\n%s\n\n",
			separator,
			expectedCode,
			res.StatusCode,
			separator)
		brexit()
	}
}
