package consumer

import (
	"encoding/json"
	"github.com/go-stomp/stomp/v3"
	"log"
	"os"
	"yt-clone-video-processing/internal/dependency"
	"yt-clone-video-processing/internal/encoder"
	"yt-clone-video-processing/internal/objectStorage"
)

type Message struct {
	FileId   int64
	FileName string
}

type FileData struct {
	FileName string `json:"fileName"`
	Quality  int    `json:"quality"`
}

type FileManagementMessage struct {
	FileId int64      `json:"fileId"`
	Files  []FileData `json:"files"`
}

type EncoderResponse struct {
	Err      error
	FileName string
	Quality  int
}

var Quality = [3]int{
	720,
	480,
	360,
}

func RunJob(msg *stomp.Message, dependency *dependency.Dependency) {
	var value Message
	err := json.Unmarshal(msg.Body, &value)
	if err != nil {
		log.Panicln(err)
	}

	object, err := objectStorage.GetObject(value.FileName, *dependency)
	if err != nil {
		log.Panicln(err)
	}

	var subProcCount = 0
	channel := make(chan EncoderResponse)

	for _, target := range Quality {
		subProcCount += 1

		go func(target int) {
			video, err2 := encoder.EncodeVideo(object, target)
			if err2 != nil {
				channel <- EncoderResponse{
					Err: err2,
				}
			}

			putObject, err2 := objectStorage.PutObject(video, *dependency)
			if err2 != nil {
				channel <- EncoderResponse{
					Err: err2,
				}
			}

			channel <- EncoderResponse{
				Err:      nil,
				FileName: putObject,
				Quality:  target,
			}
		}(target)
	}

	var response = FileManagementMessage{FileId: value.FileId}

	for i := 0; i < subProcCount; i++ {
		encoderResponse := <-channel
		if encoderResponse.Err != nil {
			log.Panicln(encoderResponse.Err)
		}
		response.Files = append(response.Files, FileData{
			FileName: encoderResponse.FileName,
			Quality:  encoderResponse.Quality,
		})
	}

	marshal, err := json.Marshal(response)
	if err != nil {
		log.Panicln(err)
	}

	err = dependency.MQConn.Send(dependency.Configs.Jobs.ManagementQueue, "text/plain", marshal)
	if err != nil {
		log.Panicln(err)
	}

	err = os.Remove(object)
	if err != nil {
		log.Panicln(err)
	}
}
