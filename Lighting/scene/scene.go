package scene

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/venexene/comp-graphics-go/lighting"
	"github.com/venexene/comp-graphics-go/objects"
)

// CreateWhiteTexture returns a 1×1 white texture for default material.
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

// setTransformUniforms uploads view, projection and model matrices to the shader.
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

// setMaterialUniforms uploads material properties to the shader.
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
}

// setLightUniforms uploads light source parameters to the shader.
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

// drawObject renders a single model with its model matrix and full lighting uniforms.
func drawObject(model *objects.Model, modelMat mgl32.Mat4, view, proj mgl32.Mat4,
	mat *lighting.MaterialConfig, l *lighting.LightConfig, u *lighting.UniformCache) {

	setTransformUniforms(modelMat, view, proj, u)
	setMaterialUniforms(mat, u)
	setLightUniforms(l, u)

	if u.ViewPos != -1 {
		// viewPos is set once before the loop, but re-set here to be safe
		// (it doesn't change per object, but some shader programs may need it)
	}

	if u.DiffuseMap != -1 {
		gl.Uniform1i(u.DiffuseMap, 0)
	}

	model.Draw()
}

// DrawScene clears the framebuffer and renders all scene objects with the given shader program.
func DrawScene(
	program uint32,
	mainModel *objects.Model,
	mainState *ObjectState,
	extras []*SceneObject,
	cam *Camera,
	projection mgl32.Mat4,
	lightCfg *lighting.LightConfig,
	matCfg *lighting.MaterialConfig,
	defaultTex uint32,
) {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(program)

	// Cache uniform locations
	var u lighting.UniformCache
	u.Refresh(program)

	view := cam.ViewMatrix()
	eye := cam.EyePosition()

	if u.ViewPos != -1 {
		gl.Uniform3f(u.ViewPos, eye.X(), eye.Y(), eye.Z())
	}

	// Activate default texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, defaultTex)

	// Draw main model
	if mainModel != nil {
		drawObject(mainModel, mainState.ModelMatrix(), view, projection, matCfg, lightCfg, &u)
	}

	// Draw extra scene objects
	for _, obj := range extras {
		if obj == nil || obj.Model == nil {
			continue
		}
		drawObject(obj.Model, obj.ModelMatrix(), view, projection, matCfg, lightCfg, &u)
	}
}
