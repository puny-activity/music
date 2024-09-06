package httprouter

import (
	"github.com/go-chi/chi/v5"
	httpcontroller "github.com/puny-activity/music/api/http/controller"
	httpmiddleware "github.com/puny-activity/music/api/http/middleware"
	"github.com/puny-activity/music/config"
	"github.com/rs/zerolog"
)

type Router struct {
	cfg        *config.HTTP
	router     *chi.Mux
	middleware *httpmiddleware.Middleware
	wrapper    *Wrapper
	controller *httpcontroller.Controller
	log        *zerolog.Logger
}

func New(cfg *config.HTTP, router *chi.Mux, middleware *httpmiddleware.Middleware,
	wrapper *Wrapper, controller *httpcontroller.Controller, log *zerolog.Logger) *Router {
	return &Router{
		cfg:        cfg,
		router:     router,
		middleware: middleware,
		wrapper:    wrapper,
		controller: controller,
		log:        log,
	}
}

func (r *Router) Setup() {
	r.router.Group(func(router chi.Router) {
		router.Route("/genres", func(router chi.Router) {
			router.Get("/", r.wrapper.Wrap(r.controller.GetGenres))
		})

		router.Route("/albums", func(router chi.Router) {
			router.Get("/", r.wrapper.Wrap(r.controller.GetAlbums))
		})

		router.Route("/artists", func(router chi.Router) {
			router.Get("/", r.wrapper.Wrap(r.controller.GetArtists))
		})
	})
}
