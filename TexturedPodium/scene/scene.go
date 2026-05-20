package scene

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/lighting"
	"github.com/venexene/comp-graphics-go/objects"
)

func CreateWhiteTexture() uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	whitePixel := []uint8{255, 255, 255, 255}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(whitePixel))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}

func setTransformUniforms(modelMat mgl32.Mat4, view, proj mgl32.Mat4, u *lighting.UniformCache) {
	if u.View != -1 {
		gl.UniformMatrix4fv(u.View, 1, false, &view[0])
	}
	if u.Projection != -1 {
		gl.UniformMatrix4fv(u.Projection, 1, false, &proj[0])
	}
	if u.Model != -1 {
		gl.UniformMatrix4fv(u.Model, 1, false, &modelMat[0])
	}
	if u.Normal != -1 {
		m := modelMat.Inv().Transpose()
		normalMat := mgl32.Mat3{
			m[0], m[1], m[2],
			m[4], m[5], m[6],
			m[8], m[9], m[10],
		}
		gl.UniformMatrix3fv(u.Normal, 1, false, &normalMat[0])
	}
}

func setMaterialUniforms(mat *lighting.MaterialConfig, u *lighting.UniformCache) {
	if u.MaterialAmbient != -1 {
		gl.Uniform3f(u.MaterialAmbient, mat.Ambient.X(), mat.Ambient.Y(), mat.Ambient.Z())
	}
	if u.MaterialDiffuse != -1 {
		gl.Uniform3f(u.MaterialDiffuse, mat.Diffuse.X(), mat.Diffuse.Y(), mat.Diffuse.Z())
	}
	if u.MaterialSpecular != -1 {
		gl.Uniform3f(u.MaterialSpecular, mat.Specular.X(), mat.Specular.Y(), mat.Specular.Z())
	}
	if u.MaterialSheen != -1 {
		gl.Uniform1f(u.MaterialSheen, mat.SheenCoef)
	}
	if u.Roughness != -1 {
		gl.Uniform1f(u.Roughness, mat.Roughness)
	}
}

func setLightUniforms(l *lighting.LightConfig, u *lighting.UniformCache) {
	if u.LightAmbient != -1 {
		gl.Uniform3f(u.LightAmbient, l.Ambient.X(), l.Ambient.Y(), l.Ambient.Z())
	}
	if u.LightDiffuse != -1 {
		gl.Uniform3f(u.LightDiffuse, l.Diffuse.X(), l.Diffuse.Y(), l.Diffuse.Z())
	}
	if u.LightSpecular != -1 {
		gl.Uniform3f(u.LightSpecular, l.Specular.X(), l.Specular.Y(), l.Specular.Z())
	}
	if u.LightPosition != -1 {
		gl.Uniform3f(u.LightPosition, l.Position.X(), l.Position.Y(), l.Position.Z())
	}
	if u.LightConstant != -1 {
		gl.Uniform1f(u.LightConstant, l.Constant)
	}
	if u.LightLinear != -1 {
		gl.Uniform1f(u.LightLinear, l.Linear)
	}
	if u.LightQuadratic != -1 {
		gl.Uniform1f(u.LightQuadratic, l.Quadratic)
	}
	if u.AmbientStrength != -1 {
		gl.Uniform1f(u.AmbientStrength, l.AmbientStrength)
	}
	if u.LinearCoef != -1 {
		gl.Uniform1f(u.LinearCoef, l.LinearCoef)
	}
	if u.QuadraticCoef != -1 {
		gl.Uniform1f(u.QuadraticCoef, l.QuadraticCoef)
	}
	if u.AttenuationMode != -1 {
		gl.Uniform1i(u.AttenuationMode, int32(l.Mode))
	}
}

func drawObject(model *objects.Model, modelMat mgl32.Mat4, view, proj mgl32.Mat4,
	mat *lighting.MaterialConfig, l *lighting.LightConfig, u *lighting.UniformCache) {

	setTransformUniforms(modelMat, view, proj, u)
	setMaterialUniforms(mat, u)
	setLightUniforms(l, u)

	if u.ViewPos != -1 {
		
		
	}

	model.Draw()
}

func setBlendUniforms(matTexID, numTexID uint32, color mgl32.Vec3, matWeight, numWeight float32, u *lighting.UniformCache) {
	
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, matTexID)
	if u.MaterialMap != -1 {
		gl.Uniform1i(u.MaterialMap, 0)
	}

	
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, numTexID)
	if u.NumberMap != -1 {
		gl.Uniform1i(u.NumberMap, 1)
	}

	if u.CubeColor != -1 {
		gl.Uniform3f(u.CubeColor, color.X(), color.Y(), color.Z())
	}
	if u.MaterialWeight != -1 {
		gl.Uniform1f(u.MaterialWeight, matWeight)
	}
	if u.NumberWeight != -1 {
		gl.Uniform1f(u.NumberWeight, numWeight)
	}
}

func DrawPodium(
	program uint32,
	p *Podium,
	cam *Camera,
	projection mgl32.Mat4,
	lightCfg *lighting.LightConfig,
	matCfg *lighting.MaterialConfig,
	materialWeight, numberWeight float32,
) {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(program)

	var u lighting.UniformCache
	u.Refresh(program)

	view := cam.ViewMatrix()
	eye := cam.EyePosition()

	if u.ViewPos != -1 {
		gl.Uniform3f(u.ViewPos, eye.X(), eye.Y(), eye.Z())
	}

	for i := range p.Cubes {
		cube := &p.Cubes[i]
		if cube.Model == nil {
			continue
		}
		setTransformUniforms(cube.ModelMatrix(), view, projection, &u)
		setMaterialUniforms(matCfg, &u)
		setLightUniforms(lightCfg, &u)
		setBlendUniforms(cube.MatTexID, cube.NumberTexID, cube.Color, materialWeight, numberWeight, &u)
		cube.Model.Draw()
	}

	
	if p.Heart != nil {
		
		
		heartY := float32(0.4 + 0.8) 
		heartScale := float32(0.2)
		heartModel := mgl32.Translate3D(0, heartY, 0).Mul4(mgl32.Scale3D(heartScale, heartScale, heartScale))
		setTransformUniforms(heartModel, view, projection, &u)
		setMaterialUniforms(matCfg, &u)
		setLightUniforms(lightCfg, &u)
		
		setBlendUniforms(p.HeartTex, p.HeartTex, p.HeartCol, 1.0, 0.0, &u)
		p.Heart.Draw()
	}
}
