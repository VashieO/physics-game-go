package main

import (
	"github.com/VashieO/physics/game"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, game.ScreenWidth, game.ScreenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	gameObj := &game.Game{}
	gameObj.Initialize(win, imd)

	for !win.Closed() {
		imd.Clear()
		win.Clear(colornames.Aliceblue)
		gameObj.Update(win)
		gameObj.Draw(win, imd)
		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
