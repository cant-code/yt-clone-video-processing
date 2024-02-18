package main

import (
	"log"
	"yt-clone-video-processing/internal/consumer"
	"yt-clone-video-processing/internal/dependency"
	"yt-clone-video-processing/internal/initializations"
)

func main() {
	dependencies, err := dependency.GetDependencies()
	if err != nil {
		log.Fatal(err)
	}

	initializations.RunMigrations(dependencies)

	consumer.Consume(dependencies)
}
