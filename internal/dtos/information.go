package dtos

// Placeholder for informational response DTOs (1xx or custom information objects).
// Define concrete structures here following any RFCs or project conventions as needed.

// Information is a generic container for informational responses.
type Information struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// ContentType returns the default media type for informational DTOs.
func (Information) ContentType() string {
	return ContentTypeJSON
}
