package objects

import (
	"unsafe"
	"math"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Вершины правильного пятиугольника
func CreatePentagon() (uint32, uint32) {
	const numTriangles = 5
	const verticesPerTriangle = 3
	const floatsPerVertex = 6 // позиция (3) + цвет (3)
	
	vertices := make([]float32, 0, numTriangles*verticesPerTriangle*floatsPerVertex)
	
	radius := float32(0.8)
	centerX, centerY := 0.0, 0.0
	
	// Монотонный цвет для пятиугольника (зеленый)
	r, g, b := float32(0.2), float32(0.8), float32(0.2)
	
	for i := 0; i < numTriangles; i++ {
		angle1 := float64(i) * 2 * math.Pi / float64(numTriangles)
		angle2 := float64(i+1) * 2 * math.Pi / float64(numTriangles)
		
		// Центр
		vertices = append(vertices, 
			float32(centerX), float32(centerY), 0.0,
			r, g, b)
		
		// Текущая вершина
		x1 := float32(centerX) + radius*float32(math.Cos(angle1))
		y1 := float32(centerY) + radius*float32(math.Sin(angle1))
		vertices = append(vertices, 
			x1, y1, 0.0,
			r, g, b)
		
		// Следующая вершина
		x2 := float32(centerX) + radius*float32(math.Cos(angle2))
		y2 := float32(centerY) + radius*float32(math.Sin(angle2))
		vertices = append(vertices, 
			x2, y2, 0.0,
			r, g, b)
	}
	
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	
	// Настройка атрибута позиции
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, nil)
	gl.EnableVertexAttribArray(0)
	
	// Настройка атрибута цвета
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, unsafe.Pointer(uintptr(3*4)))
	gl.EnableVertexAttribArray(1)
	
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	
	return vao, vbo
}