// shaders/lighting.go — Управление шейдерными программами освещения.
//
// Назначение: определяет все варианты освещения (4 модели × 2 режима шейдинга =
// 8 шейдерных программ), компилирует и линкует их, предоставляет функции
// переключения между вариантами.
//
// Ключевые типы:
//   ShadingMode — перечисление: Gouraud или Phong.
//   LightingVariant — описание одного варианта освещения: имя, модель,
//     режим шейдинга, пути к .vert/.frag, ID программы OpenGL.
//
// Ключевые функции:
//   InitLightingVariants()       — компиляция и линковка всех 8 программ.
//   GetCurrentLightingProgram()  — ID текущей шейдерной программы.
//   ToggleShadingMode()          — переключение Gouraud ↔ Phong для той же модели.
//   CycleLightingVariant(dir)    — циклическое переключение по всем вариантам.
//   CleanupLightingVariants()    — удаление всех шейдерных программ.
//
// Зависимости: shaders/shaders.go (CompileShader), вызывается из cmd/main.go.
package shaders

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// ShadingMode — режим шейдинга (интерполяции освещения).
type ShadingMode int

const (
	// ShadingGouraud — освещение вычисляется в вершинном шейдере,
	// результат интерполируется по грани треугольника.
	ShadingGouraud ShadingMode = iota
	// ShadingPhong — нормаль интерполируется, освещение вычисляется
	// во фрагментном шейдере попиксельно.
	ShadingPhong
)

// String — возвращает строковое название режима шейдинга.
func (s ShadingMode) String() string {
	if s == ShadingPhong {
		return "Phong"
	}
	return "Gouraud"
}

// LightingVariant — описание одного варианта освещения.
// Содержит все необходимые данные для идентификации, загрузки и использования
// пары шейдеров (вершинный + фрагментный) как шейдерной программы.
type LightingVariant struct {
	// Name — отображаемое имя (например, "Phong Gouraud").
	Name string
	// Model — название модели освещения ("Lambert", "Phong", "Blinn-Phong", "Toon").
	Model string
	// Mode — режим шейдинга: Gouraud или Phong.
	Mode ShadingMode
	// VertPath — путь к файлу вершинного шейдера (относительно корня проекта).
	VertPath string
	// FragPath — путь к файлу фрагментного шейдера.
	FragPath string
	// Program — идентификатор скомпилированной и слинкованной шейдерной программы OpenGL.
	Program uint32
}

// lightingVariants — полный список из 8 вариантов освещения.
// 4 модели (Lambert, Phong, Blinn-Phong, Toon) × 2 режима шейдинга (Gouraud, Phong).
// Каждая пара использует общий фрагментный шейдер basic_gouraud.frag для Gouraud-
// вариантов и общий вершинный шейдер basic_phong.vert для Phong-вариантов.
var lightingVariants = []LightingVariant{
	{Name: "Lambert Gouraud", Model: "Lambert", Mode: ShadingGouraud, VertPath: "shaders/lighting/lambert_gouraud.vert", FragPath: "shaders/lighting/basic_gouraud.frag"},
	{Name: "Lambert Phong", Model: "Lambert", Mode: ShadingPhong, VertPath: "shaders/lighting/basic_phong.vert", FragPath: "shaders/lighting/lambert_phong.frag"},
	{Name: "Phong Gouraud", Model: "Phong", Mode: ShadingGouraud, VertPath: "shaders/lighting/phong_gouraud.vert", FragPath: "shaders/lighting/basic_gouraud.frag"},
	{Name: "Phong Phong", Model: "Phong", Mode: ShadingPhong, VertPath: "shaders/lighting/basic_phong.vert", FragPath: "shaders/lighting/phong_phong.frag"},
	{Name: "Blinn-Phong Gouraud", Model: "Blinn-Phong", Mode: ShadingGouraud, VertPath: "shaders/lighting/blinn_phong_gouraud.vert", FragPath: "shaders/lighting/basic_gouraud.frag"},
	{Name: "Blinn-Phong Phong", Model: "Blinn-Phong", Mode: ShadingPhong, VertPath: "shaders/lighting/basic_phong.vert", FragPath: "shaders/lighting/blinn_phong_phong.frag"},
	{Name: "Toon Gouraud", Model: "Toon", Mode: ShadingGouraud, VertPath: "shaders/lighting/toon_gouraud.vert", FragPath: "shaders/lighting/basic_gouraud.frag"},
	{Name: "Toon Phong", Model: "Toon", Mode: ShadingPhong, VertPath: "shaders/lighting/basic_phong.vert", FragPath: "shaders/lighting/toon_phong.frag"},
}

// shaderBasePath — префикс пути к файлам шейдеров.
// Устанавливается SetBasePath() перед InitLightingVariants().
var shaderBasePath string

// SetBasePath — задаёт базовый каталог для поиска шейдерных файлов.
// Принимает: base — абсолютный путь к корню проекта (где лежит go.mod).
// Вызывается: из cmd/main.go перед InitLightingVariants().
func SetBasePath(base string) {
	shaderBasePath = base
}

// currentLightingIndex — индекс текущего выбранного варианта освещения
// в слайсе lightingVariants.
var currentLightingIndex = 0

// InitLightingVariants — компилирует и линкует все 8 шейдерных программ.
// Для каждой LightingVariant: читает .vert и .frag файлы → компилирует →
// линкует → сохраняет Program ID. В случае ошибки возвращает её немедленно.
// Побочные эффекты: создаёт 8 шейдерных программ в OpenGL.
// Вызывается: однократно в initOpenGL() при запуске.
func InitLightingVariants() error {
	for i := range lightingVariants {
		vertPath := lightingVariants[i].VertPath
		fragPath := lightingVariants[i].FragPath
		if shaderBasePath != "" {
			vertPath = shaderBasePath + "/" + vertPath
			fragPath = shaderBasePath + "/" + fragPath
		}

		vert, err := LoadShaderFile(vertPath)
		if err != nil {
			return fmt.Errorf("failed to load vertex shader %s: %w", vertPath, err)
		}

		frag, err := LoadShaderFile(fragPath)
		if err != nil {
			return fmt.Errorf("failed to load fragment shader %s: %w", fragPath, err)
		}

		vertShader, err := CompileShader(vert, gl.VERTEX_SHADER)
		if err != nil {
			return err
		}

		fragShader, err := CompileShader(frag, gl.FRAGMENT_SHADER)
		if err != nil {
			return err
		}

		program := gl.CreateProgram()
		gl.AttachShader(program, vertShader)
		gl.AttachShader(program, fragShader)
		gl.LinkProgram(program)

		var status int32
		gl.GetProgramiv(program, gl.LINK_STATUS, &status)
		if status == gl.FALSE {
			var logLength int32
			gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
			log := make([]byte, logLength)
			gl.GetProgramInfoLog(program, logLength, nil, &log[0])
			return fmt.Errorf("failed to link program for %s: %s", lightingVariants[i].Name, string(log))
		}

		gl.DeleteShader(vertShader)
		gl.DeleteShader(fragShader)

		lightingVariants[i].Program = program
	}

	return nil
}

// GetCurrentLightingProgram — возвращает ID OpenGL текущей шейдерной программы.
// Используется в scene.DrawScene() для вызова glUseProgram().
func GetCurrentLightingProgram() uint32 {
	if currentLightingIndex < 0 || currentLightingIndex >= len(lightingVariants) {
		currentLightingIndex = 0
	}
	return lightingVariants[currentLightingIndex].Program
}

// GetCurrentLightingName — возвращает отображаемое имя текущего варианта
// (например, "Phong Gouraud") для UI и заголовка окна.
func GetCurrentLightingName() string {
	if currentLightingIndex < 0 || currentLightingIndex >= len(lightingVariants) {
		currentLightingIndex = 0
	}
	return lightingVariants[currentLightingIndex].Name
}

// GetCurrentLightingModel — возвращает название модели освещения
// ("Lambert", "Phong", "Blinn-Phong", "Toon").
func GetCurrentLightingModel() string {
	if currentLightingIndex < 0 || currentLightingIndex >= len(lightingVariants) {
		currentLightingIndex = 0
	}
	return lightingVariants[currentLightingIndex].Model
}

// GetCurrentShadingMode — возвращает текущий режим шейдинга (Gouraud/Phong).
func GetCurrentShadingMode() ShadingMode {
	if currentLightingIndex < 0 || currentLightingIndex >= len(lightingVariants) {
		currentLightingIndex = 0
	}
	return lightingVariants[currentLightingIndex].Mode
}

// ToggleShadingMode — переключает между Gouraud и Phong для той же модели
// освещения. Ищет вариант с тем же Model и противоположным Mode,
// устанавливает currentLightingIndex на него.
// Вызывается: из input.ProcessInput() по клавише Y.
func ToggleShadingMode() {
	currentModel := GetCurrentLightingModel()
	currentMode := GetCurrentShadingMode()
	targetMode := ShadingGouraud
	if currentMode == ShadingGouraud {
		targetMode = ShadingPhong
	}
	for i, variant := range lightingVariants {
		if variant.Model == currentModel && variant.Mode == targetMode {
			currentLightingIndex = i
			return
		}
	}
}

// CycleLightingVariant — циклически переключает вариант освещения.
// Принимает: forward — true для перехода вперёд (T), false для назад (G).
func CycleLightingVariant(forward bool) {
	if forward {
		currentLightingIndex = (currentLightingIndex + 1) % len(lightingVariants)
	} else {
		currentLightingIndex = (currentLightingIndex - 1 + len(lightingVariants)) % len(lightingVariants)
	}
}

// GetLightingVariantCount — возвращает количество вариантов (8).
func GetLightingVariantCount() int {
	return len(lightingVariants)
}

// CleanupLightingVariants — удаляет все шейдерные программы OpenGL.
// Побочные эффекты: освобождает ресурсы GPU.
// Вызывается: при завершении программы (defer в main()).
func CleanupLightingVariants() {
	for _, variant := range lightingVariants {
		if variant.Program != 0 {
			gl.DeleteProgram(variant.Program)
		}
	}
}
