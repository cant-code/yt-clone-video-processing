package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/go-stomp/stomp/v3"
	"log"
	"os"
	"sync"
	"yt-clone-video-processing/internal/configurations"
	"yt-clone-video-processing/internal/dependency"
	"yt-clone-video-processing/internal/encoder"
	"yt-clone-video-processing/internal/objectStorage"
)

type Message struct {
	FileId   int64
	FileName string
}

var Pixels = [3]int{
	720,
	480,
	360,
}

func Consume(dependency *dependency.Dependency) {
	dial, err := stomp.Dial("tcp",
		GenerateAddress(dependency.Configs.MQ),
		stomp.ConnOpt.Login(dependency.Configs.MQ.User, dependency.Configs.MQ.Password))
	if err != nil {
		log.Fatal(err)
	}

	sub, err := dial.Subscribe(dependency.Configs.Jobs.Queue, stomp.AckAuto)
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg := <-sub.C

		if msg != nil {
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
	}

}

func GenerateAddress(conf configurations.MQConfig) string {
	return fmt.Sprintf("%s:%s", conf.Host, conf.Port)
}
