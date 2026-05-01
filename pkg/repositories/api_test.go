package repositories

import "testing"

func TestRepositoryAliases_AreUsable(t *testing.T) {
	create := RepositoryCreate[string]{Body: "a"}
	update := RepositoryUpdate[string]{Body: "b"}
	save := RepositorySave[string]{Body: "c"}
	resp := RepositoryResponse[string]{Data: "x"}
	resps := RepositoryResponses[string]{Data: []string{"x", "y"}}
	rc := RestCreate[string]{Body: "a"}
	ru := RestUpdate[string]{Body: "b"}
	rs := RestSave[string]{Body: "c"}
	rr := RestResponse[string]{Data: "x"}
	rrs := RestResponses[string]{Data: []string{"x", "y"}}

	if create.Body != "a" || update.Body != "b" || save.Body != "c" {
		t.Fatalf("unexpected repository alias payload values")
	}
	if resp.Data != "x" || len(resps.Data) != 2 {
		t.Fatalf("unexpected repository alias response values")
	}
	if rc.Body != "a" || ru.Body != "b" || rs.Body != "c" || rr.Data != "x" || len(rrs.Data) != 2 {
		t.Fatalf("unexpected rest alias values")
	}
}
