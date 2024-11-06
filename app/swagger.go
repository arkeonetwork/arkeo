package app

import (
	"net/http"

	"github.com/arkeonetwork/arkeo/docs"
	"github.com/arkeonetwork/arkeo/pkg/openapiconsole"
	"github.com/cosmos/cosmos-sdk/server/api"
)

// RegisterSwaggerAPI provides a common function which registers swagger route with API Server
func RegisterSwaggerAPI(apiSvr *api.Server) error {
	// register app's OpenAPI routes.
	apiSvr.Router.Handle("/static/swagger.min.json", http.FileServer(http.FS(docs.Docs)))
	apiSvr.Router.HandleFunc("/swagger", openapiconsole.Handler(AppName+" Swagger UI", "/static/swagger.min.json"))
	return nil
}
