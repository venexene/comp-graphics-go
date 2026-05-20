// shaders/lighting/toon_phong.frag — Фрагментный шейдер Toon (Фонг).
//
// Модель освещения: Toon (дискретные уровни + контур).
// Режим шейдинга: Phong — попиксельное освещение.
// Особенности: квантование diffuse/specular + затемнение силуэтов.

#version 330 core

in Vertex {
	vec3 normal;
	vec3 light_dir;
	vec3 view_dir;
	vec2 uv_coords;
	float distance;
} vert;

uniform sampler2D diffuse_map;
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
    vec3 norm = normalize(vert.normal);
    vec3 light_dir = normalize(vert.light_dir);
    vec3 view_dir = normalize(vert.view_dir);

    // Затухание.
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

    // Диффузное: дискретные уровни (3 ступени).
    float diff = max(dot(norm, light_dir), 0.0);
    float levels = 3.0;
    float toon_diff = floor(diff * levels) / levels;
    toon_diff = max(toon_diff, light.ambient_strength);

    // Зеркальное через половинный вектор: дискретные уровни.
    vec3 half_dir = normalize(light_dir + view_dir);
    float spec = pow(max(dot(norm, half_dir), 0.0), material.sheen_coef);
    float toon_spec = floor(spec * levels) / levels;

    // Контур: smoothstep затемняет пиксели у края объекта.
    float edge = dot(norm, view_dir);
    float edge_factor = smoothstep(0.15, 0.05, edge);

    // Итоговый цвет.
    vec3 base_color = texture(diffuse_map, vert.uv_coords).rgb;
    vec3 ambient = light.ambient * base_color * light.ambient_strength * attenuation;
    vec3 diffuse = light.diffuse * base_color * toon_diff * attenuation;
    vec3 specular = light.specular * material.specular * toon_spec * attenuation;

    frag_color = vec4((ambient + diffuse + specular) * edge_factor, 1.0);
}
