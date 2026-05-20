package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/venexene/comp-graphics-go/lighting"
	"github.com/venexene/comp-graphics-go/scene"
	"github.com/venexene/comp-graphics-go/shaders"
)

const (
	CameraRotateSpeed = 0.01
	CameraPanSpeed    = 0.01
	CameraZoomSpeed   = 0.1
	ObjectMoveSpeed   = 0.01
	ObjectScaleSpeed  = 0.001
	ParamAdjustSpeed  = 0.01
)

type State struct {
	lastT   bool
	lastG   bool
	lastY   bool
	lastTab bool
	lastM   bool
}

func ProcessInput(
	window *glfw.Window,
	cam *scene.Camera,
	mainState *scene.ObjectState,
	lightCfg *lighting.LightConfig,
	sel *scene.Selection,
	inputState *State,
) {
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		cam.Rotate(0, -CameraRotateSpeed)
	}
	if window.GetKey(glfw.KeyDown) == glfw.Press {
		cam.Rotate(0, CameraRotateSpeed)
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		cam.Rotate(-CameraRotateSpeed, 0)
	}
	if window.GetKey(glfw.KeyRight) == glfw.Press {
		cam.Rotate(CameraRotateSpeed, 0)
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		cam.PanForward(-CameraPanSpeed)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		cam.PanForward(CameraPanSpeed)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		cam.PanRight(CameraPanSpeed)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		cam.PanRight(-CameraPanSpeed)
	}
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		cam.PanUp(CameraPanSpeed)
	}
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press || window.GetKey(glfw.KeyRightShift) == glfw.Press {
		cam.PanUp(-CameraPanSpeed)
	}
	if window.GetKey(glfw.KeyKPAdd) == glfw.Press || window.GetKey(glfw.KeyEqual) == glfw.Press {
		cam.Zoom(-CameraZoomSpeed, 1.0, 50.0)
	}
	if window.GetKey(glfw.KeyKPSubtract) == glfw.Press || window.GetKey(glfw.KeyMinus) == glfw.Press {
		cam.Zoom(CameraZoomSpeed, 1.0, 50.0)
	}
	if window.GetKey(glfw.KeyQ) == glfw.Press {
		mainState.Scale -= ObjectScaleSpeed
		if mainState.Scale < 0.1 {
			mainState.Scale = 0.1
		}
	}
	if window.GetKey(glfw.KeyE) == glfw.Press {
		mainState.Scale += ObjectScaleSpeed
		if mainState.Scale > 3.0 {
			mainState.Scale = 3.0
		}
	}
	if window.GetKey(glfw.KeyR) == glfw.Press {
		mainState.RotationZ += 0.001
	}
	if window.GetKey(glfw.KeyF) == glfw.Press {
		mainState.RotationZ -= 0.001
	}
	if window.GetKey(glfw.Key1) == glfw.Press {
		mainState.RotationX += 0.001
	}
	if window.GetKey(glfw.Key2) == glfw.Press {
		mainState.RotationX -= 0.001
	}
	if window.GetKey(glfw.Key3) == glfw.Press {
		mainState.RotationY += 0.001
	}
	if window.GetKey(glfw.Key4) == glfw.Press {
		mainState.RotationY -= 0.001
	}
	altDown := window.GetKey(glfw.KeyLeftAlt) == glfw.Press || window.GetKey(glfw.KeyRightAlt) == glfw.Press
	if altDown {
		if window.GetKey(glfw.KeyI) == glfw.Press {
			lightCfg.Position[2] -= ObjectMoveSpeed
		}
		if window.GetKey(glfw.KeyK) == glfw.Press {
			lightCfg.Position[2] += ObjectMoveSpeed
		}
		if window.GetKey(glfw.KeyJ) == glfw.Press {
			lightCfg.Position[0] -= ObjectMoveSpeed
		}
		if window.GetKey(glfw.KeyL) == glfw.Press {
			lightCfg.Position[0] += ObjectMoveSpeed
		}
		if window.GetKey(glfw.KeyU) == glfw.Press {
			lightCfg.Position[1] += ObjectMoveSpeed
		}
		if window.GetKey(glfw.KeyO) == glfw.Press {
			lightCfg.Position[1] -= ObjectMoveSpeed
		}
	} else if sel != nil {
		if sel.IsMain() {
			if window.GetKey(glfw.KeyI) == glfw.Press {
				mainState.Position[2] -= ObjectMoveSpeed
			}
			if window.GetKey(glfw.KeyK) == glfw.Press {
				mainState.Position[2] += ObjectMoveSpeed
			}
			if window.GetKey(glfw.KeyJ) == glfw.Press {
				mainState.Position[0] -= ObjectMoveSpeed
			}
			if window.GetKey(glfw.KeyL) == glfw.Press {
				mainState.Position[0] += ObjectMoveSpeed
			}
			if window.GetKey(glfw.KeyU) == glfw.Press {
				mainState.Position[1] += ObjectMoveSpeed
			}
			if window.GetKey(glfw.KeyO) == glfw.Press {
				mainState.Position[1] -= ObjectMoveSpeed
			}
		} else if obj := sel.SelectedSceneObject(); obj != nil {
			if window.GetKey(glfw.KeyI) == glfw.Press {
				obj.Position[2] -= ObjectMoveSpeed
			}
			if window.GetKey(glfw.KeyK) == glfw.Press {
				obj.Position[2] += ObjectMoveSpeed
			}
			if window.GetKey(glfw.KeyJ) == glfw.Press {
				obj.Position[0] -= ObjectMoveSpeed
			}
			if window.GetKey(glfw.KeyL) == glfw.Press {
				obj.Position[0] += ObjectMoveSpeed
			}
			if window.GetKey(glfw.KeyU) == glfw.Press {
				obj.Position[1] += ObjectMoveSpeed
			}
			if window.GetKey(glfw.KeyO) == glfw.Press {
				obj.Position[1] -= ObjectMoveSpeed
			}
		}
	}
	currentT := window.GetKey(glfw.KeyT) == glfw.Press
	currentG := window.GetKey(glfw.KeyG) == glfw.Press
	currentY := window.GetKey(glfw.KeyY) == glfw.Press
	currentTab := window.GetKey(glfw.KeyTab) == glfw.Press
	currentM := window.GetKey(glfw.KeyM) == glfw.Press
	if currentT && !inputState.lastT {
		shaders.CycleLightingVariant(true)
	}
	if currentG && !inputState.lastG {
		shaders.CycleLightingVariant(false)
	}
	if currentY && !inputState.lastY {
		shaders.ToggleShadingMode()
	}
	if currentTab && !inputState.lastTab {
		if sel != nil {
			sel.CycleForward()
		}
	}
	if currentM && !inputState.lastM {
		lightCfg.CycleAttenuationMode()
	}
	inputState.lastT = currentT
	inputState.lastG = currentG
	inputState.lastY = currentY
	inputState.lastTab = currentTab
	inputState.lastM = currentM
	if window.GetKey(glfw.KeyZ) == glfw.Press {
		lightCfg.LinearCoef -= ParamAdjustSpeed
		if lightCfg.LinearCoef < 0.0 {
			lightCfg.LinearCoef = 0.0
		}
	}
	if window.GetKey(glfw.KeyX) == glfw.Press {
		lightCfg.LinearCoef += ParamAdjustSpeed
	}
	if window.GetKey(glfw.KeyC) == glfw.Press {
		lightCfg.QuadraticCoef -= ParamAdjustSpeed
		if lightCfg.QuadraticCoef < 0.0 {
			lightCfg.QuadraticCoef = 0.0
		}
	}
	if window.GetKey(glfw.KeyV) == glfw.Press {
		lightCfg.QuadraticCoef += ParamAdjustSpeed
	}
	if window.GetKey(glfw.KeyB) == glfw.Press {
		lightCfg.AmbientStrength -= ParamAdjustSpeed
		if lightCfg.AmbientStrength < 0.0 {
			lightCfg.AmbientStrength = 0.0
		}
	}
	if window.GetKey(glfw.KeyN) == glfw.Press {
		lightCfg.AmbientStrength += ParamAdjustSpeed
		if lightCfg.AmbientStrength > 1.0 {
			lightCfg.AmbientStrength = 1.0
		}
	}
	if window.GetKey(glfw.KeyR) == glfw.Press &&
		window.GetKey(glfw.KeyLeftControl) == glfw.Press {
		*mainState = scene.DefaultObjectState()
		*cam = scene.DefaultCamera()
	}
}
