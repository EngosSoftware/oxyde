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
    "strings"
)

const (
    CollectNothing                     = iota // No documentation data will be collected.
    CollectEndpointDescription                // Endpoint description will be collected.
    CollectEndpointUsageExample               // An example of endpoint usage will be collected.
    CollectEndpointDescriptionAndUsage        // Endpoint description and usage example will be collected.
)

type DocContext struct {
    mode             int        // Documentation collecting data mode.
    usageSummary     string     // Summary for the next collected usage example.
    usageDescription string     // Detailed description for the next collected usage example.
    roleName         string     // Role name of the principal for the next endpoint usage example.
    endpoint         *Endpoint  // Pointer to currently documented endpoint.
    endpoints        []Endpoint // List of documented endpoints with usage examples.
    roleNames        []string   // Role names in order they should be displayed.
    roles            roles      // Map of verified access roles for all endpoints.
}

func CreateDocContext() *DocContext {
    return &DocContext{
        mode:             CollectNothing,
        usageSummary:     "",
        usageDescription: "",
        roleName:         "",
        endpoint:         nil,
        endpoints:        make([]Endpoint, 0),
        roles:            make(roles)}
}

func (dc *DocContext) Clear() {
    dc.mode = CollectNothing
    dc.endpoint = nil
    dc.endpoints = make([]Endpoint, 0)
    dc.roles = make(map[roleKey]int)
}

func (dc *DocContext) NewEndpoint(tag string, summary string, description string) {
    dc.endpoint = &Endpoint{
        Id:          generateId(),
        Summary:     summary,
        Description: description}
    dc.endpoint.AddTag(tag)
}

func (dc *DocContext) GetEndpoint() *Endpoint {
    return dc.endpoint
}

func (dc *DocContext) CollectDescription() {
    dc.mode = CollectEndpointDescription
}

func (dc *DocContext) CollectDescriptionMode() bool {
    return dc.mode == CollectEndpointDescription || dc.mode == CollectEndpointDescriptionAndUsage
}

func (dc *DocContext) CollectUsage(summary string, description string) {
    summary = strings.TrimSpace(summary)
    description = strings.TrimSpace(description)
    if summary != "" || description != "" {
        dc.mode = CollectEndpointUsageExample
        dc.usageSummary = summary
        dc.usageDescription = description
    }
}

func (dc *DocContext) CollectRole(roleName string) {
    dc.roleName = roleName
}

func (dc *DocContext) CollectExamplesMode() bool {
    return dc.mode == CollectEndpointUsageExample || dc.mode == CollectEndpointDescriptionAndUsage
}

func (dc *DocContext) CollectAll(exampleSummary string) {
    dc.mode = CollectEndpointDescriptionAndUsage
    dc.usageSummary = exampleSummary
}

func (dc *DocContext) StopCollecting() {
    dc.mode = CollectNothing
    dc.usageSummary = ""
    dc.roleName = ""
}

func (dc *DocContext) SetRolesOrder(roleOrder []string) {
    dc.roleNames = roleOrder
}

func (dc *DocContext) SaveRole(method string, path string, status int) {
    if dc.roleName != "" {
        key := roleKey{method: method, path: path, roleName: dc.roleName}
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

func (dc *DocContext) GetRoleNames() []string {
    return dc.roleNames
}

func (dc *DocContext) GetAccess(method string, path string, roleName string) int {
    key := roleKey{
        method:   method,
        path:     path,
        roleName: roleName}
    if access, ok := dc.roles[key]; ok {
        return access
    } else {
        return AccessUnknown
    }
}

func (dc *DocContext) SaveEndpointDocumentation() {
    if dc.endpoint != nil {
        dc.endpoints = append(dc.endpoints, *dc.endpoint)
    }
}

func (dc *DocContext) GetEndpoints() []Endpoint {
    return dc.endpoints
}

func (dc *DocContext) GetExampleSummary() string {
    return dc.usageSummary
}

func (dc *DocContext) GetExampleDescription() string {
    return dc.usageDescription
}
