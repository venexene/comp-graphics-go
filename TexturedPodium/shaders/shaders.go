// Пакет shaders содержит низкоуровневые функции для работы с GLSL-шейдерами:
// загрузка исходного кода из файла (LoadShaderFile), компиляция (CompileShader).
// Зависимости: go-gl/gl — вызовы OpenGL API для компиляции шейдеров.
// Используется пакетом shaders (lighting.go) при инициализации вариантов
// освещения (InitLightingVariants).
package shaders

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// LoadShaderFile читает GLSL-файл с диска и возвращает его содержимое как строку.
// Вход: filename — путь к файлу шейдера (относительный или абсолютный).
// Возвращает: string — исходный код шейдера, error — ошибка чтения файла.
// Вызывается при компиляции шейдерных программ (InitLightingVariants).
func LoadShaderFile(filename string) (string, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return "", fmt.Errorf("failed to read shader file %s: %v", filename, err)
    }
    return string(data), nil
}

// CompileShader компилирует GLSL-код в шейдер указанного типа.
// Вход: source — исходный код шейдера, shaderType — GL_VERTEX_SHADER
// или GL_FRAGMENT_SHADER.
// Возвращает: uint32 — идентификатор скомпилированного шейдера, error — ошибка
// компиляции с логом ошибок.
// Побочные эффекты: создаёт объект шейдера в OpenGL.
// Вызывается при инициализации каждой шейдерной программы.
func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource, free := gl.Strs(source)
	length := int32(len(source))
	gl.ShaderSource(shader, 1, csource, &length)
	free()
	gl.CompileShader(shader)

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
