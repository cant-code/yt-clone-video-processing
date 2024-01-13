package main

import (
	"log"
	"yt-clone-video-processing/internal/config"
	"yt-clone-video-processing/internal/consumer"
)

func main() {
	/*	cmd := exec.Command("ffmpeg", "-i", "input.mp4", "-vf", "scale=-1:360", "-c:v", "libx264", "-crf", "18", "-preset", "veryslow", "-c:a", "copy", "output.mp4")

		cmd.Stderr = log.Writer()

		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}*/

	loadConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	_, err = consumer.Consume(loadConfig)
	if err != nil {
		log.Fatal(err)
	}
}
