package game

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/bytearena/box2d"
	"github.com/faiface/pixel"
)

func toPixelVec(v *box2d.B2Vec2) *pixel.Vec {
	return &pixel.Vec{X: v.X, Y: v.Y}
}

func toBox2dVec(v *pixel.Vec) *box2d.B2Vec2 {
	return &box2d.B2Vec2{X: v.X, Y: v.Y}
}

func worldToScreen(v *box2d.B2Vec2, cam *Camera) *pixel.Vec {
	return &pixel.Vec{X: v.X * Scale, Y: v.Y * Scale}
}

func screenToWorld(v pixel.Vec, cam *Camera) box2d.B2Vec2 {
	newVec := box2d.B2Vec2{X: v.X * InvScale, Y: v.Y * InvScale}
	newVec.X += cam.X
	return newVec
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

type BodyJson struct {
	X         float64
	Y         float64
	Angle     float64
	Hx        float64
	Hy        float64
	Radius    float64
	Density   float64
	Friction  float64
	BodyType  uint8
	BodyShape Shape
}

type ConfigData struct {
	Levels []LevelInfo
}

type LevelInfo struct {
	Name     string
	Filename string
}

type LevelData struct {
	Name   string
	Bodies []BodyJson
	Cargo  []BodyJson
}

func LoadConfig() *ConfigData {
	bytes, err := os.ReadFile("config.json")

	if err != nil {
		panic(err)
	}

	data := ConfigData{}
	err = json.Unmarshal(bytes, &data)

	if err != nil {
		panic(err)
	}

	return &data
}

func SaveToFile(g *Game) {
	data := LevelData{Name: "Dood"}

	bodies := g.Bodies
	for i := 0; i < len(bodies); i++ {
		b := bodies[i]
		pos := b.Body.GetPosition()
		friction := b.Body.GetFixtureList().GetFriction()
		density := b.Body.GetFixtureList().GetDensity()
		angle := b.Body.GetAngle()
		d := BodyJson{X: pos.X, Y: pos.Y, Angle: angle, Hx: b.HalfW, Hy: b.HalfH, Density: density, Friction: friction}
		d.BodyType = b.Body.GetType()
		d.BodyShape = b.Shape
		d.Radius = b.Radius
		d.BodyShape = b.Shape
		data.Bodies = append(data.Bodies, d)
	}

	cargo := g.CargoBodies
	for i := 0; i < len(cargo); i++ {
		b := cargo[i]
		pos := b.Body.GetPosition()
		friction := b.Body.GetFixtureList().GetFriction()
		density := b.Body.GetFixtureList().GetDensity()
		angle := b.Body.GetAngle()
		d := BodyJson{X: pos.X, Y: pos.Y, Angle: angle, Hx: b.HalfW, Hy: b.HalfH, Density: density, Friction: friction}
		d.BodyType = b.Body.GetType()
		d.BodyShape = b.Shape
		d.Radius = b.Radius
		d.BodyShape = b.Shape
		data.Cargo = append(data.Cargo, d)
	}

	fmt.Println(data.Cargo)

	file, _ := json.MarshalIndent(data, "", " ")
	err := os.WriteFile("newLevel.json", file, 0666)
	if err != nil {
		panic(err)
	}
	println("File saved!")
}

func LoadFromFile(filepath string) *LevelData {
	bytes, err := os.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	data := LevelData{}
	err = json.Unmarshal(bytes, &data)

	if err != nil {
		panic(err)
	}

	return &data
}
