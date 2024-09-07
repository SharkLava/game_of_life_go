package main

import (
	"flag"
	"fmt"
	"github.com/SharkLava/game_of_life_go/internal/automaton"
	"github.com/SharkLava/game_of_life_go/internal/video"
)

func main() {
	size := flag.Int("size", 100, "Size of the cellular automata grid")
	steps := flag.Int("steps", 100, "Number of steps to run")
	saveVideo := flag.Bool("video", false, "Save output as video")
	flag.Parse()

	ca := automaton.New(*size, automaton.GameOfLifeRule, "moore")
	frames, err := ca.Run(*steps)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if *saveVideo {
		err = video.CreateVideo(frames, "cellular_automata.mp4")
		if err != nil {
			fmt.Println("Error creating video:", err)
			return
		}
		fmt.Println("Video saved as cellular_automata.mp4")
	} else {
		err = automaton.SaveImage(frames[len(frames)-1], "cellular_automata.png")
		if err != nil {
			fmt.Println("Error saving final image:", err)
			return
		}
		fmt.Println("Final state saved as cellular_automata.png")
	}
}
