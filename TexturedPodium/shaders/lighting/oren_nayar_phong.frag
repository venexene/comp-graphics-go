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
uniform float roughness;
uniform float linear_coef;
uniform float quadratic_coef;

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
    vec3 view_dir = normalize(vert.view_dir);

    float attenuation = 1.0 / max(light.constant + 
        (light.linear * linear_coef) * vert.distance + 
        (light.quadratic * quadratic_coef) * vert.distance * vert.distance, 0.0001);

    float norm_d_light = max(dot(norm, light_dir), 0.0);
    float norm_d_view = max(dot(norm, view_dir), 0.0);

    float rough2 = roughness * roughness;
    float A = 1.0 - (0.5 * rough2 / (rough2 + 0.33));
    float B = 0.45 * rough2 / (rough2 + 0.09);

    float sin_theta_i = sqrt(1.0 - norm_d_light * norm_d_light);
    float sin_theta_r = sqrt(1.0 - norm_d_view * norm_d_view);

    float max_cos = 0.0;
    if (sin_theta_i > 0.0001 && sin_theta_r > 0.0001) {
        vec3 light_perp = normalize(light_dir - norm * norm_d_light);
        vec3 view_perp = normalize(view_dir - norm * norm_d_view);
        max_cos = max(0.0, dot(light_perp, view_perp));
    }

    float sin_alpha, tan_beta;
    if (norm_d_light > norm_d_view) {
        sin_alpha = sin_theta_r;
        tan_beta  = sin_theta_i / max(norm_d_light, 0.001);
    }
    else {
        sin_alpha = sin_theta_i;
        tan_beta  = sin_theta_r / max(norm_d_view, 0.001);
    }

    // Multi-texture blending: surface colour before lighting
    vec3 matColor = texture(u_materialTexture, vert.uv_coords).rgb;
    vec3 numColor = texture(u_numberTexture, vert.uv_coords).rgb;
    float totalWeight = u_materialWeight + u_numberWeight;
    vec3 surfaceColor = u_cubeColor * (u_materialWeight * matColor + u_numberWeight * numColor) / max(totalWeight, 0.001);

    float oren_nayar = norm_d_light * (A + B * max_cos * sin_alpha * tan_beta) * attenuation;
    vec3 ambient = light.ambient * surfaceColor * light.ambient_strength * attenuation;
    vec3 final = ambient + light.diffuse * surfaceColor * oren_nayar;

    frag_color = vec4(final, 1.0);
}
