// input/input.go — Обработка клавиатурного ввода.
//
// Назначение: читает состояние клавиш GLFW каждый кадр и обновляет камеру,
// объекты сцены, параметры освещения и выбранный шейдер.
//
// Ключевые типы:
//   State — предыдущее состояние клавиш для детекции нажатия (edge detection).
//
// Ключевые функции:
//   ProcessInput() — основной обработчик ввода (вызывается каждый кадр).
//
// Зависимости: вызывается из utils.DrawScene(); модифицирует Camera,
//   ObjectState, LightConfig, Selection, управляет shaders.CycleLightingVariant().
package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/venexene/comp-graphics-go/lighting"
	"github.com/venexene/comp-graphics-go/scene"
	"github.com/venexene/comp-graphics-go/shaders"
)

// Константы чувствительности управления.
const (
	CameraRotateSpeed = 0.01  // скорость вращения камеры (стрелки)
	CameraPanSpeed    = 0.01  // скорость панорамирования (WASD)
	CameraZoomSpeed   = 0.1   // скорость зума (+/-)
	ObjectMoveSpeed   = 0.01  // скорость перемещения объекта (IJKL/UO)
	ObjectScaleSpeed  = 0.001 // скорость масштабирования (Q/E)
	ParamAdjustSpeed  = 0.01  // скорость изменения параметров (Z/X/C/V/B/N)
)

// State — предыдущее состояние клавиш для детекции фронта (edge detection).
// Позволяет выполнить действие однократно при нажатии, а не каждый кадр.
type State struct {
	lastT   bool // предыдущее состояние клавиши T
	lastG   bool // предыдущее состояние клавиши G
	lastY   bool // предыдущее состояние клавиши Y
	lastTab bool // предыдущее состояние Tab
	lastM   bool // предыдущее состояние клавиши M
}

// ProcessInput — обрабатывает ввод с клавиатуры.
// Должен вызываться один раз в кадр перед отрисовкой сцены.
// Принимает:
//   window — окно GLFW (чтение состояния клавиш);
//   cam — камера для вращения/панорамирования/зума;
//   mainState — трансформация главного объекта;
//   lightCfg — конфигурация света (позиция, затухание, мощность);
//   sel — текущий выбор объекта;
//   inputState — предыдущее состояние клавиш (edge detection).
// Побочные эффекты: изменяет все переданные структуры.
func ProcessInput(
	window *glfw.Window,
	cam *scene.Camera,
	mainState *scene.ObjectState,
	lightCfg *lighting.LightConfig,
	sel *scene.Selection,
	inputState *State,
) {
	// ===== Вращение камеры (стрелки) =====
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

	// ===== Панорамирование камеры (WASD + Space/Shift) =====
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
	// Space — вверх по Y, Shift — вниз по Y.
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		cam.PanUp(CameraPanSpeed)
	}
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press || window.GetKey(glfw.KeyRightShift) == glfw.Press {
		cam.PanUp(-CameraPanSpeed)
	}

	// ===== Зум камеры (+/-) =====
	if window.GetKey(glfw.KeyKPAdd) == glfw.Press || window.GetKey(glfw.KeyEqual) == glfw.Press {
		cam.Zoom(-CameraZoomSpeed, 1.0, 50.0)
	}
	if window.GetKey(glfw.KeyKPSubtract) == glfw.Press || window.GetKey(glfw.KeyMinus) == glfw.Press {
		cam.Zoom(CameraZoomSpeed, 1.0, 50.0)
	}

	// ===== Масштабирование главного объекта (Q/E) =====
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

	// ===== Вращение главного объекта вокруг Z (R/F) =====
	if window.GetKey(glfw.KeyR) == glfw.Press {
		mainState.RotationZ += 0.001
	}
	if window.GetKey(glfw.KeyF) == glfw.Press {
		mainState.RotationZ -= 0.001
	}

	// ===== Вращение главного объекта вокруг X (1/2) и Y (3/4) =====
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

	// ===== Перемещение объекта/света (IJKL UO, Alt — свет) =====
	altDown := window.GetKey(glfw.KeyLeftAlt) == glfw.Press || window.GetKey(glfw.KeyRightAlt) == glfw.Press

	if altDown {
		// Alt+IJKL/UO — перемещение источника света.
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
	} else {
		// Без Alt — перемещение выбранного объекта (главного или дополнительного).
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

	// ===== Переключение шейдеров и режимов (T/G/Y/Tab/M) =====
	// Используем edge detection: действие выполняется только в момент нажатия,
	// а не каждый кадр удержания.
	currentT := window.GetKey(glfw.KeyT) == glfw.Press
	currentG := window.GetKey(glfw.KeyG) == glfw.Press
	currentY := window.GetKey(glfw.KeyY) == glfw.Press
	currentTab := window.GetKey(glfw.KeyTab) == glfw.Press
	currentM := window.GetKey(glfw.KeyM) == glfw.Press

	// T — следующий вариант освещения, G — предыдущий.
	if currentT && !inputState.lastT {
		shaders.CycleLightingVariant(true)
	}
	if currentG && !inputState.lastG {
		shaders.CycleLightingVariant(false)
	}
	// Y — переключение Gouraud ↔ Phong для той же модели.
	if currentY && !inputState.lastY {
		shaders.ToggleShadingMode()
	}
	// Tab — выбор следующего объекта (снеговик → сердце → default → ...).
	if currentTab && !inputState.lastTab {
		sel.CycleForward()
	}
	// M — циклическое переключение режима затухания (Both → Linear → Quadratic).
	if currentM && !inputState.lastM {
		lightCfg.CycleAttenuationMode()
	}

	// Сохраняем текущее состояние для следующего кадра.
	inputState.lastT = currentT
	inputState.lastG = currentG
	inputState.lastY = currentY
	inputState.lastTab = currentTab
	inputState.lastM = currentM

	// ===== Регулировка коэффициентов затухания (Z/X/C/V) =====
	// Z — уменьшение линейного коэффициента, X — увеличение.
	if window.GetKey(glfw.KeyZ) == glfw.Press {
		lightCfg.LinearCoef -= ParamAdjustSpeed
		if lightCfg.LinearCoef < 0.0 {
			lightCfg.LinearCoef = 0.0
		}
	}
	if window.GetKey(glfw.KeyX) == glfw.Press {
		lightCfg.LinearCoef += ParamAdjustSpeed
	}
	// C — уменьшение квадратичного коэффициента, V — увеличение.
	if window.GetKey(glfw.KeyC) == glfw.Press {
		lightCfg.QuadraticCoef -= ParamAdjustSpeed
		if lightCfg.QuadraticCoef < 0.0 {
			lightCfg.QuadraticCoef = 0.0
		}
	}
	if window.GetKey(glfw.KeyV) == glfw.Press {
		lightCfg.QuadraticCoef += ParamAdjustSpeed
	}

	// ===== Регулировка мощности фонового освещения (B/N) =====
	// Диапазон: [0.0, 1.0].
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

	// ===== Сброс (Ctrl+R) — возврат камеры и объекта в начальное состояние =====
	if window.GetKey(glfw.KeyR) == glfw.Press &&
		window.GetKey(glfw.KeyLeftControl) == glfw.Press {
		*mainState = scene.DefaultObjectState()
		*cam = scene.DefaultCamera()
	}
}
