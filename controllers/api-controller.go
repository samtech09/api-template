package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	config "github.com/samtech09/api-template/config"
	g "github.com/samtech09/api-template/global"
	"github.com/samtech09/api-template/web"
	"github.com/samtech09/jwtauth"

	"github.com/go-chi/render"
)

//Controller provide interface for all controllers
type Controller interface {
	New(*web.APIVersion)
	//Path(string) string
	SetRoutes() http.Handler
}

//APIController is base for all controllers
type APIController struct {
	Name string
	*web.APIVersion
}

//Register called from main to register controller
func Register(c interface{}, v *web.APIVersion) {
	c.(Controller).New(v)
}

//NewAPIController create new instance of BaseController
func NewAPIController(c interface{}, name string, v *web.APIVersion) *APIController {
	b := APIController{}
	b.Name = name
	b.APIVersion = v
	//b.initRoutes(c.(Controller))
	return &b
}

//GetClientInfo retrieve clientinfo from context of authentiated requests only
func (b *APIController) GetClientInfo(r *http.Request) jwtauth.ClientInfo {
	ac := r.Context().Value(web.KeyAppContext)
	if ac == nil {
		g.Logger.Error().Msg("Unable to get application context")
		return jwtauth.ClientInfo{}
	}
	return ac.(jwtauth.ClientInfo)
}

//
// ---------------
//

//JSON set resultant data into response that is later server as JSON using AppContexter middleware
func (b *APIController) JSON(w http.ResponseWriter, data interface{}) {
	jsondata, _ := toJSON(data)

	g.Logger.Debug().Str("data", jsondata).Msg("api-controller/JSON")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(web.NewResultFilled(jsondata, http.StatusOK, 0, ""))
}

//JSONString set resultant data into response that is later server as JSON using AppContexter middleware
func (b *APIController) JSONString(w http.ResponseWriter, jsondata string) {
	g.Logger.Debug().Str("data", jsondata).Msg("api-controller/JSONString")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(web.NewResultFilled(jsondata, http.StatusOK, 0, ""))
}

//ResultReplacer set resultant data into response that is later server as JSON using AppContexter middleware
// It also replace given string in serialized json
func (b *APIController) ResultReplacer(w http.ResponseWriter, data interface{}, find, replace string) {
	jsondata, _ := toJSON(data)
	str := strings.Replace(jsondata, find, replace, -1)

	g.Logger.Debug().Str("data", jsondata).Str("find", find).
		Str("replacewith", replace).
		Msg("api-controller/ResultReplacer")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(web.NewResultFilled(str, http.StatusOK, 0, ""))
}

//Content writes given string to response
func (b *APIController) Content(w http.ResponseWriter, content string) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, content)
}

//Error set error into response that is later served as JSON using AppContexter middleware
func (b *APIController) Error(w http.ResponseWriter, err string, code int, funcName string) {
	if !config.IsProduction {
		//Print err for debug
		g.Logger.Error().Str("Function", funcName).Msg(err)
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(web.NewResultFilled("", code, 1, err))
}

//Error500 return 500-internal server error
func (b *APIController) Error500(w http.ResponseWriter, err string, funcName string) {
	b.Error(w, err, http.StatusInternalServerError, funcName)
}

//BindJSON binds request body to JSON
func (b *APIController) BindJSON(r *http.Request, obj interface{}) error {
	defer r.Body.Close()
	return render.DecodeJSON(r.Body, obj)
}

func toJSON(d interface{}) (string, error) {
	data, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.Write(data)
	return sb.String(), nil
}
