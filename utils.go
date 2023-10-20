package oxyde

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
	"strings"
)

const (
	ApiTagName              = "api"          // Name of the tag in which documentation details are stored.
	JsonTagName             = "json"         // Name of the tag in which JSON details are stored.
	VersionPlaceholder      = "{apiVersion}" // Placeholder for API version number in request path.
	OptionalPrefix          = "?"            // Prefix used to mark the field as optional.
	TestSuitePrefix         = "ts_"          // Prefix used to name the function that runs the test suite.
	TestCasePrefix          = "tc_"          // Prefix used to name the function that runs the test case.
	TestDocumentationPrefix = "td_"          // Prefix used to name the function that documents the API.
)

var (
	reCamelBoundary = regexp.MustCompile("([a-z])([A-Z])")
	reFunctionName  = regexp.MustCompile(`\.([a-zA-Z_0-9]+)\(`)
)

// Function makeText creates a string that contains repeated text.
func makeText(text string, repeat int) string {
	var builder strings.Builder
	for i := 0; i < repeat; i++ {
		builder.WriteString(text)
	}
	return builder.String()
}

// Function nilValue checks if the value specified as parameter is nil.
func nilValue(value interface{}) bool {
	return value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil())
}

// Function prettyPrint takes JSON string as a parameter
// and returns the same string as pretty-printed JSON.
func prettyPrint(in []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, in, "", "  ")
	if err != nil {
		return string(in)
	}
	return out.String()
}

// Function panicOnError panics when the error passed as argument is not nil.
func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// Function generateId returns universally unique identifier.
func generateId() string {
	id, err := uuid.NewRandom()
	panicOnError(err)
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

// Function brexit stops the execution of the test and displays the stack trace.
// After breaking the execution flow, application returns exit code -1
// that can be utilized by test automation tools.
func brexit() {
	fmt.Printf("Stack trace:\n------------\n")
	reDeepCalls := regexp.MustCompile(`(^goroutine[^:]*:$)|(^.*oxyde.*$)`)
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

func Display(ctx *Context) {
	DisplayLevel(ctx, 3)
}

func Display2(ctx *Context) {
	DisplayLevel(ctx, 4)
}

func Ok() {
	fmt.Println("OK")
}

func DisplayLevel(ctx *Context, level int) {
	text := strings.TrimSpace(FunctionName(level))
	newline := "\n"
	if strings.HasPrefix(text, TestSuitePrefix) {
		text = " >> " + text
	} else if strings.HasPrefix(text, TestDocumentationPrefix) {
		text = "  > " + text
	} else if strings.HasPrefix(text, TestCasePrefix) {
		text = "    - " + text + " [" + ctx.UserName + "]"
		newline = ""
	} else {
		text = ">>> " + text
	}
	fmt.Printf("%-120s%-5s%s", text, ctx.Version, newline)
}

func FunctionName(level int) string {
	b := make([]byte, 8192)
	runtime.Stack(b, false)
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	index := 0
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if index == level*2+1 {
			break
		}
		index++
	}
	line = reFunctionName.FindString(line)
	line = reFunctionName.ReplaceAllString(line, "$1")
	line = reCamelBoundary.ReplaceAllString(line, `$1#$2`)
	line = strings.ToLower(strings.ReplaceAll(line, "#", "_"))
	return line
}
