package assert

import (
	"fmt"
	"github.com/wisbery/oxyde/common"
	"reflect"
)

func Nil(actual interface{}) {
	if !common.NilValue(actual) {
		displayAssertionError(nil, actual)
	}
}

func NotNil(actual interface{}) {
	if reflect.ValueOf(actual).IsNil() {
		displayAssertionError("not nil", actual)
	}
}

func NilError(e error) {
	if e != nil {
		displayAssertionError(nil, e)
	}
}

func NilString(actual *string) {
	if actual != nil {
		displayAssertionError(nil, actual)
	}
}

func NotNilString(actual *string) {
	if actual == nil {
		displayAssertionError("not nil", actual)
	}
}

func True(actual bool) {
	if !actual {
		displayAssertionError(true, actual)
	}
}

func False(actual bool) {
	if actual {
		displayAssertionError(false, actual)
	}
}

func NotNilId(actual *string) {
	NotNilString(actual)
	EqualInt(36, len(*actual))
}

func EqualString(expected string, actual string) {
	if !equalString(expected, actual) {
		displayAssertionError(expected, actual)
	}
}

func EqualStringNullable(expected *string, actual *string) {
	if !equalStringNullable(expected, actual) {
		if expected != nil && actual != nil {
			displayAssertionError(*expected, *actual)
		}
		displayAssertionError(expected, actual)
	}
}

func NilInt(actual *int) {
	if actual != nil {
		displayAssertionError(nil, actual)
	}
}

func EqualInt(expected int, actual int) {
	if !equalInt(expected, actual) {
		displayAssertionError(expected, actual)
	}
}

func EqualIntNullable(expected *int, actual *int) {
	if !equalIntNullable(expected, actual) {
		if expected != nil && actual != nil {
			displayAssertionError(*expected, *actual)
		}
		displayAssertionError(expected, actual)
	}
}

func EqualFloat64(expected float64, actual float64) {
	if !equalFloat64(expected, actual) {
		displayAssertionError(expected, actual)
	}
}

func EqualBool(expected bool, actual bool) {
	if !equalBool(expected, actual) {
		displayAssertionError(expected, actual)
	}
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
	separator := common.MakeString('-', 120)
	fmt.Printf("\n\n%s\n>     ERROR: assertion error\n>  Expected: %+v\n>    Actual: %+v\n%s\n\n",
		separator,
		expected,
		actual,
		separator)
	common.BrExit()
}
