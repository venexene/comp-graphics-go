// Файл: shaders/shaders.go
// Назначение: загрузка и компиляция GLSL-шейдеров.
//
// Ключевые функции:
//   LoadShaderFile — читает GLSL-файл в строку.
//   CompileShader — компилирует GLSL-код в шейдерный объект OpenGL.
//
// Зависимости:
//   Внутренние: — (используется shaders/lighting.go).
//   Внешние: github.com/go-gl/gl/v4.6-core/gl.

package shaders

import (
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"os"
)

// LoadShaderFile — читает содержимое файла шейдера в строку.
// Параметры:
//
//	filename — путь к GLSL-файлу (обычно .vert или .frag).
//
// Возвращает: исходный код шейдера в виде строки.
func LoadShaderFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read shader file %s: %v", filename, err)
	}
	return string(data), nil
}

// CompileShader — компилирует GLSL-код в шейдерный объект.
// Параметры:
//
//	source      — исходный код шейдера (строка).
//	shaderType  — тип шейдера (gl.VERTEX_SHADER или gl.FRAGMENT_SHADER).
//
// Возвращает: идентификатор скомпилированного шейдера.
// В случае ошибки возвращает 0 и описание ошибки компиляции.
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
