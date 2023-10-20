package oxyde

import (
	"sort"
	"strconv"
	"strings"
)

var (
	HttpMethodOrder = map[string]int{"POST": 1, "PUT": 2, "GET": 3, "DELETE": 4}
)

type PreviewModel struct {
	Groups        []PreviewGroup              // Groups of endpoints.
	Endpoints     []PreviewEndpoint           // All endpoints in preview model.
	EndpointsById map[string]*PreviewEndpoint // Endpoints indexed by identifier.
	RoleNames     []string                    // Names of access roles for endpoints.
}

type PreviewGroup struct {
	Model     *PreviewModel      // Model the group belongs to.
	Name      string             // Group name.
	Endpoints []*PreviewEndpoint // Endpoints in group.
}

type PreviewEndpoint struct {
	Id           string         // Unique endpoint identifier.
	MethodUp     string         // HTTP method name in upper-case, like GET, POST, PUT or DELETE.
	MethodLo     string         // HTTP method name in lower-case, like get, post, put or delete.
	UrlRoot      string         // Root part of request URL.
	UrlPath      string         // Request path after root part.
	Group        string         // Name of the group the endpoint belongs to.
	Summary      string         // Summary describing endpoint.
	Description  string         // Detailed endpoint description.
	Headers      []PreviewField // List of header fields.
	Parameters   []PreviewField // List of parameter fields.
	RequestBody  []PreviewField // List of request body fields.
	ResponseBody []PreviewField // List of response body fields.
	Usages       []PreviewUsage // List of examples.
	Access       []string       // List of access rights for roles.
}

type PreviewField struct {
	Name        string // Name of the field.
	Type        string // Type of the field.
	Mandatory   string // Flag indicating if field is mandatory.
	MandatoryLo string // Flag indicating if field is mandatory in lower-case.
	Description string // Description of the field.
}

type PreviewHeader struct {
	Name  string // HTTP header name.
	Value string // HTTP header value.
}

type PreviewUsage struct {
	Summary      string          // Usage example summary.
	Description  string          // Usage example description.
	MethodUp     string          // HTTP method name in upper-case.
	MethodLo     string          // HTTP method name in lower-case.
	Url          string          // Full request URL.
	Headers      []PreviewHeader // Usage headers.
	RequestBody  string          // Request body as JSON string.
	ResponseBody string          // Response body as JSON string.
	StatusCode   int             // HTTP status code.
}

func CreatePreviewModel(dc *DocContext) *PreviewModel {
	// create preview previewModel structure
	previewModel := PreviewModel{
		Groups:        make([]PreviewGroup, 0),
		Endpoints:     make([]PreviewEndpoint, 0),
		EndpointsById: make(map[string]*PreviewEndpoint),
		RoleNames:     dc.GetRoleNames()}
	// create all preview endpoints
	for _, endpoint := range dc.GetEndpoints() {
		previewEndpoint := PreviewEndpoint{
			Id:           endpoint.Id,
			MethodUp:     strings.ToUpper(endpoint.Method),
			MethodLo:     strings.ToLower(endpoint.Method),
			UrlRoot:      endpoint.RootPath,
			UrlPath:      endpoint.RequestPath,
			Group:        endpoint.Group,
			Summary:      endpoint.Summary,
			Description:  endpoint.Description,
			Headers:      prepareFields(endpoint.Headers),
			Parameters:   prepareFields(endpoint.Parameters),
			RequestBody:  prepareFields(endpoint.RequestBody),
			ResponseBody: prepareFields(endpoint.ResponseBody),
			Usages:       preparePreviewUsages(endpoint.Usages),
			Access:       previewModel.GetAccess(dc, endpoint.Method, endpoint.RequestPath)}
		previewModel.Endpoints = append(previewModel.Endpoints, previewEndpoint)
	}
	// prepare previewEndpoint mapping by identifiers
	for i, previewEndpoint := range previewModel.Endpoints {
		previewModel.EndpointsById[previewEndpoint.Id] = &previewModel.Endpoints[i]
	}
	// create groups of endpoints
	previewModel.createGroups()
	return &previewModel
}

func CreatePreviewGroup(previewModel *PreviewModel, name string, ids []string) PreviewGroup {
	previewGroup := PreviewGroup{
		Model:     previewModel,
		Name:      strings.ToUpper(name),
		Endpoints: make([]*PreviewEndpoint, 0)}
	sort.SliceStable(ids, func(i1, i2 int) bool {
		e1 := previewModel.EndpointsById[ids[i1]]
		e2 := previewModel.EndpointsById[ids[i2]]
		return compareEndpoints(e1, e2)
	})
	for _, id := range ids {
		previewGroup.Endpoints = append(previewGroup.Endpoints, previewModel.EndpointsById[id])
	}
	return previewGroup
}

func compareEndpoints(e1, e2 *PreviewEndpoint) bool {
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

func (m *PreviewModel) FindEndpointById(id string) *PreviewEndpoint {
	if endpoint, ok := m.EndpointsById[id]; ok {
		return endpoint
	}
	return nil
}

func (m *PreviewModel) createGroups() {
	groups := make(map[string][]string)
	for _, endpoint := range m.Endpoints {
		groupName := endpoint.Group
		if ids, ok := groups[groupName]; ok {
			groups[groupName] = append(ids, endpoint.Id)
		} else {
			groups[groupName] = append(make([]string, 0), endpoint.Id)
		}
	}
	groupNames := make([]string, 0)
	for groupName := range groups {
		groupNames = append(groupNames, groupName)
	}
	sort.Strings(groupNames)
	m.Groups = make([]PreviewGroup, 0)
	for _, groupName := range groupNames {
		group := CreatePreviewGroup(m, groupName, groups[groupName])
		m.Groups = append(m.Groups, group)
	}
}

func (m *PreviewModel) GetAccess(dc *DocContext, method string, path string) []string {
	access := make([]string, len(m.RoleNames))
	for i, roleName := range m.RoleNames {
		switch dc.GetAccess(method, path, roleName) {
		case AccessGranted:
			access[i] = "YES"
		case AccessDenied:
			access[i] = "NO"
		case AccessError:
			access[i] = "ERR"
		case AccessUnknown:
			access[i] = "?"
		default:
			access[i] = "-"
		}
	}
	return access
}

func prepareFields(docFields []Field) []PreviewField {
	if docFields == nil {
		return nil
	}
	return traverseFields(docFields, 0)
}

func traverseFields(fields []Field, level int) []PreviewField {
	previewFields := make([]PreviewField, 0)
	for _, field := range fields {
		mandatory := prepareMandatoryString(field.Mandatory)
		previewField := PreviewField{
			Name:        prepareFieldNameString(field.JsonName, level),
			Type:        field.JsonType,
			Mandatory:   mandatory,
			MandatoryLo: strings.ToLower(mandatory),
			Description: field.Description}
		previewFields = append(previewFields, previewField)
		if field.Children != nil {
			previewFields = append(previewFields, traverseFields(field.Children, level+1)...)
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

func preparePreviewUsages(usages []Usage) []PreviewUsage {
	previewUsages := make([]PreviewUsage, 0)
	for _, usage := range usages {
		previewExample := PreviewUsage{
			Summary:      usage.Summary,
			Description:  usage.Description,
			MethodUp:     strings.ToUpper(usage.Method),
			MethodLo:     strings.ToLower(usage.Method),
			Headers:      preparePreviewHeaders(usage.Headers),
			Url:          usage.Url,
			StatusCode:   usage.StatusCode,
			RequestBody:  usage.RequestBody,
			ResponseBody: usage.ResponseBody}
		previewUsages = append(previewUsages, previewExample)
	}
	// sort preview usages by status code in ascending order
	sort.Slice(previewUsages, func(i, j int) bool {
		return previewUsages[i].StatusCode < previewUsages[j].StatusCode
	})
	return previewUsages
}

func preparePreviewHeaders(headers headers) []PreviewHeader {
	const maxLen = 50
	previewHeaders := make([]PreviewHeader, 0)
	for name, value := range headers {
		length := len(value)
		if length > maxLen {
			value = value[:maxLen] + "[...](" + strconv.Itoa(length) + ")"
		}
		previewHeaders = append(previewHeaders, PreviewHeader{Name: name, Value: value})
	}
	return previewHeaders
}
