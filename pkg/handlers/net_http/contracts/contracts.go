// Package contracts contains public net/http handler contracts.
package contracts

import "net/http"

// NetHttpHandler defines CRUD HTTP operations exposed by the adapter.
type NetHttpHandler interface {
	// Create handles collection POST operations.
	Create(w http.ResponseWriter, r *http.Request)
	// List handles collection GET operations.
	List(w http.ResponseWriter, r *http.Request)
	// Get handles instance GET operations.
	Get(w http.ResponseWriter, r *http.Request)
	// Update handles instance PUT operations.
	Update(w http.ResponseWriter, r *http.Request)
	// Save handles instance PATCH operations.
	Save(w http.ResponseWriter, r *http.Request)
	// Delete handles instance DELETE operations.
	Delete(w http.ResponseWriter, r *http.Request)
}
