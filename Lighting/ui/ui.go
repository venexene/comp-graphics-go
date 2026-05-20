// ui/ui.go — Текстовый интерфейс пользователя (оверлей).
//
// Назначение: предоставляет простой текстовый UI в виде ASCII-панели,
// выводимой в терминал при старте, и обновляет заголовок окна в реальном
// времени. Не использует графические UI-библиотеки (Dear ImGui и т.п.).
//
// Ключевые типы:
//   LightingControlState — состояние UI (показ/скрытие панели).
//
// Ключевые функции:
//   InitializeUI()       — инициализация UI (no-op для текстового режима).
//   BeginFrame()         — начало кадра (no-op).
//   EndFrame()           — конец кадра (no-op).
//   Cleanup()            — освобождение ресурсов (no-op).
//   GetUIOverlayText()   — формирует ASCII-панель с текущими параметрами.
//
// Зависимости: вызывается из main() при инициализации и для получения текста.
package ui

import (
	"fmt"
)

// LightingControlState — состояние панели управления освещением.
type LightingControlState struct {
	ShowPanel bool // видимость панели (в будущем для переключения)
}

var state = LightingControlState{ShowPanel: true}

// InitializeUI — инициализирует UI-систему.
// В текущей реализации UI — текстовый, выводится в терминал.
func InitializeUI(window interface{}) error {
	return nil
}

// BeginFrame — начало нового кадра UI.
// No-op для текстового UI.
func BeginFrame() {}

// EndFrame — завершение кадра UI.
// No-op для текстового UI.
func EndFrame() {}

// Cleanup — освобождает ресурсы UI.
// No-op для текстового UI.
func Cleanup() {}

// GetUIOverlayText — формирует строку с ASCII-панелью управления.
// Принимает функции-геттеры для получения текущих значений:
//   getLightingName — название модели освещения;
//   getShadingMode — режим шейдинга (Gouraud/Phong);
//   getLinearCoef — линейный коэффициент затухания;
//   getQuadraticCoef — квадратичный коэффициент затухания;
//   getAmbientStrength — мощность фонового освещения;
//   getLightPosition — позиция источника света;
//   getSelectedObjectName — имя выбранного объекта;
//   getAttenuationMode — режим затухания (Both/Linear/Quadratic).
// Возвращает: строку с отформатированной ASCII-таблицей.
func GetUIOverlayText(
	getLightingName func() string,
	getShadingMode func() string,
	getLinearCoef func() float32,
	getQuadraticCoef func() float32,
	getAmbientStrength func() float32,
	getLightPosition func() (float32, float32, float32),
	getSelectedObjectName func() string,
	getAttenuationMode func() string,
) string {
	lx, ly, lz := getLightPosition()
	text := fmt.Sprintf(`
╔═══════════════════════════════════╗
║     LIGHTING CONTROLS UI          ║
╠═══════════════════════════════════╣
║ Model: %-27s ║
║ Shading: %-25s ║
║ Selected: %-23s ║
╠═══════════════════════════════════╣
║ LIGHT POSITION                    ║
║   X: %6.2f  Y: %6.2f  Z: %6.2f   ║
╠═══════════════════════════════════╣
║ ATTENUATION                       ║
║   Mode (M):      %-15s ║
║   Linear Coef (Z/X):  %.2f   ║
║   Quad. Coef (C/V):   %.2f   ║
╠═══════════════════════════════════╣
║ LIGHT POWER                       ║
║   Ambient Strength (B/N):  %.2f   ║
╠═══════════════════════════════════╣
║ KEYBOARD SHORTCUTS                ║
║   T/G     - Switch Lighting Model ║
║   Y       - Toggle Gouraud/Phong  ║
║   M       - Cycle Attenuation     ║
║   Tab     - Select Next Object    ║
║   Arrows  - Rotate Camera         ║
║   WASD    - Pan Camera            ║
║   +/-     - Zoom In/Out           ║
║   Q/E     - Scale Object          ║
║   R/F     - Rotate Object Z       ║
║   1/2     - Rotate Object X       ║
║   3/4     - Rotate Object Y       ║
║   IJKL    - Move Object XY        ║
║   U/O     - Move Object Z         ║
║   Alt+IJKL- Move Light XY         ║
║   Alt+U/O - Move Light Z          ║
║   Ctrl+R  - Reset All             ║
╚═══════════════════════════════════╝
`,
		getLightingName(),
		getShadingMode(),
		getSelectedObjectName(),
		lx, ly, lz,
		getAttenuationMode(),
		getLinearCoef(),
		getQuadraticCoef(),
		getAmbientStrength(),
	)
	return text
}
