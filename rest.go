/*
 * MIT License
 *
 * Copyright (c) 2017-2019 Dariusz Depta Engos Software
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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

// Interface for request context. Instances of this interface provide additional
// data required to successfully execute HTTP requests.
type Context interface {
    GetUrl() string     // Returns URL of the endpoint to be called.
    GetVerbose() bool   // Returns flag indicating if executing process should be more verbose.
    GetVersion() string // Returns API version to be replaced in URL.
}

// Function HttpGETString executes HTTP GET request and returns simple text result (not JSON string!)
func HttpGETString(c Context, dc *DocContext, path string, params interface{}, result interface{}, status int) {
    requestPath, err := prepareRequestPath(path, c.GetVersion(), params)
    uri := prepareUri(c, requestPath)
    displayRequestDetails(c, httpGET, uri)
    req, err := http.NewRequest(httpGET, uri, nil)
    panicOnError(err)
    setRequestHeaders(c, req, nil)
    client := http.Client{}
    res, err := client.Do(req)
    panicOnError(err)
    panicOnUnexpectedStatusCode(c, status, res)
    responseBody := readResponseBody(c, res)
    collectDocumentationData(c, dc, res, httpGET, path, requestPath, nil, params, nil, result, nil, responseBody)
    resultFields := ParseType(result)
    if len(resultFields) == 1 && resultFields[0].JsonName == "-" && resultFields[0].JsonType == "string" {
        reflect.ValueOf(result).Elem().Field(0).SetString(string(responseBody))
    }
}

// Function HttpGET executes HTTP GET request and returns JSON result.
func HttpGET(
    c Context,           /* Request context. */
    dc *DocContext,      /* Documentation context. */
    path string,         /* Request path. */
    headers interface{}, /* Request headers. */
    params interface{},  /* Request parameters. */
    result interface{},  /* Response payload (response body). */
    status int           /* Expected HTTP status code. */) {
    var responseBody []byte
    requestPath, err := prepareRequestPath(path, c.GetVersion(), params)
    uri := prepareUri(c, requestPath)
    displayRequestDetails(c, httpGET, uri)
    req, err := http.NewRequest(httpGET, uri, nil)
    panicOnError(err)
    setRequestHeaders(c, req, headers)
    client := http.Client{}
    res, err := client.Do(req)
    panicOnError(err)
    panicOnUnexpectedStatusCode(c, status, res)
    if nilValue(result) {
        responseBody = nil
    } else {
        responseBody = readResponseBody(c, res)
        err = json.Unmarshal(responseBody, result)
        panicOnError(err)
    }
    collectDocumentationData(c, dc, res, httpGET, path, requestPath, headers, params, nil, result, nil, responseBody)
}

// Function HttpPOST executes HTTP POST request.
func HttpPOST(
    c Context,           /* Request context. */
    dc *DocContext,      /* Documentation context. */
    path string,         /* Request path. */
    headers interface{}, /* Request headers. */
    params interface{},  /* Request parameters. */
    body interface{},    /* Request payload (request body). */
    result interface{},  /* Response payload (response body). */
    status int           /* Expected HTTP status code. */) {
    httpCall(c, dc, httpPOST, path, headers, params, body, result, status)
}

// Function HttpPUT executes HTTP PUT request.
func HttpPUT(
    c Context,           /* Request context. */
    dc *DocContext,      /* Documentation context. */
    path string,         /* Request path. */
    headers interface{}, /* Request headers. */
    params interface{},  /* Request parameters. */
    body interface{},    /* Request payload (request body). */
    result interface{},  /* Response payload (response body). */
    status int           /* Expected HTTP status code. */) {
    httpCall(c, dc, httpPUT, path, headers, params, body, result, status)
}

// Function HttpDELETE executes HTTP DELETE request.
func HttpDELETE(
    c Context,           /* Request context. */
    dc *DocContext,      /* Documentation context. */
    path string,         /* Request path. */
    headers interface{}, /* Request headers. */
    params interface{},  /* Request parameters. */
    body interface{},    /* Request payload (request body). */
    result interface{},  /* Response payload (response body). */
    status int           /* Expected HTTP status code. */) {
    httpCall(c, dc, httpDELETE, path, headers, params, body, result, status)
}

// Function httpCall executes HTTP request with specified HTTP method and parameters.
func httpCall(
    c Context,
    dc *DocContext,
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
    requestPath, err := prepareRequestPath(path, c.GetVersion(), params)
    panicOnError(err)
    uri := prepareUri(c, requestPath)
    displayRequestDetails(c, method, uri)
    if nilValue(body) {
        requestBody = nil
        displayRequestPayload(c, nil)
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
        displayRequestPayload(c, requestBody)
        req, err = http.NewRequest(method, uri, bytes.NewReader(requestBody))
        panicOnError(err)
        req.Header.Add("Content-Type", "application/json")
    }
    setRequestHeaders(c, req, headers)
    client := http.Client{}
    res, err := client.Do(req)
    panicOnError(err)
    panicOnUnexpectedStatusCode(c, status, res)
    if nilValue(result) {
        responseBody = nil
    } else {
        responseBody = readResponseBody(c, res)
        err = json.Unmarshal(responseBody, result)
        panicOnError(err)
    }
    collectDocumentationData(c, dc, res, method, path, requestPath, headers, params, body, result, requestBody, responseBody)
}

func collectDocumentationData(
    c Context,
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
        endpoint.RootPath = c.GetUrl()
        endpoint.RequestPath = preparePath(path, c.GetVersion())
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
            Url:          c.GetUrl() + requestPath,
            RequestBody:  prettyPrint(requestBody),
            ResponseBody: prettyPrint(responseBody),
            StatusCode:   res.StatusCode}
        endpoint.Usages = append(endpoint.Usages, usage)
    }
    dc.SaveRole(method, preparePath(path, c.GetVersion()), res.StatusCode)
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
func setRequestHeaders(c Context, req *http.Request, headers interface{}) {
    for name, value := range parseHeaders(headers) {
        req.Header.Add(name, value)
    }
}

// Function readResponseBody reads and returns the body of HTTP response.
func readResponseBody(c Context, res *http.Response) []byte {
    body, err := ioutil.ReadAll(res.Body)
    panicOnError(err)
    err = res.Body.Close()
    panicOnError(err)
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
            fmt.Printf("\n===> REQUEST PAYLOAD:\n%s\n", prettyPrint(payload))
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
            fmt.Printf("\n<=== RESPONSE BODY:\n%s\n", prettyPrint(body))
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
        separator := makeText("-", 120)
        fmt.Printf("\n\n%s\n>     ERROR: unexpected status code\n>  Expected: %d\n>    Actual: %d\n%s\n\n",
            separator,
            expectedCode,
            res.StatusCode,
            separator)
        brexit()
    }
}
