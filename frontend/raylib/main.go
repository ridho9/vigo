package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ridho9/vigo"
)

var display *vigo.Display

func main() {
	if len(os.Args) < 2 {
		fmt.Println("err missing rom file")
		os.Exit(1)
	}
	filename := os.Args[1]
	display = &vigo.Display{}

	run(filename)
}

func run(filename string) {
	rl.InitWindow(640, 320, "raylib [core] example - basic window")
	rl.SetTargetFPS(60)

	go runCPU(filename)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		drawFrame()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func runCPU(filename string) {
	startTime := time.Now()
	rom, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("err ", err)
		os.Exit(2)
	}

	cpu := vigo.NewCPU(display)
	cpu.LoadRom(rom)

	cpu.Run()
	fmt.Println("elapsed time", time.Since(startTime))
}

func drawFrame() {
	rl.ClearBackground(rl.Black)
	for y := int32(0); y < 32; y++ {
		for x := int32(0); x < 64; x++ {
			if display[x][y] {
				rl.DrawRectangle(x*10, y*10, 10, 10, rl.White)
			}
		}
	}
}
