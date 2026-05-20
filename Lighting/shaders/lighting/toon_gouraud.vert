#version 330 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal_in;
layout (location = 3) in vec2 uv_coords_in;

uniform struct Transform {
    mat4 model;
    mat4 view;
    mat4 projection;
    mat3 normal_mat;
    vec3 view_pos;
} transform;

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

uniform sampler2D diffuse_map;
uniform float linear_coef;
uniform float quadratic_coef;
uniform int attenuation_mode;

out vec3 vert_color;
out vec2 uv_coords;

void main() {
    vec4 world_pos = transform.model * vec4(position, 1.0);

    vec3 normal = normalize(transform.normal_mat * normal_in);
    vec3 light_dir = light.position - world_pos.xyz;
    vec3 view_dir = transform.view_pos - world_pos.xyz;
    float distance = length(light_dir);
    light_dir = normalize(light_dir);
    view_dir = normalize(view_dir);

    float attenuation;
    if (attenuation_mode == 1) {
        attenuation = 1.0 / max(light.constant + (light.linear * linear_coef) * distance, 0.0001);
    } else if (attenuation_mode == 2) {
        attenuation = 1.0 / max(light.constant + (light.quadratic * quadratic_coef) * distance * distance, 0.0001);
    } else {
        attenuation = 1.0 / max(light.constant + 
            (light.linear * linear_coef) * distance + 
            (light.quadratic * quadratic_coef) * distance * distance, 0.0001);
    }

    // Diffuse: discrete levels
    float diff = max(dot(normal, light_dir), 0.0);
    float levels = 3.0;
    float toon_diff = floor(diff * levels) / levels;
    toon_diff = max(toon_diff, light.ambient_strength);

    // Specular: discrete highlight (Blinn-Phong half-vector approach)
    vec3 half_dir = normalize(light_dir + view_dir);
    float spec = pow(max(dot(normal, half_dir), 0.0), material.sheen_coef);
    float toon_spec = floor(spec * levels) / levels;

    // Silhouette edge: darken near grazing angles
    float edge = dot(normal, view_dir);
    float edge_factor = smoothstep(0.15, 0.05, edge);

    vec3 base_color = texture(diffuse_map, uv_coords_in).rgb;
    vec3 ambient = light.ambient * base_color * light.ambient_strength * attenuation;
    vec3 diffuse = light.diffuse * base_color * toon_diff * attenuation;
    vec3 specular = light.specular * material.specular * toon_spec * attenuation;

    vert_color = (ambient + diffuse + specular) * edge_factor;
    uv_coords = vec2(uv_coords_in.x, 1.0 - uv_coords_in.y);

    gl_Position = transform.projection * transform.view * world_pos;
}