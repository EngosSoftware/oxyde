package oxyde

import (
	"testing"
)

func TestEqualStrings(t *testing.T) {
	if !equalString("string", "string") {
		t.Error("strings are equal but test shows they are not")
	}
	if equalString("string", "string1") {
		t.Error("strings are not equal but test shows they are")
	}
	if equalString("Alfa", "Alfb") {
		t.Error("strings are not equal but test shows they are")
	}
	if equalString("string", "") {
		t.Error("strings are not equal but test shows they are")
	}
	if equalString("", "string") {
		t.Error("strings are not equal but test shows they are")
	}
}

func TestEqualStringsNullable(t *testing.T) {
	s1 := "string"
	s2 := "string"
	if !equalStringNullable(&s1, &s2) {
		t.Error("strings are equal but test shows they are not")
	}
	if !equalStringNullable(nil, nil) {
		t.Error("strings are equal but test shows they are not")
	}
	if equalStringNullable(&s1, nil) {
		t.Error("strings are not equal but test shows they are")
	}
	if equalStringNullable(nil, &s2) {
		t.Error("strings are not equal but test shows they are")
	}
	s1 = "Alfa"
	s2 = "Alfb"
	if equalStringNullable(&s1, &s2) {
		t.Error("strings are not equal but test shows they are")
	}
}
