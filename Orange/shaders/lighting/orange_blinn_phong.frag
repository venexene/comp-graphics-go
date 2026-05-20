// Файл: shaders/lighting/orange_blinn_phong.frag
// Назначение: фрагментный шейдер для модели освещения Блинна-Фонга (Blinn-Phong)
// с bump mapping (normal mapping) в tangent space.
//
// Роль в пайплайне: получает интерполированные направления в tangent space,
// сэмплирует карту нормалей, вычисляет Blinn-Phong-освещение.
//
// Работает в паре с: orange_phong.vert.
//
// Модель освещения: Blinn-Phong.
// Отличие от Phong: specular вычисляется через half-vector H = normalize(L + V)
// вместо отражённого вектора R. Это даёт более мягкие и реалистичные блики
// и slightly different (обычно более физически корректное) распределение.
//
//   Ambient:  I_a = light.ambient * surfaceColor * material.ambient * strength * att * AO
//   Diffuse:  I_d = light.diffuse * surfaceColor * material.diffuse * max(N·L, 0) * att
//   Specular: I_s = light.specular * material.specular * max(N·H, 0)^shininess * att
//              где H = normalize(L + V) — half-vector между L и V
//   Итог:    I = I_a + I_d + I_s

#version 330 core

// Интерполированные входные данные из вершинного шейдера (в tangent space)
in TangentSpace {
	vec3 light_dir;
	vec3 view_dir;
	vec2 uv_coords;
	float distance;
} vert;

// Текстурные карты
uniform sampler2D u_diffuseMap;   // GL_TEXTURE0 — диффузная карта
uniform sampler2D u_normalMap;    // GL_TEXTURE1 — карта нормалей (tangent space)
uniform sampler2D u_aoMap;        // GL_TEXTURE2 — карта ambient occlusion
uniform sampler2D u_roughnessMap; // GL_TEXTURE3 — карта шероховатости

// Коэффициенты затухания
uniform float linear_coef;
uniform float quadratic_coef;
uniform int attenuation_mode;     // 0=оба, 1=линейное, 2=квадратичное

// Параметры материала
uniform struct Material {
	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
	float sheen_coef;   // Коэффициент блеска [1, 256]
} material;

// Параметры источника света
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
	// ===== 1. ВЫБОРКА НОРМАЛИ ИЗ КАРТЫ НОРМАЛЕЙ =====
	// Преобразование из [0, 1] (RGB текстуры) в [-1, 1] (вектор нормали)
	vec3 normal = texture(u_normalMap, vert.uv_coords).rgb;
	normal = normalize(normal * 2.0 - 1.0);

	// ===== 2. НОРМАЛИЗАЦИЯ НАПРАВЛЕНИЙ =====
	vec3 light_dir = normalize(vert.light_dir);
	vec3 view_dir  = normalize(vert.view_dir);

	// ===== 3. ЗАТУХАНИЕ =====
	float attenuation;
	if (attenuation_mode == 1) {
		// Только линейное затухание: att = 1 / (Kc + Kl * d)
		attenuation = 1.0 / max(light.constant + (light.linear * linear_coef) * vert.distance, 0.0001);
	} else if (attenuation_mode == 2) {
		// Только квадратичное затухание: att = 1 / (Kc + Kq * d²)
		attenuation = 1.0 / max(light.constant + (light.quadratic * quadratic_coef) * vert.distance * vert.distance, 0.0001);
	} else {
		// Полное затухание: att = 1 / (Kc + Kl * d + Kq * d²)
		attenuation = 1.0 / max(light.constant +
			(light.linear * linear_coef) * vert.distance +
			(light.quadratic * quadratic_coef) * vert.distance * vert.distance, 0.0001);
	}

	// ===== 4. ЦВЕТ ПОВЕРХНОСТИ ИЗ ДИФФУЗНОЙ КАРТЫ =====
	vec3 surface_color = texture(u_diffuseMap, vert.uv_coords).rgb;

	// ===== 5. AMBIENT OCCLUSION =====
	float ao = texture(u_aoMap, vert.uv_coords).r;

	// ===== 6. ОСВЕЩЕНИЕ БЛИННА-ФОНГА =====
	//
	// N — возмущённая нормаль из карты нормалей (tangent space).
	// L — направление на источник, V — направление на камеру.
	//
	// Diffuse (закон Ламберта):
	float NdotL = max(dot(normal, light_dir), 0.0);

	// Specular (Блинн-Фонг):
	// Вместо отражённого вектора R используется half-vector H:
	//   H = normalize(L + V)
	// Это вектор, направленный ровно посередине между L и V.
	// spec = max(N·H, 0)^shininess
	//
	// Преимущество Blinn-Phong перед Phong:
	// - Более гладкое распределение блика
	// - Не требует вычисления reflect (меньше инструкций)
	// - При высоком shininess даёт более реалистичный результат
	vec3 half_dir = normalize(light_dir + view_dir);
	float spec = pow(max(dot(normal, half_dir), 0.0), material.sheen_coef);

	// ===== 7. ИТОГОВЫЙ ЦВЕТ =====
	vec3 ambient  = light.ambient * surface_color * material.ambient * light.ambient_strength * attenuation * ao;
	vec3 diffuse  = light.diffuse * surface_color * material.diffuse * NdotL * attenuation;
	vec3 specular = light.specular * spec * material.specular * attenuation;

	frag_color = vec4(ambient + diffuse + specular, 1.0);
}
