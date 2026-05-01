package net_http

import (
	"regexp"
	"testing"
)

func TestExtractParams(t *testing.T) {
	r1 := regexp.MustCompile(`^\d+$`)
	r2 := regexp.MustCompile(`^[a-z]+$`)
	params := extractParams("users/{userId}/profiles/{profileId}", []*regexp.Regexp{r1, r2})
	if len(params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(params))
	}
	if params[0].name != "userId" || params[0].format != r1 {
		t.Fatalf("unexpected first param")
	}
	if params[1].name != "profileId" || params[1].format != r2 {
		t.Fatalf("unexpected second param")
	}
}

func TestSanitizePath(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"users", "users"},
		{"/users", "users"},
		{"users/", "users"},
		{"/users/", "users"},
		{"//users//", "users"},
		{"users/{id}", "users/{id}"},
		{"/users/{id}/", "users/{id}"},
		{"//users//{id}//", "users/{id}"},
		{"", ""},
	}
	for _, tc := range cases {
		got := sanitizePath(tc.input)
		if got != tc.want {
			t.Errorf("sanitizePath(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
