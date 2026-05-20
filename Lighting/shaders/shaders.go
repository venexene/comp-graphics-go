// shaders/shaders.go — Загрузка и компиляция GLSL-шейдеров.
//
// Назначение: предоставляет функции для чтения исходного кода шейдеров из
// файлов и их компиляции через OpenGL.
//
// Ключевые функции:
//   LoadShaderFile(filename string) (string, error) — читает GLSL-файл в строку.
//   CompileShader(source, shaderType) (uint32, error) — компилирует шейдер.
//
// Зависимости: вызывается из shaders/lighting.go:InitLightingVariants().
package shaders

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// LoadShaderFile — читает содержимое файла шейдера.
// Принимает: filename — путь к .vert или .frag файлу.
// Возвращает: исходный код шейдера как строку.
// Ошибка: если файл не найден или не читается.
func LoadShaderFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read shader file %s: %v", filename, err)
	}
	return string(data), nil
}

// CompileShader — компилирует GLSL-шейдер.
// Принимает:
//   source — исходный код шейдера (строка).
//   shaderType — тип шейдера: gl.VERTEX_SHADER или gl.FRAGMENT_SHADER.
// Возвращает: идентификатор скомпилированного шейдера (uint32).
// Побочные эффекты: создаёт объект шейдера в OpenGL.
// Ошибка: если компиляция не удалась (синтаксическая ошибка в GLSL).
func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	// Передача исходного кода в OpenGL.
	// gl.Strs возвращает C-строку с нулевым окончанием; free() освобождает её.
	csource, free := gl.Strs(source)
	length := int32(len(source))
	gl.ShaderSource(shader, 1, csource, &length)
	free()
	gl.CompileShader(shader)

	// Проверка статуса компиляции.
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength)
		gl.GetShaderInfoLog(shader, logLength, nil, &log[0])

		return 0, fmt.Errorf("failed to compile shader: %s", string(log))
	}

	return shader, nil
}