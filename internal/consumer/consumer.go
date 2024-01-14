package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/go-stomp/stomp/v3"
	"log"
	"yt-clone-video-processing/internal/config"
	"yt-clone-video-processing/internal/encoder"
)

type Message struct {
	Test string
}

func Consume(config config.Config) {
	dial, err := stomp.Dial("tcp", GenerateAddress(config.MQ), stomp.ConnOpt.Login(config.MQ.User, config.MQ.Password))
	if err != nil {
		log.Fatal(err)
	}

	sub, err := dial.Subscribe(config.Jobs.Queue, stomp.AckAuto)
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

			encoder.EncodeVideo(value.Test)
		}
	}

}

func GenerateAddress(conf config.MQConfig) string {
	return fmt.Sprintf("%s:%s", conf.Host, conf.Port)
}
