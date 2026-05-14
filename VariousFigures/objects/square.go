package objects

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Вершины квадрата
// Формат: позиция (x, y, z) и цвет (r, g, b) - цвет будет переопределен в шейдере
var StripedSquareVertices = []float32{
	// Первый треугольник
	-0.8, -0.8, 0.0, 1.0, 1.0, 1.0, // временный белый цвет
	0.8, -0.8, 0.0, 1.0, 1.0, 1.0,
	0.8, 0.8, 0.0, 1.0, 1.0, 1.0,
	
	// Второй треугольник
	-0.8, -0.8, 0.0, 1.0, 1.0, 1.0,
	0.8, 0.8, 0.0, 1.0, 1.0, 1.0,
	-0.8, 0.8, 0.0, 1.0, 1.0, 1.0,
}

// Создание квадрата для полосатой заливки
func CreateStripedSquare() (uint32, uint32) {
	// Создание и настройка VAO
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Создание и настройка VBO
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(StripedSquareVertices)*4, gl.Ptr(StripedSquareVertices), gl.STATIC_DRAW)

	// Настройка атрибута позиции (location = 0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, nil)
	gl.EnableVertexAttribArray(0)

	// Настройка атрибута цвета (location = 1) - но он будет переопределен в шейдере
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, unsafe.Pointer(uintptr(3*4)))
	gl.EnableVertexAttribArray(1)

	// Очистка биндингов
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return vao, vbo
}