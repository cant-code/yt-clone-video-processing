package main

import (
	"log"
	"yt-clone-video-processing/internal/config"
	"yt-clone-video-processing/internal/consumer"
)

func main() {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	consumer.Consume(loadConfig)
}
