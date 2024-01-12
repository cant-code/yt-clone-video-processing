package main

import (
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("ffmpeg", "-i", "input.mp4", "-vf", "scale=-1:360", "-c:v", "libx264", "-crf", "18", "-preset", "veryslow", "-c:a", "copy", "output.mp4")

	cmd.Stderr = log.Writer()

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
