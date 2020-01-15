package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	g "github.com/samtech09/api-template/global"
)

//APIVersionCtx adds api version to request Context
func APIVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), KeyAPIVersion, version))
			next.ServeHTTP(w, r)
		})
	}
}

//AuthRequired check whether request contais
//valid JWT token
func AuthRequired() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Check tokens only in production
			if g.TestEnv {
				next.ServeHTTP(w, r)
				return
			}

			ci, err := g.JWTval.ValidateRequest(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, err.Error())
				return
			}

			//g.Logger.Debugf("clientinfo: %q", ci)

			if ci.Scopes != nil && len(ci.Scopes) < 1 {
				w.WriteHeader(http.StatusNotAcceptable)
				fmt.Fprint(w, "scope(s) not set")
				return
			}

			// JWT and scopes are fine
			// set client info to context
			r = r.WithContext(context.WithValue(r.Context(), KeyAppContext, *ci))

			// Extract route info
			s := strings.Split(r.URL.Path, "/")
			ln := len(s)
			if ln < 4 {
				//invalid URI, it should have
				//  /version/controller/action
				//e.g.  /v1/speedtests/getstlist
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "invalid resource uri")
				return
			}

			g.Logger.Debug().Str("endpoint", r.URL.Path).Msg("Requested endpoint")

			// s[0] will be blank
			//version := s[1]
			controller := s[2]
			endpoint := s[3]
			method := r.Method
			//controller := ""
			// if ln > 2 {
			// 	controller = s[2]
			// }

			//Check client's granted scopes has access to this route
			scopes, err := g.Mgosesion.GetScopesFromRoute(controller, endpoint, method)
			if (err != nil) || (len(*scopes) < 1) {
				// failed to get scope for given route
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, fmt.Sprintf("no route: %v", err))
				return
			}

			//g.Logger.Debugm("AuthRequired", "Found Scopes: %q\n", scopes)

			//scopeFound := false
			//check if client has any of allowed scope
			hasScope := false
			for _, scop := range *scopes {
				for _, cscop := range ci.Scopes {

					//g.Logger.Debugm("AuthRequired", "      %s\n", scop+" == "+cscop)

					if scop == cscop {
						//g.Logger.Debugm("AuthRequired", " Scope matched calling Next")
						// scope found, so allow access to route
						//scopeFound = true
						hasScope = true
						break
					}
				}
				if hasScope {
					break
				}
			}

			if hasScope {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprint(w, "no route")
				return
			}

		})
	}
}
