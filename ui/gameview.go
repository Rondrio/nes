package ui

import (
	"image"

	"github.com/fogleman/nes/nes"
)

const padding = 0

type GameView struct {
	director *Director
	console  *nes.Console
	title    string
	hash     string
	record   bool
	frames   []image.Image
}

func NewGameView(director *Director, console *nes.Console, title, hash string) View {
	return &GameView{director, console, title, hash, false, nil}
}

func (view *GameView) Enter() {
	view.console.SetAudioChannel(view.director.audio.channel)
	view.console.SetAudioSampleRate(view.director.audio.sampleRate)
	// load state
	if err := view.console.LoadState(savePath(view.hash)); err == nil {
		return
	} else {
		view.console.Reset()
	}
	// load sram
	cartridge := view.console.Cartridge
	if cartridge.Battery != 0 {
		if sram, err := readSRAM(sramPath(view.hash)); err == nil {
			cartridge.SRAM = sram
		}
	}
}

func (view *GameView) Exit() {
	view.console.SetAudioChannel(nil)
	view.console.SetAudioSampleRate(0)
	// save sram
	cartridge := view.console.Cartridge
	if cartridge.Battery != 0 {
		writeSRAM(sramPath(view.hash), cartridge.SRAM)
	}
	// save state
	view.console.SaveState(savePath(view.hash))
}

func (view *GameView) Update() {
	console := view.console
	console.StepSeconds(0.0166)

	updateControllers(console)

	setTexture(console.Buffer())
	if view.record {
		view.frames = append(view.frames, copyImage(console.Buffer()))
	}
}

func updateControllers(console *nes.Console) {
	console.SetButtons1(keystate)
}
