// Пакет objects: загрузка OBJ-моделей (LoadOBJ), программное создание
// геометрии куба (CreateCube), управление VAO/VBO (CreateModelFromVertices).
// Формат вершины: Position (vec3, location 0), Normal (vec3, location 1),
// UV (vec2, location 3). Interleaved VBO, 8 float на вершину, DrawArrays.
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

type Vertex struct {
	Position [3]float32
	Normal   [3]float32
	UV       [2]float32
}

type Model struct {
	VAO         uint32
	VBO         uint32
	VertexCount int32
}

func LoadOBJ(path string) (*Model, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open OBJ file %s: %w", path, err)
	}
	defer file.Close()

	positions := make([][3]float32, 0)
	texCoords := make([][2]float32, 0)
	normals := make([][3]float32, 0)
	vertices := make([]Vertex, 0)
	normalAccum := make(map[int][3]float32)

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
			normals = append(normals, [3]float32{float32(nx), float32(ny), float32(nz)})
		case "f":
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

			
			for i := 1; i < len(faceIndices)-1; i++ {
				v0 := faceIndices[0]
				v1 := faceIndices[i]
				v2 := faceIndices[i+1]

				
				vertices = append(vertices, createVertexFromOBJ(positions, texCoords, normals, v0, len(vertices)))
				vertices = append(vertices, createVertexFromOBJ(positions, texCoords, normals, v1, len(vertices)))
				vertices = append(vertices, createVertexFromOBJ(positions, texCoords, normals, v2, len(vertices)))

				
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

	
	for i := range vertices {
		if norm, ok := normalAccum[i]; ok {
			vertices[i].Normal = normalize3f(norm)
		}
	}

	return CreateModelFromVertices(vertices)
}

type OBJFaceIndex struct {
	Position int
	TexCoord int
	Normal   int
}

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
		
		v.Normal = [3]float32{0.0, 1.0, 0.0}
	}

	return v
}

func normalize3f(v [3]float32) [3]float32 {
	len := math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2]))
	if len < 0.0001 {
		return [3]float32{0, 1, 0}
	}
	invLen := float32(1.0 / len)
	return [3]float32{v[0] * invLen, v[1] * invLen, v[2] * invLen}
}

func add3f(a, b [3]float32) [3]float32 {
	return [3]float32{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

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

func (m *Model) Draw() {
	gl.BindVertexArray(m.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, m.VertexCount)
	gl.BindVertexArray(0)
}

func (m *Model) Delete() {
	gl.DeleteVertexArrays(1, &m.VAO)
	gl.DeleteBuffers(1, &m.VBO)
}

func CreateCube(size float32) *Model {
	h := size * 0.5

	
	type face struct {
		normal [3]float32
		verts  [4][3]float32 
	}

	faces := []face{
		{ 
			normal: [3]float32{1, 0, 0},
			verts: [4][3]float32{
				{h, -h, -h}, {h, -h, h}, {h, h, h}, {h, h, -h},
			},
		},
		{ 
			normal: [3]float32{-1, 0, 0},
			verts: [4][3]float32{
				{-h, -h, -h}, {-h, -h, h}, {-h, h, h}, {-h, h, -h},
			},
		},
		{ 
			normal: [3]float32{0, 1, 0},
			verts: [4][3]float32{
				{-h, h, -h}, {h, h, -h}, {h, h, h}, {-h, h, h},
			},
		},
		{ 
			normal: [3]float32{0, -1, 0},
			verts: [4][3]float32{
				{-h, -h, -h}, {h, -h, -h}, {h, -h, h}, {-h, -h, h},
			},
		},
		{ 
			normal: [3]float32{0, 0, 1},
			verts: [4][3]float32{
				{-h, -h, h}, {h, -h, h}, {h, h, h}, {-h, h, h},
			},
		},
		{ 
			normal: [3]float32{0, 0, -1},
			verts: [4][3]float32{
				{h, -h, -h}, {-h, -h, -h}, {-h, h, -h}, {h, h, -h},
			},
		},
	}

	
	uvs := [4][2]float32{
		{0, 0}, {1, 0}, {1, 1}, {0, 1},
	}

	
	vertices := make([]Vertex, 36)
	idx := 0
	for _, f := range faces {
		
		
		triIdx := [6]int{0, 1, 2, 0, 2, 3}
		triUV := [6]int{0, 1, 2, 0, 2, 3}
		for k := 0; k < 6; k++ {
			vi := triIdx[k]
			vertices[idx] = Vertex{
				Position: f.verts[vi],
				Normal:   f.normal,
				UV:       uvs[triUV[k]],
			}
			idx++
		}
	}

	model, _ := CreateModelFromVertices(vertices)
	return model
}
