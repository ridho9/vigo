package main

import rl "github.com/gen2brain/raylib-go/raylib"

func initWindow() {
	rl.InitWindow(640, 320, "Vigo")
	rl.SetTargetFPS(120)
}

func deinitWindow() {
	rl.CloseWindow()
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
