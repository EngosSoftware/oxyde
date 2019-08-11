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
    "errors"
    "fmt"
    "reflect"
    "strings"
)

const (
    CollectNone        = iota // Collect no documentation data.
    CollectDescription        // Collect endpoint description.
    CollectExamples           // Collect an example of endpoint usage.
    CollectAll                // Collect endpoint description and usage example.
)

const (
    AccessGranted = iota // Access to endpoint is granted.
    AccessDenied         // Access to endpoint is denied.
    AccessUnknown        // No info about endpoint access rights.
    AccessError          // Error occurred during collecting access rights.
)

type RoleKey struct {
    method   string
    path     string
    roleName string
}

type DocumentContext struct {
    mode               int             // Documentation collecting data mode.
    exampleSummary     string          // Summary for the next collected example.
    exampleDescription string          // Detailed description for the next collected example.
    roleName           string          // Role name of the principal for the next endpoint example.
    endpoint           *Endpoint       // Currently documented endpoint data.
    endpoints          []Endpoint      // List of documented endpoints.
    roleNames          []string        // Role names in order they should be displayed.
    roles              map[RoleKey]int // Map of access roles tested for endpoints.
}

func CreateDocumentContext() *DocumentContext {
    return &DocumentContext{
        mode:      CollectNone,
        endpoint:  nil,
        endpoints: make([]Endpoint, 0),
        roles:     make(map[RoleKey]int)}
}

func (dc *DocumentContext) ClearDocumentation() {
    dc.mode = CollectNone
    dc.endpoint = nil
    dc.endpoints = make([]Endpoint, 0)
    dc.roles = make(map[RoleKey]int)
}

func (dc *DocumentContext) PublishDocumentation() {
    for _, endpoint := range dc.endpoints {
        PrintEndpoint(endpoint)
    }
}

func (dc *DocumentContext) NewEndpointDocumentation(id string, tag string, summary string) {
    if id == "" {
        id = GenerateId()
    }
    dc.endpoint = &Endpoint{
        Id:      id,
        Summary: summary}
    dc.endpoint.AddTag(tag)
}

func (dc *DocumentContext) GetEndpoint() *Endpoint {
    return dc.endpoint
}

func (dc *DocumentContext) CollectDescription() {
    dc.mode = CollectDescription
}

func (dc *DocumentContext) CollectDescriptionMode() bool {
    return dc.mode == CollectDescription || dc.mode == CollectAll
}

func (dc *DocumentContext) CollectExamples(exampleSummary string, exampleDescription string) {
    exampleSummary = strings.TrimSpace(exampleSummary)
    exampleDescription = strings.TrimSpace(exampleDescription)
    if exampleSummary != "" || exampleDescription != "" {
        dc.mode = CollectExamples
        dc.exampleSummary = exampleSummary
        dc.exampleDescription = exampleDescription
    }
}

func (dc *DocumentContext) CollectRole(roleName string) {
    dc.roleName = roleName
}

func (dc *DocumentContext) CollectExamplesMode() bool {
    return dc.mode == CollectExamples || dc.mode == CollectAll
}

func (dc *DocumentContext) CollectAll(exampleSummary string) {
    dc.mode = CollectAll
    dc.exampleSummary = exampleSummary
}

func (dc *DocumentContext) StopCollecting() {
    dc.mode = CollectNone
    dc.exampleSummary = ""
    dc.roleName = ""
}

func (dc *DocumentContext) SetRolesOrder(roleOrder []string) {
    dc.roleNames = roleOrder
}

func (dc *DocumentContext) SaveRole(method string, path string, status int) {
    if dc.roleName != "" {
        key := RoleKey{method: method, path: path, roleName: dc.roleName}
        switch status {
        case 200:
            dc.roles[key] = AccessGranted
        case 401:
            dc.roles[key] = AccessDenied
        default:
            dc.roles[key] = AccessError
        }
    }
}

func (dc *DocumentContext) GetRoleNames() []string {
    return dc.roleNames
}

func (dc *DocumentContext) GetAccess(method string, path string, roleName string) int {
    key := RoleKey{
        method:   method,
        path:     path,
        roleName: roleName}
    if access, ok := dc.roles[key]; ok {
        return access
    } else {
        return AccessUnknown
    }
}

func (dc *DocumentContext) SaveEndpointDocumentation() {
    if dc.endpoint != nil {
        dc.endpoints = append(dc.endpoints, *dc.endpoint)
    }
}

func (dc *DocumentContext) GetEndpoints() []Endpoint {
    return dc.endpoints
}

func (dc *DocumentContext) GetExampleSummary() string {
    return dc.exampleSummary
}

func (dc *DocumentContext) GetExampleDescription() string {
    return dc.exampleDescription
}

type Endpoint struct {
    Id           string    // Unique endpoint identifier.
    Tags         []string  // List of tags of endpoint.
    Method       string    // HTTP method name, like GET, POST, PUT or DELETE.
    UrlRoot      string    // Request URL root.
    UrlPath      string    // Request URL path after root.
    Summary      string    // Summary text describing endpoint.
    Parameters   []Field   // Description of request parameters.
    RequestBody  []Field   // Description of request body.
    ResponseBody []Field   // Description of results.
    Examples     []Example // Description of usage examples.
}

func (e *Endpoint) AddTag(tag string) {
    if e.Tags == nil {
        e.Tags = make([]string, 0)
    }
    e.Tags = append(e.Tags, tag)
}

type Field struct {
    JsonName    string  // Name of the field in JSON.
    JsonType    string  // Type of the field in JSON.
    Mandatory   bool    // Flag indicating if field is mandatory in JSON.
    Description string  // Description of the field.
    Children    []Field // List of child fields (may be empty).
}

func CreateField(typ reflect.Type, structField reflect.StructField) Field {
    jsonType := jsonType(typ)
    jsonName := structField.Tag.Get(JsonTagName)
    apiTagContent := structField.Tag.Get(ApiTagName)
    mandatory := !strings.HasPrefix(apiTagContent, OptionalPrefix)
    apiTagContent = strings.TrimPrefix(apiTagContent, OptionalPrefix)
    return Field{
        JsonName:    jsonName,
        JsonType:    jsonType,
        Mandatory:   mandatory,
        Description: apiTagContent,
        Children:    make([]Field, 0)}
}

type Example struct {
    Summary      string // Example summary.
    Description  string // Detailed example description.
    Method       string // HTTP method name.
    Uri          string // Request URI.
    StatusCode   int    // HTTP status code.
    RequestBody  string // Request body as JSON string.
    ResponseBody string // Response body as JSON string.
}

func ParseObject(o interface{}) []Field {
    typ := reflect.TypeOf(o)
    return ParseFields(typ)
}

func ParseFields(typ reflect.Type) []Field {
    switch typ.Kind() {
    case reflect.Ptr:
        return ParseFields(typ.Elem())
    case reflect.Struct:
        fields := make([]Field, 0)
        for i := 0; i < typ.NumField(); i++ {
            childField := typ.Field(i)
            childType := childField.Type
            field := CreateField(childType, childField)
            switch field.JsonType {
            case "object":
                field.Children = append(field.Children, ParseFields(childType)...)
            case "array":
                field.Children = append(field.Children, ParseFields(childType.Elem())...)
            }
            fields = append(fields, field)
        }
        return fields
    }
    return []Field{}
}

func jsonType(typ reflect.Type) string {
    switch typ.Kind() {
    case reflect.Ptr:
        return jsonType(typ.Elem())
    case reflect.Struct:
        return "object"
    case reflect.Slice:
        return "array"
    case reflect.String:
        return "string"
    case reflect.Bool:
        return "boolean"
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
        reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
        reflect.Float32, reflect.Float64:
        return "number"
    default:
        panic(errors.New("unsupported type: " + typ.Kind().String()))
    }
}

func PrintEndpoint(endpoint Endpoint) {
    fmt.Printf("\n\n%s\n%s\n\n", endpoint.Method, endpoint.UrlRoot)
    fmt.Println("Parameters:")
    fmt.Println(makeText("-", 120))
    PrintFields(endpoint.RequestBody, "  ", 0)
    fmt.Println()
    fmt.Println(makeText("-", 120))
    fmt.Printf("\n")
    fmt.Println("ResponseBody:")
    fmt.Println(makeText("-", 120))
    PrintFields(endpoint.ResponseBody, "  ", 0)
    fmt.Println()
    fmt.Println(makeText("-", 120))
    fmt.Printf("\n")
    for _, usage := range endpoint.Examples {
        PrintExample(usage)
    }
}

func PrintFields(fields []Field, indent string, level int) {
    for i, field := range fields {
        nameText := indent + field.JsonName
        nameText = nameText + strings.Repeat(" ", 30-len(nameText))
        typeText := field.JsonType
        typeText = typeText + strings.Repeat(" ", 20-len(typeText))
        mandatoryText := "N"
        if field.Mandatory {
            mandatoryText = "Y"
        }
        if i > 0 || level > 0 {
            fmt.Println()
        }
        fmt.Printf("|%s | %s | %s | %s", nameText, typeText, mandatoryText, field.Description)
        if field.JsonType == "object" || field.JsonType == "array" {
            PrintFields(field.Children, indent+indent, level+1)
        }
    }
}

func PrintExample(usage Example) {
    fmt.Printf("\nExample:\n")
    fmt.Printf("%d %s %s\n", usage.StatusCode, usage.Method, usage.Uri)
    fmt.Printf("Parameters:\n%s\n", usage.RequestBody)
    fmt.Printf("ResponseBody:\n%s\n", usage.ResponseBody)
}
