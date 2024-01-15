package consumer

import (
	"encoding/json"
	"github.com/go-stomp/stomp/v3"
	"log"
	"os"
	"sync"
	"yt-clone-video-processing/internal/dependency"
	"yt-clone-video-processing/internal/encoder"
	"yt-clone-video-processing/internal/objectStorage"
)

func RunJob(msg *stomp.Message, dependency *dependency.Dependency) {
	var value Message
	err := json.Unmarshal(msg.Body, &value)
	if err != nil {
		log.Println(err)
	}

	object, err := objectStorage.GetObject(value.FileName, *dependency)
	if err != nil {
		log.Println(err)
	}

	var waitGroup sync.WaitGroup

	for _, target := range Pixels {
		waitGroup.Add(1)

		go func(target int) {
			defer waitGroup.Done()
			encoder.EncodeVideo(object, target)
		}(target)
	}

	waitGroup.Wait()
	err = os.Remove(object)
	if err != nil {
		log.Panicln(err)
	}
}
