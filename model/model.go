package model

import (
	d "github.com/wisbery/oxyde/doc"
	"sort"
	"strings"
)

var (
	HttpMethodOrder = map[string]int{"POST": 1, "PUT": 2, "GET": 3, "DELETE": 4}
)

type Model struct {
	Groups        []Group              // Groups of endpoints.
	Endpoints     []Endpoint           // List of all endpoints in model.
	EndpointsById map[string]*Endpoint // Pointers to endpoints by endpoint identifier.
	RoleNames     []string             // List of tested role names for endpoints.
}

func CreateModel(dc *d.Context) *Model {
	// create model structure
	model := Model{
		Groups:        make([]Group, 0),
		Endpoints:     make([]Endpoint, 0),
		EndpointsById: make(map[string]*Endpoint),
		RoleNames:     dc.GetRoleNames()}
	// create all preview endpoints
	for _, docEndpoint := range dc.GetEndpoints() {
		endpoint := Endpoint{
			Id:           docEndpoint.Id,
			MethodUp:     strings.ToUpper(docEndpoint.Method),
			MethodLo:     strings.ToLower(docEndpoint.Method),
			UrlRoot:      docEndpoint.UrlRoot,
			UrlPath:      docEndpoint.UrlPath,
			Tags:         append(make([]string, 0), docEndpoint.Tags...),
			Summary:      docEndpoint.Summary,
			Parameters:   prepareFields(docEndpoint.Parameters),
			RequestBody:  prepareFields(docEndpoint.RequestBody),
			ResponseBody: prepareFields(docEndpoint.ResponseBody),
			Examples:     prepareExamples(docEndpoint.Examples),
			Access:       model.GetAccess(dc, docEndpoint.Method, docEndpoint.UrlPath)}
		model.Endpoints = append(model.Endpoints, endpoint)
	}
	// prepare endpoint mapping by identifiers
	for i, endpoint := range model.Endpoints {
		model.EndpointsById[endpoint.Id] = &model.Endpoints[i]
	}
	// create groups of endpoints
	model.createGroups()
	return &model
}

func (m *Model) FindEndpointById(id string) *Endpoint {
	if endpoint, ok := m.EndpointsById[id]; ok {
		return endpoint
	}
	return nil
}

func (m *Model) createGroups() {
	tags := make(map[string][]string)
	for _, docEndpoint := range m.Endpoints {
		if docEndpoint.Tags != nil {
			for _, tag := range docEndpoint.Tags {
				if ids, ok := tags[tag]; ok {
					tags[tag] = append(ids, docEndpoint.Id)
				} else {
					tags[tag] = append(make([]string, 0), docEndpoint.Id)
				}
			}
		}
	}
	groupNames := make([]string, 0)
	for tag := range tags {
		groupNames = append(groupNames, tag)
	}
	sort.Strings(groupNames)
	m.Groups = make([]Group, 0)
	for _, groupName := range groupNames {
		group := CreateGroup(m, groupName, tags[groupName])
		m.Groups = append(m.Groups, group)
	}
}

func (m *Model) GetAccess(dc *d.Context, method string, path string) []string {
	access := make([]string, len(m.RoleNames))
	for i, roleName := range m.RoleNames {
		switch dc.GetAccess(method, path, roleName) {
		case d.AccessGranted:
			access[i] = "YES"
		case d.AccessDenied:
			access[i] = "NO"
		case d.AccessError:
			access[i] = "ERR"
		case d.AccessUnknown:
			access[i] = "?"
		default:
			access[i] = "-"
		}
	}
	return access
}

type Group struct {
	Model     *Model      // Mode the group belong to.
	Name      string      // Group name.
	Endpoints []*Endpoint // List of endpoint identifiers in group.
}

func CreateGroup(model *Model, name string, ids []string) Group {
	group := Group{
		Model:     model,
		Name:      strings.ToUpper(name),
		Endpoints: make([]*Endpoint, 0)}
	sort.SliceStable(ids, func(i1, i2 int) bool {
		e1 := model.EndpointsById[ids[i1]]
		e2 := model.EndpointsById[ids[i2]]
		return compareEndpoints(e1, e2)
	})
	for _, id := range ids {
		group.Endpoints = append(group.Endpoints, model.EndpointsById[id])
	}
	return group
}

type Endpoint struct {
	Id           string    // Unique endpoint identifier.
	MethodUp     string    // HTTP method name in uppercase, like GET, POST, PUT or DELETE.
	MethodLo     string    // HTTP method name in lowercase, like get, post, put or delete.
	UrlRoot      string    // Root part of request URL.
	UrlPath      string    // Request path after root part.
	Tags         []string  // List of tags for endpoint.
	Summary      string    // Summary describing endpoint.
	Parameters   []Field   // List of parameter fields.
	RequestBody  []Field   // List of request body fields.
	ResponseBody []Field   // List of response body fields.
	Examples     []Example // List of examples.
	Access       []string  // List of access rights for roles.
}

type Field struct {
	Name        string // Name of the field.
	Type        string // Type of the field.
	Mandatory   string // Flag indicating if field is mandatory.
	MandatoryLo string // Flag indicating if field is mandatory in lowercase.
	Description string // Description of the field.
}

type Example struct {
	Summary      string // Example summary.
	Description  string // Example detailed description.
	Method       string // HTTP method name.
	MethodLo     string // HTTP method name in lowercase.
	Uri          string // Request URI.
	StatusCode   int    // HTTP status code.
	RequestBody  string // Request body as JSON string.
	ResponseBody string // Response body as JSON string.
}

func compareEndpoints(e1, e2 *Endpoint) bool {
	if i1, ok1 := HttpMethodOrder[e1.MethodUp]; ok1 {
		if i2, ok2 := HttpMethodOrder[e2.MethodUp]; ok2 {
			if i1 < i2 {
				return true
			} else if i1 == i2 {
				return len(e1.UrlPath) < len(e2.UrlPath)
			}
		}
	}
	return false
}

func prepareFields(docFields []d.Field) []Field {
	if docFields == nil {
		return nil
	}
	return traverseFields(docFields, 0)
}

func traverseFields(docFields []d.Field, level int) []Field {
	previewFields := make([]Field, 0)
	for _, docField := range docFields {
		mandatory := prepareMandatoryString(docField.Mandatory)
		previewField := Field{
			Name:        prepareFieldNameString(docField.JsonName, level),
			Type:        docField.JsonType,
			Mandatory:   mandatory,
			MandatoryLo: strings.ToLower(mandatory),
			Description: docField.Description}
		previewFields = append(previewFields, previewField)
		if docField.Children != nil {
			previewFields = append(previewFields, traverseFields(docField.Children, level+1)...)
		}
	}
	return previewFields
}

func prepareFieldNameString(name string, level int) string {
	indentString := "&nbsp;&nbsp;&nbsp;&nbsp;"
	indent := ""
	for i := 0; i < level; i++ {
		indent = indent + indentString
	}
	return indent + name
}

func prepareMandatoryString(mandatory bool) string {
	if mandatory {
		return "Yes"
	} else {
		return "No"
	}
}

func prepareExamples(docExamples []d.Example) []Example {
	examples := make([]Example, 0)
	for _, docExample := range docExamples {
		example := Example{
			Summary:      docExample.Summary,
			Description:  docExample.Description,
			Method:       docExample.Method,
			MethodLo:     strings.ToLower(docExample.Method),
			Uri:          docExample.Uri,
			StatusCode:   docExample.StatusCode,
			RequestBody:  docExample.RequestBody,
			ResponseBody: docExample.ResponseBody}
		examples = append(examples, example)
	}
	// sort examples by status code in ascending order
	sort.Slice(examples, func(i, j int) bool {
		return examples[i].StatusCode < examples[j].StatusCode
	})
	return examples
}
