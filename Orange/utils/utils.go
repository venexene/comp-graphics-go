package utils

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/venexene/comp-graphics-go/input"
	"github.com/venexene/comp-graphics-go/lighting"
	"github.com/venexene/comp-graphics-go/scene"
	"github.com/venexene/comp-graphics-go/shaders"
)

var (
	cam         = scene.DefaultCamera()
	lightCfg    = lighting.DefaultLight()
	matCfg      = lighting.DefaultMaterial()
	inputState  input.State
	orangeScene *scene.OrangeScene
	heartScene  *scene.HeartScene
	orangeState scene.ObjectState
	heartState  scene.ObjectState
)

func initObjectStates() {
	orangeState = scene.ObjectState{
		Position:  mgl32.Vec3{-1.56, -1.51, 0.91},
		Scale:     0.5,
		RotationX: 0.0,
		RotationY: 0.0,
		RotationZ: 0.0,
	}
	heartState = scene.ObjectState{
		Position:  mgl32.Vec3{-1.8, -2.495 - 0.25, 0.0},
		Scale:     0.18,
		RotationX: 0.0,
		RotationY: 0.25,
		RotationZ: 0.0,
	}
}
func InitScene(basePath string) error {
	var err error
	orangeScene, err = scene.NewOrangeScene(basePath)
	if err != nil {
		return err
	}
	heartScene, err = scene.NewHeartScene(basePath)
	if err != nil {
		return err
	}
	initObjectStates()
	return nil
}
func DrawScene(window *glfw.Window, projection mgl32.Mat4) {
	input.ProcessInput(window, &cam, &orangeState, &lightCfg, nil, &inputState)
	program := shaders.GetCurrentLightingProgram()
	scene.DrawSceneObjects(program, orangeScene, heartScene, &cam, projection, &lightCfg, &matCfg, &orangeState, &heartState)
	glfw.PollEvents()
	window.SwapBuffers()
}
func Cleanup() {
	if orangeScene != nil {
		orangeScene.Cleanup()
	}
	if heartScene != nil {
		heartScene.Cleanup()
	}
}
func GetLightingName() string     { return shaders.GetCurrentLightingName() }
func GetShadingMode() string      { return shaders.GetCurrentShadingMode().String() }
func GetLinearCoef() float32      { return lightCfg.LinearCoef }
func GetQuadraticCoef() float32   { return lightCfg.QuadraticCoef }
func GetAmbientStrength() float32 { return lightCfg.AmbientStrength }
func GetAttenuationMode() string  { return lightCfg.Mode.String() }
func GetLightPosition() (float32, float32, float32) {
	return lightCfg.Position.X(), lightCfg.Position.Y(), lightCfg.Position.Z()
}
func GetSelectedObjectName() string { return "Orange + Heart" }
