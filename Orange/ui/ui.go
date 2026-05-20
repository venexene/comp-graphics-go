// Файл: ui/ui.go
// Назначение: интерфейс пользователя (overlay-текст с информацией о сцене и управлении).
//
// Ключевые типы:
//   LightingControlState — состояние панели управления (показ/скрытие).
//
// Ключевые функции:
//   GetUIOverlayText — возвращает строку с информацией о текущем состоянии сцены.
//   InitializeUI — заглушка инициализации UI.
//   BeginFrame/EndFrame — заглушки для начала/конца кадра UI.
//
// Зависимости:
//   Внутренние: — (используется cmd/main.go через функции-колбеки).
//   Внешние: нет (только стандартная библиотека).

package ui

import (
	"fmt"
)

// LightingControlState — состояние панели управления.
// Поле ShowPanel управляет видимостью информационного overlay.
type LightingControlState struct {
	ShowPanel bool
}

var state = LightingControlState{ShowPanel: true}

// InitializeUI — инициализация UI (заглушка).
// В текущей реализации UI — это только текстовый overlay, без графических элементов.
// Параметр window зарезервирован для будущей интеграции с ImGui или аналогами.
func InitializeUI(window interface{}) error {
	return nil
}

// BeginFrame — начало кадра UI (заглушка для будущей интеграции).
func BeginFrame() {
}

// EndFrame — конец кадра UI (заглушка для будущей интеграции).
func EndFrame() {
}

// Cleanup — очистка ресурсов UI (заглушка).
func Cleanup() {
}

// GetUIOverlayText — возвращает форматированную строку с информацией о сцене.
// Параметры — функции-колбеки для получения текущих значений из utils:
//
//	getLightingName     — имя текущей модели освещения.
//	getShadingMode      — режим затенения (Gouraud/Phong).
//	getLinearCoef       — множитель линейного затухания.
//	getQuadraticCoef    — множитель квадратичного затухания.
//	getAmbientStrength  — сила ambient-освещения.
//	getLightPosition    — позиция источника света.
//	getAttenuationMode  — режим затухания (Both/Linear/Quadratic).
//
// Возвращает: отформатированный текст для вывода в консоль.
func GetUIOverlayText(
	getLightingName func() string,
	getShadingMode func() string,
	getLinearCoef func() float32,
	getQuadraticCoef func() float32,
	getAmbientStrength func() float32,
	getLightPosition func() (float32, float32, float32),
	getAttenuationMode func() string,
) string {
	lx, ly, lz := getLightPosition()
	text := fmt.Sprintf(`
╔═══════════════════════════════════╗
║  ORANGE + HEART — BUMP MAPPING    ║
╠═══════════════════════════════════╣
║ Model: %-27s ║
║ Shading: %-25s ║
╠═══════════════════════════════════╣
║ LIGHT                             ║
║   Pos:  %5.2f %5.2f %5.2f          ║
║   Atten (M):   %-15s ║
║   Linear (Z/X):    %.2f         ║
║   Quad   (C/V):    %.2f         ║
║   Ambient (B/N):   %.2f         ║
╠═══════════════════════════════════╣
║ CAMERA                            ║
║   Arrows  - Rotate               ║
║   WASD    - Pan                  ║
║   Space/Shift - Up/Down          ║
║   +/-     - Zoom                 ║
╠═══════════════════════════════════╣
║ LIGHT & SHADER                    ║
║   Alt+IJKLUO - Move Light        ║
║   T/G     - Switch Model         ║
║   Y       - Toggle Shading       ║
║   M       - Cycle Attenuation    ║
║   Ctrl+R  - Reset All            ║
╚═══════════════════════════════════╝
`,
		getLightingName(),
		getShadingMode(),
		lx, ly, lz,
		getAttenuationMode(),
		getLinearCoef(),
		getQuadraticCoef(),
		getAmbientStrength(),
	)
	return text
}
