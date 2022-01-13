package ui

import (
	"log"
	"os/exec"
	"runtime"

	"azul3d.org/engine/gfx"
	"azul3d.org/engine/gfx/window"
	"azul3d.org/engine/keyboard"
	"github.com/fogleman/nes/nes"
	"github.com/gordonklaus/portaudio"
)

const (
	width  = 256
	height = 240
	scale  = 3
	title  = "NES"
)

var keystate [8]bool

func init() {

	// we need a parallel OS thread to avoid audio stuttering
	runtime.GOMAXPROCS(3)

	// we need to keep OpenGL calls on a single thread
	runtime.LockOSThread()

	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func gfxLoop(w window.Window, d gfx.Device) {
	defer w.Close()
	for {

		keystate[nes.ButtonA] = w.Keyboard().Down(keyboard.X)
		keystate[nes.ButtonB] = w.Keyboard().Down(keyboard.Y)
		keystate[nes.ButtonSelect] = w.Keyboard().Down(keyboard.LeftShift)
		keystate[nes.ButtonStart] = w.Keyboard().Down(keyboard.Enter)
		keystate[nes.ButtonUp] = w.Keyboard().Down(keyboard.ArrowUp)
		keystate[nes.ButtonDown] = w.Keyboard().Down(keyboard.ArrowDown)
		keystate[nes.ButtonLeft] = w.Keyboard().Down(keyboard.ArrowLeft)
		keystate[nes.ButtonRight] = w.Keyboard().Down(keyboard.ArrowRight)
	}
}

func Run(paths []string) {

	props := window.NewProps()
	props.Minimized()
	props.SetVisible(false)
	props.SetResizeRenderSync(false)
	go window.Run(gfxLoop, props)

	// initialize audio
	portaudio.Initialize()
	defer portaudio.Terminate()

	audio := NewAudio()
	if err := audio.Start(); err != nil {
		log.Fatalln(err)
	}
	defer audio.Stop()

	// run director
	director := NewDirector(audio)
	director.Start(paths)
}
