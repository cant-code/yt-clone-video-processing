package consumer

import (
	"fmt"
	"github.com/go-stomp/stomp/v3"
	"yt-clone-video-processing/internal/config"
)

type Message struct {
	Test string
}

func Consume(config config.Config) (*stomp.Subscription, error) {
	dial, err := stomp.Dial("tcp", GenerateAddress(config.MQ), stomp.ConnOpt.Login(config.MQ.User, config.MQ.Password))
	if err != nil {
		return nil, err
	}

	sub, err := dial.Subscribe(config.Jobs.Queue, stomp.AckAuto)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func GenerateAddress(conf config.MQConfig) string {
	return fmt.Sprintf("%s:%s", conf.Host, conf.Port)
}
