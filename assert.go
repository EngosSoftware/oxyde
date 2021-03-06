/*
 * MIT License
 *
 * Copyright (c) 2017-2020 Engos Software
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

// Function AssertNil asserts that actual value is nil.
// When actual value is not nil, an error is reported.
func AssertNil(actual interface{}) {
    if !isNil(actual) {
        displayAssertionError(nil, actual)
    }
}

func AssertNotNil(actual interface{}) {
    if isNil(actual) {
        displayAssertionError("not nil", actual)
    }
}

func AssertNilError(e error) {
    if e != nil {
        displayAssertionError(nil, e)
    }
}

func AssertNilString(actual *string) {
    if actual != nil {
        displayAssertionError(nil, actual)
    }
}

func AssertNotNilString(actual *string) {
    if actual == nil {
        displayAssertionError("not nil", actual)
    }
}

func AssertTrue(actual bool) {
    if !actual {
        displayAssertionError(true, actual)
    }
}

func AssertFalse(actual bool) {
    if actual {
        displayAssertionError(false, actual)
    }
}

func AssertNotNilId(actual *string) {
    AssertNotNilString(actual)
    AssertEqualInt(36, len(*actual))
}

func AssertEqualString(expected string, actual string) {
    if !equalString(expected, actual) {
        displayAssertionError(expected, actual)
    }
}

func AssertEqualStringNullable(expected *string, actual *string) {
    if !equalStringNullable(expected, actual) {
        if expected != nil && actual != nil {
            displayAssertionError(*expected, *actual)
        }
        displayAssertionError(expected, actual)
    }
}

func AssertNilInt(actual *int) {
    if actual != nil {
        displayAssertionError(nil, actual)
    }
}

func AssertEqualInt(expected int, actual int) {
    if !equalInt(expected, actual) {
        displayAssertionError(expected, actual)
    }
}

func AssertEqualIntNullable(expected *int, actual *int) {
    if !equalIntNullable(expected, actual) {
        if expected != nil && actual != nil {
            displayAssertionError(*expected, *actual)
        }
        displayAssertionError(expected, actual)
    }
}

func AssertEqualInt64Nullable(expected *int64, actual *int64) {
    if !equalInt64Nullable(expected, actual) {
        if expected != nil && actual != nil {
            displayAssertionError(*expected, *actual)
        }
        displayAssertionError(expected, actual)
    }
}

func AssertEqualFloat64(expected float64, actual float64) {
    if !equalFloat64(expected, actual) {
        displayAssertionError(expected, actual)
    }
}

func AssertEqualBool(expected bool, actual bool) {
    if !equalBool(expected, actual) {
        displayAssertionError(expected, actual)
    }
}

// Function isNil checks if the value specified as parameter is nil.
func isNil(value interface{}) bool {
    return value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil())
}

// Function equalString checks if two string values are equal.
func equalString(expected string, actual string) bool {
    return expected == actual
}

// Function equalStringNullable checks if two pointers to string values are equal.
func equalStringNullable(expected *string, actual *string) bool {
    if expected != nil && actual != nil {
        return equalString(*expected, *actual)
    }
    if expected != nil || actual != nil {
        return false
    }
    return true
}

// Function equalInt checks if two int values are equal.
func equalInt(expected int, actual int) bool {
    return expected == actual
}

// Function equalIntNullable checks if two pointers to int values are equal.
func equalIntNullable(expected *int, actual *int) bool {
    if expected != nil && actual != nil {
        return equalInt(*expected, *actual)
    }
    if expected != nil || actual != nil {
        return false
    }
    return true
}

// Function equalInt64 checks if two int64 values are equal.
func equalInt64(expected int64, actual int64) bool {
    return expected == actual
}

// Function equalInt64Nullable checks if two pointers to int64 values are equal.
func equalInt64Nullable(expected *int64, actual *int64) bool {
    if expected != nil && actual != nil {
        return equalInt64(*expected, *actual)
    }
    if expected != nil || actual != nil {
        return false
    }
    return true
}

// Function equalFloat64 checks if two float64 values are equal.
func equalFloat64(expected float64, actual float64) bool {
    return expected == actual
}

// Function equalBool checks if two boolean values are equal.
func equalBool(expected bool, actual bool) bool {
    return expected == actual
}

// Function displayAssertionError displays assertion error details.
func displayAssertionError(expected interface{}, actual interface{}) {
    separator := makeText("-", 120)
    fmt.Printf("\n\n%s\n>     ERROR: assertion error\n>  Expected: %+v\n>    Actual: %+v\n%s\n\n",
        separator,
        expected,
        actual,
        separator)
    brexit()
}
