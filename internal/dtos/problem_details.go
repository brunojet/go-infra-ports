package dtos

import "encoding/json"

// ProblemDetails implements RFC 7807 (Problem Details for HTTP APIs).
// Reference: https://datatracker.ietf.org/doc/html/rfc7807
type ProblemDetails struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	// Extensions holds extension members; not emitted as a single nested field.
	Extensions map[string]any `json:"-"`
}

func (p *ProblemDetails) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	if p.Type != "" {
		m["type"] = p.Type
	}
	if p.Title != "" {
		m["title"] = p.Title
	}
	if p.Status != 0 {
		m["status"] = p.Status
	}
	if p.Detail != "" {
		m["detail"] = p.Detail
	}
	if p.Instance != "" {
		m["instance"] = p.Instance
	}
	for k, v := range p.Extensions {
		if _, exists := m[k]; !exists {
			m[k] = v
		}
	}
	return json.Marshal(m)
}

func (p *ProblemDetails) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if err := unmarshalRawField(raw, "type", &p.Type); err != nil {
		return err
	}
	if err := unmarshalRawField(raw, "title", &p.Title); err != nil {
		return err
	}
	if err := unmarshalRawField(raw, "status", &p.Status); err != nil {
		return err
	}
	if err := unmarshalRawField(raw, "detail", &p.Detail); err != nil {
		return err
	}
	if err := unmarshalRawField(raw, "instance", &p.Instance); err != nil {
		return err
	}
	extensions, err := unmarshalExtensions(raw)
	p.Extensions = extensions
	return err
}

func unmarshalRawField[T any](raw map[string]json.RawMessage, key string, dest *T) error {
	v, ok := raw[key]
	if !ok {
		return nil
	}
	if err := json.Unmarshal(v, dest); err != nil {
		return err
	}
	delete(raw, key)
	return nil
}

func unmarshalExtensions(raw map[string]json.RawMessage) (map[string]any, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	extensions := make(map[string]any, len(raw))
	for k, v := range raw {
		var vv any
		if err := json.Unmarshal(v, &vv); err != nil {
			return nil, err
		}
		extensions[k] = vv
	}
	return extensions, nil
}

// ContentType returns the official media type for ProblemDetails (RFC 7807).
func (ProblemDetails) ContentType() string {
	return ContentTypeProblemJSON
}
