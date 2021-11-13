package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_SAMPLES = 512
const SAMPLE_RATE = 44100
const MAX_SAMPLES_PER_UPDATE = SAMPLE_RATE / 30

var stream rl.AudioStream
var data []float32
var writeBuf []float32
var wavelength int
var readCursor = 0
var playSound bool

func initAudio() {
	rl.InitAudioDevice()
	stream = rl.LoadAudioStream(SAMPLE_RATE, 32, 1)
	rl.SetAudioStreamBufferSizeDefault(MAX_SAMPLES_PER_UPDATE)
	rl.SetAudioStreamVolume(stream, 0.3)

	writeBuf = make([]float32, MAX_SAMPLES_PER_UPDATE)
	data = make([]float32, MAX_SAMPLES)
	generateWaveData()
}

func generateWaveData() {
	frequency := 440

	// amount of sample needed for one wave of specified frequency over sample rate
	wavelength = SAMPLE_RATE / frequency
	if wavelength > (MAX_SAMPLES / 2) {
		wavelength = MAX_SAMPLES / 2
	}
	if wavelength < 1 {
		wavelength = 1
	}
	for i := 0; i < wavelength*2; i++ {
		rad := 2 * rl.Pi * float32(i) / float32(wavelength)
		data[i] = float32(math.Sin(float64(rad))) * 1
		fmt.Println(i, data[i])
	}
}

func updateAudio() {
	if playSound {
		rl.PlayAudioStream(stream)
	} else {
		rl.StopAudioStream(stream)
	}
	if rl.IsAudioStreamProcessed(stream) {
		updateWriteBuffer()
	}
}

func updateWriteBuffer() {
	writeCursor := 0

	for writeCursor < MAX_SAMPLES_PER_UPDATE {
		writeLength := MAX_SAMPLES_PER_UPDATE - writeCursor
		readLength := wavelength - readCursor
		if writeLength > readLength {
			writeLength = readLength
		}

		// write data[rc:rc+wl] to wb[wc:wc+wl]
		for offset := 0; offset < writeLength; offset += 1 {
			writeBuf[writeCursor+offset] = data[readCursor+offset]
		}

		// loop readcursor
		readCursor = (readCursor + writeLength) % wavelength

		writeCursor += writeLength
	}

	fmt.Println("WB LEN", len(writeBuf))
	rl.UpdateAudioStream(stream, writeBuf, MAX_SAMPLES_PER_UPDATE)
}

func deinitAudio() {
	defer rl.CloseAudioDevice()
	defer rl.UnloadAudioStream(stream)
}
