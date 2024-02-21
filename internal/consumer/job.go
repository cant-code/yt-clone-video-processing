package consumer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-stomp/stomp/v3"
	"log"
	"os"
	"strconv"
	"yt-clone-video-processing/internal/dependency"
	"yt-clone-video-processing/internal/encoder"
	"yt-clone-video-processing/internal/objectStorage"
	"yt-clone-video-processing/pkg/model"
)

type EncoderResponse struct {
	Err      error
	FileName string
	Quality  int
	Size     int64
}

var Quality = [3]int{
	720,
	480,
	360,
}

const (
	contentType            = "text/plain"
	insertQuery            = "INSERT INTO FILE_STATUS (id, vid, status) VALUES (nextval('file_status_seq'), $1, $2) RETURNING ID"
	updateQuery            = "UPDATE FILE_STATUS SET STATUS = $1 WHERE ID = $2"
	insertJobStatWithError = "INSERT INTO JOB_STATUS (vid, quality, success, error) VALUES ($1, $2, false, $3)"
	insertJobStat          = "INSERT INTO JOB_STATUS (vid, quality, success) VALUES ($1, $2, true)"
)

type STATUS int

const (
	Started STATUS = iota
	Success
	Failed
	PartialSuccess
)

func (status STATUS) string() string {
	return [...]string{"STARTED", "SUCCESS", "FAILED", "PARTIAL_SUCCESS"}[status]
}

func RunJob(msg *stomp.Message, dependency *dependency.Dependency) {

	var value model.Message
	err := json.Unmarshal(msg.Body, &value)
	if err != nil {
		log.Println(err)
		return
	}

	var id int64
	err = dependency.DBConn.QueryRow(insertQuery, value.FileId, Started.string()).Scan(&id)
	if err != nil {
		log.Println("Error while adding entry into the DB", err)
	}

	object, err := objectStorage.GetObject(value.FileName, *dependency)
	if err != nil {
		log.Println(err)
	}

	channel := make(chan EncoderResponse)

	for _, target := range Quality {
		go encodeVideoAndUploadToS3(target, object, channel, dependency)
	}

	var response = model.FileManagementMessage{FileId: value.FileId}

	var failCount = 0

	for range Quality {
		encoderResponse := <-channel
		if encoderResponse.Err != nil {
			log.Println(encoderResponse.Err)
			failCount++

			addJobStatus(value.FileId, strconv.Itoa(encoderResponse.Quality), encoderResponse.Err, dependency.DBConn)

			response.Files = append(response.Files, model.FileData{
				Quality: encoderResponse.Quality,
				Success: false,
				Error:   encoderResponse.Err,
			})
		} else {
			addJobStatus(value.FileId, strconv.Itoa(encoderResponse.Quality), nil, dependency.DBConn)

			response.Files = append(response.Files, model.FileData{
				FileName: encoderResponse.FileName,
				Size:     encoderResponse.Size,
				Quality:  encoderResponse.Quality,
				Success:  true,
			})
		}
	}

	switch failCount {
	case len(Quality):
		updateStatusForId(id, Failed, dependency.DBConn)
	case 0:
		updateStatusForId(id, Success, dependency.DBConn)
	default:
		updateStatusForId(id, PartialSuccess, dependency.DBConn)
	}

	marshal, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}

	err = dependency.MQConn.Send(dependency.Configs.Jobs.ManagementQueue, contentType, marshal)
	if err != nil {
		log.Println(err)
	}

	err = os.Remove(object)
	if err != nil {
		log.Println(err)
	}
}

func encodeVideoAndUploadToS3(target int, object string, channel chan EncoderResponse, dependency *dependency.Dependency) {
	video, err := encoder.EncodeVideo(object, target)
	if err != nil {
		channel <- EncoderResponse{
			Err:     err,
			Quality: target,
		}
	}

	var size int64 = 0
	stat, err := os.Stat(fmt.Sprintf("./files/%s", video))
	if err != nil {
		log.Println(err)
	}
	size = stat.Size()

	key, err := objectStorage.PutObject(video, *dependency)
	if err != nil {
		channel <- EncoderResponse{
			Err:     err,
			Quality: target,
		}
	}

	channel <- EncoderResponse{
		FileName: key,
		Quality:  target,
		Size:     size,
	}
}

func addJobStatus(id int64, quality string, message error, DBConn *sql.DB) {
	var err error
	if message != nil {
		_, err = DBConn.Exec(insertJobStatWithError, id, quality, message)
	} else {
		_, err = DBConn.Exec(insertJobStat, id, quality)
	}
	if err != nil {
		log.Println("Error while inserting job status", err)
	}
}

func updateStatusForId(id int64, status STATUS, DBConn *sql.DB) {
	_, err := DBConn.Exec(updateQuery, status.string(), id)
	if err != nil {
		log.Println("Error while updating status", err)
	}
}
