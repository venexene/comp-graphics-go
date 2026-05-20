package lighting

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

type UniformCache struct {
	Model            int32
	View             int32
	Projection       int32
	Normal           int32
	ViewPos          int32
	MaterialAmbient  int32
	MaterialDiffuse  int32
	MaterialSpecular int32
	MaterialSheen    int32
	Roughness        int32
	LightAmbient     int32
	LightDiffuse     int32
	LightSpecular    int32
	LightPosition    int32
	LightConstant    int32
	LightLinear      int32
	LightQuadratic   int32
	AmbientStrength  int32
	LinearCoef       int32
	QuadraticCoef    int32
	AttenuationMode  int32
	DiffuseMap       int32
	NormalMap        int32
	AOMap            int32
	RoughnessMap     int32
}

func (c *UniformCache) Refresh(program uint32) {
	c.Model = gl.GetUniformLocation(program, gl.Str("transform.model\x00"))
	c.View = gl.GetUniformLocation(program, gl.Str("transform.view\x00"))
	c.Projection = gl.GetUniformLocation(program, gl.Str("transform.projection\x00"))
	c.Normal = gl.GetUniformLocation(program, gl.Str("transform.normal_mat\x00"))
	c.ViewPos = gl.GetUniformLocation(program, gl.Str("transform.view_pos\x00"))
	c.MaterialAmbient = gl.GetUniformLocation(program, gl.Str("material.ambient\x00"))
	c.MaterialDiffuse = gl.GetUniformLocation(program, gl.Str("material.diffuse\x00"))
	c.MaterialSpecular = gl.GetUniformLocation(program, gl.Str("material.specular\x00"))
	c.MaterialSheen = gl.GetUniformLocation(program, gl.Str("material.sheen_coef\x00"))
	c.Roughness = gl.GetUniformLocation(program, gl.Str("roughness\x00"))
	c.LightAmbient = gl.GetUniformLocation(program, gl.Str("light.ambient\x00"))
	c.LightDiffuse = gl.GetUniformLocation(program, gl.Str("light.diffuse\x00"))
	c.LightSpecular = gl.GetUniformLocation(program, gl.Str("light.specular\x00"))
	c.LightPosition = gl.GetUniformLocation(program, gl.Str("light.position\x00"))
	c.LightConstant = gl.GetUniformLocation(program, gl.Str("light.constant\x00"))
	c.LightLinear = gl.GetUniformLocation(program, gl.Str("light.linear\x00"))
	c.LightQuadratic = gl.GetUniformLocation(program, gl.Str("light.quadratic\x00"))
	c.AmbientStrength = gl.GetUniformLocation(program, gl.Str("light.ambient_strength\x00"))
	c.LinearCoef = gl.GetUniformLocation(program, gl.Str("linear_coef\x00"))
	c.QuadraticCoef = gl.GetUniformLocation(program, gl.Str("quadratic_coef\x00"))
	c.AttenuationMode = gl.GetUniformLocation(program, gl.Str("attenuation_mode\x00"))
	c.DiffuseMap = gl.GetUniformLocation(program, gl.Str("u_diffuseMap\x00"))
	c.NormalMap = gl.GetUniformLocation(program, gl.Str("u_normalMap\x00"))
	c.AOMap = gl.GetUniformLocation(program, gl.Str("u_aoMap\x00"))
	c.RoughnessMap = gl.GetUniformLocation(program, gl.Str("u_roughnessMap\x00"))
}
