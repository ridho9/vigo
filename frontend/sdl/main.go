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
	window      *sdl.Window
	renderer    *sdl.Renderer
	screenPixel [64][32]sdl.Rect

	display vigo.Display
	keyDown [16]bool

	COLOR_WHITE uint32
)

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
		start := time.Now()
		stop := pollEvent()
		if stop {
			break
		}
		update()
		draw()
		fmt.Println("looptime", time.Since(start))
		sdl.Delay(16)
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

func setupWindow() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	window, err = sdl.CreateWindow("Vigo", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 640, 320, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	for y := int32(0); y < 32; y += 1 {
		for x := int32(0); x < 64; x += 1 {
			screenPixel[x][y] = sdl.Rect{X: x * 10, Y: y * 10, W: 10, H: 10}
		}
	}
}

func pollEvent() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			return true
		}
	}
	return false
}

func destroyWindow() {
	defer sdl.Quit()
	defer window.Destroy()
	defer renderer.Destroy()
}
