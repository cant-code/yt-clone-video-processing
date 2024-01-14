package encoder

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func EncodeVideo(input string) {
	cmd := exec.Command("ffmpeg", "-i", input, "-vf", "scale=-1:360", "-c:v", "libx264", "-crf", "18",
		"-preset", "veryslow",
		"-c:a", "copy", fmt.Sprintf("./files/%v-output.mp4", time.Now().Unix()))

	cmd.Stderr = log.Writer()

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
