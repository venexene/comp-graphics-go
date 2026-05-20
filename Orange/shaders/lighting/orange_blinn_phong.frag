#version 330 core

in TangentSpace {
	vec3 light_dir;
	vec3 view_dir;
	vec2 uv_coords;
	float distance;
} vert;

// Texture maps
uniform sampler2D u_diffuseMap;
uniform sampler2D u_normalMap;
uniform sampler2D u_aoMap;
uniform sampler2D u_roughnessMap;

uniform float linear_coef;
uniform float quadratic_coef;
uniform int attenuation_mode;

uniform struct Material {
	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
    float sheen_coef;
} material;

uniform struct PointLight {
	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
    vec3 position;

    float constant;
    float linear;
    float quadratic;

    float ambient_strength;
} light;

out vec4 frag_color;

void main() {
	// --- Sample normal from normal map, transform from [0,1] to [-1,1] ---
	vec3 normal = texture(u_normalMap, vert.uv_coords).rgb;
	normal = normalize(normal * 2.0 - 1.0);

	// Directions in tangent space
	vec3 light_dir = normalize(vert.light_dir);
	vec3 view_dir  = normalize(vert.view_dir);

	// --- Attenuation ---
	float attenuation;
	if (attenuation_mode == 1) {
		attenuation = 1.0 / max(light.constant + (light.linear * linear_coef) * vert.distance, 0.0001);
	} else if (attenuation_mode == 2) {
		attenuation = 1.0 / max(light.constant + (light.quadratic * quadratic_coef) * vert.distance * vert.distance, 0.0001);
	} else {
		attenuation = 1.0 / max(light.constant +
			(light.linear * linear_coef) * vert.distance +
			(light.quadratic * quadratic_coef) * vert.distance * vert.distance, 0.0001);
	}

	// --- Surface colour from diffuse map ---
	vec3 surface_color = texture(u_diffuseMap, vert.uv_coords).rgb;

	// Optional: apply ambient occlusion
	float ao = texture(u_aoMap, vert.uv_coords).r;

	// --- Blinn-Phong lighting with perturbed normal ---
	float NdotL = max(dot(normal, light_dir), 0.0);

	// Blinn-Phong specular: (N · H)^shininess  where H = normalize(L + V)
	vec3 half_dir = normalize(light_dir + view_dir);
	float spec = pow(max(dot(normal, half_dir), 0.0), material.sheen_coef);

	vec3 ambient  = light.ambient * surface_color * material.ambient * light.ambient_strength * attenuation * ao;
	vec3 diffuse  = light.diffuse * surface_color * material.diffuse * NdotL * attenuation;
	vec3 specular = light.specular * spec * material.specular * attenuation;

	frag_color = vec4(ambient + diffuse + specular, 1.0);
}
