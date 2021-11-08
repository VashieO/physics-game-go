package game

import (
	"github.com/bytearena/box2d"
	"github.com/faiface/pixel/pixelgl"
)

func handleForce(g *Game) {
	// Force applying
	if g.Window.JustPressed(pixelgl.MouseButton1) && !g.EditMode {
		pos := g.Window.MousePosition()

		var bodies []*GameBody
		bodies = append(bodies, g.car.body)
		bodies = append(bodies, g.car.wheel1)
		bodies = append(bodies, g.car.wheel2)
		bodies = append(bodies, g.Bodies...)

		for i := 0; i < len(bodies); i++ {
			body := bodies[i].Body
			worldPos := screenToWorld(pos, g.camera)
			collided := body.GetFixtureList().TestPoint(worldPos)
			if collided {
				g.isDragging = true
				localPos := body.GetLocalPoint(worldPos)
				g.forceDrag = &ForceDrag{body, &localPos}
				break
			}
		}
	}

	// Force applying
	if g.isDragging && g.Window.JustReleased(pixelgl.MouseButton1) {
		g.isDragging = false
		worldPos := g.forceDrag.body.GetWorldPoint(*g.forceDrag.localPos)
		mass := g.forceDrag.body.GetMass()
		mousePos := g.Window.MousePosition()
		mouseWorld := screenToWorld(mousePos, g.camera)
		acc := box2d.B2Vec2Sub(mouseWorld, worldPos)
		force := box2d.B2Vec2MulScalar(100*mass, acc)
		g.forceDrag.body.ApplyForce(force, worldPos, true)
	}
}
