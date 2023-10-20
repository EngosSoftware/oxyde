package doc

import (
	"errors"
	"fmt"
	"github.com/wisbery/oxyde/common"
	"reflect"
	"strings"
)

const (
	CollectNone = iota
	CollectDescription
	CollectExamples
	CollectAll
)

const (
	AccessGranted = iota
	AccessDenied
	AccessUnknown
	AccessError
)

type RoleKey struct {
	method   string
	path     string
	roleName string
}

type Context struct {
	mode               int             // Documentation collecting data mode.
	exampleSummary     string          // Summary for the next collected example.
	exampleDescription string          // Detailed description for the next collected example.
	roleName           string          // Role name of the principal for the next endpoint example.
	endpoint           *Endpoint       // Currently documented endpoint data.
	endpoints          []Endpoint      // List of documented endpoints.
	roleNames          []string        // Role names in order they should be displayed.
	roles              map[RoleKey]int // Map of access roles tested for endpoints.
}

func CreateDocContext() *Context {
	return &Context{
		mode:      CollectNone,
		endpoint:  nil,
		endpoints: make([]Endpoint, 0),
		roles:     make(map[RoleKey]int)}
}

func (dc *Context) ClearDocumentation() {
	dc.mode = CollectNone
	dc.endpoint = nil
	dc.endpoints = make([]Endpoint, 0)
	dc.roles = make(map[RoleKey]int)
}

func (dc *Context) PublishDocumentation() {
	for _, endpoint := range dc.endpoints {
		PrintEndpoint(endpoint)
	}
}

func (dc *Context) NewEndpointDocumentation(id string, tag string, summary string) {
	if id == "" {
		id = common.GenerateId()
	}
	dc.endpoint = &Endpoint{
		Id:      id,
		Summary: summary}
	dc.endpoint.AddTag(tag)
}

func (dc *Context) GetEndpoint() *Endpoint {
	return dc.endpoint
}

func (dc *Context) CollectDescription() {
	dc.mode = CollectDescription
}

func (dc *Context) CollectDescriptionMode() bool {
	return dc.mode == CollectDescription || dc.mode == CollectAll
}

func (dc *Context) CollectExamples(exampleSummary string, exampleDescription string) {
	exampleSummary = strings.TrimSpace(exampleSummary)
	exampleDescription = strings.TrimSpace(exampleDescription)
	if exampleSummary != "" || exampleDescription != "" {
		dc.mode = CollectExamples
		dc.exampleSummary = exampleSummary
		dc.exampleDescription = exampleDescription
	}
}

func (dc *Context) CollectRole(roleName string) {
	dc.roleName = roleName
}

func (dc *Context) CollectExamplesMode() bool {
	return dc.mode == CollectExamples || dc.mode == CollectAll
}

func (dc *Context) CollectAll(exampleSummary string) {
	dc.mode = CollectAll
	dc.exampleSummary = exampleSummary
}

func (dc *Context) StopCollecting() {
	dc.mode = CollectNone
	dc.exampleSummary = ""
	dc.roleName = ""
}

func (dc *Context) SetRolesOrder(roleOrder []string) {
	dc.roleNames = roleOrder
}

func (dc *Context) SaveRole(method string, path string, status int) {
	if dc.roleName != "" {
		key := RoleKey{method: method, path: path, roleName: dc.roleName}
		switch status {
		case 200:
			dc.roles[key] = AccessGranted
		case 401:
			dc.roles[key] = AccessDenied
		default:
			dc.roles[key] = AccessError
		}
	}
}

func (dc *Context) GetRoleNames() []string {
	return dc.roleNames
}

func (dc *Context) GetAccess(method string, path string, roleName string) int {
	key := RoleKey{
		method:   method,
		path:     path,
		roleName: roleName}
	if access, ok := dc.roles[key]; ok {
		return access
	} else {
		return AccessUnknown
	}
}

func (dc *Context) SaveEndpointDocumentation() {
	if dc.endpoint != nil {
		dc.endpoints = append(dc.endpoints, *dc.endpoint)
	}
}

func (dc *Context) GetEndpoints() []Endpoint {
	return dc.endpoints
}

func (dc *Context) GetExampleSummary() string {
	return dc.exampleSummary
}

func (dc *Context) GetExampleDescription() string {
	return dc.exampleDescription
}

type Endpoint struct {
	Id           string    // Unique endpoint identifier.
	Tags         []string  // List of tags of endpoint.
	Method       string    // HTTP method name, like GET, POST, PUT or DELETE.
	UrlRoot      string    // Request URL root.
	UrlPath      string    // Request URL path after root.
	Summary      string    // Summary text describing endpoint.
	Parameters   []Field   // Description of request parameters.
	RequestBody  []Field   // Description of request body.
	ResponseBody []Field   // Description of results.
	Examples     []Example // Description of usage examples.
}

func (e *Endpoint) AddTag(tag string) {
	if e.Tags == nil {
		e.Tags = make([]string, 0)
	}
	e.Tags = append(e.Tags, tag)
}

type Field struct {
	JsonName    string  // Name of the field in JSON.
	JsonType    string  // Type of the field in JSON.
	Mandatory   bool    // Flag indicating if field is mandatory in JSON.
	Description string  // Description of the field.
	Children    []Field // List of child fields (may be empty).
}

func CreateField(typ reflect.Type, structField reflect.StructField) Field {
	jsonType := jsonType(typ)
	jsonName := structField.Tag.Get(common.JsonTagName)
	apiTagContent := structField.Tag.Get(common.ApiTagName)
	mandatory := !strings.HasPrefix(apiTagContent, common.OptionalPrefix)
	apiTagContent = strings.TrimPrefix(apiTagContent, common.OptionalPrefix)
	return Field{
		JsonName:    jsonName,
		JsonType:    jsonType,
		Mandatory:   mandatory,
		Description: apiTagContent,
		Children:    make([]Field, 0)}
}

type Example struct {
	Summary      string // Example summary.
	Description  string // Detailed example description.
	Method       string // HTTP method name.
	Uri          string // Request URI.
	StatusCode   int    // HTTP status code.
	RequestBody  string // Request body as JSON string.
	ResponseBody string // Response body as JSON string.
}

func ParseObject(o interface{}) []Field {
	typ := reflect.TypeOf(o)
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
			field := CreateField(childType, childField)
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

func PrintEndpoint(endpoint Endpoint) {
	fmt.Printf("\n\n%s\n%s\n\n", endpoint.Method, endpoint.UrlRoot)
	fmt.Println("Parameters:")
	fmt.Println(common.MakeString('-', 120))
	PrintFields(endpoint.RequestBody, "  ", 0)
	fmt.Println()
	fmt.Println(common.MakeString('-', 120))
	fmt.Printf("\n")
	fmt.Println("ResponseBody:")
	fmt.Println(common.MakeString('-', 120))
	PrintFields(endpoint.ResponseBody, "  ", 0)
	fmt.Println()
	fmt.Println(common.MakeString('-', 120))
	fmt.Printf("\n")
	for _, usage := range endpoint.Examples {
		PrintExample(usage)
	}
}

func PrintFields(fields []Field, indent string, level int) {
	for i, field := range fields {
		nameText := indent + field.JsonName
		nameText = nameText + strings.Repeat(" ", 30-len(nameText))
		typeText := field.JsonType
		typeText = typeText + strings.Repeat(" ", 20-len(typeText))
		mandatoryText := "N"
		if field.Mandatory {
			mandatoryText = "Y"
		}
		if i > 0 || level > 0 {
			fmt.Println()
		}
		fmt.Printf("|%s | %s | %s | %s", nameText, typeText, mandatoryText, field.Description)
		if field.JsonType == "object" || field.JsonType == "array" {
			PrintFields(field.Children, indent+indent, level+1)
		}
	}
}

func PrintExample(usage Example) {
	fmt.Printf("\nExample:\n")
	fmt.Printf("%d %s %s\n", usage.StatusCode, usage.Method, usage.Uri)
	fmt.Printf("Parameters:\n%s\n", usage.RequestBody)
	fmt.Printf("ResponseBody:\n%s\n", usage.ResponseBody)
}
