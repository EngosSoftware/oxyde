package oxyde

import (
	"testing"
)

func TestSingleParameterInjection(t *testing.T) {
	path := "/users/{userId}"
	params := struct {
		UserId string `json:"userId"`
	}{
		UserId: "c4f63c4f-e66b-4cd4-b0a7-f7a5e2bc6edd",
	}
	requestPath, err := prepareRequestPath(path, "v1", params)
	if requestPath != "/users/c4f63c4f-e66b-4cd4-b0a7-f7a5e2bc6edd" || err != nil {
		t.Error("single parameter not injected")
	}
}

func TestMultipleParameterInjection(t *testing.T) {
	path := "/users/{userId}/{userName}"
	params := struct {
		UserId   string `json:"userId"`
		UserName string `json:"userName"`
	}{
		UserId:   "b494fd53-10c8-43bf-b585-334a2cac0995",
		UserName: "John",
	}
	requestPath, err := prepareRequestPath(path, "v1", params)
	if requestPath != "/users/b494fd53-10c8-43bf-b585-334a2cac0995/John" || err != nil {
		t.Error("multiple parameters not injected")
	}
}

func TestSingleRepeatedParameterInjection(t *testing.T) {
	path := "/users/{userId}/{userId}"
	params := struct {
		UserId string `json:"userId"`
	}{
		UserId: "ee90021b-15ce-4d3e-bd2c-6ce023503fff",
	}
	requestPath, err := prepareRequestPath(path, "v1", params)
	if requestPath != "/users/ee90021b-15ce-4d3e-bd2c-6ce023503fff/ee90021b-15ce-4d3e-bd2c-6ce023503fff" || err != nil {
		t.Error("single repeated parameters not injected")
	}
}

func TestSingleParameterAppend(t *testing.T) {
	path := "/users"
	params := struct {
		UserId string `json:"userId"`
	}{
		UserId: "2b4ca889-7ed0-41ca-b832-222a9ecaf183",
	}
	requestPath, err := prepareRequestPath(path, "v1", params)
	if requestPath != "/users?userId=2b4ca889-7ed0-41ca-b832-222a9ecaf183" || err != nil {
		t.Error("single parameters not appended")
	}
}

func TestMultipleParameterAppend(t *testing.T) {
	path := "/users"
	params := struct {
		UserId   string `json:"userId"`
		UserName string `json:"userName"`
		Age      int    `json:"age"`
	}{
		UserId:   "2b4ca889-7ed0-41ca-b832-222a9ecaf183",
		UserName: "Matthew",
		Age:      32,
	}
	requestPath, err := prepareRequestPath(path, "v1", params)
	if requestPath != "/users?userId=2b4ca889-7ed0-41ca-b832-222a9ecaf183&userName=Matthew&age=32" || err != nil {
		t.Error("multiple parameters not appended")
	}
}

func TestEmptyParameterInjection(t *testing.T) {
	path := "/users/empty{userId}"
	params := struct {
		UserId string `json:"userId"`
	}{
		UserId: "",
	}
	requestPath, err := prepareRequestPath(path, "v1", params)
	if requestPath != "/users/empty" || err != nil {
		t.Error("empty parameter not injected")
	}
}

func TestEmptyParameterAppend(t *testing.T) {
	path := "/users"
	params := struct {
		UserId string `json:"userId"`
	}{
		UserId: "",
	}
	requestPath, err := prepareRequestPath(path, "v1", params)
	if requestPath != "/users?userId=" || err != nil {
		t.Error("empty parameter not appended")
	}
}
