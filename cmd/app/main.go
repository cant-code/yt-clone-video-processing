package main

import (
	"log"
	"yt-clone-video-processing/internal/consumer"
	"yt-clone-video-processing/internal/dependency"
)

func main() {
	dependencies, err := dependency.GetDependencies()
	if err != nil {
		log.Fatal(err)
	}

	consumer.Consume(dependencies)
}
