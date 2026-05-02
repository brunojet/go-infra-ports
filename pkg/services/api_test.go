package services

import "testing"

func TestServiceAliases_AreUsable(t *testing.T) {
	create := ServiceCreate[string]{Body: "a"}
	update := ServiceUpdate[string]{Body: "b"}
	save := ServiceSave[string]{Body: "c"}
	meta := ServiceMeta{Message: "ok", Code: "CODE"}
	resp := ServiceResponse[string]{Meta: meta, Data: "x"}
	resps := ServiceResponses[string]{Meta: meta, Data: []string{"x", "y"}}

	if create.Body != "a" || update.Body != "b" || save.Body != "c" {
		t.Fatalf("unexpected alias payload values")
	}
	if resp.Meta.Code != "CODE" || len(resps.Data) != 2 {
		t.Fatalf("unexpected alias response values")
	}
}
