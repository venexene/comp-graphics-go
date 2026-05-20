// shaders/lighting/basic_phong.vert — Вершинный шейдер для Phong-шейдинга.
//
// Назначение: передача интерполируемых данных (нормаль, направления света
// и обзора, UV, расстояние) во фрагментный шейдер.
// Сам расчёт освещения выполняется во фрагментном шейдере (Phong shading).

#version 330 core

layout (location = 0) in vec3 position;    // позиция в object space
layout (location = 1) in vec3 normal_in;   // нормаль в object space
layout (location = 3) in vec2 uv_coords_in; // UV-координаты

uniform struct Transform {
	mat4 model;
	mat4 view;
	mat4 projection;
	mat3 normal_mat;
	vec3 view_pos;
} transform;

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

// Интерфейсный блок Vertex — данные, интерполируемые между вершинами и
// фрагментами. Фрагментный шейдер получит интерполированные значения.
out Vertex {
	vec3 normal;      // нормаль в world space (интерполируется)
	vec3 light_dir;   // направление на источник (интерп.)
	vec3 view_dir;    // направление на камеру (интерп.)
	vec2 uv_coords;   // текстурные координаты (интерп.)
	float distance;   // расстояние до источника (интерп.)
} vert;

void main() {	
	// Позиция вершины в world space.
	vec4 world_pos = transform.model * vec4(position, 1.0);

	// Нормаль (без нормализации — нормализуется во фрагментном шейдере).
	vert.normal = transform.normal_mat * normal_in;
	
	// Направления на свет и камеру.
	vert.light_dir = light.position - world_pos.xyz;
	vert.view_dir = transform.view_pos - world_pos.xyz;
	
	// UV с инвертированным Y.
	vert.uv_coords = vec2(uv_coords_in.x, 1.0 - uv_coords_in.y);
	
	// Расстояние для затухания.
	vert.distance = length(vert.light_dir);
	
	// Финальная позиция.
	gl_Position = transform.projection * transform.view * world_pos;
}