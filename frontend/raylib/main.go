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

func run(filename string) {
	rl.InitWindow(640, 320, "Vigo")
	rl.SetTargetFPS(60)

	go runCPU(filename)

	for !rl.WindowShouldClose() {
		pollKeyDown()
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

	cpu := vigo.NewCPU(&display, &keyDown)
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

var KEYMAP = [0x10]int32{
	rl.KeyX,     // 0
	rl.KeyOne,   // 1
	rl.KeyTwo,   // 2
	rl.KeyThree, // 3
	rl.KeyQ,     // 4
	rl.KeyW,     // 5
	rl.KeyE,     // 6
	rl.KeyA,     // 7
	rl.KeyS,     // 8
	rl.KeyD,     // 9
	rl.KeyZ,     // A
	rl.KeyC,     // B
	rl.KeyFour,  // C
	rl.KeyR,     // D
	rl.KeyF,     // E
	rl.KeyV,     // F
}

func pollKeyDown() {
	for k := 0; k <= 0xF; k++ {
		keyDown[k] = rl.IsKeyDown(KEYMAP[k])
	}
}
