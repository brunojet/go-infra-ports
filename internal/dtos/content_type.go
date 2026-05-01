package dtos

// Content type constants and helper interface for DTOs.
const (
	ContentTypeProblemJSON = "application/problem+json"
	ContentTypeJSON        = "application/json"
	ContentTypeURIList     = "text/uri-list"
)

// ContentTyper is implemented by DTOs that can report their preferred Content-Type.
type ContentTyper interface {
	ContentType() string
}
