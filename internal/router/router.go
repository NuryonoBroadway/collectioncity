// Package router
package router

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"

	"gitlab.privy.id/collection/collection-city/internal/appctx"
	"gitlab.privy.id/collection/collection-city/internal/bootstrap"
	"gitlab.privy.id/collection/collection-city/internal/consts"
	"gitlab.privy.id/collection/collection-city/internal/handler"
	"gitlab.privy.id/collection/collection-city/internal/middleware"
	"gitlab.privy.id/collection/collection-city/internal/provider"
	"gitlab.privy.id/collection/collection-city/internal/provider/collection_core_provider"
	"gitlab.privy.id/collection/collection-city/internal/ucase"
	"gitlab.privy.id/collection/collection-city/pkg/logger"
	"gitlab.privy.id/collection/collection-city/pkg/routerkit"

	//"gitlab.privy.id/collection/collection-city/pkg/mariadb"
	//"gitlab.privy.id/collection/collection-city/internal/repositories"
	//"gitlab.privy.id/collection/collection-city/internal/ucase/example"

	"gitlab.privy.id/collection/collection-city/internal/ucase/collection_city_ucase"
	ucaseContract "gitlab.privy.id/collection/collection-city/internal/ucase/contract"
)

type router struct {
	config *appctx.Config
	router *routerkit.Router
}

// NewRouter initialize new router wil return Router Interface
func NewRouter(cfg *appctx.Config) Router {
	bootstrap.RegistryMessage()
	bootstrap.RegistryLogger(cfg)

	return &router{
		config: cfg,
		router: routerkit.NewRouter(routerkit.WithServiceName(cfg.App.AppName)),
	}
}

func (rtr *router) handle(hfn httpHandlerFunc, svc ucaseContract.UseCase, mdws ...middleware.MiddlewareFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
				w.WriteHeader(http.StatusInternalServerError)
				res := appctx.Response{
					Code: consts.CodeInternalServerError,
				}

				res.GenerateMessage()
				logger.Error(logger.MessageFormat("error %v", string(debug.Stack())))
				json.NewEncoder(w).Encode(res)
				return
			}
		}()

		ctx := context.WithValue(r.Context(), "access", map[string]interface{}{
			"path":      r.URL.Path,
			"remote_ip": r.RemoteAddr,
			"method":    r.Method,
		})

		req := r.WithContext(ctx)

		if status := middleware.FilterFunc(rtr.config, req, mdws); status != 200 {
			rtr.response(w, appctx.Response{
				Code: status,
			})

			return
		}

		resp := hfn(req, svc, rtr.config)
		resp.Lang = rtr.defaultLang(req.Header.Get(consts.HeaderLanguageKey))
		rtr.response(w, resp)
	}
}

// response prints as a json and formatted string for DGP legacy
func (rtr *router) response(w http.ResponseWriter, resp appctx.Response) {

	w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)

	defer func() {
		resp.GenerateMessage()
		w.WriteHeader(resp.GetCode())
		json.NewEncoder(w).Encode(resp)
	}()

	return

}

// Route preparing http router and will return mux router object
func (rtr *router) Route() *routerkit.Router {

	root := rtr.router.PathPrefix("/").Subrouter()
	//in := root.PathPrefix("/in/").Subrouter()
	liveness := root.PathPrefix("/").Subrouter()
	//inV1 := in.PathPrefix("/v1/").Subrouter()

	// open tracer setup
	bootstrap.RegistryOpenTracing(rtr.config)

	// db := bootstrap.RegistryMariaMasterSlave(rtr.config.WriteDB, rtr.config.ReadDB)

	// use case
	healthy := ucase.NewHealthCheck()

	// healthy
	liveness.HandleFunc("/liveness", rtr.handle(
		handler.HttpRequest,
		healthy,
	)).Methods(http.MethodGet)

	// this is use case for example purpose, please delete
	//repoExample := repositories.NewExample(db)
	//el := example.NewExampleList(repoExample)
	//ec := example.NewPartnerCreate(repoExample)
	//ed := example.NewExampleDelete(repoExample)

	// TODO: create your route here

	// this route for example rest, please delete
	// example list
	//inV1.HandleFunc("/example", rtr.handle(
	//    handler.HttpRequest,
	//    el,
	//)).Methods(http.MethodGet)

	//inV1.HandleFunc("/example", rtr.handle(
	//    handler.HttpRequest,
	//    ec,
	//)).Methods(http.MethodPost)

	//inV1.HandleFunc("/example/{id:[0-9]+}", rtr.handle(
	//    handler.HttpRequest,
	//    ed,
	//)).Methods(http.MethodDelete)
	collection_core := bootstrap.RegistryCollectionCore(rtr.config)
	prov := provider.NewProvider(collection_core)
	collection_core_provider := collection_core_provider.NewCollecionCoreProvider(prov)

	collection_book := root.PathPrefix("/").Subrouter()
	collection_book.HandleFunc("/v1/internal/collection/book", rtr.handle(
		handler.HttpRequest,
		collection_city_ucase.NewCreateCollectionBook(collection_core_provider),
	)).Methods(http.MethodPost)

	collection_book.HandleFunc("/v1/internal/collection/book/all", rtr.handle(
		handler.HttpRequest,
		collection_city_ucase.NewGetCollectionBook(collection_core_provider),
	)).Methods(http.MethodGet)

	collection_book.HandleFunc("/v1/internal/collection/book/detail/{document_id}", rtr.handle(
		handler.HttpRequest,
		collection_city_ucase.NewGetDetailCollectionBook(collection_core_provider),
	)).Methods(http.MethodGet)

	return rtr.router

}

func (rtr *router) defaultLang(l string) string {

	if len(l) == 0 {
		return rtr.config.App.DefaultLang
	}

	return l
}
