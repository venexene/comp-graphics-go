// Package utils provides a compatibility layer that delegates to specialized
// packages (lighting, scene, input, shaders). New code should import those
// packages directly.
package utils

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/input"
	"github.com/venexene/comp-graphics-go/lighting"
	"github.com/venexene/comp-graphics-go/objects"
	"github.com/venexene/comp-graphics-go/scene"
	"github.com/venexene/comp-graphics-go/shaders"
)

// Re-exported types for backward compatibility.
type SceneObject = scene.SceneObject

var (
	defaultTex uint32

	cam        = scene.DefaultCamera()
	mainState  = scene.DefaultObjectState()
	lightCfg   = lighting.DefaultLight()
	matCfg     = lighting.DefaultMaterial()
	sel        = scene.NewSelection("Main")
	inputState input.State
)

// InitScene creates the default white texture used for untextured models.
func InitScene() {
	defaultTex = scene.CreateWhiteTexture()
}

// RegisterSceneObjects stores extra scene objects for selection cycling.
func RegisterSceneObjects(objs ...*SceneObject) {
	sel.RegisterObjects(objs...)
}

// SetMainObjectName sets the display name of the primary model.
func SetMainObjectName(name string) {
	sel.SetMainName(name)
}

// DrawScene processes input, clears the framebuffer and renders all objects.
func DrawScene(window *glfw.Window, model *objects.Model, view, projection mgl32.Mat4, extras ...*SceneObject) {
	input.ProcessInput(window, &cam, &mainState, &lightCfg, sel, &inputState)
	program := shaders.GetCurrentLightingProgram()
	scene.DrawScene(program, model, &mainState, extras, &cam, projection, &lightCfg, &matCfg, defaultTex)
	glfw.PollEvents()
	window.SwapBuffers()
}

// Getters for UI / title bar.

func GetLightingName() string       { return shaders.GetCurrentLightingName() }
func GetShadingMode() string        { return shaders.GetCurrentShadingMode().String() }
func GetLinearCoef() float32        { return lightCfg.LinearCoef }
func GetQuadraticCoef() float32     { return lightCfg.QuadraticCoef }
func GetAmbientStrength() float32   { return lightCfg.AmbientStrength }
func GetSelectedObjectName() string { return sel.SelectedName() }
func GetAttenuationMode() string    { return lightCfg.Mode.String() }

func GetLightPosition() (float32, float32, float32) {
	return lightCfg.Position.X(), lightCfg.Position.Y(), lightCfg.Position.Z()
}
