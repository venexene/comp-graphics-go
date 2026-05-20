// Файл: objects/model.go
// Назначение: загрузка 3D-моделей из формата OBJ, вычисление касательных (tangents)
// для normal mapping, создание VAO/VBO для OpenGL.
//
// Ключевые структуры:
//   Vertex — вершина модели с позицией, нормалью, UV-координатами и касательной.
//   Model — OpenGL-объект (VAO + VBO), готовый к отрисовке.
//   OBJFaceIndex — индекс вершины в грани (позиция, текстурная координата, нормаль).
//
// Ключевые функции:
//   LoadOBJ — парсит OBJ-файл, вычисляет нормали и касательные, создаёт Model.
//   parseOBJFaceIndex — разбирает токен грани вида "v/vt/vn".
//   createVertexFromOBJ — собирает Vertex по OBJFaceIndex.
//   smoothNormals — усредняет нормали треугольника по вершинам для сглаживания.
//   ComputeTangents — вычисляет касательные (tangents) для каждого треугольника.
//   CreateModelFromVertices — создаёт VAO/VBO из массива вершин.
//   Model.Draw — отрисовывает модель через glDrawArrays(GL_TRIANGLES).
//   Model.Delete — удаляет VAO и VBO.
//
// Зависимости:
//   Внутренние: — (базовый модуль, не зависит от других пакетов проекта).
//   Внешние: github.com/go-gl/gl/v4.6-core/gl — OpenGL-функции.

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

// Vertex — структура одной вершины модели.
// Поля хранят атрибуты, которые будут переданы в VBO как interleaved массив:
//
//	Offset 0:  Position [3]float32 — координаты вершины (x, y, z), единицы — единицы модели.
//	Offset 12: Normal   [3]float32 — нормаль вершины (nx, ny, nz), единичный вектор.
//	Offset 24: UV       [2]float32 — текстурные координаты (u, v), диапазон [0,1].
//	Offset 32: Tangent  [3]float32 — касательная (tx, ty, tz), единичный вектор.
//	Общий stride вершины = 11 * 4 = 44 байта.
type Vertex struct {
	Position [3]float32
	Normal   [3]float32
	UV       [2]float32
	Tangent  [3]float32
}

// Model — OpenGL-представление модели.
// Поля:
//
//	VAO         — Vertex Array Object (инкапсулирует атрибутные указатели).
//	VBO         — Vertex Buffer Object (буфер с interleaved вершинами).
//	VertexCount — количество вершин для отрисовки (всегда кратно 3, т.к. GL_TRIANGLES).
type Model struct {
	VAO         uint32
	VBO         uint32
	VertexCount int32
}

// LoadOBJ — загрузка OBJ-файла и создание Model.
// Параметры:
//
//	path — путь к OBJ-файлу относительно корня проекта или абсолютный.
//
// Возвращает:
//
//	*Model — VAO/VBO с загруженной геометрией, готовый к отрисовке.
//
// Алгоритм:
// 1. Чтение строк v (позиции), vt (текстурные координаты), vn (нормали).
// 2. Разбор граней f (формат v/vt/vn), триангуляция (fan из первой вершины).
// 3. Сглаживание нормалей (smoothNormals) — усреднение по уникальным позициям.
// 4. Вычисление касательных (ComputeTangents) для normal mapping.
// 5. Создание VAO/VBO через CreateModelFromVertices.
func LoadOBJ(path string) (*Model, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open OBJ file %s: %w", path, err)
	}
	defer file.Close()
	positions := make([][3]float32, 0)
	texCoords := make([][2]float32, 0)
	objNormals := make([][3]float32, 0)
	vertices := make([]Vertex, 0)
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
			if len(fields) < 4 {
				continue
			}
			nx, err1 := strconv.ParseFloat(fields[1], 32)
			ny, err2 := strconv.ParseFloat(fields[2], 32)
			nz, err3 := strconv.ParseFloat(fields[3], 32)
			if err1 != nil || err2 != nil || err3 != nil {
				continue
			}
			objNormals = append(objNormals, [3]float32{float32(nx), float32(ny), float32(nz)})
		case "f":
			if len(fields) < 4 {
				continue
			}
			faceIndices := make([]OBJFaceIndex, 0, len(fields)-1)
			for _, part := range fields[1:] {
				idx, err := parseOBJFaceIndex(part, len(positions), len(texCoords), len(objNormals))
				if err != nil {
					return nil, err
				}
				faceIndices = append(faceIndices, idx)
			}
			for i := 1; i < len(faceIndices)-1; i++ {
				v0 := faceIndices[0]
				v1 := faceIndices[i]
				v2 := faceIndices[i+1]
				vertices = append(vertices, createVertexFromOBJ(positions, texCoords, objNormals, v0))
				vertices = append(vertices, createVertexFromOBJ(positions, texCoords, objNormals, v1))
				vertices = append(vertices, createVertexFromOBJ(positions, texCoords, objNormals, v2))
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan OBJ file %s: %w", path, err)
	}
	if len(vertices) == 0 {
		return nil, fmt.Errorf("no geometry found in OBJ file %s", path)
	}
	smoothNormals(vertices)
	ComputeTangents(vertices)
	return CreateModelFromVertices(vertices)
}

// OBJFaceIndex — индексы одной вершины внутри грани OBJ.
// Поля:
//
//	Position — индекс позиции в массиве positions (0-based).
//	TexCoord — индекс текстурной координаты в texCoords (-1 если отсутствует).
//	Normal   — индекс нормали в objNormals (-1 если отсутствует).
type OBJFaceIndex struct {
	Position int
	TexCoord int
	Normal   int
}

// parseOBJFaceIndex — разбирает строку "v/vt/vn" (или "v//vn", "v") в OBJFaceIndex.
// Параметры:
//
//	token     — строка вида "1/2/3", "1//3" или "1".
//	posCount  — общее количество позиций в модели (для отрицательных индексов).
//	texCount  — общее количество текстурных координат.
//	normCount — общее количество нормалей.
//
// OBJ использует 1-based индексы; отрицательные — отсчёт с конца.
// Функция преобразует их в 0-based.
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
// createVertexFromOBJ — создаёт Vertex из индексов OBJFaceIndex.
// Если индекс отрицательный (не задан), используются значения по умолчанию
// (нулевая позиция, UV=(0,0), нормаль=(0,1,0)).
func createVertexFromOBJ(positions [][3]float32, texCoords [][2]float32, normals [][3]float32, idx OBJFaceIndex) Vertex {
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
		v.Normal = [3]float32{0.0, 1.0, 0.0}
	}
	return v
}
// smoothNormals — усреднение нормалей по вершинам с одинаковой позицией.
// Алгоритм:
// 1. Для каждого треугольника вычисляется его нормаль (cross(edge1, edge2)).
// 2. Нормаль треугольника добавляется к аккумулятору для каждой из 3 вершин.
// 3. В конце для каждой вершины аккумулированная нормаль нормализуется.
// Эффект — сглаженное освещение (Gouraud/Phong shading) без видимых рёбер.
func smoothNormals(vertices []Vertex) {
	if len(vertices)%3 != 0 {
		return
	}
	type posKey struct {
		x, y, z float32
	}
	accum := make(map[posKey][3]float32)
	count := make(map[posKey]int)
	for i := 0; i < len(vertices); i += 3 {
		p0 := vertices[i].Position
		p1 := vertices[i+1].Position
		p2 := vertices[i+2].Position
		fn := normalize3f(cross3f(sub3f(p1, p0), sub3f(p2, p0)))
		for k := 0; k < 3; k++ {
			pos := vertices[i+k].Position
			key := posKey{pos[0], pos[1], pos[2]}
			accum[key] = add3f(accum[key], fn)
			count[key]++
		}
	}
	for i := range vertices {
		pos := vertices[i].Position
		key := posKey{pos[0], pos[1], pos[2]}
		if cnt := count[key]; cnt > 0 {
			vertices[i].Normal = normalize3f(accum[key])
		}
	}
}
// normalize3f — нормализует вектор [3]float32.
// Возвращает единичный вектор. Для нулевого вектора возвращает (0, 1, 0).
func normalize3f(v [3]float32) [3]float32 {
	len := math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2]))
	if len < 0.0001 {
		return [3]float32{0, 1, 0}
	}
	invLen := float32(1.0 / len)
	return [3]float32{v[0] * invLen, v[1] * invLen, v[2] * invLen}
}

// add3f — покомпонентное сложение двух векторов.
func add3f(a, b [3]float32) [3]float32 {
	return [3]float32{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

// sub3f — покомпонентное вычитание векторов (a - b).
func sub3f(a, b [3]float32) [3]float32 {
	return [3]float32{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

// mul3fScalar — умножение вектора на скаляр.
func mul3fScalar(v [3]float32, s float32) [3]float32 {
	return [3]float32{v[0] * s, v[1] * s, v[2] * s}
}

// dot3f — скалярное произведение двух векторов.
func dot3f(a, b [3]float32) float32 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

// cross3f — векторное произведение a × b.
func cross3f(a, b [3]float32) [3]float32 {
	return [3]float32{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}
// ComputeTangents — вычисляет касательные (tangents) для каждой вершины.
// Алгоритм (стандартный для normal mapping):
//
// Для каждого треугольника с вершинами (p0, p1, p2) и UV (uv0, uv1, uv2):
//
//	deltaPos1 = p1 - p0,  deltaPos2 = p2 - p0
//	deltaUV1  = uv1 - uv0, deltaUV2 = uv2 - uv0
//
//	Решаем систему: [deltaPos1, deltaPos2] = T * [deltaUV1.x, deltaUV2.x]
//	                 (только для T, бикасательная B не вычисляется)
//
//	f = 1 / (deltaUV1.x * deltaUV2.y - deltaUV1.y * deltaUV2.x)
//	T = f * (deltaPos1 * deltaUV2.y - deltaPos2 * deltaUV1.y)
//
// После накопления T для каждой вершины применяется ортогонализация
// Грама — Шмидта относительно нормали:
//
//	T' = normalize(T - N * dot(N, T))
//
// Результат — касательная, ортогональная нормали (лежит в касательной плоскости).
func ComputeTangents(vertices []Vertex) {
	if len(vertices)%3 != 0 {
		return
	}
	acc := make([][3]float32, len(vertices))
	for i := 0; i < len(vertices); i += 3 {
		p0 := vertices[i].Position
		p1 := vertices[i+1].Position
		p2 := vertices[i+2].Position
		uv0 := vertices[i].UV
		uv1 := vertices[i+1].UV
		uv2 := vertices[i+2].UV
		deltaPos1 := sub3f(p1, p0)
		deltaPos2 := sub3f(p2, p0)
		duv1 := [2]float32{uv1[0] - uv0[0], uv1[1] - uv0[1]}
		duv2 := [2]float32{uv2[0] - uv0[0], uv2[1] - uv0[1]}
		denom := duv1[0]*duv2[1] - duv1[1]*duv2[0]
		f := float32(1.0)
		if math.Abs(float64(denom)) > 0.000001 {
			f = 1.0 / denom
		}
		t := sub3f(
			mul3fScalar(deltaPos1, duv2[1]),
			mul3fScalar(deltaPos2, duv1[1]),
		)
		t = mul3fScalar(t, f)
		for k := 0; k < 3; k++ {
			acc[i+k] = add3f(acc[i+k], t)
		}
	}
	for i := range vertices {
		n := vertices[i].Normal
		t := acc[i]
		proj := dot3f(n, t)
		tOrtho := sub3f(t, mul3fScalar(n, proj))
		vertices[i].Tangent = normalize3f(tOrtho)
	}
}

// vertexStride — размер одной вершины в байтах.
// 11 компонент float32 * 4 байта = 44 байта на вершину.
// Раскладка: Position(3) + Normal(3) + UV(2) + Tangent(3) = 11 float32.
const vertexStride = 11 * 4

// CreateModelFromVertices — создаёт VAO и VBO из массива Vertex.
// Атрибутные указатели (locations) соответствуют шейдеру orange_phong.vert:
//
//	location 0 — Position (vec3, stride 44, offset 0)
//	location 1 — Normal   (vec3, stride 44, offset 12)
//	location 2 — UV       (vec2, stride 44, offset 24)
//	location 3 — Tangent  (vec3, stride 44, offset 32)
//
// Использует gl.DrawArrays с GL_TRIANGLES (без EBO).
func CreateModelFromVertices(vertices []Vertex) (*Model, error) {
	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*vertexStride, gl.Ptr(&vertices[0].Position[0]), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, vertexStride, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, vertexStride, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, vertexStride, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, vertexStride, gl.PtrOffset(8*4))
	gl.EnableVertexAttribArray(3)
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return &Model{
		VAO:         vao,
		VBO:         vbo,
		VertexCount: int32(len(vertices)),
	}, nil
}
// Draw — отрисовывает модель через glDrawArrays(GL_TRIANGLES).
// Вызывается в главном цикле рендера для каждого кадра.
// Предусловие: шейдерная программа должна быть активна (gl.UseProgram),
// текстуры и uniform-переменные — установлены.
func (m *Model) Draw() {
	gl.BindVertexArray(m.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, m.VertexCount)
	gl.BindVertexArray(0)
}

// Delete — удаляет OpenGL-ресурсы (VAO, VBO).
// Вызывается при завершении программы (cleanup).
func (m *Model) Delete() {
	gl.DeleteVertexArrays(1, &m.VAO)
	gl.DeleteBuffers(1, &m.VBO)
}
