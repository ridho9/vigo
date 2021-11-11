package main

import "github.com/veandco/go-sdl2/sdl"

var (
	window      *sdl.Window
	renderer    *sdl.Renderer
	screenPixel [64][32]sdl.Rect
)

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
