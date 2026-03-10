package server

type AdmissionRequest struct {
	Operation string         `json:"operation,omitempty"`
	Kind      string         `json:"kind,omitempty"`
	Name      string         `json:"name,omitempty"`
	Namespace string         `json:"namespace,omitempty"`
	Object    map[string]any `json:"object,omitempty"`
	OldObject map[string]any `json:"old_object,omitempty"`
}

type AdmissionResponse struct {
	Allowed  bool           `json:"allowed"`
	Message  string         `json:"message,omitempty"`
	Patch    map[string]any `json:"patch,omitempty"`
	Warnings []string       `json:"warnings,omitempty"`
}
