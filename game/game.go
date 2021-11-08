package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/bytearena/box2d"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type PlayerTurn int

const (
	PlayerOne PlayerTurn = 0
	PlayerTwo PlayerTurn = 1
)

const TimeStep = 1.0 / 60.0
const VelocityIterations = 8
const PositionIterations = 3

const (
	ScreenWidth          = 1200
	ScreenHeight         = 900
	Scale        float64 = 1200.0 / 10.0
	InvScale     float64 = 1 / Scale
)

type Shape int

const (
	Rectangle Shape = 0
	Circle    Shape = 1
)

type BoxDef struct {
	x        float64
	y        float64
	hx       float64
	hy       float64
	density  float64
	friction float64
	isSensor bool
	bodyType uint8
}

type BallDef struct {
	x        float64
	y        float64
	r        float64
	density  float64
	friction float64
	isSensor bool
	bodyType uint8
}

const BoxMode = 0
const BallMode = 1

type PlaceMode int

type Game struct {
	World        *box2d.B2World
	Window       *pixelgl.Window
	imDraw       *imdraw.IMDraw
	camera       *Camera
	ground       *GameBody
	Bodies       []*GameBody
	CargoBodies  []*GameBody
	score        int
	forceDrag    *ForceDrag
	EditMode     bool
	isDragging   bool
	toggleGrid   bool
	car          *Car
	groundSprite *pixel.Sprite
	newBody      *GameBody
	text         *text.Text
	sideText     *text.Text
	infoText     *text.Text
	finishedText *text.Text
	startText    *text.Text
	scoreText    *text.Text
	goalBody     *GameBody
	states       GameStateStack
	editStates   EditModeStateStack
	config       *ConfigData
	levelData    *LevelData
	levelInfo    *LevelInfo
	levelIndex   int
	placeMode    PlaceMode
}

type Camera struct {
	X float64
	Y float64
}

type GameBody struct {
	Body       *box2d.B2Body
	HalfW      float64
	HalfH      float64
	Radius     float64
	Density    float64
	Friction   float64
	Shape      Shape
	IsSelected bool
	IsCargo    bool
}

type Car struct {
	carAcc      float64
	body        *GameBody
	wheel1      *GameBody
	wheel2      *GameBody
	wheelJoint1 *box2d.B2WheelJoint
	wheelJoint2 *box2d.B2WheelJoint
}

type ForceDrag struct {
	body     *box2d.B2Body
	localPos *box2d.B2Vec2
}

func (car *Car) Forward() {
	car.carAcc = -20.0
	car.wheel1.Body.SetAngularDamping(0.0)
	car.wheel2.Body.SetAngularDamping(0.0)
	car.wheelJoint1.EnableMotor(true)
	car.wheelJoint1.SetMotorSpeed(car.carAcc)
	car.wheelJoint2.EnableMotor(true)
	car.wheelJoint2.SetMotorSpeed(car.carAcc)
}

func (car *Car) Backwards() {
	car.carAcc = 20.0
	car.wheel1.Body.SetAngularDamping(0.0)
	car.wheel2.Body.SetAngularDamping(0.0)
	car.wheelJoint1.EnableMotor(true)
	car.wheelJoint1.SetMotorSpeed(car.carAcc)
	car.wheelJoint2.EnableMotor(true)
	car.wheelJoint2.SetMotorSpeed(car.carAcc)
}

func (car *Car) Stop() {
	car.carAcc = 0.0
	car.wheel1.Body.SetAngularDamping(1.0)
	car.wheel2.Body.SetAngularDamping(1.0)
	car.wheelJoint1.EnableMotor(false)
	car.wheelJoint1.SetMotorSpeed(car.carAcc)
	car.wheelJoint2.EnableMotor(false)
	car.wheelJoint2.SetMotorSpeed(car.carAcc)
}

func (car *Car) Break() {
	car.carAcc = 0.0
	car.wheel1.Body.SetAngularDamping(1.0)
	car.wheel2.Body.SetAngularDamping(1.0)
	car.wheelJoint1.EnableMotor(true)
	car.wheelJoint1.SetMotorSpeed(-car.wheelJoint1.GetJointAngularSpeed())
	car.wheelJoint2.EnableMotor(true)
	car.wheelJoint2.SetMotorSpeed(-car.wheelJoint2.GetJointAngularSpeed())

	// angularVel := car.wheel1.GetAngularVelocity()
	// fmt.Println(car.wheelJoint1.GetJointAngularSpeed())
}

func (g *Game) resetCar() {
	g.car.body.Body.SetTransform(box2d.B2Vec2{X: 3.5, Y: 1.5}, 0)
	g.car.wheel1.Body.SetTransform(box2d.B2Vec2{X: 2.7, Y: 1.3}, 0)
	g.car.wheel2.Body.SetTransform(box2d.B2Vec2{X: 4.3, Y: 1.5}, 0)
	g.car.body.Body.SetLinearVelocity(box2d.B2Vec2{X: 0, Y: 0})
	g.car.body.Body.SetAngularVelocity(0)
	g.car.wheel1.Body.SetLinearVelocity(box2d.B2Vec2{X: 0, Y: 0})
	g.car.wheel1.Body.SetAngularVelocity(0)
	g.car.wheel2.Body.SetLinearVelocity(box2d.B2Vec2{X: 0, Y: 0})
	g.car.wheel2.Body.SetAngularVelocity(0)
}

func (g *Game) Initialize(win *pixelgl.Window, imd *imdraw.IMDraw) {
	g.Window = win
	g.imDraw = imd
	g.states.game = g
	g.camera = &Camera{}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 870), basicAtlas)
	basicTxt.Color = colornames.Black
	g.text = basicTxt
	fmt.Fprintln(basicTxt, "")

	sideTxt := text.New(pixel.V(980, 800), basicAtlas)
	sideTxt.Color = colornames.Black
	g.sideText = sideTxt

	infoTxt := text.New(pixel.V(980, 600), basicAtlas)
	infoTxt.Color = colornames.Black
	g.infoText = infoTxt
	fmt.Fprintln(infoTxt, "")

	finishedTxt := text.New(pixel.V(160, 780), basicAtlas)
	finishedTxt.Color = colornames.Black
	g.finishedText = finishedTxt
	fmt.Fprintln(finishedTxt, "Congrats you reached the goal")

	startTxt := text.New(pixel.V(140, 730), basicAtlas)
	startTxt.Color = color.RGBA{0, 0, 0, 255}
	g.startText = startTxt
	fmt.Fprintln(startTxt, "Carry the payload to the finish line")

	g.scoreText = text.New(pixel.V(380, 860), basicAtlas)
	g.scoreText.Color = colornames.Black
	fmt.Fprintf(g.scoreText, "Score: %d", g.score)

	// Create world
	world := box2d.MakeB2World(box2d.B2Vec2{X: 0.0, Y: -3.0})
	g.World = &world

	path := "./resources/grassLongPlatform.png"
	picture, err := loadPicture(path)
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(picture, pixel.R(0, 0, 300, 100))
	g.groundSprite = sprite

	// Load first level
	g.config = LoadConfig()
	g.states.Push(LoadingState{levelInfo: g.config.Levels[0]})
}

func (g *Game) Update(win *pixelgl.Window) error {
	g.scoreText.Clear()
	fmt.Fprintf(g.scoreText, "Score: %d", g.score)

	g.states.Top().Update(g)

	handleInput(g, win)

	return nil
}

func (g *Game) checkGoal() bool {
	goalPos := g.goalBody.Body.GetPosition()
	carPos := g.car.body.Body.GetPosition()

	// Back of car and a bit extra
	if carPos.X-g.car.body.HalfW-0.3 > goalPos.X {
		return true
	}
	return false
}

func (g *Game) CalcScore() int {
	score := 0
	goalPos := g.goalBody.Body.GetPosition()
	for i := 0; i < len(g.CargoBodies); i++ {
		body := g.CargoBodies[i]
		pos := body.Body.GetPosition()
		if pos.X > goalPos.X {
			score += int(4000 * body.HalfW * body.HalfH)
		}
	}
	return score
}

func handleInput(g *Game, win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyM) {
		g.toggleGrid = !g.toggleGrid
	}
}

func checkCollision(body *box2d.B2Body, pos pixel.Vec, cam *Camera) bool {
	worldPos := screenToWorld(pos, cam)
	return body.GetFixtureList().TestPoint(worldPos)
}

func HandleEditModeSelect(g *Game) *GameBody {
	// Select body we clicked on
	pos := g.Window.MousePosition()
	bodies := append(g.Bodies[:], g.CargoBodies...)

	for i := 0; i < len(bodies); i++ {
		bodies[i].IsSelected = false
	}

	for i := 0; i < len(bodies); i++ {
		worldPos := box2d.B2Vec2{X: pos.X / Scale, Y: pos.Y / Scale}
		worldPos.X = worldPos.X + g.camera.X
		collided := bodies[i].Body.GetFixtureList().TestPoint(worldPos)
		if collided {
			bodies[i].IsSelected = true
			return bodies[i]
		}
	}
	return nil
}

func handleEditCreateNew(g *Game) *GameBody {
	mousePos := g.Window.MousePosition()
	worldPos := screenToWorld(mousePos, g.camera)

	fmt.Println(g.placeMode)
	if g.newBody == nil {
		if g.placeMode == BoxMode {
			boxDef := BoxDef{x: worldPos.X, y: worldPos.Y, hx: 1, hy: 0.2, density: 1.0, friction: 0.8, isSensor: true}
			boxDef.bodyType = box2d.B2BodyType.B2_staticBody
			g.newBody = createBox(boxDef, g.World)
			fmt.Println("Created box")
		} else {
			ballDef := BallDef{x: worldPos.X, y: worldPos.Y, r: 0.2, density: 1.0, friction: 0.8, isSensor: true}
			ballDef.bodyType = box2d.B2BodyType.B2_staticBody
			g.newBody = createBall(ballDef, g.World)
			fmt.Println("Created ball")
		}
		return g.newBody
	}

	return nil
}

func handleEditCreateNewCargo(g *Game) *GameBody {
	mousePos := g.Window.MousePosition()
	worldPos := screenToWorld(mousePos, g.camera)
	boxDef := BoxDef{x: worldPos.X, y: worldPos.Y, hx: 0.2, hy: 0.2, density: 1.0, friction: 1.0, isSensor: true}
	boxDef.bodyType = box2d.B2BodyType.B2_dynamicBody
	if g.newBody == nil {
		g.newBody = createBox(boxDef, g.World)
		g.newBody.IsCargo = true
		return g.newBody
	}

	return nil
}

func handleEditShape(g *Game) {
	needNewShape := false
	if g.placeMode == BoxMode {
		if g.Window.JustPressed(pixelgl.KeyRight) && g.newBody.HalfW+0.1 > 0 {
			g.newBody.HalfW += 0.1
			needNewShape = true
		} else if g.Window.JustPressed(pixelgl.KeyLeft) && g.newBody.HalfW-0.1 > 0 {
			g.newBody.HalfW -= 0.1
			needNewShape = true
		} else if g.Window.JustPressed(pixelgl.KeyUp) && g.newBody.HalfH+0.1 > 0 {
			g.newBody.HalfH += 0.1
			needNewShape = true
		} else if g.Window.JustPressed(pixelgl.KeyDown) && g.newBody.HalfH-0.1 > 0 {
			g.newBody.HalfH -= 0.1
			needNewShape = true
		}
	} else {
		if g.Window.JustPressed(pixelgl.KeyRight) && g.newBody.Radius+0.1 > 0 {
			g.newBody.Radius += 0.1
			needNewShape = true
		} else if g.Window.JustPressed(pixelgl.KeyLeft) && g.newBody.Radius+0.1 > 0 {
			g.newBody.Radius -= 0.1
			needNewShape = true
		}
	}
	if g.Window.Pressed(pixelgl.KeyPeriod) {
		angle := g.newBody.Body.GetAngle() - math.Pi/160
		g.newBody.Body.SetTransform(g.newBody.Body.GetPosition(), angle)
	} else if g.Window.Pressed(pixelgl.KeyComma) {
		angle := g.newBody.Body.GetAngle() + math.Pi/160
		g.newBody.Body.SetTransform(g.newBody.Body.GetPosition(), angle)
	} else if g.Window.JustPressed(pixelgl.KeyR) {
		g.newBody.Density -= 0.1
		needNewShape = true
	} else if g.Window.JustPressed(pixelgl.KeyT) {
		g.newBody.Density += 0.1
		needNewShape = true
	} else if g.Window.JustPressed(pixelgl.KeyF) {
		g.newBody.Friction -= 0.1
		needNewShape = true
	} else if g.Window.JustPressed(pixelgl.KeyG) {
		g.newBody.Friction += 0.1
		needNewShape = true
	}

	if needNewShape {
		shape := g.newBody.Body.GetFixtureList().GetShape()
		shape.Destroy()
		if g.placeMode == BoxMode {
			def := createBoxFixureDef(g.newBody.HalfW, g.newBody.HalfH, g.newBody.Density, g.newBody.Friction, true)
			g.newBody.Body.CreateFixtureFromDef(&def)
		} else {
			def := createBallFixureDef(g.newBody.Radius, g.newBody.Density, g.newBody.Friction, true)
			g.newBody.Body.CreateFixtureFromDef(&def)
		}
	}
}

func handleEditSwapShape(g *Game) {
	if g.placeMode == BoxMode {
		g.placeMode = BallMode
	} else {
		g.placeMode = BoxMode
	}

	if g.newBody != nil {
		mousePos := g.Window.MousePosition()
		worldPos := screenToWorld(mousePos, g.camera)
		g.World.DestroyBody(g.newBody.Body)

		if g.placeMode == BoxMode {
			boxDef := BoxDef{x: worldPos.X, y: worldPos.Y, hx: 1, hy: 0.2, density: 1.0, friction: 0.8, isSensor: true}
			boxDef.bodyType = box2d.B2BodyType.B2_staticBody
			g.newBody = createBox(boxDef, g.World)
			fmt.Println("Created box")
		} else {
			ballDef := BallDef{x: worldPos.X, y: worldPos.Y, r: 0.2, density: 1.0, friction: 0.8, isSensor: true}
			ballDef.bodyType = box2d.B2BodyType.B2_staticBody
			g.newBody = createBall(ballDef, g.World)
			fmt.Println("Created ball")
		}
	}
}

func handleEditMode(g *Game) {
	// Save to file
	if g.Window.JustPressed(pixelgl.KeyS) {
		SaveToFile(g)
	}
}

func handleCarControls(g *Game) {
	// Car controls
	if !g.Window.Pressed(pixelgl.KeyLeft) && !g.Window.Pressed(pixelgl.KeyLeft) {
		g.car.Stop()
	}
	if g.Window.Pressed(pixelgl.KeyLeft) {
		g.car.Backwards()
	}

	if g.Window.Pressed(pixelgl.KeyRight) {
		g.car.Forward()
	}

	if g.Window.Pressed(pixelgl.KeySpace) {
		g.car.Break()
	}

	if g.Window.JustPressed(pixelgl.Key1) {
		g.resetCar()
	}
}
