package common

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"reflect"
	"regexp"
	"runtime"
)

const (
	ApiTagName     = "api"  // Name of the tag in which documentation details are stored.
	JsonTagName    = "json" // Name of the tag in which JSON details are stored.
	OptionalPrefix = "?"    // Prefix used to mark th field as optional.
)

// Function MakeString creates a string of length 'len' containing the same character 'ch'.
func MakeString(ch byte, len int) string {
	b := make([]byte, len)
	for i := 0; i < len; i++ {
		b[i] = ch
	}
	return string(b)
}

// Function NilValue checks if the value specified as parameter is nil.
func NilValue(value interface{}) bool {
	return value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil())
}

// Function PrettyPrint takes JSON string as an argument
// and returns the same JSON but pretty-printed.
func PrettyPrint(in []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, in, "", "  ")
	if err != nil {
		return string(in)
	}
	return out.String()
}

// Function PanicOnError panics when the error passed as argument is not nil.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// Function GenerateId generated unique identifier in form of UUID string.
func GenerateId() string {
	id, err := uuid.NewRandom()
	PanicOnError(err)
	return id.String()
}

// Function TypeOfValue returns the type of value specified as parameter.
// If specified value is a pointer, then is first dereferenced.
func TypeOfValue(value interface{}) reflect.Type {
	t := reflect.TypeOf(value)
	if t.Kind().String() == "ptr" {
		v := reflect.ValueOf(value)
		t = reflect.Indirect(v).Type()
	}
	return t
}

// Function ValueOfValue returns the value of specified parameter.
// If specified value is a pointer, then is first dereferenced.
func ValueOfValue(value interface{}) reflect.Value {
	t := reflect.TypeOf(value)
	if t.Kind().String() == "ptr" {
		v := reflect.ValueOf(value)
		return reflect.Indirect(v)
	}
	return reflect.ValueOf(value)
}

// Function BrExit breaks the execution of test and displays stack trace.
// After breaking the execution flow, application returns exit code -1
// that can be utilized by test automation tools.
func BrExit() {
	fmt.Printf("Stack trace:\n------------\n")
	reDeepCalls := regexp.MustCompile(`(^goroutine[^:]*:$)|(^.*/oxyde/.*$)`)
	reFuncParams := regexp.MustCompile(`([a-zA-Z_0-9]+)\([^\)]+\)`)
	reFuncOffset := regexp.MustCompile(`\s+\+.*$`)
	b := make([]byte, 100000)
	runtime.Stack(b, false)
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	for scanner.Scan() {
		line := scanner.Text()
		if reDeepCalls.MatchString(line) {
			continue
		}
		line = reFuncParams.ReplaceAllString(line, "$1()")
		line = reFuncOffset.ReplaceAllString(line, "")
		fmt.Println(line)
	}
	os.Exit(-1)
}
