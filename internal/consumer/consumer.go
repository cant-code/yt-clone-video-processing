package consumer

import (
	"fmt"
	"github.com/go-stomp/stomp/v3"
	"log"
	"yt-clone-video-processing/internal/configurations"
	"yt-clone-video-processing/internal/dependency"
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
			go RunJob(msg, dependency)
		}
	}

}

func GenerateAddress(conf configurations.MQConfig) string {
	return fmt.Sprintf("%s:%s", conf.Host, conf.Port)
}
