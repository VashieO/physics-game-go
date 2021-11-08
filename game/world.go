package game

import "github.com/bytearena/box2d"

func CreateCar(g *Game) *Car {
	carBodyDef := BoxDef{x: 3.5, y: 1.3, hx: 1.3, hy: 0.2, density: 0.5, friction: 0.8}
	carBodyDef.bodyType = box2d.B2BodyType.B2_dynamicBody
	carBody := createBox(carBodyDef, g.World)

	wheelDef1 := BallDef{x: 2.6, y: 1.1, r: 0.3, density: 1.0, friction: 1.0}
	wheelDef1.bodyType = box2d.B2BodyType.B2_dynamicBody
	wheel1 := createBall(wheelDef1, g.World)

	wheelDef2 := BallDef{x: 4.4, y: 1.1, r: 0.3, density: 1.0, friction: 1.0}
	wheelDef2.bodyType = box2d.B2BodyType.B2_dynamicBody
	wheel2 := createBall(wheelDef2, g.World)

	motorDef := box2d.MakeB2WheelJointDef()
	motorDef.Initialize(carBody.Body, wheel1.Body, wheel1.Body.GetWorldCenter(), box2d.B2Vec2{X: 0, Y: 1})
	motorDef.MaxMotorTorque = 2
	motorDef.DampingRatio = 0.7
	motorDef.FrequencyHz = 4
	joint1 := g.World.CreateJoint(&motorDef)
	wheelJoint1, ok := joint1.(*box2d.B2WheelJoint)
	if !ok {
		panic("Could not convert joint")
	}

	motorDef2 := box2d.MakeB2WheelJointDef()
	motorDef2.Initialize(carBody.Body, wheel2.Body, wheel2.Body.GetWorldCenter(), box2d.B2Vec2{X: 0, Y: 1})
	motorDef2.MaxMotorTorque = 2
	motorDef2.DampingRatio = 0.7
	motorDef2.FrequencyHz = 4
	joint2 := g.World.CreateJoint(&motorDef2)
	wheelJoint2, ok := joint2.(*box2d.B2WheelJoint)

	if !ok {
		panic("Could not convert joint")
	}

	car := &Car{}
	car.body = carBody
	car.wheel1 = wheel1
	car.wheel2 = wheel2
	car.wheelJoint1 = wheelJoint1
	car.wheelJoint2 = wheelJoint2

	return car
}

func CreateGroundAndGoal(g *Game) (*GameBody, *GameBody) {
	groundDef := BoxDef{x: 35, y: 0.3, hx: 50, hy: 0.5, density: 1.0, friction: 0.8}
	ground := createBox(groundDef, g.World)

	goalDef := BoxDef{x: 50, y: 2.4, hx: 0.1, hy: 1.8, density: 1.0, friction: 0.8, isSensor: true}
	goal := createBox(goalDef, g.World)

	return ground, goal
}

func CreateBodies(g *Game, bodies []BodyJson) []*GameBody {
	// Create bodies
	var newBodies []*GameBody

	for i := 0; i < len(bodies); i++ {
		body := bodies[i]
		if body.BodyShape == Rectangle {
			boxDef := BoxDef{x: body.X, y: body.Y, hx: body.Hx, hy: body.Hy, density: body.Density, friction: body.Friction}
			boxDef.bodyType = body.BodyType
			box := createBox(boxDef, g.World)
			box.Body.SetTransform(box.Body.GetPosition(), body.Angle)
			newBodies = append(newBodies, box)
		} else if body.BodyShape == Circle {
			ballDef := BallDef{x: body.X, y: body.Y, r: body.Radius, density: body.Density, friction: body.Friction}
			ballDef.bodyType = body.BodyType
			ball := createBall(ballDef, g.World)
			ball.Body.SetTransform(ball.Body.GetPosition(), body.Angle)
			newBodies = append(newBodies, ball)
		}
	}
	return newBodies
}

func createBox(def BoxDef, world *box2d.B2World) *GameBody {
	boxDef := box2d.MakeB2BodyDef()
	boxDef.Type = def.bodyType

	boxDef.Position.Set(def.x, def.y)
	boxBody := world.CreateBody(&boxDef)

	boxBox := box2d.B2PolygonShape{}
	boxBox.SetAsBox(float64(def.hx), float64(def.hy))
	boxFixDef := box2d.MakeB2FixtureDef()
	boxFixDef.Shape = &boxBox
	boxFixDef.Density = def.density
	boxFixDef.Friction = def.friction
	boxFixDef.IsSensor = def.isSensor
	boxBody.CreateFixtureFromDef(&boxFixDef)

	return &GameBody{Body: boxBody, HalfW: def.hx, HalfH: def.hy, Density: def.density, Friction: def.friction, Shape: Rectangle}
}

func createBall(def BallDef, world *box2d.B2World) *GameBody {
	ballDef := box2d.MakeB2BodyDef()
	ballDef.Position.Set(def.x, def.y)
	ballDef.Type = def.bodyType
	ballBody := world.CreateBody(&ballDef)

	ballShape := box2d.B2CircleShape{}
	ballShape.SetRadius(def.r)

	ballFixDef := box2d.MakeB2FixtureDef()
	ballFixDef.Shape = &ballShape
	ballFixDef.Density = def.density
	ballFixDef.Friction = def.friction
	ballFixDef.IsSensor = def.isSensor
	ballBody.CreateFixtureFromDef(&ballFixDef)

	return &GameBody{Body: ballBody, Radius: def.r, Shape: Circle, Density: def.density, Friction: def.friction}
}

func createBallFixureDef(radius, density, friction float64, isSensor bool) box2d.B2FixtureDef {
	shape := box2d.B2CircleShape{}
	shape.SetRadius(radius)
	boxFixDef := box2d.MakeB2FixtureDef()
	boxFixDef.Shape = &shape
	boxFixDef.Density = density
	boxFixDef.Friction = friction
	boxFixDef.IsSensor = isSensor

	return boxFixDef
}

func createBoxFixureDef(hx, hy, density, friction float64, isSensor bool) box2d.B2FixtureDef {
	boxBox := box2d.B2PolygonShape{}
	boxBox.SetAsBox(float64(hx), float64(hy))
	boxFixDef := box2d.MakeB2FixtureDef()
	boxFixDef.Shape = &boxBox
	boxFixDef.Density = density
	boxFixDef.Friction = friction
	boxFixDef.IsSensor = isSensor

	return boxFixDef
}
