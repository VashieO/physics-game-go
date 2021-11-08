package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type GameStateStack struct {
	states []GameState
	game   *Game
}

func (s *GameStateStack) Push(state GameState) {
	s.states = append(s.states, state)
	state.Init(s.game)
}

func (s *GameStateStack) Pop() GameState {
	if s.isEmpty() {
		return nil
	}
	i := len(s.states) - 1
	state := s.states[i]
	s.states = append(s.states[:i])
	return state
}

func (s *GameStateStack) Top() GameState {
	if s.isEmpty() {
		return nil
	}
	return s.states[len(s.states)-1]
}

func (s *GameStateStack) isEmpty() bool {
	return len(s.states) == 0
}

type GameState interface {
	Init(g *Game)
	Update(g *Game)
	Render(g *Game)
}

type GameStartState struct {
	startTime time.Time
}
type PlayState struct{}
type FinishedState struct{}
type PauseState struct{}
type EditState struct{}
type RestartState struct{}
type LoadingState struct {
	levelInfo LevelInfo
}

func (state GameStartState) Init(g *Game) {
	g.text.Clear()
	fmt.Fprintln(g.text, "Normal mode")

	g.startText.Clear()
	fmt.Fprintln(g.startText, g.levelInfo.Name)
	fmt.Fprintln(g.startText, "Carry the payload to the finish line")

	g.sideText.Clear()
	fmt.Fprintln(g.sideText, "Accelerate with <- and -> keys")
	fmt.Fprintln(g.sideText, "Break with space")
	fmt.Fprintln(g.sideText, "Restart with Enter")
	fmt.Println("GameStartState")
}

func (state GameStartState) Update(g *Game) {
	now := time.Now()
	elapsed := now.Sub(state.startTime)
	if elapsed.Seconds() > 4 {
		g.states.Pop()
		g.states.Push(PlayState{})
		return
	}

	if elapsed.Seconds() > 2 {
		g.startText.Clear()
		ratio := float64((4 - elapsed.Seconds()) / 2.0)
		g.startText.Color = color.RGBA{0, 0, 0, uint8(255 * ratio)}
		fmt.Fprintln(g.startText, g.levelInfo.Name)
		fmt.Fprintln(g.startText, "Carry the payload to the finish line")
	}

	g.World.Step(TimeStep, VelocityIterations, PositionIterations)

	pos := g.car.body.Body.GetPosition()
	g.camera.X = pos.X - 5.0 // Follow car, 5.0 is half the screen

	handleCarControls(g)

	if g.Window.JustPressed(pixelgl.KeyM) {
		g.toggleGrid = !g.toggleGrid
	}

	if g.Window.JustPressed(pixelgl.KeyP) {
		g.states.Pop()
		g.states.Push(PauseState{})
	}

	if g.Window.JustPressed(pixelgl.KeyE) {
		g.states.Pop()
		g.states.Push(EditState{})
	}
	if g.Window.JustPressed(pixelgl.KeyEnter) {
		g.states.Pop()
		g.states.Push(RestartState{})
	}
	if g.Window.JustPressed(pixelgl.KeyL) {
		g.states.Pop()
		info := LevelInfo{Name: "New level", Filename: "newLevel.json"}
		g.states.Push(LoadingState{levelInfo: info})
	}
}

func (state GameStartState) Render(g *Game) {
	g.text.Draw(g.Window, pixel.IM.Scaled(g.text.Orig, 2))
	g.sideText.Draw(g.Window, pixel.IM.Scaled(g.text.Orig, 1))
	g.startText.Draw(g.Window, pixel.IM.Scaled(g.text.Orig, 3))
}

func (state PlayState) Init(g *Game) {
	g.text.Clear()
	fmt.Fprintln(g.text, "Normal mode")
	fmt.Println("Playstate")

	g.sideText.Clear()
	fmt.Fprintln(g.sideText, "Accelerate with <- and -> keys")
	fmt.Fprintln(g.sideText, "Break with space")
	fmt.Fprintln(g.sideText, "Restart with Enter")
}

func (state PlayState) Update(g *Game) {
	if g.checkGoal() {
		g.states.Pop()
		g.states.Push(FinishedState{})
		return
	}
	g.World.Step(TimeStep, VelocityIterations, PositionIterations)

	pos := g.car.body.Body.GetPosition()
	g.camera.X = pos.X - 5.0 // Follow car, 5.0 is half the screen

	handleCarControls(g)

	if g.Window.JustPressed(pixelgl.KeyM) {
		g.toggleGrid = !g.toggleGrid
	}

	if g.Window.JustPressed(pixelgl.KeyP) {
		g.states.Pop()
		g.states.Push(PauseState{})
	}

	if g.Window.JustPressed(pixelgl.KeyE) {
		g.states.Pop()
		g.states.Push(EditState{})
	}

	if g.Window.JustPressed(pixelgl.KeyEnter) {
		g.states.Pop()
		g.states.Push(RestartState{})
	}

	if g.Window.JustPressed(pixelgl.KeyL) {
		g.states.Pop()
		info := LevelInfo{Name: "New level", Filename: "newLevel.json"}
		g.states.Push(LoadingState{levelInfo: info})
	}

	handleForce(g)
}

func (state PlayState) Render(g *Game) {
	g.text.Draw(g.Window, pixel.IM.Scaled(g.text.Orig, 2))
	g.sideText.Draw(g.Window, pixel.IM.Scaled(g.text.Orig, 1))

	if g.isDragging {
		g.imDraw.SetMatrix(pixel.IM)
		g.imDraw.Color = colornames.Blueviolet
		worldPos := g.forceDrag.body.GetWorldPoint(*g.forceDrag.localPos)
		worldPos.X = worldPos.X - g.camera.X
		screenPos := worldToScreen(&worldPos, g.camera)
		g.imDraw.Push(*screenPos, g.Window.MousePosition())
		g.imDraw.Line(3)
	}
}

func (state PauseState) Init(g *Game) {
	fmt.Println("PauseState")
	g.text.Clear()
	fmt.Fprintln(g.text, "Paused")

	g.sideText.Clear()
	fmt.Fprintln(g.sideText, "Accelerate with <- and -> keys")
	fmt.Fprintln(g.sideText, "Break with space")
}

func (state PauseState) Update(g *Game) {
	if g.Window.JustPressed(pixelgl.KeyP) {
		g.states.Pop()
		g.states.Push(PlayState{})
	}
}

func (state PauseState) Render(g *Game) {
	g.text.Draw(g.Window, pixel.IM.Scaled(g.text.Orig, 2))
	g.sideText.Draw(g.Window, pixel.IM.Scaled(g.text.Orig, 1))
}

func (state EditState) Init(g *Game) {
	g.text.Clear()
	fmt.Fprintln(g.text, "Edit mode")
	fmt.Println("EditState")

	g.sideText.Clear()
	fmt.Fprintln(g.sideText, "Press N for new body")
	fmt.Fprintln(g.sideText, "Press Esc to cancel placement")
	fmt.Fprintln(g.sideText, "Press S to save")
	fmt.Fprintln(g.sideText, "Arrows to change size")
	fmt.Fprintln(g.sideText, "A and D to use camera")

	g.editStates.Push(&MainEditState{})
}

func (state EditState) Update(g *Game) {
	if g.Window.JustPressed(pixelgl.KeyE) {
		for g.editStates.isEmpty() {
			g.editStates.Pop()
		}

		g.states.Pop()
		g.states.Push(PlayState{})
		return
	}

	g.editStates.Top().Update(g)

	g.infoText.Clear()
	if g.newBody != nil {
		pos := g.newBody.Body.GetPosition()
		angle := g.newBody.Body.GetAngle()
		fmt.Fprintln(g.infoText, "Body info")
		fmt.Fprintf(g.infoText, "X: %.2f\n", pos.X)
		fmt.Fprintf(g.infoText, "Y: %.2f\n", pos.Y)
		fmt.Fprintf(g.infoText, "Angle: %.2f\n", angle)
		fmt.Fprintf(g.infoText, "Width: %.2f\n", g.newBody.HalfW*2)
		fmt.Fprintf(g.infoText, "Height: %.2f\n", g.newBody.HalfH*2)
		fmt.Fprintf(g.infoText, "Radius: %.2f\n", g.newBody.Radius)
		fmt.Fprintf(g.infoText, "Density: %.1f\n", g.newBody.Density)
		fmt.Fprintf(g.infoText, "Friction: %.1f\n", g.newBody.Friction)
	} else {
		for i := 0; i < len(g.Bodies); i++ {
			if g.Bodies[i].IsSelected {
				body := g.Bodies[i]
				pos := body.Body.GetPosition()
				angle := body.Body.GetAngle()
				fmt.Fprintln(g.infoText, "Body info")
				fmt.Fprintf(g.infoText, "X: %.2f\n", pos.X)
				fmt.Fprintf(g.infoText, "Y: %.2f\n", pos.Y)
				fmt.Fprintf(g.infoText, "Angle: %.2f\n", angle)
				fmt.Fprintf(g.infoText, "Width: %.2f\n", body.HalfW*2)
				fmt.Fprintf(g.infoText, "Height: %.2f\n", body.HalfH*2)
				fmt.Fprintf(g.infoText, "Radius: %.2f\n", body.Radius)
				fmt.Fprintf(g.infoText, "Density %.1f\n", body.Density)
				fmt.Fprintf(g.infoText, "Friction %.1f\n", body.Friction)
			}
		}
	}

	mousePos := g.Window.MousePosition()
	pos := screenToWorld(mousePos, g.camera)
	if g.newBody != nil {
		g.newBody.Body.SetTransform(pos, g.newBody.Body.GetAngle())
	}

	if g.Window.Pressed(pixelgl.KeyD) {
		g.camera.X += 0.1
	}
	if g.Window.Pressed(pixelgl.KeyA) {
		g.camera.X -= 0.1
	}

	handleEditMode(g)
}

func (state EditState) Render(g *Game) {
	g.text.Draw(g.Window, pixel.IM.Scaled(g.text.Orig, 2))
	g.sideText.Draw(g.Window, pixel.IM.Scaled(g.text.Orig, 1))
	g.infoText.Draw(g.Window, pixel.IM.Scaled(g.text.Orig, 1))

	if g.newBody != nil {
		g.newBody.Render(g, g.Window, g.imDraw)
	}

}

func (state FinishedState) Init(g *Game) {
	g.levelIndex += 1
	g.score += g.CalcScore()
	g.finishedText.Clear()
	fmt.Fprintln(g.finishedText, "Congrats you reached the goal")
	fmt.Fprintf(g.finishedText, "Score: %d\n", g.score)

	if g.levelIndex < len(g.config.Levels) {
		fmt.Fprintln(g.finishedText, "Continue with Enter")
	} else {
		fmt.Fprintln(g.finishedText, "You have beaten the game")
	}

}

func (state FinishedState) Update(g *Game) {
	if g.Window.JustPressed(pixelgl.KeyEnter) {
		if g.levelIndex < len(g.config.Levels) {
			g.states.Pop()
			g.states.Push(LoadingState{g.config.Levels[g.levelIndex]})
		}
	}
}

func (state FinishedState) Render(g *Game) {
	g.text.Draw(g.Window, pixel.IM.Scaled(g.text.Orig, 2))
	g.finishedText.Draw(g.Window, pixel.IM.Scaled(g.text.Orig, 3))
}

func (state RestartState) Init(g *Game) {
	handleRestart(g)
}

func (state RestartState) Update(g *Game) {
	g.states.Pop()
	g.states.Push(PlayState{})
}

func (state RestartState) Render(g *Game) {
}

func handleRestart(g *Game) {
	for i := 0; i < len(g.Bodies); i++ {
		g.World.DestroyBody(g.Bodies[i].Body)
	}
	for i := 0; i < len(g.CargoBodies); i++ {
		g.World.DestroyBody(g.CargoBodies[i].Body)
	}
	g.Bodies = CreateBodies(g, g.levelData.Bodies)
	g.CargoBodies = CreateBodies(g, g.levelData.Cargo)
	g.resetCar()
}

func DestroyWorld(g *Game) {
	if g.ground != nil {
		g.World.DestroyBody(g.ground.Body)
	}
	if g.goalBody != nil {
		g.World.DestroyBody(g.goalBody.Body)
	}
	if g.car != nil {
		if g.car.wheelJoint1 != nil {
			g.World.DestroyJoint(g.car.wheelJoint1)
		}
		if g.car.wheelJoint2 != nil {
			g.World.DestroyJoint(g.car.wheelJoint2)
		}
		if g.car.wheel1.Body != nil {
			g.World.DestroyBody(g.car.wheel1.Body)
		}
		if g.car.wheel2.Body != nil {
			g.World.DestroyBody(g.car.wheel2.Body)
		}
		if g.car.body.Body != nil {
			g.World.DestroyBody(g.car.body.Body)
		}
	}

	for i := 0; i < len(g.Bodies); i++ {
		g.World.DestroyBody(g.Bodies[i].Body)
	}
	for i := 0; i < len(g.CargoBodies); i++ {
		g.World.DestroyBody(g.CargoBodies[i].Body)
	}
}

func (state LoadingState) Init(g *Game) {
	data := LoadFromFile(state.levelInfo.Filename)
	g.levelData = data
	g.levelInfo = &state.levelInfo

	DestroyWorld(g)

	g.ground, g.goalBody = CreateGroundAndGoal(g)
	g.car = CreateCar(g)
	g.Bodies = CreateBodies(g, data.Bodies)
	g.CargoBodies = CreateBodies(g, data.Cargo)
}

func (state LoadingState) Update(g *Game) {
	g.states.Pop()
	g.states.Push(GameStartState{startTime: time.Now()})
}

func (state LoadingState) Render(g *Game) {
}
