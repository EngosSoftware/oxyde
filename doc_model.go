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
    AccessGranted = iota // Access to endpoint for selected role is granted.
    AccessDenied         // Access to endpoint for selected role is denied.
    AccessUnknown        // No info about endpoint access rights for selected role.
    AccessError          // Error occurred during collecting access rights for selected role.
)

// Properties of documented API endpoint.
type DocEndpoint struct {
    Id           string  // Unique endpoint identifier.
    Group        string  // Name of the group this endpoint belongs to.
    Version      string  // API version number for endpoint.
    Method       string  // HTTP method name, like GET, POST, PUT or DELETE.
    RootPath     string  // Request root path.
    RequestPath  string  // Request path after root path.
    Summary      string  // Endpoint summary.
    Description  string  // Endpoint detailed description.
    Headers      []Field // Description of request headers.
    Parameters   []Field // Description of request parameters.
    RequestBody  []Field // Description of request body.
    ResponseBody []Field // Description of response body.
    Usages       []Usage // Description of usage examples.
}

// Function createEndpoint creates new API endpoint.
func createEndpoint(group string, version string, summary string, description string) *DocEndpoint {
    return &DocEndpoint{
        Id:          generateId(),
        Group:       group,
        Version:     version,
        Summary:     summary,
        Description: description}
}

type Field struct {
    FieldName   string  // Field name in struct or in array.
    JsonName    string  // Name of the field in JSON object.
    JsonType    string  // Type of the field in JSON object.
    Mandatory   bool    // Flag indicating if field is mandatory in JSON object.
    Description string  // Description of the field.
    Children    []Field // List of child fields (may be empty).
}

func ParseType(i interface{}) []Field {
    typ := reflect.TypeOf(i)
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
            field := createField(childType, childField)
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

func createField(typ reflect.Type, structField reflect.StructField) Field {
    fieldName := structField.Name
    jsonType := jsonType(typ)
    jsonName := structField.Tag.Get(JsonTagName)
    apiTagContent := structField.Tag.Get(ApiTagName)
    mandatory := !strings.HasPrefix(apiTagContent, OptionalPrefix)
    apiTagContent = strings.TrimPrefix(apiTagContent, OptionalPrefix)
    return Field{
        FieldName:   fieldName,
        JsonName:    jsonName,
        JsonType:    jsonType,
        Mandatory:   mandatory,
        Description: apiTagContent,
        Children:    make([]Field, 0)}
}

type Usage struct {
    Summary      string  // Usage summary.
    Description  string  // Usage detailed description.
    Method       string  // HTTP method name.
    Headers      headers // Request headers.
    Url          string  // Request URL.
    RequestBody  string  // Request body as JSON string.
    ResponseBody string  // Response body as JSON string.
    StatusCode   int     // HTTP status code.
}

// Type headers is a map that defines names and values of HTTP request headers.
// Keys are header names and values are header values. This is a convenient way
// to pass any number of headers to functions that call REST endpoints.
type headers map[string]string

// Function parseHeaders traverses the interface given in parameter and retrieves
// names and values of request headers. All request headers required in endpoint call
// should be defined as a struct having string fields (or pointers to strings).
// Each field in such a struct should have a tag named 'json' with the name of the header.
// This way allows to define and document headers and pass header values in one
// single (and simple) structure.
func parseHeaders(any interface{}) headers {
    headersMap := make(headers)
    if any == nil {
        return headersMap
    }
    typ := reflect.TypeOf(any)
    value := reflect.ValueOf(any)
    if typ.Kind() == reflect.Ptr {
        typ = typ.Elem()
        if value.IsNil() || !value.IsValid() {
            return headersMap
        }
        value = reflect.Indirect(value)
    }
    if typ.Kind() != reflect.Struct {
        return headersMap
    }
    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        fieldType := field.Type
        fieldValue := value.Field(i)
        if fieldType.Kind() == reflect.Ptr {
            fieldType = fieldType.Elem()
            if fieldValue.IsNil() {
                continue
            }
            fieldValue = reflect.Indirect(fieldValue)
        }
        if fieldType.Kind() != reflect.String {
            continue
        }
        fieldName := field.Tag.Get(JsonTagName)
        headersMap[fieldName] = fmt.Sprintf("%s", fieldValue)
    }
    return headersMap
}

type roleKey struct {
    method   string // HTTP method name.
    path     string // Request path.
    roleName string // Name of the role.
}

type roles map[roleKey]int