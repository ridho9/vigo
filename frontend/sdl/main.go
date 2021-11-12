package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ridho9/vigo"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	display vigo.Display
	keyDown [16]bool
)

var KEYMAP = [0x10]int{
	sdl.SCANCODE_X, // 0
	sdl.SCANCODE_1, // 1
	sdl.SCANCODE_2, // 2
	sdl.SCANCODE_3, // 3
	sdl.SCANCODE_Q, // 4
	sdl.SCANCODE_W, // 5
	sdl.SCANCODE_E, // 6
	sdl.SCANCODE_A, // 7
	sdl.SCANCODE_S, // 8
	sdl.SCANCODE_D, // 9
	sdl.SCANCODE_Z, // A
	sdl.SCANCODE_C, // B
	sdl.SCANCODE_4, // C
	sdl.SCANCODE_R, // D
	sdl.SCANCODE_F, // E
	sdl.SCANCODE_V, // F
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("err missing rom file")
		os.Exit(1)
	}
	filename := os.Args[1]

	setupWindow()
	defer destroyWindow()

	go runCPU(filename)

	for {
		// start := time.Now()
		stop := pollEvent()
		if stop {
			break
		}
		update()
		// fmt.Println("UPDATE TIME", time.Since(start))
		draw()
		sdl.Delay(8)
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
	cpu.LoadRom(rom)

	cpu.Run()
	fmt.Println("elapsed time", time.Since(startTime))
}

func update() {
	curKeyState := sdl.GetKeyboardState()
	for i := 0; i <= 0xF; i += 1 {
		keyDown[i] = curKeyState[KEYMAP[i]] == 1
	}
}

func draw() {
	renderer.SetDrawColor(0, 0, 0, 0)
	renderer.Clear()
	renderer.SetDrawColor(255, 255, 255, 0)
	for y := int32(0); y < 32; y++ {
		for x := int32(0); x < 64; x++ {
			if display[x][y] {
				renderer.FillRect(&screenPixel[x][y])
			}
		}
	}
	renderer.Present()
}
