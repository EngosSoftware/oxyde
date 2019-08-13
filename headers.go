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
    "fmt"
    "reflect"
)

// Type headers is a map that defines names and values of HTTP request headers.
// Keys are header names and values are header values. This is a convenient way
// to pass any number of headers to functions that call REST endpoints.
type headers map[string]string

// Function parseHeaders traverses the interface given in parameter and retrieves
// names and values of request headers. All request headers required in endpoint call
// should be defined as a struct having string fields (or pointers to strings).
// Each field in such a struct should have a tag named 'json' with the name of header.
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
