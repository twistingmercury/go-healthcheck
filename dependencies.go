package healthcheck

import "encoding/json"

// DependencyDescriptor defines a resource to be checked during a heartbeat request.
type DependencyDescriptor struct {
	Name        string                          `json:"name"`
	Type        string                          `json:"type"`
	Connection  string                          `json:"connection"`
	HandlerFunc func() (hsr HealthStatusResult) `json:"-"`
}

func (d *DependencyDescriptor) String() string {
	text, _ := json.MarshalIndent(d, "", "  ")
	return string(text)
}

// HealthStatusResult represents another process or API that this service relies upon to be considered healthy.
type HealthStatusResult struct {
	Status          HealthStatus `json:"status"`
	Name            string       `json:"name,omitempty"`
	Resource        string       `json:"resource"`
	RequestDuration float64      `json:"request_duration_ms"`
	StatusCode      int          `json:"http_status_code"`
	Message         string       `json:"message,omitempty"`
}

func (dep *HealthStatusResult) String() string {
	text, _ := json.MarshalIndent(dep, "", "  ")
	return string(text)
}
