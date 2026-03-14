package dto

type UpsertQueueRequest struct {
	TenantID    string  `json:"tenant_id"`
	Name        string  `json:"name"`
	Weight      int     `json:"weight"`
	Priority    int     `json:"priority"`
	Enabled     *bool   `json:"enabled,omitempty"`
	Description *string `json:"description,omitempty"`
}

type QueueResponse struct {
	ID          string  `json:"id"`
	TenantID    string  `json:"tenant_id"`
	Name        string  `json:"name"`
	Weight      int     `json:"weight"`
	Priority    int     `json:"priority"`
	Enabled     bool    `json:"enabled"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
