// Package router
package router

import (
	"net/http"

	"collection-squad/collection/collection-city/internal/appctx"
	"collection-squad/collection/collection-city/internal/ucase/contract"
	"collection-squad/collection/collection-city/pkg/routerkit"
)

// httpHandlerFunc is a contract http handler for router
type httpHandlerFunc func(request *http.Request, svc contract.UseCase, conf *appctx.Config) appctx.Response

// Router is a contract router and must implement this interface
type Router interface {
	Route() *routerkit.Router
}
