package objects

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

var (
	// Вершины куба с монотонным цветом
	// Первые три значения (x, y, z) - позиция вершины
	// Следующие три значения (r, g, b) - цвет вершины
	ColoredCubeVertices = []float32{
		// Передняя грань 
		-0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0, 0.0,

		// Задняя грань
		-0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 0.0, 0.0,

		// Левая грань 
		-0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		-0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 0.0, 0.0,

		// Правая грань 
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0,

		// Верхняя грань 
		-0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0,

		// Нижняя грань
		-0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		-0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
	}

	// Индексы для отрисовки треугольников
	ColoredCubeIndices = []uint32{
		// Передняя грань
		0, 1, 2,
		2, 3, 0,

		// Задняя грань
		4, 5, 6,
		6, 7, 4,

		// Левая грань
		8, 9, 10,
		10, 11, 8,

		// Правая грань
		12, 13, 14,
		14, 15, 12,

		// Верхняя грань
		16, 17, 18,
		18, 19, 16,

		// Нижняя грань
		20, 21, 22,
		22, 23, 20,
	}
)

// Создание цветного куба
func CreateColoredCube() (uint32, uint32, uint32) {
	// Создание и настройка VAO
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Создание и настройка VBO
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(ColoredCubeVertices)*4, gl.Ptr(ColoredCubeVertices), gl.STATIC_DRAW)

	// Создание и настройка EBO
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(ColoredCubeIndices)*4, gl.Ptr(ColoredCubeIndices), gl.STATIC_DRAW)

	// Настройка атрибута позиции (location = 0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, nil)
	gl.EnableVertexAttribArray(0)

	// Настройка атрибута цвета (location = 1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, unsafe.Pointer(uintptr(3*4)))
	gl.EnableVertexAttribArray(1)

	// Очистка биндингов
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return vao, vbo, ebo
}