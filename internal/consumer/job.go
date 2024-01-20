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

type EncoderResponse struct {
	Err      error
	FilePath string
}

var Pixels = [3]int{
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

	for _, target := range Pixels {
		subProcCount += 1

		go func(target int) {
			video, err2 := encoder.EncodeVideo(object, target)
			if err2 != nil {
				channel <- EncoderResponse{
					Err:      err2,
					FilePath: "",
				}
			}

			putObject, err2 := objectStorage.PutObject(video, *dependency)
			if err2 != nil {
				channel <- EncoderResponse{
					Err:      err2,
					FilePath: "",
				}
			}

			channel <- EncoderResponse{
				Err:      nil,
				FilePath: putObject,
			}
		}(target)
	}

	for i := 0; i < subProcCount; i++ {
		log.Println(<-channel)
	}

	err = os.Remove(object)
	if err != nil {
		log.Panicln(err)
	}
}
