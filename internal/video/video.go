package video

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
)

func CreateVideo(frames []image.Image, outputFile string) error {
	tmpDir, err := os.MkdirTemp("", "ca_frames")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	for i, frame := range frames {
		framePath := filepath.Join(tmpDir, fmt.Sprintf("frame_%04d.png", i))
		if err := saveFrame(frame, framePath); err != nil {
			return fmt.Errorf("failed to save frame %d: %v", i, err)
		}
	}

	return runFFmpeg(tmpDir, outputFile)
}

func saveFrame(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func runFFmpeg(frameDir, outputFile string) error {
	cmd := exec.Command("ffmpeg",
		"-framerate", "10",
		"-pattern_type", "glob",
		"-i", filepath.Join(frameDir, "frame_*.png"),
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-loglevel", "quiet",
		outputFile)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
