package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ridho9/vigo"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("err missing rom file")
		os.Exit(1)
	}
	filename := os.Args[1]
	rom, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("err ", err)
		os.Exit(2)
	}

	cpu := vigo.NewCPU()
	cpu.LoadRom(rom)

	cpu.SetDisplayCallback(func(d vigo.Display) {
		fmt.Println("")
		for y := 0; y < 32; y++ {
			for x := 0; x < 64; x++ {
				if d[x][y] {
					fmt.Print("#")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Println("|")
		}
	})

	cpu.Run()
	// for {
	// 	cpu.Step()
	// }
}
