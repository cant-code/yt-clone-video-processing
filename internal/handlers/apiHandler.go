package handlers

import (
	"database/sql"
	"net/http"
	"yt-clone-video-processing/internal/auth"
)

type Dependencies struct {
	DBConn *sql.DB
}

func (apiConfig *Dependencies) ApiHandler() *http.ServeMux {
	mux := http.NewServeMux()

	handler := auth.HandleJwtAuthMiddleware()

	mux.Handle("GET /videos/errors/{id}", handler(http.HandlerFunc(apiConfig.errorHandler)))

	return mux
}
