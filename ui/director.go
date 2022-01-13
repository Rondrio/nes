package ui

import (
	"log"

	"github.com/fogleman/nes/nes"
)

type View interface {
	Enter()
	Exit()
	Update()
}

type Director struct {
	audio     *Audio
	view      View
	menuView  View
	timestamp float64
}

func NewDirector(audio *Audio) *Director {
	director := Director{}
	director.audio = audio
	return &director
}

func (d *Director) SetView(view View) {
	if d.view != nil {
		d.view.Exit()
	}
	d.view = view
	if d.view != nil {
		d.view.Enter()
	}
}

func (d *Director) Step() {
	if d.view != nil {
		d.view.Update()
	}
}

func (d *Director) Start(paths []string) {
	d.PlayGame(paths[0])
	d.Run()
}

func (d *Director) Run() {
	for {
		d.Step()
	}
}

func (d *Director) PlayGame(path string) {
	hash, err := hashFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	console, err := nes.NewConsole(path)
	if err != nil {
		log.Fatalln(err)
	}
	d.SetView(NewGameView(d, console, path, hash))
}

func (d *Director) ShowMenu() {
	d.SetView(d.menuView)
}
