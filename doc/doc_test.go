package doc

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type TestLoginParams struct {
	Login    string `json:"login" api:"Login"`
	Password string `json:"password" api:"Password"`
}

func TestSimpleJSONObject(t *testing.T) {

	fields := ParseObject(TestLoginParams{})
	if len(fields) != 2 {
		t.Error("expected two fields in object")
	}
	if fields[0].JsonName != "login" {
		t.Error("expected field with name 'login'")
	}
	if fields[1].JsonName != "password" {
		t.Error("expected field with name 'password'")
	}
}

func Traverse(o interface{}) {
	t := reflect.TypeOf(o)
	TraverseFields(t)
}

func TraverseFields(t reflect.Type) {
	switch t.Kind() {
	case reflect.Ptr:
		t := t.Elem()
		TraverseFields(t)
	case reflect.Bool:
		fmt.Printf("%s\n", "boolean")
	case reflect.String:
		fmt.Printf("%s\n", "string")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		fmt.Printf("%s\n", "number")
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			TraverseFields(t.Field(i).Type)
		}
	case reflect.Slice:
		TraverseFields(t.Elem())
	default:
		panic(errors.New(fmt.Sprintf("'%s' is not supported in static type analysis", t.Kind().String())))
	}
}

func TestA(t *testing.T) {
	fmt.Println()
	str := "test"
	fmt.Println("---a")
	Traverse(str)
	fmt.Println("---b")
	Traverse(&str)
	ps := &str
	fmt.Println("---c")
	Traverse(&ps)
	type Tk struct {
		Code int
	}
	k := Tk{Code: 120}
	type Tp struct {
		Info float64
	}
	// p := Tp{Info: 64.56}
	type Arr []struct {
		Bobo float64
		Kupo int64
	}

	s := struct {
		Name       string
		Address    *string
		Company    **string
		Details    Tk
		SubDetails *Tp
		Tags       Arr
	}{
		Name:       "John",
		Address:    nil,
		Company:    &ps,
		Details:    k,
		SubDetails: nil,
		Tags:       nil}
	fmt.Println("---d")
	ParseObject(s)
	// Traverse(&s)
	fmt.Println()
}

func TestSimpleTypes(t *testing.T) {
	type Data struct {
		Name    string   `json:"name" api:"Name."`
		Age     int      `json:"age" api:"Age."`
		Salary  int64    `json:"salary" api:"Salary."`
		Married bool     `json:"married" api:"Is married?"`
		Height  float64  `json:"height" api:"Height."`
		Tags    []string `json:"tags" api:"Tags."`
	}
	d := Data{}
	fields := ParseObject(d)
	PrintFields(fields, " ", 0)
	fmt.Println()
}

func TestStructures(t *testing.T) {
	type Address struct {
		Country string `json:"country" api:"Country name."`
		Street  string `json:"street" api:"Street name."`
	}
	type Child struct {
		Name string `json:"name" api:"Name of the child."`
		Age  int8   `json:"age" api:"Age of the child."`
	}
	type Data struct {
		Name     *string  `json:"name" api:"Name."`
		Age      int      `json:"age" api:"Age."`
		Salary   int64    `json:"salary" api:"Salary."`
		Married  *bool    `json:"married" api:"Is married?"`
		Height   float64  `json:"height" api:"Height."`
		Tags     []string `json:"tags" api:"Tags."`
		Children []Child  `json:"children" api:"Children."`
		Address  *Address `json:"address" api:"Address details."`
	}
	d := &Data{}
	fields := ParseObject(d)
	PrintFields(fields, "   ", 0)
	fmt.Println()
}
