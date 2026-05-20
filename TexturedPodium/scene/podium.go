
package scene

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/objects"
)

type PodiumCube struct {
	Model       *objects.Model
	Position    mgl32.Vec3
	Scale       float32
	NumberTexID uint32    
	MatTexID    uint32    
	Color       mgl32.Vec3 
}

func (c *PodiumCube) ModelMatrix() mgl32.Mat4 {
	s := mgl32.Scale3D(c.Scale, c.Scale, c.Scale)
	t := mgl32.Translate3D(c.Position.X(), c.Position.Y(), c.Position.Z())
	return t.Mul4(s)
}

type Podium struct {
	Cubes    [4]PodiumCube
	Heart    *objects.Model 
	HeartTex uint32         
	HeartCol mgl32.Vec3     
}

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

func (p *Podium) SetHeart(model *objects.Model, tex uint32, color mgl32.Vec3) {
	p.Heart = model
	p.HeartTex = tex
	p.HeartCol = color
}

func (p *Podium) Delete() {
	if p.Cubes[0].Model != nil {
		p.Cubes[0].Model.Delete()
	}
	if p.Heart != nil {
		p.Heart.Delete()
	}
}
