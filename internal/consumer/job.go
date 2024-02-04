package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/go-stomp/stomp/v3"
	"log"
	"os"
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

const contentType = "text/plain"

func RunJob(msg *stomp.Message, dependency *dependency.Dependency) {

	var value model.Message
	err := json.Unmarshal(msg.Body, &value)
	if err != nil {
		log.Println(err)
		return
	}

	object, err := objectStorage.GetObject(value.FileName, *dependency)
	if err != nil {
		log.Println(err)
	}

	var subProcCount = 0
	channel := make(chan EncoderResponse)

	for _, target := range Quality {
		subProcCount += 1

		go EncodeVideoAndUploadToS3(target, object, channel, dependency)
	}

	var response = model.FileManagementMessage{FileId: value.FileId}

	for i := 0; i < subProcCount; i++ {
		encoderResponse := <-channel
		if encoderResponse.Err != nil {
			log.Println(encoderResponse.Err)
			response.Files = append(response.Files, model.FileData{
				Quality: encoderResponse.Quality,
				Success: false,
				Error:   encoderResponse.Err,
			})
		} else {
			response.Files = append(response.Files, model.FileData{
				FileName: encoderResponse.FileName,
				Size:     encoderResponse.Size,
				Quality:  encoderResponse.Quality,
				Success:  true,
			})
		}
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

func EncodeVideoAndUploadToS3(target int, object string, channel chan EncoderResponse, dependency *dependency.Dependency) {
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
