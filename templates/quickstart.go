package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

/*All `Game` code is suggested to be separeted from main.go*/
func main() {
	ebiten.SetWindowResizable(true)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func NewGame() *Game {
	return &Game{}
}

func (self *Game) Update() error {
	return nil
}

func (self *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.White)
}

func (self *Game) Layout(w, h int) (int, int) {
	f := ebiten.DeviceScaleFactor()
	return w * int(f), h * int(f)
}
