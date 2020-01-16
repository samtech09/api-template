package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
func (b *APIController) GetClientInfo(r *http.Request, callerName string) jwtauth.ClientInfo {
	// check request has Auth header, otherwise there will be not client info
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		g.Logger.Error().Str("Fn", callerName).Msg("GetClientInfo called for request without authorization header")
		return jwtauth.ClientInfo{}
	}

	ac := r.Context().Value(web.KeyAppContext)
	if ac == nil {
		g.Logger.Error().Str("Fn", callerName).Msg("GetClientInfo failed to get application context")
		return jwtauth.ClientInfo{}
	}
	return ac.(jwtauth.ClientInfo)
}

//
// ---------------
//

//JSON set resultant data into response that is injected into Data property of APiResult as JSONString
func (b *APIController) JSON(w http.ResponseWriter, data interface{}) {
	jsondata, _ := toJSON(data)

	g.Logger.Debug().Str("data", jsondata).
		Msg("api-controller/JSON")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(web.NewResultFilled(jsondata, http.StatusOK, 0, ""))
}

//JSONStr set resultant data into response that is injected into Data property of APiResult as-it-is without any conversion
func (b *APIController) JSONStr(w http.ResponseWriter, jsondata string) {
	g.Logger.Debug().Str("data", jsondata).Msg("api-controller/JSONStr")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(web.NewResultFilled(jsondata, http.StatusOK, 0, ""))
}

//JSONRaw set resultant data as JSON of passed data itself. Data not encapsulated in APiResult.
//Use with caution as called must know what data it will return for parsing.
func (b *APIController) JSONRaw(w http.ResponseWriter, data interface{}) {
	g.Logger.Debug().Str("data", "raw-data").Msg("api-controller/JSONRaw")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

//JSONReplacer set resultant data into response APIResult, but before sending response,
//it makes replacements in JSONString of data as per supplied replacers
func (b *APIController) JSONReplacer(w http.ResponseWriter, data interface{}, r *strings.Replacer) {
	jsondata, _ := toJSON(data)
	jsondata = r.Replace(jsondata)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(web.NewResultFilled(jsondata, http.StatusOK, 0, ""))
}

//Content writes given string to response
func (b *APIController) Content(w http.ResponseWriter, content string) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, content)
}

//Error set error into response that is later served as JSON
func (b *APIController) Error(w http.ResponseWriter, err string, code int, callerName string) {
	g.Logger.Error().Str("Fn", callerName).Msg(err)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(web.NewResultFilled("", code, 1, err))
}

//Error500 return 500-internal server error
func (b *APIController) Error500(w http.ResponseWriter, err string, callerName string) {
	// pc, _, _, ok := runtime.Caller(1)
	// fn := runtime.FuncForPC(pc)
	// if ok && fn != nil {
	// 	dotName := filepath.Ext(fn.Name())
	// 	fnName := strings.TrimLeft(dotName, ".") + "()"
	// 	fmt.Println("Function: ", fnName)
	// }

	g.Logger.Error().Str("Fn", callerName).Msg(err)
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(web.NewResultFilled("", http.StatusInternalServerError, 1, err))
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
