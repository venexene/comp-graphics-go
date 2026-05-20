// Пакет scene: Podium — четыре кубика пьедестала + сердце.
// PodiumCube — один кубик с материалом, номером и цветом.
// Все кубики разделяют один VAO, созданный через objects.CreateCube.
package scene

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/objects"
)

// PodiumCube — один кубик пьедестала.
// Model — общая геометрия (один VAO на все 4 кубика).
// Position — смещение центра куба в world space.
// Scale — масштаб (1.0 = без изменений).
// NumberTexID — OpenGL ID текстуры номера (1.jpg, 2.png, 3.png).
// MatTexID — OpenGL ID текстуры материала (metal, marble, wood).
// Color — базовый цвет кубика (умножается на смешанные текстуры).
type PodiumCube struct {
	Model       *objects.Model
	Position    mgl32.Vec3
	Scale       float32
	NumberTexID uint32
	MatTexID    uint32
	Color       mgl32.Vec3
}

// ModelMatrix возвращает матрицу Model = Translate(Position) * Scale(Scale).
func (c *PodiumCube) ModelMatrix() mgl32.Mat4 {
	s := mgl32.Scale3D(c.Scale, c.Scale, c.Scale)
	t := mgl32.Translate3D(c.Position.X(), c.Position.Y(), c.Position.Z())
	return t.Mul4(s)
}

// Podium управляет четырьмя кубиками и опциональным сердцем.
// Cubes — массив из 4 кубиков (нижний ряд: центр, лево, право; верх: центр).
// Heart — модель сердца (OBJ), размещается на нижнем центральном кубике.
// HeartTex — текстура материала для сердца (onyx.jpg).
// HeartCol — цвет сердца (красный).
type Podium struct {
	Cubes    [4]PodiumCube
	Heart    *objects.Model
	HeartTex uint32
	HeartCol mgl32.Vec3
}

// NewPodium создаёт пьедестал из 4 кубиков.
// cubeSize — сторона куба, spacing — расстояние между центрами.
// numberTexIDs — карта номер→текстура: 1→1.jpg, 2→2.png, 3→3.png.
// matTexIDs — карта номер→материал: 1→metal, 2→marble, 3→wood.
// cubeColors — карта номер→цвет: 1→жёлтый, 2→серый, 3→оранжевый.
func NewPodium(cubeSize, spacing float32, numberTexIDs, matTexIDs map[int]uint32, cubeColors map[int]mgl32.Vec3) *Podium {
	cubeModel := objects.CreateCube(cubeSize)

	p := &Podium{}

	p.Cubes[0] = PodiumCube{
		Model:       cubeModel,
		Position:    mgl32.Vec3{0, 0, 0},
		Scale:       1.0,
		NumberTexID: numberTexIDs[1],
		MatTexID:    matTexIDs[1],
		Color:       cubeColors[1],
	}

	p.Cubes[1] = PodiumCube{
		Model:       cubeModel,
		Position:    mgl32.Vec3{-spacing, 0, 0},
		Scale:       1.0,
		NumberTexID: numberTexIDs[2],
		MatTexID:    matTexIDs[2],
		Color:       cubeColors[2],
	}

	p.Cubes[2] = PodiumCube{
		Model:       cubeModel,
		Position:    mgl32.Vec3{spacing, 0, 0},
		Scale:       1.0,
		NumberTexID: numberTexIDs[3],
		MatTexID:    matTexIDs[3],
		Color:       cubeColors[3],
	}

	p.Cubes[3] = PodiumCube{
		Model:       cubeModel,
		Position:    mgl32.Vec3{0, cubeSize, 0},
		Scale:       1.0,
		NumberTexID: numberTexIDs[1],
		MatTexID:    matTexIDs[1],
		Color:       cubeColors[1],
	}

	return p
}

// SetHeart размещает сердце на нижнем центральном кубике.
func (p *Podium) SetHeart(model *objects.Model, tex uint32, color mgl32.Vec3) {
	p.Heart = model
	p.HeartTex = tex
	p.HeartCol = color
}

// Delete освобождает VAO кубиков и сердца.
func (p *Podium) Delete() {
	if p.Cubes[0].Model != nil {
		p.Cubes[0].Model.Delete()
	}
	if p.Heart != nil {
		p.Heart.Delete()
	}
}
