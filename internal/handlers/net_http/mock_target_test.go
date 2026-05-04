package net_http

import "testing"

func TestMockTargetTypes_BasicFields(t *testing.T) {
	c := userCreate{Name: "alice"}
	u := userUpdate{Name: "bob"}
	r := userResp{ID: "1", Name: "charlie"}

	if c.Name != "alice" {
		t.Fatalf("expected create name alice, got %s", c.Name)
	}
	if u.Name != "bob" {
		t.Fatalf("expected update name bob, got %s", u.Name)
	}
	if r.ID != "1" || r.Name != "charlie" {
		t.Fatalf("unexpected response value: %+v", r)
	}
}
