package consumer

import (
	"github.com/go-stomp/stomp/v3"
	"log"
	"yt-clone-video-processing/internal/dependency"
)

func Consume(dependency *dependency.Dependency) {

	sub, err := dependency.MQConn.Subscribe(dependency.Configs.Jobs.TranscodingQueue, stomp.AckAuto)
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
