package handlers

import (
	"database/sql"
	"net/http"
)

type Dependencies struct {
	DBConn *sql.DB
}

func (dependencies *Dependencies) ApiHandler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /videos/errors/{id}", dependencies.errorHandler)

	return mux
}
