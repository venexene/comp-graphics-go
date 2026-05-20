// shaders/lighting/phong_gouraud.vert — Вершинный шейдер модели Фонга (Гуро).
//
// Модель освещения: Phong (Блинн-Фонг с отражённым вектором).
// Режим шейдинга: Gouraud — освещение вычисляется на вершинах,
// результат интерполируется по грани.
// Формула: I = I_a + I_d * max(N·L, 0) + I_s * (max(V·R, 0))^sh
// Где:
//   I_a — фоновое освещение (ambient)
//   I_d — диффузное освещение (Lambert)
//   I_s — зеркальное освещение (Phong reflection)
//   N — нормаль поверхности
//   L — направление на источник света
//   V — направление на камеру
//   R — отражённый вектор: reflect(-L, N)
//   sh — экспонента блеска (material.sheen_coef)

#version 330 core

// Входные атрибуты вершины.
layout (location = 0) in vec3 position;    // позиция в object space
layout (location = 1) in vec3 normal_in;   // нормаль в object space
layout (location = 3) in vec2 uv_coords_in; // текстурные координаты

// Uniform-структура трансформаций.
uniform struct Transform {
	mat4 model;       // модельная матрица (Object → World)
	mat4 view;        // видовая матрица (World → View)
	mat4 projection;  // матрица проекции (View → Clip)
	mat3 normal_mat;  // матрица нормалей ((M⁻¹)ᵀ)
	vec3 view_pos;    // позиция камеры в world space
} transform;

// Uniform-структура материала поверхности.
uniform struct Material {
	vec3 ambient;     // коэффициент фонового отражения k_a
	vec3 diffuse;     // коэффициент диффузного отражения k_d
	vec3 specular;    // коэффициент зеркального отражения k_s
    float sheen_coef; // экспонента блеска (shininess)
} material;

// Uniform-структура точечного источника света.
uniform struct PointLight {
	vec3 ambient;    // цвет фоновой составляющей света
	vec3 diffuse;    // цвет диффузной составляющей
	vec3 specular;   // цвет зеркальной составляющей
    vec3 position;   // позиция источника в world space

    float constant;  // постоянный коэффициент затухания c
    float linear;    // линейный коэффициент k_l
    float quadratic; // квадратичный коэффициент k_q

	float ambient_strength; // множитель фонового освещения [0,1]
} light;

uniform sampler2D diffuse_map;   // диффузная текстура
uniform float linear_coef;       // множитель линейного затухания
uniform float quadratic_coef;    // множитель квадратичного затухания
uniform int attenuation_mode;    // режим затухания: 0=Both,1=Linear,2=Quadratic

// Выходные данные для фрагментного шейдера (интерполируются).
out vec3 vert_color; // вычисленный цвет вершины
out vec2 uv_coords;  // текстурные координаты (Y-инвертированные)

void main() {	
	// 1. Преобразование позиции вершины в мировое пространство.
	vec4 world_pos = transform.model * vec4(position, 1.0);

	// 2. Нормаль в world space через матрицу нормалей (сохраняет длину).
	vec3 normal = normalize(transform.normal_mat * normal_in);
	
	// 3. Вектор от точки поверхности к источнику света.
	vec3 light_dir = light.position - world_pos.xyz;
	
	// 4. Вектор от точки поверхности к камере.
	vec3 view_dir = transform.view_pos - world_pos.xyz;
	
	// 5. Расстояние до источника для затухания.
	float distance = length(light_dir);
	
	// 6. Нормализация направлений.
	light_dir = normalize(light_dir);
    view_dir = normalize(view_dir);
	
	// 7. Отражённый вектор: R = reflect(-L, N).
	vec3 refl_dir = reflect(-light_dir, normal);

	// 8. Вычисление затухания в зависимости от режима.
	//    Формула: attenuation = 1 / (c + k_l * d + k_q * d²)
	float attenuation;
	if (attenuation_mode == 1) {
		// Только линейное.
		attenuation = 1.0 / max(light.constant + (light.linear * linear_coef) * distance, 0.0001);
	} else if (attenuation_mode == 2) {
		// Только квадратичное.
		attenuation = 1.0 / max(light.constant + (light.quadratic * quadratic_coef) * distance * distance, 0.0001);
	} else {
		// Оба (по умолчанию).
		attenuation = 1.0 / max(light.constant + 
			(light.linear * linear_coef) * distance + 
			(light.quadratic * quadratic_coef) * distance * distance, 0.0001);
	}

	// 9. Диффузное затенение: max(N·L, 0).
	float norm_d_light = max(dot(normal, light_dir), 0.0);
	
	// 10. Зеркальное затенение: max(V·R, 0)^sh.
    float view_d_refl = max(dot(view_dir, refl_dir), 0.0);

	// 11. Чтение цвета из текстуры.
    vec3 base_color = texture(diffuse_map, uv_coords_in).rgb;
    
	// 12. Три составляющие освещения по Фонгу.
    vec3 ambient = light.ambient * base_color * material.ambient * light.ambient_strength * attenuation; 
    vec3 diffuse = light.diffuse * base_color * material.diffuse * norm_d_light * attenuation;
    vec3 specular = light.specular * (pow(view_d_refl, material.sheen_coef) * material.specular) * attenuation;
	
	// 13. Итоговый цвет вершины (сумма).
	vert_color = ambient + diffuse + specular;
	
	// 14. Инвертирование Y текстурных координат (OpenGL vs OBJ).
	uv_coords = vec2(uv_coords_in.x, 1.0 - uv_coords_in.y);

	// 15. Финальная позиция в Clip Space.
	gl_Position = transform.projection * transform.view * world_pos;
}