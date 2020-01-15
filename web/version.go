package web

import (
	"github.com/go-chi/chi"
)

//APIVersion struct to hold api version information
type APIVersion struct {
	Router chi.Router
	V      string
}

//NewAPIVersion create new version path for APIs
// like /v1  or  /v2 etc.
func NewAPIVersion(version string, r *chi.Mux) *APIVersion {
	v := APIVersion{}
	v.V = version
	r.Use(APIVersionCtx(version))
	v.Router = r.Route("/"+version, func(r chi.Router) {})
	return &v
}
