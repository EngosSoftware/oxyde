package oxyde

import (
	"fmt"
	"testing"
)

func TestHeadersNil(t *testing.T) {
	out := parseHeaders(nil)
	assertTestResultEmptyHeadersMap(t, out)
}

func TestHeadersNilPointerToStruct(t *testing.T) {
	type Headers struct {
		Authorization string
	}
	type PHeaders = *Headers
	var in PHeaders
	out := parseHeaders(in)
	assertTestResultEmptyHeadersMap(t, out)
}

func TestHeadersString(t *testing.T) {
	in := "Authorization"
	out := parseHeaders(in)
	assertTestResultEmptyHeadersMap(t, out)
}

func TestHeadersStruct(t *testing.T) {
	type Headers struct {
		H1 string `json:"h1"`
		H2 string `json:"h2"`
		H3 string `json:"h3"`
	}
	in := Headers{
		H1: "v1",
		H2: "v2",
		H3: "v3"}
	out := parseHeaders(in)
	assertTestResultHeaders3(t, out, in.H1, in.H2, in.H3)
}

func TestHeadersPointerToStruct(t *testing.T) {
	type Headers struct {
		H1 string `json:"h1"`
		H2 string `json:"h2"`
		H3 string `json:"h3"`
	}
	in := Headers{
		H1: "v1",
		H2: "v2",
		H3: "v3"}
	out := parseHeaders(&in)
	assertTestResultHeaders3(t, out, in.H1, in.H2, in.H3)
}

func TestHeadersStructWithPointers(t *testing.T) {
	type Headers struct {
		H1 *string `json:"h1"`
		H2 *string `json:"h2"`
		H3 *string `json:"h3"`
	}
	v1 := "v1"
	v2 := "v2"
	v3 := "v3"
	in := Headers{
		H1: &v1,
		H2: &v2,
		H3: &v3}
	out := parseHeaders(in)
	assertTestResultHeaders3(t, out, *in.H1, *in.H2, *in.H3)
}

func TestHeadersStructWithPointersAndNilValues(t *testing.T) {
	type Headers struct {
		H1 *string `json:"h1"`
		H2 *string `json:"h2"`
		H3 *string `json:"h3"`
	}
	v1 := "v1"
	v3 := "v3"
	in := Headers{
		H1: &v1,
		H2: nil,
		H3: &v3}
	out := parseHeaders(in)
	if len(out) != 2 {
		t.Error(fmt.Sprintf("expected exactly 2 headers, but %d found", len(out)))
	}
	key := "h1"
	if value, ok := out[key]; !ok || value != *in.H1 {
		t.Error(fmt.Sprintf("expected key '%s' with value '%s' not found", key, *in.H1))
	}
	key = "h2"
	if value, ok := out[key]; ok {
		t.Error(fmt.Sprintf("expected no key '%s', but value '%s' found", key, value))
	}
	key = "h3"
	if value, ok := out[key]; !ok || value != *in.H3 {
		t.Error(fmt.Sprintf("expected key '%s' with value '%s' not found", key, *in.H3))
	}
}

func assertTestResultEmptyHeadersMap(t *testing.T, actual headers) {
	if len(actual) != 0 {
		t.Error(fmt.Sprintf("expected empty map, found: %+v", actual))
	}
}

func assertTestResultHeaders3(t *testing.T, actual headers, h1, h2, h3 string) {
	if len(actual) != 3 {
		t.Error(fmt.Sprintf("expected 3 headers, found %d", len(actual)))
	}
	key := "h1"
	if value, ok := actual[key]; !ok || value != h1 {
		t.Error(fmt.Sprintf("expected key '%s' with value '%s', not found", key, h1))
	}
	key = "h2"
	if value, ok := actual[key]; !ok || value != h2 {
		t.Error(fmt.Sprintf("expected key '%s' with value '%s' not found", key, h2))
	}
	key = "h3"
	if value, ok := actual[key]; !ok || value != h3 {
		t.Error(fmt.Sprintf("expected key '%s' with value '%s' not found", key, h3))
	}
}
