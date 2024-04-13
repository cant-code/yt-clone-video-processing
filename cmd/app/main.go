package main

import (
	"log"
	"net/http"
	"yt-clone-video-processing/internal/consumer"
	"yt-clone-video-processing/internal/dependency"
	"yt-clone-video-processing/internal/handlers"
	"yt-clone-video-processing/internal/initializations"
)

func main() {
	dependencies, err := dependency.GetDependencies()
	if err != nil {
		log.Fatal(err)
	}

	initializations.RunMigrations(dependencies)

	go consumer.Consume(dependencies)

	h := handlers.Dependencies{DBConn: dependencies.DBConn, Auth: dependencies.Configs.Auth}
	if err := http.ListenAndServe(":8080", h.ApiHandler()); err != nil {
		log.Println(err)
	}
}
