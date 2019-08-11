/*
 * MIT License
 *
 * Copyright (c) 2017-2019 Dariusz Depta Engos Software
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package oxyde

import (
    "sort"
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
    Id           string           // Unique endpoint identifier.
    MethodUp     string           // HTTP method name in uppercase, like GET, POST, PUT or DELETE.
    MethodLo     string           // HTTP method name in lowercase, like get, post, put or delete.
    UrlRoot      string           // Root part of request URL.
    UrlPath      string           // Request path after root part.
    Tags         []string         // List of tags for endpoint.
    Summary      string           // Summary describing endpoint.
    Parameters   []PreviewField   // List of parameter fields.
    RequestBody  []PreviewField   // List of request body fields.
    ResponseBody []PreviewField   // List of response body fields.
    Examples     []PreviewExample // List of examples.
    Access       []string         // List of access rights for roles.
}

type PreviewField struct {
    Name        string // Name of the field.
    Type        string // Type of the field.
    Mandatory   string // Flag indicating if field is mandatory.
    MandatoryLo string // Flag indicating if field is mandatory in lowercase.
    Description string // Description of the field.
}

type PreviewExample struct {
    Summary      string // Usage example summary.
    Description  string // Usage example description.
    Method       string // HTTP method name.
    MethodLo     string // HTTP method name in lowercase.
    Uri          string // Full request URI.
    StatusCode   int    // HTTP status code.
    RequestBody  string // Request body as JSON string.
    ResponseBody string // Response body as JSON string.
}

func CreatePreviewModel(dc *DocumentContext) *PreviewModel {
    // create preview previewModel structure
    previewModel := PreviewModel{
        Groups:        make([]PreviewGroup, 0),
        Endpoints:     make([]PreviewEndpoint, 0),
        EndpointsById: make(map[string]*PreviewEndpoint),
        RoleNames:     dc.GetRoleNames()}
    // create all preview endpoints
    for _, docEndpoint := range dc.GetEndpoints() {
        previewEndpoint := PreviewEndpoint{
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
            Access:       previewModel.GetAccess(dc, docEndpoint.Method, docEndpoint.UrlPath)}
        previewModel.Endpoints = append(previewModel.Endpoints, previewEndpoint)
    }
    // prepare endpoint mapping by identifiers
    for i, endpoint := range previewModel.Endpoints {
        previewModel.EndpointsById[endpoint.Id] = &previewModel.Endpoints[i]
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
    m.Groups = make([]PreviewGroup, 0)
    for _, groupName := range groupNames {
        group := CreatePreviewGroup(m, groupName, tags[groupName])
        m.Groups = append(m.Groups, group)
    }
}

func (m *PreviewModel) GetAccess(dc *DocumentContext, method string, path string) []string {
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

func traverseFields(docFields []Field, level int) []PreviewField {
    previewFields := make([]PreviewField, 0)
    for _, docField := range docFields {
        mandatory := prepareMandatoryString(docField.Mandatory)
        previewField := PreviewField{
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

func prepareExamples(docExamples []Example) []PreviewExample {
    examples := make([]PreviewExample, 0)
    for _, docExample := range docExamples {
        example := PreviewExample{
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
