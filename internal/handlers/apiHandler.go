package handlers

import (
	"database/sql"
	"net/http"
	"yt-clone-video-processing/internal/auth"
	"yt-clone-video-processing/internal/configurations"
)

type Dependencies struct {
	DBConn *sql.DB
	Auth   configurations.Auth
}

func (apiConfig *Dependencies) ApiHandler() *http.ServeMux {
	mux := http.NewServeMux()

	handler := auth.HandleJwtAuthMiddleware(&apiConfig.Auth)

	mux.Handle("GET /videos/errors/{id}", handler(http.HandlerFunc(apiConfig.errorHandler)))

	return mux
}
