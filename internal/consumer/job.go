package consumer

import (
	"encoding/json"
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
}

var Quality = [3]int{
	720,
	480,
	360,
}

func RunJob(msg *stomp.Message, dependency *dependency.Dependency) {

	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic occurred while running job:", err)
		}
	}()

	var value model.Message
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

		go EncodeVideoAndUploadToS3(target, object, channel, dependency)
	}

	var response = model.FileManagementMessage{FileId: value.FileId}

	for i := 0; i < subProcCount; i++ {
		encoderResponse := <-channel
		if encoderResponse.Err != nil {
			log.Panicln(encoderResponse.Err)
		}
		response.Files = append(response.Files, model.FileData{
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

func EncodeVideoAndUploadToS3(target int, object string, channel chan EncoderResponse, dependency *dependency.Dependency) {
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
}
