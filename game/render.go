package game

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func (body *GameBody) Render(g *Game, win *pixelgl.Window, imd *imdraw.IMDraw) {
	pos := body.Body.GetPosition()
	pos.OperatorScalarMulInplace(float64(Scale))
	m := pixel.IM.Rotated(pixel.V(0, 0), body.Body.GetAngle())
	m = m.Moved(pixel.V(pos.X-g.camera.X*Scale, pos.Y))
	imd.SetMatrix(m)
	imd.EndShape = imdraw.RoundEndShape

	switch body.Shape {
	case Rectangle:
		imd.Color = colornames.Blueviolet
		p1 := pixel.V(-body.HalfW*float64(Scale), -body.HalfH*float64(Scale))
		p2 := pixel.V(body.HalfW*float64(Scale), body.HalfH*float64(Scale))
		imd.Push(p1, p2)
		imd.Rectangle(3)

		if body.IsSelected {
			p1.X -= 3
			p1.Y -= 3
			p2.X += 3
			p2.Y += 3
			imd.Color = colornames.Darkorange
			imd.Push(p1, p2)
			imd.Rectangle(3)
		}
	case Circle:
		imd.Color = colornames.Brown
		imd.Push(pixel.V(0, 0))
		imd.Circle(body.Radius*Scale, 3)
		imd.Push(pixel.V(0, 0), pixel.V(0, body.Radius*Scale))
		imd.Line(3)
	}
}

func (g *Game) Draw(win *pixelgl.Window, imd *imdraw.IMDraw) {
	if g.toggleGrid {
		DrawGrid(imd)
	}

	g.goalBody.Render(g, win, imd)

	g.scoreText.Draw(g.Window, pixel.IM.Scaled(g.text.Orig, 3))

	// Render bodies
	for i := 0; i < len(g.Bodies); i++ {
		g.Bodies[i].Render(g, win, imd)
	}

	for i := 0; i < len(g.CargoBodies); i++ {
		g.CargoBodies[i].Render(g, win, imd)
	}

	// Render Car
	g.car.body.Render(g, win, imd)
	g.car.wheel1.Render(g, win, imd)
	g.car.wheel2.Render(g, win, imd)

	// Render floor
	g.ground.Render(g, win, imd)
	pos := g.ground.Body.GetPosition()
	startX := (pos.X-g.ground.HalfW)*Scale + g.groundSprite.Frame().W()/2
	y := (pos.Y+g.ground.HalfH)*Scale - g.groundSprite.Frame().H()/2
	spriteW := g.groundSprite.Frame().W()
	repeat := int(g.ground.HalfW * 2 * Scale / spriteW)

	for i := 0; i < repeat; i++ {
		g.groundSprite.Draw(win, pixel.IM.Moved(pixel.V(startX+float64(i)*spriteW-g.camera.X*Scale, y)))
	}

	g.states.Top().Render(g)
}

func DrawGrid(imd *imdraw.IMDraw) {
	imd.SetMatrix(pixel.IM)
	for i := 1; i < 10; i++ {
		imd.Push(pixel.V(float64(i)*Scale, 0), pixel.V(float64(i)*Scale, 7.5*Scale))
		imd.Line(3)
	}
	for i := 1; i < 8; i++ {
		imd.Push(pixel.V(0, float64(i)*Scale), pixel.V(float64(10)*Scale, float64(i)*Scale))
		imd.Line(3)
	}
}
