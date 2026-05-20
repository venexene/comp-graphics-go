// Файл: shaders/lighting/orange_phong.frag
// Назначение: фрагментный шейдер для модели освещения Фонга (Phong)
// с bump mapping (normal mapping) в tangent space.
//
// Роль в пайплайне: получает интерполированные направления света и обзора
// в tangent space от вершинного шейдера, сэмплирует карту нормалей,
// вычисляет Phong-освещение с возмущённой нормалью.
//
// Работает в паре с: orange_phong.vert.
//
// Модель освещения: Phong.
//   Ambient:  I_a = light.ambient * surfaceColor * material.ambient * strength * attenuation * AO
//   Diffuse:  I_d = light.diffuse * surfaceColor * material.diffuse * max(N·L, 0) * attenuation
//   Specular: I_s = light.specular * material.specular * max(R·V, 0)^shininess * attenuation
//              где R = reflect(-L, N)
//   Итог:    I = I_a + I_d + I_s
//
// Все вычисления производятся в tangent space. Нормаль берётся из карты
// нормалей (tangent space), L и V переданы из вершинного шейдера уже
// в tangent space.

#version 330 core

// Интерполированные входные данные из вершинного шейдера
in TangentSpace {
	vec3 light_dir;   // Направление на источник в tangent space
	vec3 view_dir;    // Направление на камеру в tangent space
	vec2 uv_coords;   // Текстурные координаты
	float distance;   // Расстояние до источника
} vert;

// Текстурные карты
uniform sampler2D u_diffuseMap;   // GL_TEXTURE0 — диффузная карта (цвет апельсина)
uniform sampler2D u_normalMap;    // GL_TEXTURE1 — карта нормалей (tangent space)
uniform sampler2D u_aoMap;        // GL_TEXTURE2 — карта ambient occlusion (опционально)
uniform sampler2D u_roughnessMap; // GL_TEXTURE3 — карта шероховатости (опционально)

// Коэффициенты затухания (регулируются через GUI)
uniform float linear_coef;        // Множитель линейного затухания [0, ∞)
uniform float quadratic_coef;     // Множитель квадратичного затухания [0, ∞)
uniform int attenuation_mode;     // Режим: 0=оба, 1=линейное, 2=квадратичное

// Параметры материала
uniform struct Material {
	vec3 ambient;       // Цвет ambient-отражения материала
	vec3 diffuse;       // Цвет diffuse-отражения материала
	vec3 specular;      // Цвет specular-отражения материала
	float sheen_coef;   // Коэффициент блеска (shininess) [1, 256]
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

	float ambient_strength; // Сила ambient-освещения [0, 1]
} light;

// Выходной цвет фрагмента
out vec4 frag_color;

void main() {
	// ===== 1. ВЫБОРКА И ПРЕОБРАЗОВАНИЕ НОРМАЛИ ИЗ КАРТЫ НОРМАЛЕЙ =====
	//
	// Карта нормалей хранит векторы в формате RGB, где каждый компонент
	// находится в диапазоне [0, 1]. Направление нормали в tangent space
	// имеет компоненты в [-1, 1]. Преобразование:
	//
	//   normal = texture.rgb * 2.0 - 1.0
	//
	// В карте нормалей:
	//   R = tangent axis (X)  — отклонение вдоль касательной
	//   G = bitangent axis (Y) — отклонение вдоль бикасательной
	//   B = normal axis (Z)    — отклонение вдоль нормали (обычно близко к 1)
	//
	// Для OpenGL-нормалей (NormalGL): ось Y направлена вверх (стандарт OpenGL).
	vec3 normal = texture(u_normalMap, vert.uv_coords).rgb;
	normal = normalize(normal * 2.0 - 1.0);

	// ===== 2. НОРМАЛИЗАЦИЯ НАПРАВЛЕНИЙ В TANGENT SPACE =====
	// Интерполяция между вершинами изменяет длину векторов,
	// поэтому необходима повторная нормализация.
	vec3 light_dir = normalize(vert.light_dir);
	vec3 view_dir  = normalize(vert.view_dir);

	// ===== 3. ЗАТУХАНИЕ СВЕТА (ATTENUATION) =====
	//
	// Формула затухания точечного источника:
	//
	//   att = 1 / (Kc + Kl * d + Kq * d²)
	//
	// где d — расстояние от источника до точки поверхности.
	// В зависимости от attenuation_mode:
	//   0 — полная формула (линейное + квадратичное)
	//   1 — только линейное (Kq = 0)
	//   2 — только квадратичное (Kl = 0)
	//
	// linear_coef и quadratic_coef — множители, регулируемые пользователем.
	// max(..., 0.0001) предотвращает деление на ноль.
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

	// ===== 4. ДИФФУЗНЫЙ ЦВЕТ ПОВЕРХНОСТИ =====
	vec3 surface_color = texture(u_diffuseMap, vert.uv_coords).rgb;

	// ===== 5. AMBIENT OCCLUSION =====
	// AO-карта (красный канал) модулирует ambient-составляющую:
	// в углублениях (AO ~ 0) ambient ослабляется.
	float ao = texture(u_aoMap, vert.uv_coords).r;

	// ===== 6. ОСВЕЩЕНИЕ ФОНГА (PHONG) =====
	//
	// N = возмущённая нормаль из карты нормалей (в tangent space)
	// L = направление на источник (в tangent space)
	// V = направление на камеру (в tangent space)
	//
	// Diffuse (закон Ламберта):
	//   I_diffuse = light.diffuse * surfaceColor * material.diffuse * max(N·L, 0) * att
	//
	float NdotL = max(dot(normal, light_dir), 0.0);

	// Specular (Фонг):
	//   R = reflect(-L, N) — отражённый вектор
	//   I_specular = light.specular * material.specular * max(R·V, 0)^shininess * att
	//
	vec3 refl_dir = reflect(-light_dir, normal);
	float spec = pow(max(dot(view_dir, refl_dir), 0.0), material.sheen_coef);

	// ===== 7. СУММИРОВАНИЕ КОМПОНЕНТ =====
	//
	// Ambient: фоновое освещение, не зависит от направления.
	//   умножается на ambient_strength (регулировка), attenuation и AO.
	//
	// Diffuse: рассеянный свет, зависит от угла падения.
	//
	// Specular: блик, зависит от угла между отражённым лучом и взглядом.
	//
	vec3 ambient  = light.ambient * surface_color * material.ambient * light.ambient_strength * attenuation * ao;
	vec3 diffuse  = light.diffuse * surface_color * material.diffuse * NdotL * attenuation;
	vec3 specular = light.specular * spec * material.specular * attenuation;

	frag_color = vec4(ambient + diffuse + specular, 1.0);
}
