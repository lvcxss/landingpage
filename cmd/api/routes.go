package main

import (
	"net/http"

	"landingpage/internal/ui"
)

func (app *application) routes() http.Handler {
	syncHandler := http.HandlerFunc(app.syncPost)
	adminToken := app.config.AdminToken
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.FS(ui.StaticFiles))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/post/{id}", app.showPost)
	mux.Handle("/api/admin/sync", requireAdminToken(adminToken, syncHandler))
	return rateLimiter(mux)
}
