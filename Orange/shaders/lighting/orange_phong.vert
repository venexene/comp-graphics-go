// Файл: shaders/lighting/orange_phong.vert
// Назначение: вершинный шейдер для bump mapping (normal mapping).
// Роль в пайплайне: преобразует вершины в clip space, строит TBN-матрицу,
// переводит направления света и обзора в tangent space для фрагментного шейдера.
//
// Работает в паре с:
//   - orange_phong.frag (модель Фонга)
//   - orange_blinn_phong.frag (модель Блинна-Фонга)
//
// Атрибуты (см. objects/model.go, CreateModelFromVertices):
//   location 0 — position   (vec3) — координаты вершины в объектном пространстве
//   location 1 — normal_in  (vec3) — нормаль вершины
//   location 2 — uv_coords_in (vec2) — текстурные координаты
//   location 3 — tangent_in (vec3) — касательная (вычислена на CPU)
//
// Uniform-структуры:
//   transform — матрицы Model, View, Projection, normal_mat, позиция камеры
//   light     — параметры источника света (позиция в world space)

#version 330 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal_in;
layout (location = 2) in vec2 uv_coords_in;
layout (location = 3) in vec3 tangent_in;

uniform struct Transform {
	mat4 model;         // Object space → World space
	mat4 view;          // World space → View space
	mat4 projection;    // View space → Clip space
	mat3 normal_mat;    // transpose(inverse(model)) — для нормалей
	vec3 view_pos;      // Позиция камеры в world space
} transform;

uniform struct PointLight {
	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
	vec3 position;      // Позиция источника в world space

	float constant;
	float linear;
	float quadratic;

	float ambient_strength;
} light;

// Выходные данные для фрагментного шейдера.
// Все направления передаются в tangent space, чтобы фрагментный шейдер
// мог использовать нормаль из карты нормалей (которая тоже в tangent space)
// без дополнительных преобразований.
out TangentSpace {
	vec3 light_dir;   // Направление на источник в tangent space
	vec3 view_dir;    // Направление на камеру в tangent space
	vec2 uv_coords;   // Текстурные координаты (v перевёрнута: 1.0 - y)
	float distance;   // Расстояние от вершины до источника (для затухания)
} vert;

void main() {
	// 1. Преобразование вершины в мировое пространство
	vec4 world_pos = transform.model * vec4(position, 1.0);

	// 2. Построение TBN-матрицы (Tangent-Bitangent-Normal)
	//
	// Нормаль и касательная преобразуются в world space через normal_mat.
	// normal_mat = transpose(inverse(model)) — обеспечивает корректное
	// преобразование нормалей даже при неравномерном масштабе.
	// Бикасательная B вычисляется как cross(N, T) — не требуется
	// передавать как атрибут, т.к. для ортонормированного базиса B = N × T.
	vec3 N = normalize(transform.normal_mat * normal_in);
	vec3 T = normalize(transform.normal_mat * tangent_in);

	// 3. Ортогонализация Грама — Шмидта:
	//    T' = normalize(T - N * dot(N, T))
	//    Обеспечивает, что T ортогонален N (касательная лежит в касательной
	//    плоскости, перпендикулярной нормали).
	T = normalize(T - dot(T, N) * N);
	vec3 B = cross(N, T);

	// 4. TBN-матрица переводит векторы из tangent space в world space.
	//    Её транспонирование (равное обратной для ортонормированной матрицы)
	//    переводит векторы из world space в tangent space.
	//
	//    Зачем переводить в tangent space, а не наоборот?
	//    Потому что карта нормалей хранит нормали в tangent space.
	//    Если перевести саму нормаль в world space (через TBN * normal),
	//    то при интерполяции между вершинами возникнут артефакты.
	//    Вместо этого мы переводим light_dir и view_dir в tangent space
	//    на этапе вершинного шейдера, а во фрагментном работаем
	//    полностью в tangent space.
	mat3 TBN = transpose(mat3(T, B, N));

	// 5. Направления на источник и камеру в world space
	vec3 world_light_dir = light.position - world_pos.xyz;
	vec3 world_view_dir  = transform.view_pos - world_pos.xyz;

	// 6. Перевод направлений в tangent space через TBN
	vert.light_dir = TBN * world_light_dir;
	vert.view_dir  = TBN * world_view_dir;

	// Переворачиваем V-координату текстуры (OpenGL vs. изображение)
	vert.uv_coords = vec2(uv_coords_in.x, 1.0 - uv_coords_in.y);

	// Расстояние до источника (для затухания)
	vert.distance = length(world_light_dir);

	// 7. Стандартное преобразование вершины в clip space
	gl_Position = transform.projection * transform.view * world_pos;
}
