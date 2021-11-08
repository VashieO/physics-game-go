package game

import (
	"fmt"
	"math"

	"github.com/faiface/pixel/pixelgl"
)

type EditModeStateStack struct {
	states []EditModeState
}

func (s *EditModeStateStack) Push(state EditModeState) {
	s.states = append(s.states, state)
}

func (s *EditModeStateStack) Pop() EditModeState {
	if s.isEmpty() {
		return nil
	}
	i := len(s.states) - 1
	state := s.states[i]
	s.states = append(s.states[:i])
	return state
}

func (s *EditModeStateStack) Top() EditModeState {
	if s.isEmpty() {
		return nil
	}
	return s.states[len(s.states)-1]
}

func (s *EditModeStateStack) isEmpty() bool {
	return len(s.states) == 0
}

type EditModeState interface {
	Update(g *Game)
}

type MainEditState struct{}
type PlacementState struct{}
type SelectedState struct {
	gBody *GameBody
}

func (state *MainEditState) Update(g *Game) {
	fmt.Println("MainState")

	if g.Window.JustPressed(pixelgl.MouseButton1) {
		body := HandleEditModeSelect(g)
		if body != nil {
			g.editStates.Pop()
			g.editStates.Push(&SelectedState{gBody: body})
			return
		}
	}

	// Press N to create new body
	if g.Window.JustPressed(pixelgl.KeyN) {
		handleEditCreateNew(g)
		g.editStates.Pop()
		g.editStates.Push(&PlacementState{})
	}
	// Press N to create new cargo body
	if g.Window.JustPressed(pixelgl.KeyC) {
		handleEditCreateNewCargo(g)
		g.editStates.Pop()
		g.editStates.Push(&PlacementState{})
	}
}

func (state *PlacementState) Update(g *Game) {
	fmt.Println("PlacementState")
	handleEditShape(g)

	// Press V to to swap between Box and Ball
	if g.Window.JustPressed(pixelgl.KeyV) {
		handleEditSwapShape(g)
	}

	// Add new body to world
	if g.Window.JustPressed(pixelgl.MouseButton1) {
		g.newBody.Body.GetFixtureList().SetSensor(false)
		if g.newBody.IsCargo {
			g.CargoBodies = append(g.CargoBodies, g.newBody)
		} else {
			g.Bodies = append(g.Bodies, g.newBody)
		}
		g.newBody = nil
		g.editStates.Pop()
		g.editStates.Push(&MainEditState{})
	}

	// Press Esc to cancel body placement
	if g.Window.JustPressed(pixelgl.KeyEscape) && g.newBody != nil {
		g.World.DestroyBody(g.newBody.Body)
		g.newBody = nil
		g.editStates.Pop()
		g.editStates.Push(&MainEditState{})
	}
}

func (state *SelectedState) Update(g *Game) {
	fmt.Println("SelectedState")
	// Delete body
	if g.Window.JustPressed(pixelgl.KeyDelete) {
		fmt.Println("Delete")
		index := 0
		for i := 0; i < len(g.Bodies); i++ {
			if g.Bodies[i].IsSelected {
				index = i
				g.World.DestroyBody(g.Bodies[i].Body)
				g.Bodies = append(g.Bodies[:index], g.Bodies[index+1:]...)
				break
			}
		}

		for i := 0; i < len(g.CargoBodies); i++ {
			if g.CargoBodies[i].IsSelected {
				index = i
				g.World.DestroyBody(g.CargoBodies[i].Body)
				g.CargoBodies = append(g.CargoBodies[:index], g.CargoBodies[index+1:]...)
				break
			}
		}
		g.editStates.Pop()
		g.editStates.Push(&MainEditState{})
	}

	if g.Window.Pressed(pixelgl.KeyRight) {
		p := state.gBody.Body.GetPosition()
		p.X += 0.01
		state.gBody.Body.SetTransform(p, state.gBody.Body.GetAngle())
	}
	if g.Window.Pressed(pixelgl.KeyLeft) {
		p := state.gBody.Body.GetPosition()
		p.X -= 0.01
		state.gBody.Body.SetTransform(p, state.gBody.Body.GetAngle())
	}
	if g.Window.Pressed(pixelgl.KeyUp) {
		p := state.gBody.Body.GetPosition()
		p.Y += 0.01
		state.gBody.Body.SetTransform(p, state.gBody.Body.GetAngle())
	}
	if g.Window.Pressed(pixelgl.KeyDown) {
		p := state.gBody.Body.GetPosition()
		p.Y -= 0.01
		state.gBody.Body.SetTransform(p, state.gBody.Body.GetAngle())
	}
	if g.Window.Pressed(pixelgl.KeyComma) {
		p := state.gBody.Body.GetPosition()
		a := state.gBody.Body.GetAngle()
		a += math.Pi / 160
		state.gBody.Body.SetTransform(p, a)
	}
	if g.Window.Pressed(pixelgl.KeyPeriod) {
		p := state.gBody.Body.GetPosition()
		a := state.gBody.Body.GetAngle()
		a -= math.Pi / 160
		state.gBody.Body.SetTransform(p, a)
	}

	if g.Window.JustPressed(pixelgl.MouseButton1) {
		body := HandleEditModeSelect(g)
		state.gBody = body
		if body == nil {
			g.editStates.Pop()
			g.editStates.Push(&MainEditState{})
		}
	}
}
