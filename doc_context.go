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
	mode             int           // Documentation collecting data mode.
	usageSummary     string        // Summary for the next collected usage example.
	usageDescription string        // Detailed description for the next collected usage example.
	roleName         string        // Principal's role name for the next endpoint usage example.
	endpoint         *DocEndpoint  // Pointer to currently documented endpoint.
	endpoints        []DocEndpoint // List of all documented endpoints with usage examples.
	roleNames        []string      // Role names in order they should be displayed in preview.
	roles            roles         // Map of verified access roles for all endpoints.
}

func CreateDocContext() *DocContext {
	return &DocContext{
		mode:             CollectNothing,
		usageSummary:     "",
		usageDescription: "",
		roleName:         "",
		endpoint:         nil,
		endpoints:        make([]DocEndpoint, 0),
		roles:            make(roles)}
}

func (dc *DocContext) Clear() {
	dc.mode = CollectNothing
	dc.usageSummary = ""
	dc.usageDescription = ""
	dc.roleName = ""
	dc.endpoint = nil
	dc.endpoints = make([]DocEndpoint, 0)
	dc.roles = make(map[roleKey]int)
}

func (dc *DocContext) NewEndpoint(version string, group string, summary string, description string) {
	if version == "" {
		version = "v1"
	}
	summary = strings.TrimSpace(summary)
	description = strings.TrimSpace(description)
	dc.endpoint = createEndpoint(group, version, summary, description)
}

func (dc *DocContext) GetEndpoint() *DocEndpoint {
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

func (dc *DocContext) CollectAll(summary string, description string) {
	dc.mode = CollectEndpointDescriptionAndUsage
	dc.usageSummary = summary
	dc.usageDescription = description
}

func (dc *DocContext) StopCollecting() {
	dc.mode = CollectNothing
	dc.usageSummary = ""
	dc.usageDescription = ""
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

func (dc *DocContext) SaveEndpoint() {
	if dc.endpoint != nil {
		dc.endpoints = append(dc.endpoints, *dc.endpoint)
		dc.endpoint = nil
	}
}

func (dc *DocContext) GetEndpoints() []DocEndpoint {
	return dc.endpoints
}

func (dc *DocContext) GetExampleSummary() string {
	return dc.usageSummary
}

func (dc *DocContext) GetExampleDescription() string {
	return dc.usageDescription
}
