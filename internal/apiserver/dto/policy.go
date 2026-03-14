package dto

type UpsertPolicyRequest struct {
	TenantID           string         `json:"tenant_id"`
	Name               string         `json:"name"`
	Queue              string         `json:"queue"`
	Priority           int            `json:"priority"`
	Enabled            *bool          `json:"enabled,omitempty"`
	Preemptible        bool           `json:"preemptible"`
	RequireHealthy     bool           `json:"require_healthy"`
	RequireMIG         bool           `json:"require_mig"`
	MaxGPUCount        int            `json:"max_gpu_count"`
	RequiredGPUModel   string         `json:"required_gpu_model"`
	RequiredNodeLabels map[string]any `json:"required_node_labels"`
	Selector           map[string]any `json:"selector"`
	Description        *string        `json:"description,omitempty"`
}

type PolicyResponse struct {
	ID                 string         `json:"id"`
	TenantID           string         `json:"tenant_id"`
	Name               string         `json:"name"`
	Queue              string         `json:"queue"`
	Priority           int            `json:"priority"`
	Enabled            bool           `json:"enabled"`
	Preemptible        bool           `json:"preemptible"`
	RequireHealthy     bool           `json:"require_healthy"`
	RequireMIG         bool           `json:"require_mig"`
	MaxGPUCount        int            `json:"max_gpu_count"`
	RequiredGPUModel   string         `json:"required_gpu_model"`
	RequiredNodeLabels map[string]any `json:"required_node_labels,omitempty"`
	Selector           map[string]any `json:"selector,omitempty"`
	Description        *string        `json:"description,omitempty"`
	CreatedAt          string         `json:"created_at"`
	UpdatedAt          string         `json:"updated_at"`
}
