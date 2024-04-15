package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Auth(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")

		ctx := context.WithValue(r.Context(), "token", token)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Routes(bannerHandler *bannerHandler) chi.Router {
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	root.Use(middleware.RequestID)

	r := chi.NewRouter()
	r.Use(Auth)

	root.Mount("/api", r)
	r.Post("/banner", bannerHandler.CreateBanner)
	r.Get("/banner", bannerHandler.GetBannersWithFilter)
	r.Get("/user-banner", bannerHandler.GetUserBanner)
	r.Put("/banner", nil)
	r.Delete("/banner", nil)

	return root
}
