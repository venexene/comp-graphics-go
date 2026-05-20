# ARCHITECTURE.md — Bump Mapping на апельсине

## 1. Обзор проекта

**Полное название:** Лабораторная работа «Bump Mapping (Normal Mapping) на модели апельсина с использованием tangent space и TBN-матрицы».

**Сцена:** программа рендерит 3D-сцену, содержащую модель апельсина (`orange.obj`) и модель сердца (`heart.obj`), с применением **bump mapping** (normal mapping) для реалистичного отображения неровностей кожуры апельсина и поверхности ёлочной игрушки (сердце). Модели освещены точечным источником света.

**Демонстрируемые технологии:**
- **Bump mapping (normal mapping)** — имитация мелких неровностей поверхности без увеличения геометрической сложности.
- **Tangent space** — локальная система координат каждой вершины, в которой хранится карта нормалей.
- **TBN-матрица** — преобразование между tangent space и world space.
- **Модель освещения Фонга (Phong)** — классическая модель с reflect-вектором.
- **Модель освещения Блинна-Фонга (Blinn-Phong)** — улучшенная модель с half-vector.
- **Затухание точечного источника** — линейное, квадратичное или оба.

**Что видит пользователь при запуске:**
- Окно 1280×720 с тёмно-серым фоном.
- Два 3D-объекта: апельсин слева и сердце справа, оба с bump mapping.
- В консоли выводится информационная панель с параметрами сцены.

**Управление:**
- Вращение камеры, панорамирование, зум.
- Перемещение объекта и источника света.
- Переключение модели освещения (Фонг / Блинн-Фонг).
- Регулировка затухания и ambient-освещения.

---

## 2. Дерево проекта

```
Orange/
├── ARCHITECTURE.md              # Настоящий файл — архитектурная документация
├── go.mod                       # Go-модуль (github.com/venexene/comp-graphics-go)
├── cmd/
│   └── main.go                  # Точка входа: инициализация, цикл рендера, cleanup
├── input/
│   └── input.go                 # Обработка клавиатурного ввода
├── lighting/
│   ├── light.go                 # Структура LightConfig (точечный источник)
│   ├── material.go              # Структура MaterialConfig (материал поверхности)
│   └── uniforms.go              # Кеш uniform-переменных (UniformCache)
├── models/
│   ├── orange.obj               # 3D-модель апельсина
│   └── heart.obj                # 3D-модель сердца
├── objects/
│   └── model.go                 # Загрузка OBJ, вычисление касательных, VAO/VBO
├── scene/
│   ├── camera.go                # Орбитальная камера (target-distance-yaw-pitch)
│   ├── object.go                # Состояние объектов, матрицы Model, Selection
│   └── scene.go                 # Сцены OrangeScene/HeartScene, привязка текстур, отрисовка
├── shaders/
│   ├── shaders.go               # Загрузка и компиляция GLSL-шейдеров
│   ├── lighting.go              # Управление вариантами освещения (LightingVariant)
│   └── lighting/
│       ├── orange_phong.vert    # Вершинный шейдер (TBN → tangent space)
│       ├── orange_phong.frag    # Фрагментный шейдер (Phong + normal map)
│       └── orange_blinn_phong.frag  # Фрагментный шейдер (Blinn-Phong + normal map)
├── textures/
│   ├── texture.go               # Загрузка текстур (JPEG/PNG → OpenGL texture)
│   ├── orange/                  # Текстуры апельсина (color, normal, ao, roughness)
│   └── ornament2/               # Текстуры сердца (color, normal, metalness, roughness)
├── ui/
│   └── ui.go                    # UI-оверлей (текстовая информационная панель)
└── utils/
    └── utils.go                 # Глобальное состояние, инициализация, цикл рендера
```

---

## 3. Графический пайплайн

### Инициализация (один раз при запуске)

```
main.go
  ├── findProjectRoot() — поиск корня проекта (go.mod)
  ├── shaders.SetBasePath() — установка пути к шейдерам
  ├── initGlfw()
  │     ├── glfw.Init()
  │     ├── glfw.WindowHint(...) — OpenGL 4.1 Core Profile
  │     └── glfw.CreateWindow(1280×720)
  ├── initOpenGL()
  │     ├── gl.Init()
  │     └── shaders.InitLightingVariants()
  │           ├── Загрузка orange_phong.vert + orange_phong.frag → программа 1
  │           └── Загрузка orange_phong.vert + orange_blinn_phong.frag → программа 2
  ├── ui.InitializeUI()
  ├── utils.InitScene()
  │     ├── scene.NewOrangeScene() — загрузка orange.obj + текстуры апельсина
  │     ├── scene.NewHeartScene() — загрузка heart.obj + текстуры сердца
  │     └── initObjectStates() — начальные позиции/масштаб
  ├── gl.Enable(DEPTH_TEST)
  ├── mgl32.Perspective(45°, 1280/720, 0.1, 100.0) — матрица проекции
  └── Главный цикл рендера
```

### Цикл рендера (каждый кадр)

```
for !window.ShouldClose():
  ├── ui.BeginFrame()
  ├── utils.DrawScene(window, projection):
  │     ├── input.ProcessInput() — клавиатура
  │     ├── shaders.GetCurrentLightingProgram() — активная программа
  │     ├── scene.DrawSceneObjects():
  │     │     ├── gl.ClearColor(0.15, 0.15, 0.18) + gl.Clear(COLOR | DEPTH)
  │     │     ├── gl.UseProgram(program)
  │     │     ├── UniformCache.Refresh(program) — location-ы
  │     │     ├── cam.ViewMatrix() + cam.EyePosition()
  │     │     ├── gl.Uniform3f(view_pos)
  │     │     ├── Для апельсина:
  │     │     │     ├── orangeState.ModelMatrix()
  │     │     │     ├── setTransformUniforms(model, view, proj, normal_mat)
  │     │     │     ├── setMaterialUniforms(matCfg)
  │     │     │     ├── setLightUniforms(lightCfg)
  │     │     │     ├── bindOrangeTextures() — TEXTURE0..3
  │     │     │     └── orange.Model.Draw()
  │     │     └── Для сердца:
  │     │           ├── heartState.ModelMatrix()
  │     │           ├── setTransformUniforms / setMaterial / setLight
  │     │           ├── bindHeartTextures() — TEXTURE0..3
  │     │           └── heart.Model.Draw()
  │     ├── glfw.PollEvents()
  │     └── window.SwapBuffers()
  ├── Обновление заголовка окна
  └── ui.EndFrame()
```

### Завершение

```
defer (при выходе из main):
  ├── shaders.CleanupLightingVariants() — gl.DeleteProgram
  ├── ui.Cleanup()
  ├── utils.Cleanup()
  │     ├── orangeScene.Cleanup() — удаление модели и текстур
  │     └── heartScene.Cleanup() — удаление модели и текстур
  └── glfw.Terminate()
```

---

## 4. Геометрия модели

### Формат OBJ

Парсер в `objects/model.go` обрабатывает следующие типы строк:

| Тип строки | Описание | Пример |
|-----------|----------|--------|
| `v` | Вершина (позиция) | `v 0.123 0.456 0.789` |
| `vt` | Текстурная координата | `vt 0.500 0.750` |
| `vn` | Нормаль | `vn 0.0 1.0 0.0` |
| `f` | Грань (треугольник) | `f 1/1/1 2/2/2 3/3/3` |

### Формат граней

Грани парсятся в формате `v/vt/vn` (позиция/текстурная координата/нормаль), где все три индекса могут присутствовать. Поддерживаются отрицательные индексы (отсчёт с конца списка). Триангуляция выполняется методом «fan»: первая вершина + каждая последующая пара образуют треугольник.

### Interleaved массив вершин

После парсинга вершины укладываются в VBO как interleaved массив со следующим расположением атрибутов (stride = 44 байта = 11 × float32):

| Location | Атрибут | Тип | Размер (байт) | Смещение |
|----------|---------|-----|---------------|----------|
| 0 | Position | vec3 | 12 | 0 |
| 1 | Normal | vec3 | 12 | 12 |
| 2 | UV | vec2 | 8 | 24 |
| 3 | Tangent | vec3 | 12 | 32 |

**Бикасательная не передаётся как атрибут.** Она вычисляется в вершинном шейдере как `B = cross(N, T)`. Это экономит 12 байт на вершину и упрощает формат данных. Для ортонормированного базиса после ортогонализации Грама — Шмидта бикасательная однозначно определяется нормалью и касательной.

### Размер модели

Модель `orange.obj` содержит несколько тысяч треугольников (типичная модель апельсина среднего разрешения). Точное количество вер/треугольников можно получить из лога парсинга.

---

## 5. Bump Mapping (Normal Mapping)

### Что такое bump mapping

Bump mapping — техника имитации мелких неровностей поверхности (бугорков, вмятин, пор) без изменения геометрии модели. Вместо добавления тысяч дополнительных треугольников используется текстура (карта нормалей), которая на этапе освещения **возмущает** нормаль поверхности.

### Normal mapping

Normal mapping — разновидность bump mapping, где в текстуре хранятся не высоты (как в height map), а непосредственно векторы нормалей, закодированные в RGB-каналы.

### Формат карты нормалей

Карта нормалей хранит отклонения нормали в **tangent space**:
- R-канал (красный) — отклонение вдоль оси T (касательная), диапазон [0, 1] соответствует [-1, 1].
- G-канал (зелёный) — отклонение вдоль оси B (бикасательная), диапазон [0, 1] соответствует [-1, 1].
- B-канал (синий) — отклонение вдоль оси N (нормаль), обычно близок к 1 (т.е. 255 в текстуре).

Преобразование во фрагментном шейдере:

```glsl
vec3 normal = texture(u_normalMap, uv).rgb;      // [0, 1]
normal = normalize(normal * 2.0 - 1.0);           // → [-1, 1]
```

Карта нормалей **не должна загружаться как sRGB**, т.к. данные нормали — это векторы, а не цвета. sRGB-преобразование исказило бы направления.

Для данной работы используется карта нормалей в формате OpenGL (`food_0022_normal_opengl_1k.png`), где ось Y направлена вверх.

### Tangent space

Tangent space — локальная декартова система координат каждой вершины:
- **Ось T (касательная)** — направление увеличения U-текстурной координаты.
- **Ось B (бикасательная)** — направление увеличения V-текстурной координаты.
- **Ось N (нормаль)** — перпендикуляр к поверхности вершины.

Tangent space необходим потому, что карта нормалей хранит отклонения относительно локальной ориентации поверхности. Без tangent space одна и та же карта нормалей не могла бы использоваться на разных гранях модели.

### Вычисление касательных на CPU (objects/model.go — ComputeTangents)

Для каждого треугольника:

```
p0, p1, p2 — позиции вершин
uv0, uv1, uv2 — текстурные координаты

deltaPos1 = p1 - p0
deltaPos2 = p2 - p0
duv1 = uv1 - uv0 = (du1, dv1)
duv2 = uv2 - uv0 = (du2, dv2)

Решаем систему:
  deltaPos1 = T * du1 + B * dv1
  deltaPos2 = T * du2 + B * dv2

Отсюда (только для T):
  f = 1 / (du1 * dv2 - dv1 * du2)
  T = f * (deltaPos1 * dv2 - deltaPos2 * dv1)
```

Касательная накапливается для каждой вершины из всех треугольников, в которых она участвует. После накопления применяется ортогонализация Грама — Шмидта относительно нормали:

```
T' = normalize(T - N * dot(N, T))
```

### Построение TBN-матрицы (вершинный шейдер — orange_phong.vert)

```glsl
// Преобразование normal_mat — transpose(inverse(model))
// для корректного преобразования нормалей в world space
vec3 N = normalize(transform.normal_mat * normal_in);
vec3 T = normalize(transform.normal_mat * tangent_in);

// Ортогонализация Грама — Шмидта
T = normalize(T - dot(T, N) * N);
vec3 B = cross(N, T);

// TBN: tangent space → world space
// transpose(TBN): world space → tangent space
mat3 TBN = transpose(mat3(T, B, N));
```

### Преобразование нормали (фрагментный шейдер)

Во фрагментном шейдере направления L и V уже находятся в tangent space (переданы из вершинного шейдера). Нормаль из карты нормалей также в tangent space. Поэтому освещение вычисляется полностью в tangent space без дополнительных преобразований.

Важно: в отличие от подхода «перевести нормаль в world space», здесь мы переводим **направления** (L и V) в tangent space. Это делается в вершинном шейдере для производительности — вместо преобразования нормали для каждого фрагмента (миллионы раз) преобразуются только два направления на вершину.

---

## 6. Модели освещения

### Модель Фонга (Phong) — orange_phong.frag

Формула освещения (все векторы в tangent space):

```
N — возмущённая нормаль из карты нормалей (normalized)
L — направление на источник (normalized)
V — направление на камеру (normalized)
R = reflect(-L, N) — отражённый вектор

NdotL = max(N·L, 0)

Diffuse  = light.diffuse * surfaceColor * material.diffuse * NdotL * attenuation

Specular = light.specular * material.specular * max(R·V, 0)^shininess * attenuation

Ambient  = light.ambient * surfaceColor * material.ambient * ambientStrength * attenuation * AO

Итог: fragmentColor = Ambient + Diffuse + Specular
```

### Модель Блинна-Фонга (Blinn-Phong) — orange_blinn_phong.frag

Отличается только specular-компонентой:

```
H = normalize(L + V) — half-vector между L и V
Specular = light.specular * material.specular * max(N·H, 0)^shininess * attenuation
```

Преимущества Blinn-Phong:
- Более гладкое распределение блика.
- Не требует вычисления `reflect()`.
- При высоких значениях shininess даёт более физически корректный результат.

### Переключение между моделями

Uniform-переменная `u_lightModel` не используется напрямую — переключение происходит через выбор шейдерной программы. В `shaders/lighting.go` определены два варианта:

| Индекс | Имя | Вершинный шейдер | Фрагментный шейдер |
|--------|-----|------------------|-------------------|
| 0 | Phong + Normal Map | orange_phong.vert | orange_phong.frag |
| 1 | Blinn-Phong + Normal Map | orange_phong.vert | orange_blinn_phong.frag |

Переключение: клавиши **T** (вперёд) и **G** (назад).

---

## 7. Система освещения

### Структура LightConfig (lighting/light.go)

| Поле | Тип | Описание | Uniform в шейдере |
|------|-----|----------|-------------------|
| Position | mgl32.Vec3 | Позиция источника в world space | `light.position` |
| Ambient | mgl32.Vec3 | Цвет фонового освещения (RGB) | `light.ambient` |
| Diffuse | mgl32.Vec3 | Цвет диффузного освещения (RGB) | `light.diffuse` |
| Specular | mgl32.Vec3 | Цвет зеркального освещения (RGB) | `light.specular` |
| Constant | float32 | Константа затухания (Kc) | `light.constant` |
| Linear | float32 | Коэффициент линейного затухания (Kl) | `light.linear` |
| Quadratic | float32 | Коэффициент квадратичного затухания (Kq) | `light.quadratic` |
| LinearCoef | float32 | Множитель линейного затухания | `linear_coef` |
| QuadraticCoef | float32 | Множитель квадратичного затухания | `quadratic_coef` |
| AmbientStrength | float32 | Сила ambient-освещения [0, 1] | `light.ambient_strength` |
| Mode | AttenuationMode | Режим затухания (0=оба, 1=лин, 2=квадр) | `attenuation_mode` |

### Тип источника: точечный (PointLight)

Свет распространяется во все стороны из одной точки. Затухание зависит от расстояния.

### Формула затухания

```
attenuation = 1.0 / (Kc + Kl * linearCoef * d + Kq * quadraticCoef * d²)
```

где:
- `d` — расстояние от источника до точки поверхности.
- `Kc = 1.0` — константа.
- `Kl = 0.09` — коэффициент линейного затухания.
- `Kq = 0.032` — коэффициент квадратичного затухания.
- `linearCoef` — множитель, регулируемый пользователем (Z/X).
- `quadraticCoef` — множитель, регулируемый пользователем (C/V).

В зависимости от `attenuation_mode`:
- **0 (Both):** полная формула.
- **1 (Linear):** только линейный член (`Kq = 0`).
- **2 (Quadratic):** только квадратичный член (`Kl = 0`).

### Положение источника

Начальная позиция: `(3.0, 5.0, 2.5)` — справа и сверху от апельсина. Может быть изменена клавишами Alt+IJKLUO.

---

## 8. GUI и управление

В текущей реализации **нет графического GUI-фреймворка** (ImGui не используется). UI представлен текстовым оверлеем, который выводится в консоль при запуске (функция `GetUIOverlayText` в `ui/ui.go`). Управление осуществляется исключительно с клавиатуры.

### Полный перечень элементов управления

| Действие | Тип | Клавиша | Переменная в Go | Uniform в шейдере | Диапазон |
|----------|-----|---------|-----------------|-------------------|----------|
| Вращение камеры (орбита) | Клавиши | Стрелки | cam.Yaw / cam.Pitch | — | Неогр. |
| Панорамирование (вперёд/назад) | Клавиши | W / S | cam.Target | — | Неогр. |
| Панорамирование (влево/вправо) | Клавиши | A / D | cam.Target | — | Неогр. |
| Панорамирование (вверх/вниз) | Клавиши | Space / Shift | cam.Target | — | Неогр. |
| Зум (приближение/отдаление) | Клавиши | + / - | cam.Distance | — | [1.0, 50.0] |
| Масштаб объекта (увел/умен) | Клавиши | E / Q | mainState.Scale | transform.model | [0.1, 3.0] |
| Вращение объекта по Z | Клавиши | R / F | mainState.RotationZ | transform.model | Неогр. |
| Вращение объекта по X | Клавиши | 1 / 2 | mainState.RotationX | transform.model | Неогр. |
| Вращение объекта по Y | Клавиши | 3 / 4 | mainState.RotationY | transform.model | Неогр. |
| Перемещение объекта (IJKL) | Клавиши | I/K/J/L/U/O | mainState.Position | transform.model | Неогр. |
| Перемещение источника света | Клавиши | Alt+IJKLUO | lightCfg.Position | light.position | Неогр. |
| Следующий вариант освещения | Клавиша | T | currentLightingIndex | — | [0, 1] |
| Предыдущий вариант освещения | Клавиша | G | currentLightingIndex | — | [0, 1] |
| Переключение Gouraud/Phong | Клавиша | Y | currentLightingIndex | — | 0/1 |
| Цикл затухания | Клавиша | M | lightCfg.Mode | attenuation_mode | 0/1/2 |
| Линейное затухание (умен/увел) | Клавиши | Z / X | lightCfg.LinearCoef | linear_coef | [0, ∞) |
| Квадратичное затухание (умен/увел) | Клавиши | C / V | lightCfg.QuadraticCoef | quadratic_coef | [0, ∞) |
| Ambient strength (умен/увел) | Клавиши | B / N | lightCfg.AmbientStrength | light.ambient_strength | [0, 1] |
| Сброс всех параметров | Комбинация | Ctrl+R | — | — | — |

---

## 9. Система координат и матрицы

### Матрицы

| Матрица | Назначение | Uniform | Описание |
|---------|-----------|---------|----------|
| Model | Object → World | `transform.model` | Включает масштабирование, вращение (Z×Y×X), перенос |
| View | World → View | `transform.view` | Построена через `mgl32.LookAt` |
| Projection | View → Clip | `transform.projection` | Перспективная, FOV=45°, near=0.1, far=100 |
| Normal | Нормали | `transform.normal_mat` | `mat3(transpose(inverse(model)))` |

### Библиотека матричных вычислений

Используется **github.com/go-gl/mathgl/mgl32** — библиотека математики для OpenGL на Go.

### Положение камеры

- **Позиция:** вычисляется из `cam.EyePosition()` — сферические координаты вокруг Target.
- **Начальный Target:** `(0.0, 0.05, 0.3)` — между апельсином и сердцем.
- **Начальное расстояние:** 5.0 единиц.
- **Up-вектор:** `(0, 1, 0)` — ось Y направлена вверх.

### Положение источника света

- **Начальная позиция:** `(3.0, 5.0, 2.5)` — справа и сверху от апельсина.

### Порядок умножения матриц в вершинном шейдере

```glsl
gl_Position = transform.projection * transform.view * transform.model * vec4(position, 1.0);
```

Порядок умножения **справа налево**: сначала Model, затем View, затем Projection.

### Построение матрицы Model (scene/object.go)

```go
Model = T * Rz * Ry * Rx * S
```

Порядок: сначала масштаб, потом вращения (X, затем Y, затем Z), потом перенос.

### Матрица нормалей

```glsl
mat3 normal_mat = transpose(inverse(model))
```

Используется для преобразования нормалей и касательных в world space. Обеспечивает корректное преобразование даже при неравномерном масштабе (хотя в данной работе масштаб равномерный).

---

## 10. Текстуры и текстурные блоки

### Текстуры апельсина (OrangeScene)

| Файл | Тип | Блок | Uniform sampler2D | Внутренний формат | Фильтрация | Wrap |
|------|-----|------|-------------------|-------------------|------------|------|
| `food_0022_color_1k.jpg` | Diffuse | GL_TEXTURE0 | u_diffuseMap | GL_RGBA | GL_LINEAR | GL_REPEAT |
| `food_0022_normal_opengl_1k.png` | Normal map | GL_TEXTURE1 | u_normalMap | GL_RGBA | GL_LINEAR | GL_REPEAT |
| `food_0022_ao_1k.jpg` | Ambient occlusion | GL_TEXTURE2 | u_aoMap | GL_RGBA | GL_LINEAR | GL_REPEAT |
| `food_0022_roughness_1k.jpg` | Roughness | GL_TEXTURE3 | u_roughnessMap | GL_RGBA | GL_LINEAR | GL_REPEAT |

### Текстуры сердца (HeartScene)

| Файл | Тип | Блок | Uniform sampler2D |
|------|-----|------|-------------------|
| `ChristmasTreeOrnament021_1K-JPG_Color.jpg` | Diffuse | GL_TEXTURE0 | u_diffuseMap |
| `ChristmasTreeOrnament021_1K-JPG_NormalGL.jpg` | Normal map | GL_TEXTURE1 | u_normalMap |
| `ChristmasTreeOrnament021_1K-JPG_Metalness.jpg` | Metalness | GL_TEXTURE2 | u_aoMap (замена) |
| `ChristmasTreeOrnament021_1K-JPG_Roughness.jpg` | Roughness | GL_TEXTURE3 | u_roughnessMap |

### Важно о карте нормалей

Карта нормалей **не должна загружаться как sRGB** (внутренний формат `GL_RGB` или `GL_RGBA`, **не** `GL_SRGB`). Если бы она загружалась как sRGB, значения цветов были бы гамма-скорректированы, что привело бы к неправильным направлениям нормалей.

### Mipmap-уровни

В текущей реализации mipmap-уровни **не генерируются** (фильтрация `GL_LINEAR` без `GL_LINEAR_MIPMAP_LINEAR`). Это может вызывать мерцание текстур на больших расстояниях. Улучшение: добавить `gl.GenerateMipmap(GL_TEXTURE_2D)` после загрузки и использовать `GL_LINEAR_MIPMAP_LINEAR`.

---

## 11. Известные ограничения и возможные улучшения

### Ограничения

1. **Нет теней.** Источник света точечный, но объекты не отбрасывают тени друг на друга.
2. **Нет ambient occlusion в реальном времени.** AO используется из текстуры (pre-baked), что не даёт динамических теней в углублениях.
3. **Bump mapping не изменяет геометрию.** На силуэте апельсин остаётся идеально гладким. Parallax occlusion mapping мог бы частично решить эту проблему.
4. **Касательные на CPU.** Вычисляются при загрузке и хранятся в VBO. Для очень больших моделей (миллионы треугольников) это увеличивает потребление памяти и время загрузки.
5. **Нет PBR.** Реализованы классические модели Фонга и Блинна-Фонга, а не физически корректный рендеринг.
6. **Отсутствует mipmapping.** Текстуры загружаются с фильтрацией `GL_LINEAR`, что может вызывать aliasing на больших расстояниях.
7. **Нет HDR-освещения.** Цвета ограничены диапазоном [0, 1] (LDR).
8. **Нет анимации.** Модели и источник статичны (вращение модели только по команде пользователя).

### Возможные улучшения

1. **Parallax occlusion mapping** — более продвинутая техника bump mapping, которая смещает текстурные координаты в зависимости от угла обзора, создавая иллюзию рельефа на силуэте.
2. **Динамические тени** (shadow mapping) — точечный источник может отбрасывать тени через cube map shadow map.
3. **HDR-освещение и tone mapping** — рендеринг в HDR-буфер с последующим tone mapping для более реалистичного диапазона яркостей.
4. **Отражения окружения** (cubemap) — фон в виде cubemap и отражения на блестящих поверхностях.
5. **SSAO (Screen Space Ambient Occlusion)** — динамическая ambient occlusion на основе глубины.
6. **PBR (Physically Based Rendering)** — переход на модель Кука — Торренса с использованием roughness/metalness карт.
7. **Mipmap-генерация** — `gl.GenerateMipmap` для всех текстур.
8. **Анимация вращения** — автоматическое вращение апельсина (переключаемое).
9. **Графический GUI** (ImGui) — полноценные слайдеры и кнопки вместо клавиатурного управления.
10. **Загрузка нескольких моделей** — расширение на произвольное количество объектов сцены.
