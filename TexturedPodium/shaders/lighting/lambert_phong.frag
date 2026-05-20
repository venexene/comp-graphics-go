#version 330 core

in Vertex {
	vec3 normal;
	vec3 light_dir;
    vec3 view_dir;
	vec2 uv_coords;
	float distance;
} vert;

uniform sampler2D u_materialTexture;
uniform sampler2D u_numberTexture;
uniform vec3 u_cubeColor;
uniform float u_materialWeight;
uniform float u_numberWeight;
uniform float linear_coef;
uniform float quadratic_coef;
uniform int attenuation_mode;

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
    vec3 norm = normalize(vert.normal);
    vec3 light_dir = normalize(vert.light_dir);
    float norm_d_light = max(dot(norm, light_dir), 0.0); 
    
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
    
    // Multi-texture blending: surface colour before lighting
    vec3 matColor = texture(u_materialTexture, vert.uv_coords).rgb;
    vec3 numColor = texture(u_numberTexture, vert.uv_coords).rgb;
    float totalWeight = u_materialWeight + u_numberWeight;
    vec3 surfaceColor = u_cubeColor * (u_materialWeight * matColor + u_numberWeight * numColor) / max(totalWeight, 0.001);

    vec3 ambient = light.ambient * surfaceColor * light.ambient_strength * attenuation;
    vec3 diffuse = light.diffuse * norm_d_light * surfaceColor * attenuation;

    frag_color = vec4(ambient + diffuse, 1.0);
}