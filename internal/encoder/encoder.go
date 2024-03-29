package encoder

import (
	"fmt"
	"os/exec"
	"time"
)

func EncodeVideo(input string, target int) (string, error) {
	var key = fmt.Sprintf("%v-%v-output.mp4", time.Now().Unix(), target)

	cmd := exec.Command("ffmpeg", "-i", input, "-vf", fmt.Sprintf("scale=-2:%v", target),
		"-c:v", "libx264", "-crf", "18", "-preset", "veryslow",
		"-c:a", "copy", fmt.Sprintf("./files/%s", key))

	cmd.Stdout = nil

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return key, nil
}
