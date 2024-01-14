package encoder

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func EncodeVideo(input string, target int) {
	cmd := exec.Command("ffmpeg", "-i", input, "-vf", fmt.Sprintf("scale=-2:%v", target),
		"-c:v", "libx264", "-crf", "18", "-preset", "veryslow",
		"-c:a", "copy", fmt.Sprintf("./files/%v-%v-output.mp4", time.Now().Unix(), target))

	cmd.Stderr = log.Writer()

	err := cmd.Run()
	if err != nil {
		log.Panicln(err)
	}
}
