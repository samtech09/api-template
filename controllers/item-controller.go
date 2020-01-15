package controllers

import (
	"github.com/samtech09/api-template/web"
	"net/http"

	"github.com/go-chi/chi"
)

//Item is type for ItemController
type Item struct {
	*APIController
}

// New creates New DataSync and initialize it with BaseController
func (t *Item) New(v *web.APIVersion) {
	ac := NewAPIController(t, "items", v)

	//Add mount and add routes
	v.Router.Mount("/items", t.SetRoutes())

	t.APIController = ac
}

//SetRoutes create routes for controller
func (t *Item) SetRoutes() http.Handler {
	r := chi.NewRouter()

	// Set anonumous routes
	r.Get("/", t.Listitems)

	// set routes with Authentication
	r.Group(func(r chi.Router) {
		r.Use(web.AuthRequired())

		r.Delete("/{id}", t.Deleteitem)
	})
	return r
}

//Listitems will list all items
func (t *Item) Listitems(w http.ResponseWriter, r *http.Request) {
	t.Content(w, "it will list items")
}

//Deleteitem will delete given item by id
func (t *Item) Deleteitem(w http.ResponseWriter, r *http.Request) {
	t.Content(w, "it will delete item")
}
