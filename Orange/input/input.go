// Файл: input/input.go
// Назначение: обработка пользовательского ввода с клавиатуры (GLFW).
//
// Ключевые константы:
//   CameraRotateSpeed — скорость вращения камеры (радиан/кадр).
//   CameraPanSpeed — скорость панорамирования (единиц/кадр).
//   CameraZoomSpeed — скорость зума (единиц/кадр).
//   ObjectMoveSpeed — скорость перемещения объекта/света (единиц/кадр).
//   ObjectScaleSpeed — скорость масштабирования (единиц/кадр).
//   ParamAdjustSpeed — скорость регулировки параметров (единиц/кадр).
//
// Ключевые типы:
//   State — состояние клавиш для edge-триггеров (однократное срабатывание).
//
// Ключевые функции:
//   ProcessInput — обрабатывает все клавиши и обновляет состояние сцены.
//
// Зависимости:
//   Внутренние: lighting, scene, shaders.
//   Внешние: github.com/go-gl/glfw/v3.3/glfw.

package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/venexene/comp-graphics-go/lighting"
	"github.com/venexene/comp-graphics-go/scene"
	"github.com/venexene/comp-graphics-go/shaders"
)

// Константы скорости для различных действий.
const (
	CameraRotateSpeed = 0.01   // радиан/кадр — вращение камеры
	CameraPanSpeed    = 0.01   // единиц модели/кадр — панорамирование
	CameraZoomSpeed   = 0.1    // единиц/кадр — приближение/отдаление
	ObjectMoveSpeed   = 0.01   // единиц/кадр — перемещение объекта или света
	ObjectScaleSpeed  = 0.001  // единиц/кадр — изменение масштаба
	ParamAdjustSpeed  = 0.01   // единиц/кадр — регулировка linear/quadratic/ambient
)

// State — хранит предыдущее состояние клавиш для edge-триггеров.
// Поле lastX = true, если в предыдущем кадре клавиша X была нажата.
// Используется для однократного срабатывания (например, переключение режима).
type State struct {
	lastT   bool
	lastG   bool
	lastY   bool
	lastTab bool
	lastM   bool
}

// ProcessInput — обрабатывает нажатия клавиш и обновляет состояние сцены.
// Параметры:
//
//	window    — GLFW-окно (для получения состояния клавиш).
//	cam       — камера (вращение, панорамирование, зум).
//	mainState — состояние основного объекта (перемещение, масштаб, вращение).
//	lightCfg  — конфигурация источника света (позиция, коэффициенты, режим).
//	sel       — механизм выбора объекта (может быть nil).
//	inputState — состояние клавиш для edge-триггеров.
//
// Управление камерой:
//   Стрелки — вращение вокруг target.
//   WASD — панорамирование (вперёд/назад/вправо/влево).
//   Space/Shift — вверх/вниз.
//   +/- — зум.
//
// Управление объектом:
//   Q/E — уменьшение/увеличение масштаба.
//   R/F — вращение вокруг Z.
//   1/2 — вращение вокруг X.
//   3/4 — вращение вокруг Y.
//   IJKL (без Alt) — перемещение выбранного объекта.
//
// Управление источником света:
//   Alt+IJKLUO — перемещение источника.
//   Z/X — изменение линейного затухания.
//   C/V — изменение квадратичного затухания.
//   B/N — изменение ambient strength.
//
// Переключение режимов:
//   T/G — следующий/предыдущий вариант освещения.
//   Y — переключение Gouraud/Phong.
//   M — циклическое переключение затухания.
//   Tab — переключение выбранного объекта.
//   Ctrl+R — сброс всех параметров.
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
