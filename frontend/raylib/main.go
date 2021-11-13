package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ridho9/vigo"
)

var display vigo.Display
var keyDown [16]bool

func main() {
	if len(os.Args) < 2 {
		fmt.Println("err missing rom file")
		os.Exit(1)
	}
	filename := os.Args[1]

	run(filename)
}

func initt() {
	initWindow()
	initAudio()
}

func deinit() {
	defer deinitWindow()
	defer deinitAudio()
}

func run(filename string) {
	initt()
	defer deinit()

	go runCPU(filename)

	for !rl.WindowShouldClose() {
		start := time.Now()

		updateAudio()
		fmt.Println("UPDATE AUDIO", time.Since(start))

		pollKeyDown()
		fmt.Println("POLLING", time.Since(start))

		rl.BeginDrawing()
		fmt.Println("BEGIN DRAW", time.Since(start))

		drawFrame()
		fmt.Println("FINISH DRAW FRAME", time.Since(start))

		rl.EndDrawing()
		fmt.Println("END DRAW", time.Since(start))
	}
}

func runCPU(filename string) {
	startTime := time.Now()
	rom, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("err ", err)
		os.Exit(2)
	}

	cpu := vigo.NewCPU(&display, &keyDown)
	cpu.SetSoundFlag(&playSound)
	cpu.LoadRom(rom)

	cpu.Run()
	fmt.Println("elapsed time", time.Since(startTime))
}
