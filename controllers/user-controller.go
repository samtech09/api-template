package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	g "github.com/samtech09/api-template/global"
	"github.com/samtech09/api-template/sqls"
	"github.com/samtech09/api-template/viewmodels"
	"github.com/samtech09/api-template/web"

	"github.com/go-chi/chi"
)

//
//#region base initialization

//User is type for UserController
type User struct {
	*APIController
}

// New creates New UserController and initialize it with BaseController
func (t *User) New(v *web.APIVersion) {
	ac := NewAPIController(t, "users", v)

	//Create mount and add routes
	v.Router.Mount("/users", t.SetRoutes())

	t.APIController = ac
}

//SetRoutes create routes for controller
func (t *User) SetRoutes() http.Handler {
	r := chi.NewRouter()

	// Set anonumous routes
	r.Get("/", t.Listusers)

	r.Post("/", t.Createuser)

	// set routes with Authentication
	r.Group(func(r chi.Router) {
		r.Use(web.AuthRequired())
		//r.Post("/", t.Createuser)
		r.Delete("/{id}", t.Deleteuser)
	})
	return r
}

//#endregion
//

//Listusers will list all users
func (t *User) Listusers(w http.ResponseWriter, r *http.Request) {
	//get list of users
	stmt := sqls.UserListAll()

	// request will be cancelled after 10 seconds if takes more than 10 sec to execute
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	rows, err := g.Db.Conn(stmt.ReadOnly).Query(ctx, stmt.SQL)
	if err != nil {
		t.Error500(w, err.Error(), "user.ListUsers")
		return
	}
	// make sure to always close rows
	defer rows.Close()

	var list []viewmodels.DbUser
	count := 0
	isfirst := true
	for rows.Next() {
		// scan rows to struct
		u := viewmodels.DbUser{}
		err := rows.Scan(&u.ID, &u.Name, &count)
		if err != nil {
			t.Error500(w, err.Error(), "user.ListUsers")
			return
		}

		if isfirst {
			isfirst = false
			// instead of append, make slice of fixed length to avoid reallocations
			list = make([]viewmodels.DbUser, 0, count)
		}
		list = append(list, u)
	}

	t.JSON(w, list)
}

//Createuser will create new user in database
func (t *User) Createuser(w http.ResponseWriter, r *http.Request) {
	usr := viewmodels.DbUser{}
	err := t.BindJSON(r, &usr)
	if err != nil {
		t.Error500(w, err.Error(), "user.Createuser")
		return
	}

	// request will be cancelled after 10 seconds if takes more than 10 sec to execute
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// create user
	stmt := sqls.UserCreate()
	rows, err := g.Db.Conn(stmt.ReadOnly).Query(ctx, stmt.SQL, usr.ID, usr.Name)
	if err != nil {
		t.Error500(w, err.Error(), "user.Createuser")
		return
	}

	if rows.Err() != nil {
		fmt.Println("rows.err")
		t.Error500(w, rows.Err().Error(), "user.Createuser")
		return
	}

	// make sure to always close rows
	defer rows.Close()
	var id int
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			t.Error500(w, err.Error(), "user.Createuser")
			return
		}
	}

	t.JSON(w, id)
}

//Deleteuser will delete given user by id
func (t *User) Deleteuser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// request will be cancelled after 10 seconds if takes more than 10 sec to execute
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// delete user
	stmt := sqls.UserDelete()
	ctag, err := g.Db.Conn(stmt.ReadOnly).Exec(ctx, stmt.SQL, id)
	if err != nil {
		t.Error500(w, err.Error(), "user.Deleteuser")
		return
	}
	t.JSON(w, ctag.RowsAffected())
}
