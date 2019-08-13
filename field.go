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
    "reflect"
    "strings"
)

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
