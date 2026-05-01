package net_http

// Concrete types used by tests and mocks.
type userCreate struct{ Name string }
type userResp struct {
	ID   string
	Name string
}
type userUpdate struct{ Name string }
