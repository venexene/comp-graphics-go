package scene

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/venexene/comp-graphics-go/lighting"
	"github.com/venexene/comp-graphics-go/objects"
	"github.com/venexene/comp-graphics-go/textures"
)

type OrangeScene struct {
	Model        *objects.Model
	DiffuseTex   uint32
	NormalTex    uint32
	AOTex        uint32
	RoughnessTex uint32
}

func NewOrangeScene(basePath string) (*OrangeScene, error) {
	s := &OrangeScene{}
	var err error
	s.Model, err = objects.LoadOBJ(basePath + "/models/orange.obj")
	if err != nil {
		return nil, err
	}
	s.DiffuseTex, err = textures.LoadTexture(basePath + "/textures/orange/food_0022_color_1k.jpg")
	if err != nil {
		return nil, err
	}
	s.NormalTex, err = textures.LoadTexture(basePath + "/textures/orange/food_0022_normal_opengl_1k.png")
	if err != nil {
		return nil, err
	}
	s.AOTex, err = textures.LoadTexture(basePath + "/textures/orange/food_0022_ao_1k.jpg")
	if err != nil {
		return nil, err
	}
	s.RoughnessTex, err = textures.LoadTexture(basePath + "/textures/orange/food_0022_roughness_1k.jpg")
	if err != nil {
		return nil, err
	}
	return s, nil
}
func (s *OrangeScene) Cleanup() {
	if s.Model != nil {
		s.Model.Delete()
	}
	for _, tex := range []uint32{s.DiffuseTex, s.NormalTex, s.AOTex, s.RoughnessTex} {
		if tex != 0 {
			gl.DeleteTextures(1, &tex)
		}
	}
}
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
func bindOrangeTextures(s *OrangeScene, u *lighting.UniformCache) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.DiffuseTex)
	if u.DiffuseMap != -1 {
		gl.Uniform1i(u.DiffuseMap, 0)
	}
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, s.NormalTex)
	if u.NormalMap != -1 {
		gl.Uniform1i(u.NormalMap, 1)
	}
	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, s.AOTex)
	if u.AOMap != -1 {
		gl.Uniform1i(u.AOMap, 2)
	}
	gl.ActiveTexture(gl.TEXTURE3)
	gl.BindTexture(gl.TEXTURE_2D, s.RoughnessTex)
	if u.RoughnessMap != -1 {
		gl.Uniform1i(u.RoughnessMap, 3)
	}
}
func DrawOrangeScene(
	program uint32,
	s *OrangeScene,
	cam *Camera,
	projection mgl32.Mat4,
	lightCfg *lighting.LightConfig,
	matCfg *lighting.MaterialConfig,
	objectState *ObjectState,
) {
	DrawSceneObjects(program, s, nil, cam, projection, lightCfg, matCfg, objectState, nil)
}

type HeartScene struct {
	Model        *objects.Model
	DiffuseTex   uint32
	NormalTex    uint32
	MetalnessTex uint32
	RoughnessTex uint32
}

func NewHeartScene(basePath string) (*HeartScene, error) {
	s := &HeartScene{}
	var err error
	s.Model, err = objects.LoadOBJ(basePath + "/models/heart.obj")
	if err != nil {
		return nil, err
	}
	s.DiffuseTex, err = textures.LoadTexture(basePath + "/textures/ornament2/ChristmasTreeOrnament021_1K-JPG_Color.jpg")
	if err != nil {
		return nil, err
	}
	s.NormalTex, err = textures.LoadTexture(basePath + "/textures/ornament2/ChristmasTreeOrnament021_1K-JPG_NormalGL.jpg")
	if err != nil {
		return nil, err
	}
	s.MetalnessTex, err = textures.LoadTexture(basePath + "/textures/ornament2/ChristmasTreeOrnament021_1K-JPG_Metalness.jpg")
	if err != nil {
		return nil, err
	}
	s.RoughnessTex, err = textures.LoadTexture(basePath + "/textures/ornament2/ChristmasTreeOrnament021_1K-JPG_Roughness.jpg")
	if err != nil {
		return nil, err
	}
	return s, nil
}
func (s *HeartScene) Cleanup() {
	if s.Model != nil {
		s.Model.Delete()
	}
	for _, tex := range []uint32{s.DiffuseTex, s.NormalTex, s.MetalnessTex, s.RoughnessTex} {
		if tex != 0 {
			gl.DeleteTextures(1, &tex)
		}
	}
}
func bindHeartTextures(s *HeartScene, u *lighting.UniformCache) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.DiffuseTex)
	if u.DiffuseMap != -1 {
		gl.Uniform1i(u.DiffuseMap, 0)
	}
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, s.NormalTex)
	if u.NormalMap != -1 {
		gl.Uniform1i(u.NormalMap, 1)
	}
	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, s.MetalnessTex)
	if u.AOMap != -1 {
		gl.Uniform1i(u.AOMap, 2)
	}
	gl.ActiveTexture(gl.TEXTURE3)
	gl.BindTexture(gl.TEXTURE_2D, s.RoughnessTex)
	if u.RoughnessMap != -1 {
		gl.Uniform1i(u.RoughnessMap, 3)
	}
}
func DrawSceneObjects(
	program uint32,
	orange *OrangeScene,
	heart *HeartScene,
	cam *Camera,
	projection mgl32.Mat4,
	lightCfg *lighting.LightConfig,
	matCfg *lighting.MaterialConfig,
	orangeState *ObjectState,
	heartState *ObjectState,
) {
	gl.ClearColor(0.15, 0.15, 0.18, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)
	var u lighting.UniformCache
	u.Refresh(program)
	view := cam.ViewMatrix()
	eye := cam.EyePosition()
	if u.ViewPos != -1 {
		gl.Uniform3f(u.ViewPos, eye.X(), eye.Y(), eye.Z())
	}
	if orange != nil && orangeState != nil {
		modelMat := orangeState.ModelMatrix()
		setTransformUniforms(modelMat, view, projection, &u)
		setMaterialUniforms(matCfg, &u)
		setLightUniforms(lightCfg, &u)
		bindOrangeTextures(orange, &u)
		orange.Model.Draw()
	}
	if heart != nil && heartState != nil {
		modelMat := heartState.ModelMatrix()
		setTransformUniforms(modelMat, view, projection, &u)
		setMaterialUniforms(matCfg, &u)
		setLightUniforms(lightCfg, &u)
		bindHeartTextures(heart, &u)
		heart.Model.Draw()
	}
}
