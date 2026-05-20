package shaders

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Загрузка шейдера из файла
func LoadShaderFile(filename string) (string, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return "", fmt.Errorf("failed to read shader file %s: %v", filename, err)
    }
    return string(data), nil
}

// Компиляция шейдера
func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	// Convert Go string to null-terminated C string and pass with explicit length
	csource, free := gl.Strs(source)
	length := int32(len(source))
	gl.ShaderSource(shader, 1, csource, &length)
	free()
	gl.CompileShader(shader)

	// Проверка результата компиляции (примерно как в initOpenGL)
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