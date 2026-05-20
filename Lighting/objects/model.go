// objects/model.go — Загрузка OBJ-файлов и управление геометрией OpenGL.
//
// Назначение: парсинг Wavefront OBJ-файлов (позиции, нормали, текстурные
// координаты, грани), создание буферов OpenGL (VAO/VBO), отрисовка и удаление.
//
// Ключевые типы:
//   Vertex       — вершина: позиция (3), нормаль (3), UV (2).
//   Model        — OpenGL-представление модели (VAO, VBO, количество вершин).
//   OBJFaceIndex — парсерный тип: индексы позиции/текстуры/нормали для грани OBJ.
//
// Ключевые функции:
//   LoadOBJ(path)                 — полный парсинг OBJ-файла → Model.
//   CreateModelFromVertices()     — создание Model из массива Vertex.
//   CreateModel()                 — создание Model из плоского массива float32.
//   Model.Draw()                  — отрисовка через glDrawArrays.
//   Model.Delete()                — удаление VAO/VBO.
//
// Зависимости: используется в main() для загрузки снеговика, сердца и default.
package objects

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// Vertex — структура вершины для передачи в OpenGL.
// Упакована для соответствия атрибутам шейдера:
//   location 0: Position (vec3) — смещение 0
//   location 1: Normal (vec3)   — смещение 12 байт (3*4)
//   location 3: UV (vec2)       — смещение 24 байта (6*4)
// Общий размер: 8 * 4 = 32 байта на вершину.
type Vertex struct {
	Position [3]float32 // координаты вершины (object space)
	Normal   [3]float32 // нормаль вершины (object space)
	UV       [2]float32 // текстурные координаты (u, v)
}

// Model — OpenGL-представление 3D-модели.
type Model struct {
	VAO         uint32 // Vertex Array Object — хранит состояние атрибутов
	VBO         uint32 // Vertex Buffer Object — буфер с данными вершин
	VertexCount int32  // количество вершин для glDrawArrays
}

// LoadOBJ — загружает и парсит Wavefront OBJ-файл.
// Принимает: path — путь к .obj файлу.
// Возвращает: указатель на Model (VAO/VBO уже созданы).
// Алгоритм:
// 1. Чтение строк файла: v (позиции), vt (текстурные координаты), vn (нормали).
// 2. Для каждой грани f — триангуляция (разбиение на треугольники),
//    создание вершин с интерполяцией нормалей.
// 3. Нормали накапливаются и усредняются для сглаженного шейдинга.
// 4. Создание VAO/VBO через CreateModelFromVertices().
func LoadOBJ(path string) (*Model, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open OBJ file %s: %w", path, err)
	}
	defer file.Close()

	positions := make([][3]float32, 0)   // сырые позиции из строк v
	texCoords := make([][2]float32, 0)   // текстурные координаты из vt
	normals := make([][3]float32, 0)     // нормали из vn
	vertices := make([]Vertex, 0)        // результирующий массив вершин
	normalAccum := make(map[int][3]float32) // накопление нормалей для усреднения

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}

		switch fields[0] {
		case "v":
			// Вершина: v x y z
			if len(fields) < 4 {
				continue
			}
			x, err1 := strconv.ParseFloat(fields[1], 32)
			y, err2 := strconv.ParseFloat(fields[2], 32)
			z, err3 := strconv.ParseFloat(fields[3], 32)
			if err1 != nil || err2 != nil || err3 != nil {
				return nil, fmt.Errorf("invalid vertex data in OBJ: %s", line)
			}
			positions = append(positions, [3]float32{float32(x), float32(y), float32(z)})
		case "vt":
			// Текстурная координата: vt u v
			if len(fields) < 3 {
				continue
			}
			u, err1 := strconv.ParseFloat(fields[1], 32)
			v, err2 := strconv.ParseFloat(fields[2], 32)
			if err1 != nil || err2 != nil {
				continue
			}
			texCoords = append(texCoords, [2]float32{float32(u), float32(v)})
		case "vn":
			// Нормаль: vn nx ny nz
			if len(fields) < 4 {
				continue
			}
			nx, err1 := strconv.ParseFloat(fields[1], 32)
			ny, err2 := strconv.ParseFloat(fields[2], 32)
			nz, err3 := strconv.ParseFloat(fields[3], 32)
			if err1 != nil || err2 != nil || err3 != nil {
				continue
			}
			normals = append(normals, [3]float32{float32(nx), float32(ny), float32(nz)})
		case "f":
			// Грань: f v1/vt1/vn1 v2/vt2/vn2 v3/vt3/vn3 ...
			if len(fields) < 4 {
				continue
			}

			faceIndices := make([]OBJFaceIndex, 0, len(fields)-1)
			for _, part := range fields[1:] {
				idx, err := parseOBJFaceIndex(part, len(positions), len(texCoords), len(normals))
				if err != nil {
					return nil, err
				}
				faceIndices = append(faceIndices, idx)
			}

			// Триангуляция: веер треугольников от первой вершины.
			for i := 1; i < len(faceIndices)-1; i++ {
				v0 := faceIndices[0]
				v1 := faceIndices[i]
				v2 := faceIndices[i+1]

				vertices = append(vertices, createVertexFromOBJ(positions, texCoords, normals, v0, len(vertices)))
				vertices = append(vertices, createVertexFromOBJ(positions, texCoords, normals, v1, len(vertices)))
				vertices = append(vertices, createVertexFromOBJ(positions, texCoords, normals, v2, len(vertices)))

				// Накопление нормалей для последующего усреднения (smooth shading).
				vertIdx0 := len(vertices) - 3
				vertIdx1 := len(vertices) - 2
				vertIdx2 := len(vertices) - 1

				if v0.Normal >= 0 {
					normalAccum[vertIdx0] = add3f(normalAccum[vertIdx0], normals[v0.Normal])
				}
				if v1.Normal >= 0 {
					normalAccum[vertIdx1] = add3f(normalAccum[vertIdx1], normals[v1.Normal])
				}
				if v2.Normal >= 0 {
					normalAccum[vertIdx2] = add3f(normalAccum[vertIdx2], normals[v2.Normal])
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan OBJ file %s: %w", path, err)
	}

	if len(vertices) == 0 {
		return nil, fmt.Errorf("no geometry found in OBJ file %s", path)
	}

	// Нормализация накопленных нормалей (усреднение).
	for i := range vertices {
		if norm, ok := normalAccum[i]; ok {
			vertices[i].Normal = normalize3f(norm)
		}
	}

	return CreateModelFromVertices(vertices)
}

// OBJFaceIndex — индексы позиции, текстурной координаты и нормали для одной
// вершины грани OBJ-файла. Формат: v/vt/vn, где vt и vn опциональны.
type OBJFaceIndex struct {
	Position int // индекс позиции (обязателен)
	TexCoord int // индекс текстурной координаты (-1 если отсутствует)
	Normal   int // индекс нормали (-1 если отсутствует)
}

// parseOBJFaceIndex — разбирает один токен грани OBJ (например "1/2/3").
// Поддерживает отрицательные индексы (относительные) и частичное отсутствие
// текстурных координат или нормалей (например "1//3").
func parseOBJFaceIndex(token string, posCount, texCount, normCount int) (OBJFaceIndex, error) {
	parts := strings.Split(token, "/")
	var idx OBJFaceIndex

	if len(parts) < 1 || parts[0] == "" {
		return idx, fmt.Errorf("invalid face token: %s", token)
	}

	posIdx, err := strconv.Atoi(parts[0])
	if err != nil {
		return idx, fmt.Errorf("invalid OBJ position index: %s", token)
	}
	if posIdx < 0 {
		posIdx = posCount + posIdx + 1
	}
	idx.Position = posIdx - 1

	if len(parts) > 1 && parts[1] != "" {
		texIdx, err := strconv.Atoi(parts[1])
		if err == nil {
			if texIdx < 0 {
				texIdx = texCount + texIdx + 1
			}
			idx.TexCoord = texIdx - 1
		} else {
			idx.TexCoord = -1
		}
	} else {
		idx.TexCoord = -1
	}

	if len(parts) > 2 && parts[2] != "" {
		normIdx, err := strconv.Atoi(parts[2])
		if err == nil {
			if normIdx < 0 {
				normIdx = normCount + normIdx + 1
			}
			idx.Normal = normIdx - 1
		} else {
			idx.Normal = -1
		}
	} else {
		idx.Normal = -1
	}

	return idx, nil
}

// createVertexFromOBJ — создаёт Vertex из данных OBJ-файла по индексам грани.
// Если нормаль или UV отсутствуют — используются значения по умолчанию.
func createVertexFromOBJ(positions [][3]float32, texCoords [][2]float32, normals [][3]float32, idx OBJFaceIndex, vertexCount int) Vertex {
	v := Vertex{}

	if idx.Position >= 0 && idx.Position < len(positions) {
		v.Position = positions[idx.Position]
	}

	if idx.TexCoord >= 0 && idx.TexCoord < len(texCoords) {
		v.UV = texCoords[idx.TexCoord]
	} else {
		v.UV = [2]float32{0.0, 0.0}
	}

	if idx.Normal >= 0 && idx.Normal < len(normals) {
		v.Normal = normals[idx.Normal]
	} else {
		// Если нормаль не указана — временное значение (будет усреднено).
		v.Normal = [3]float32{0.0, 1.0, 0.0}
	}

	return v
}

// normalize3f — нормализует вектор из трёх float32.
// Возвращает единичный вектор. Если длина < 0.0001 — возвращает (0,1,0).
func normalize3f(v [3]float32) [3]float32 {
	length := math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2]))
	if length < 0.0001 {
		return [3]float32{0, 1, 0}
	}
	invLen := float32(1.0 / length)
	return [3]float32{v[0] * invLen, v[1] * invLen, v[2] * invLen}
}

// add3f — покомпонентное сложение двух векторов [3]float32.
func add3f(a, b [3]float32) [3]float32 {
	return [3]float32{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

// CreateModelFromVertices — создаёт VAO и VBO из массива структур Vertex.
// Формат атрибутов:
//   location 0: vec3 position (12 байт)
//   location 1: vec3 normal   (12 байт, смещение 12)
//   location 3: vec2 uv       (8 байт, смещение 24)
//   stride: 32 байта (8 * 4)
// Возвращает: указатель на Model (VAO/VBO).
func CreateModelFromVertices(vertices []Vertex) (*Model, error) {
	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*8*4, gl.Ptr(&vertices[0].Position[0]), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(3, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(3)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return &Model{
		VAO:         vao,
		VBO:         vbo,
		VertexCount: int32(len(vertices)),
	}, nil
}

// CreateModel — создаёт Model из плоского массива float32.
// Формат: позиция (3) + нормаль (3) на каждую вершину (без UV).
// Используется для процедурно сгенерированной геометрии.
func CreateModel(vertices []float32) (*Model, error) {
	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return &Model{
		VAO:         vao,
		VBO:         vbo,
		VertexCount: int32(len(vertices) / 6),
	}, nil
}

// Draw — отрисовывает модель через glDrawArrays(GL_TRIANGLES).
// Побочные эффекты: изменяет привязанный VAO.
func (m *Model) Draw() {
	gl.BindVertexArray(m.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, m.VertexCount)
	gl.BindVertexArray(0)
}

// Delete — удаляет VAO и VBO модели в OpenGL.
// Побочные эффекты: освобождает ресурсы GPU.
func (m *Model) Delete() {
	gl.DeleteVertexArrays(1, &m.VAO)
	gl.DeleteBuffers(1, &m.VBO)
}
