package dtos

// Placeholder for redirection DTOs (3xx semantics).
// Add fields as required by RFCs or internal conventions.

type Redirection struct {
	Location string `json:"location,omitempty"`
	Code     int    `json:"code,omitempty"`
}

// ContentType returns the default media type for redirection DTOs.
// This DTO is encoded as JSON by the handler.
func (Redirection) ContentType() string {
	return ContentTypeJSON
}
