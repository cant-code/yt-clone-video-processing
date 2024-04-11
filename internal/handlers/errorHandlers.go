package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type FileStatus struct {
	Vid          int64  `json:"vid"`
	Quality      string `json:"quality"`
	ErrorMessage string `json:"errorMessage"`
}

const (
	selectErrorsUsingVid = "SELECT vid, quality, error FROM job_status WHERE vid = $1 AND success = false"
	id                   = "id"
	contentType          = "Content-Type"
	applicationJson      = "application/json"
)

func (dependencies *Dependencies) errorHandler(w http.ResponseWriter, r *http.Request) {
	exec, err := dependencies.DBConn.Query(selectErrorsUsingVid, r.PathValue(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprintf(w, err.Error())

		if err != nil {
			writeInternalServerError(w, err)
		}
		return
	}

	defer func() {
		err = exec.Close()
		if err != nil {
			writeInternalServerError(w, err)
		}
	}()

	statuses := make([]FileStatus, 0)
	for exec.Next() {
		var vid int64
		var quality string
		var errorMessage string
		err := exec.Scan(&vid, &quality, &errorMessage)
		if err != nil {
			log.Println("Error scanning job status:", err)
		}
		statuses = append(statuses, FileStatus{
			Vid:          vid,
			Quality:      quality,
			ErrorMessage: errorMessage,
		})
	}

	encoder, err := json.Marshal(statuses)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}
	w.Header().Set(contentType, applicationJson)
	_, err = w.Write(encoder)

	if err != nil {
		writeInternalServerError(w, err)
	}
}

func writeInternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
}
