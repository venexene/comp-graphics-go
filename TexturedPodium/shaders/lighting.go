// Пакет shaders управляет вариантами освещения: определяет список
// LightingVariant (10 комбинаций модель x шейдинг), компилирует шейдерные
// программы (InitLightingVariants), переключает активный вариант
// (CycleLightingVariant, ToggleShadingMode).
// Модели: Lambert, Phong, Blinn-Phong, Toon, Oren-Nayar.
// Шейдинг: Gouraud (освещение в вершинном шейдере), Phong (во фрагментном).
// Зависимости: shaders/shaders.go (CompileShader, LoadShaderFile).
package shaders

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// ShadingMode определяет, где вычисляется освещение:
// ShadingGouraud — в вершинном шейдере (быстрее, но артефакты),
// ShadingPhong — во фрагментном шейдере (точнее, но медленнее).
type ShadingMode int

const (
	ShadingGouraud ShadingMode = iota
	ShadingPhong
)

func (s ShadingMode) String() string {
	if s == ShadingPhong {
		return "Phong"
	}
	return "Gouraud"
}

// LightingVariant описывает одну шейдерную программу: модель освещения,
// тип шейдинга, пути к GLSL-файлам, идентификатор скомпилированной программы.
type LightingVariant struct {
	Name     string      // отображаемое имя (напр. "Lambert Gouraud")
	Model    string      // модель освещения (Lambert, Phong, Blinn-Phong, Toon)
	Mode     ShadingMode // Gouraud или Phong
	VertPath string      // путь к вершинному шейдеру
	FragPath string      // путь к фрагментному шейдеру
	Program  uint32      // OpenGL-идентификатор шейдерной программы
}

// lightingVariants — полный перечень из 10 шейдерных программ.
// Для каждой модели освещения существует Gouraud- и Phong-вариант.
var lightingVariants = []LightingVariant{
	{Name: "Lambert Gouraud", Model: "Lambert", Mode: ShadingGouraud, VertPath: "shaders/lighting/lambert_gouraud.vert", FragPath: "shaders/lighting/basic_gouraud.frag"},
	{Name: "Lambert Phong", Model: "Lambert", Mode: ShadingPhong, VertPath: "shaders/lighting/basic_phong.vert", FragPath: "shaders/lighting/lambert_phong.frag"},
	{Name: "Phong Gouraud", Model: "Phong", Mode: ShadingGouraud, VertPath: "shaders/lighting/phong_gouraud.vert", FragPath: "shaders/lighting/basic_gouraud.frag"},
	{Name: "Phong Phong", Model: "Phong", Mode: ShadingPhong, VertPath: "shaders/lighting/basic_phong.vert", FragPath: "shaders/lighting/phong_phong.frag"},
	{Name: "Blinn-Phong Gouraud", Model: "Blinn-Phong", Mode: ShadingGouraud, VertPath: "shaders/lighting/blinn_phong_gouraud.vert", FragPath: "shaders/lighting/basic_gouraud.frag"},
	{Name: "Blinn-Phong Phong", Model: "Blinn-Phong", Mode: ShadingPhong, VertPath: "shaders/lighting/basic_phong.vert", FragPath: "shaders/lighting/blinn_phong_phong.frag"},
	{Name: "Toon Gouraud", Model: "Toon", Mode: ShadingGouraud, VertPath: "shaders/lighting/toon_gouraud.vert", FragPath: "shaders/lighting/basic_gouraud.frag"},
	{Name: "Toon Phong", Model: "Toon", Mode: ShadingPhong, VertPath: "shaders/lighting/basic_phong.vert", FragPath: "shaders/lighting/toon_phong.frag"},
	{Name: "Oren-Nayar Gouraud", Model: "Oren-Nayar", Mode: ShadingGouraud, VertPath: "shaders/lighting/oren_nayar_gouraud.vert", FragPath: "shaders/lighting/basic_gouraud.frag"},
	{Name: "Oren-Nayar Phong", Model: "Oren-Nayar", Mode: ShadingPhong, VertPath: "shaders/lighting/basic_phong.vert", FragPath: "shaders/lighting/oren_nayar_phong.frag"},
}

// shaderBasePath — префикс для путей к шейдерным файлам.
// Устанавливается через SetBasePath перед InitLightingVariants.
var shaderBasePath string

// SetBasePath задаёт корневой каталог для поиска GLSL-файлов.
func SetBasePath(base string) {
	shaderBasePath = base
}

var currentLightingIndex = 0

// InitLightingVariants компилирует и линкует все 10 шейдерных программ.
// Для каждого LightingVariant читает .vert и .frag, компилирует, линкует,
// сохраняет programID. Удаляет объекты шейдеров после линковки.
// Возвращает ошибку при неудачной компиляции или линковке.
// Вызывается один раз при запуске (initOpenGL в cmd/main.go).
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

func GetCurrentLightingProgram() uint32 {
	if currentLightingIndex < 0 || currentLightingIndex >= len(lightingVariants) {
		currentLightingIndex = 0
	}
	return lightingVariants[currentLightingIndex].Program
}

func GetCurrentLightingName() string {
	if currentLightingIndex < 0 || currentLightingIndex >= len(lightingVariants) {
		currentLightingIndex = 0
	}
	return lightingVariants[currentLightingIndex].Name
}

func GetCurrentLightingModel() string {
	if currentLightingIndex < 0 || currentLightingIndex >= len(lightingVariants) {
		currentLightingIndex = 0
	}
	return lightingVariants[currentLightingIndex].Model
}

func GetCurrentShadingMode() ShadingMode {
	if currentLightingIndex < 0 || currentLightingIndex >= len(lightingVariants) {
		currentLightingIndex = 0
	}
	return lightingVariants[currentLightingIndex].Mode
}

// ToggleShadingMode переключает между Gouraud и Phong для текущей модели.
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

// CycleLightingVariant переключает модель освещения.
func CycleLightingVariant(forward bool) {
	if forward {
		currentLightingIndex = (currentLightingIndex + 1) % len(lightingVariants)
	} else {
		currentLightingIndex = (currentLightingIndex - 1 + len(lightingVariants)) % len(lightingVariants)
	}
}

func GetLightingVariantCount() int {
	return len(lightingVariants)
}

func CleanupLightingVariants() {
	for _, variant := range lightingVariants {
		if variant.Program != 0 {
			gl.DeleteProgram(variant.Program)
		}
	}
}
