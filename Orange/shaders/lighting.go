// Файл: shaders/lighting.go
// Назначение: управление вариантами моделей освещения (lighting variants).
// Каждый вариант — пара (вершинный шейдер, фрагментный шейдер), скомпилированная
// в программу OpenGL. Поддерживается переключение между вариантами.
//
// Ключевые типы:
//   ShadingMode — режим затенения: Gouraud (вершинное) или Phong (пиксельное).
//   LightingVariant — один вариант освещения: имя, модель, пути к шейдерам, программа.
//
// Ключевые функции:
//   InitLightingVariants — компилирует и линкует все варианты освещения.
//   GetCurrentLightingProgram — возвращает программу текущего варианта.
//   GetCurrentLightingName — возвращает имя текущего варианта.
//   ToggleShadingMode — переключает между Gouraud и Phong для текущей модели.
//   CycleLightingVariant — циклически переключает вариант освещения.
//   CleanupLightingVariants — удаляет все шейдерные программы.
//
// Зависимости:
//   Внутренние: shaders/shaders.go (LoadShaderFile, CompileShader).
//   Внешние: github.com/go-gl/gl/v4.6-core/gl.

package shaders

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// ShadingMode — режим затенения (интерполяции освещения).
// Gouraud — освещение вычисляется в вершинном шейдере, интерполируется.
// Phong  — нормаль интерполируется, освещение вычисляется во фрагментном шейдере.
type ShadingMode int

const (
	ShadingGouraud ShadingMode = iota
	ShadingPhong
)

// String — возвращает строковое представление режима затенения.
func (s ShadingMode) String() string {
	if s == ShadingPhong {
		return "Phong"
	}
	return "Gouraud"
}

// LightingVariant — один вариант модели освещения.
// Поля:
//
//	Name     — отображаемое имя (например, "Phong + Normal Map").
//	Model    — название модели освещения ("Phong" или "Blinn-Phong").
//	Mode     — режим затенения (Gouraud или Phong).
//	VertPath — путь к файлу вершинного шейдера (.vert).
//	FragPath — путь к файлу фрагментного шейдера (.frag).
//	Program  — идентификатор слинкованной шейдерной программы OpenGL.
type LightingVariant struct {
	Name     string
	Model    string
	Mode     ShadingMode
	VertPath string
	FragPath string
	Program  uint32
}

// lightingVariants — список всех доступных вариантов освещения.
// Текущая реализация содержит два варианта:
// 1. Фонг (Phong) + карта нормалей — orange_phong.{vert,frag}.
// 2. Блинн-Фонг (Blinn-Phong) + карта нормалей — orange_phong.vert + orange_blinn_phong.frag.
// Оба варианта используют Phong Shading (пиксельное освещение с bump mapping).
var lightingVariants = []LightingVariant{
	{Name: "Phong + Normal Map", Model: "Phong", Mode: ShadingPhong,
		VertPath: "shaders/lighting/orange_phong.vert", FragPath: "shaders/lighting/orange_phong.frag"},
	{Name: "Blinn-Phong + Normal Map", Model: "Blinn-Phong", Mode: ShadingPhong,
		VertPath: "shaders/lighting/orange_phong.vert", FragPath: "shaders/lighting/orange_blinn_phong.frag"},
}

// shaderBasePath — корневой путь для поиска файлов шейдеров.
// Устанавливается из main.go через SetBasePath.
var shaderBasePath string

// SetBasePath — устанавливает корневой путь для загрузки шейдерных файлов.
// Вызывается из main.go после определения корня проекта.
func SetBasePath(base string) {
	shaderBasePath = base
}

// currentLightingIndex — индекс текущего активного варианта освещения.
var currentLightingIndex = 0

// InitLightingVariants — компилирует и линкует все варианты освещения.
// Для каждого варианта в lightingVariants:
// 1. Загружает исходный код вершинного и фрагментного шейдера.
// 2. Компилирует каждый шейдер.
// 3. Создаёт программу, прикрепляет шейдеры, линкует.
// 4. Удаляет шейдерные объекты (после линковки не нужны).
// 5. Сохраняет идентификатор программы в variant.Program.
//
// Вызывается однократно при инициализации OpenGL.
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
// GetCurrentLightingProgram — возвращает программу текущего варианта освещения.
// Используется в цикле рендера для активации шейдера.
func GetCurrentLightingProgram() uint32 {
	if currentLightingIndex < 0 || currentLightingIndex >= len(lightingVariants) {
		currentLightingIndex = 0
	}
	return lightingVariants[currentLightingIndex].Program
}

// GetCurrentLightingName — возвращает имя текущего варианта освещения.
// Используется для отображения в заголовке окна и в UI.
func GetCurrentLightingName() string {
	if currentLightingIndex < 0 || currentLightingIndex >= len(lightingVariants) {
		currentLightingIndex = 0
	}
	return lightingVariants[currentLightingIndex].Name
}

// GetCurrentLightingModel — возвращает название модели освещения ("Phong" или "Blinn-Phong").
func GetCurrentLightingModel() string {
	if currentLightingIndex < 0 || currentLightingIndex >= len(lightingVariants) {
		currentLightingIndex = 0
	}
	return lightingVariants[currentLightingIndex].Model
}

// GetCurrentShadingMode — возвращает текущий режим затенения (Gouraud или Phong).
func GetCurrentShadingMode() ShadingMode {
	if currentLightingIndex < 0 || currentLightingIndex >= len(lightingVariants) {
		currentLightingIndex = 0
	}
	return lightingVariants[currentLightingIndex].Mode
}

// ToggleShadingMode — переключает режим затенения (Gouraud ↔ Phong)
// для текущей модели освещения. Ищет вариант с той же Model и противоположным Mode.
// Вызывается по нажатию клавиши Y.
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

// CycleLightingVariant — переключает вариант освещения по циклу.
// forward=true — следующий вариант (клавиша T).
// forward=false — предыдущий вариант (клавиша G).
func CycleLightingVariant(forward bool) {
	if forward {
		currentLightingIndex = (currentLightingIndex + 1) % len(lightingVariants)
	} else {
		currentLightingIndex = (currentLightingIndex - 1 + len(lightingVariants)) % len(lightingVariants)
	}
}

// GetLightingVariantCount — возвращает количество загруженных вариантов освещения.
func GetLightingVariantCount() int {
	return len(lightingVariants)
}

// CleanupLightingVariants — удаляет все шейдерные программы.
// Вызывается при завершении программы (defer в main.go).
func CleanupLightingVariants() {
	for _, variant := range lightingVariants {
		if variant.Program != 0 {
			gl.DeleteProgram(variant.Program)
		}
	}
}
