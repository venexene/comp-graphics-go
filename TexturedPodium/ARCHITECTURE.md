# ARCHITECTURE.md — Пьедестал почёта с мультитекстурированием

## 1. Обзор проекта

**Полное название лабораторной работы:** Пьедестал почёта с мультитекстурированием.

**Сцена:** Программа рендерит трёхмерный пьедестал почёта, состоящий из четырёх одинаковых кубиков. Компоновка: нижний ряд — три кубика вплотную (центр, левый, правый), верхний ряд — один кубик по центру над нижним центральным. На нижнем центральном кубике стоит декоративное сердце. 

**Графические технологии:**
- **Мультитекстурирование:** на каждый кубик накладываются две текстуры (материала и номера), которые смешиваются с индивидуальным цветом кубика.
- **Модели освещения:** Ламберт, Фонг, Блинн-Фонг, Тун (Toon), Орен-Наяр (Oren-Nayar).
- **Шейдинг:** Гуро (Gouraud) — освещение в вершинном шейдере; Фонг (Phong) — освещение во фрагментном шейдере.
- **Затухание света:** точечный источник с настраиваемым линейным/квадратичным затуханием.

**Что видит пользователь при запуске:** четыре кубика (жёлтый металл по центру, серый мрамор слева, оранжевое дерево справа) с наложенными текстурами номеров 1, 2, 3, красное сердце с текстурой оникса на первом пьедестале. Пользователь может вращать камеру, перемещать источник света, переключать модели освещения, регулировать веса смешивания текстур.

## 2. Дерево проекта

```
TexturedPodium/
├── .deepseek/
│   ├── planner.agent.md       # Спецификация планировщика
│   └── implementer.agent.md   # Спецификация имплементатора
├── ARCHITECTURE.md            # Данный файл
├── cmd/
│   └── main.go                # Точка входа: инициализация и главный цикл
├── go.mod                     # Go-модуль, зависимости
├── go.sum                     # Хеши зависимостей
├── input/
│   └── input.go               # Обработка клавиатурного ввода
├── lighting/
│   ├── light.go               # LightConfig — источник света
│   ├── material.go            # MaterialConfig — свойства материала
│   └── uniforms.go            # UniformCache — кеш uniform-переменных
├── models/
│   └── heart.obj              # 3D-модель сердца (OBJ-формат)
├── objects/
│   └── model.go               # Загрузка OBJ, CreateCube, VAO/VBO
├── scene/
│   ├── camera.go              # Орбитальная камера
│   ├── object.go              # ObjectState, SceneObject, Selection (устаревшие)
│   ├── podium.go              # Podium, PodiumCube, NewPodium
│   └── scene.go               # DrawPodium, set*Uniforms, CreateWhiteTexture
├── shaders/
│   ├── lighting.go            # 10 вариантов освещения (переключение)
│   ├── shaders.go             # LoadShaderFile, CompileShader
│   └── lighting/              # GLSL-файлы (см. ниже)
│       ├── basic_gouraud.frag        # Фрагментный шейдер для Gouraud
│       ├── basic_phong.vert          # Вершинный шейдер для Phong (передаёт данные)
│       ├── lambert_gouraud.vert      # Ламберт Gouraud: N·L в вершинах
│       ├── lambert_phong.frag        # Ламберт Phong: N·L в пикселях
│       ├── phong_gouraud.vert        # Фонг Gouraud: R·V в вершинах
│       ├── phong_phong.frag          # Фонг Phong: R·V в пикселях
│       ├── blinn_phong_gouraud.vert  # Блинн-Фонг Gouraud: N·H в вершинах
│       ├── blinn_phong_phong.frag    # Блинн-Фонг Phong: N·H в пикселях
│       ├── toon_gouraud.vert        # Тун Gouraud: ступенчатый свет + обводка
│       ├── toon_phong.frag          # Тун Phong: ступенчатый свет + обводка
│       ├── oren_nayar_gouraud.vert  # Орен-Наяр Gouraud
│       └── oren_nayar_phong.frag    # Орен-Наяр Phong
├── textures/
│   ├── texture.go             # LoadTexture — загрузка PNG/JPEG в OpenGL
│   ├── imgs/                  # Текстуры номеров
│   │   ├── 1.jpg              # Номер 1 (центральные кубики)
│   │   ├── 2.png              # Номер 2 (левый кубик)
│   │   └── 3.png              # Номер 3 (правый кубик)
│   └── materials/             # Текстуры материалов
│       ├── marble.jpg         # Мрамор (левый кубик, номер 2)
│       ├── metal.jpg          # Металл (центральные кубики, номер 1)
│       ├── onyx.jpg           # Оникс (сердце)
│       └── wood.jpg           # Дерево (правый кубик, номер 3)
├── ui/
│   └── ui.go                  # Текстовый UI (вывод в консоль)
└── utils/
    └── utils.go               # Инициализация сцены и связка компонентов
```

## 3. Графический пайплайн

### Инициализация (один раз при запуске)

```
cmd/main.go: main()
  └─ findProjectRoot() — поиск go.mod от CWD вверх
  └─ initGlfw() — создание окна 1280x720, контекст OpenGL 4.1 Core
  └─ initOpenGL()
      └─ gl.Init() — инициализация OpenGL
      └─ shaders.InitLightingVariants() — компиляция 10 шейдерных программ
          ├─ LoadShaderFile(lambert_gouraud.vert + basic_gouraud.frag)
          ├─ LoadShaderFile(basic_phong.vert + lambert_phong.frag)
          ├─ LoadShaderFile(phong_gouraud.vert + basic_gouraud.frag)
          ├─ LoadShaderFile(basic_phong.vert + phong_phong.frag)
          ├─ LoadShaderFile(blinn_phong_gouraud.vert + basic_gouraud.frag)
          ├─ LoadShaderFile(basic_phong.vert + blinn_phong_phong.frag)
          ├─ LoadShaderFile(toon_gouraud.vert + basic_gouraud.frag)
          ├─ LoadShaderFile(basic_phong.vert + toon_phong.frag)
          ├─ LoadShaderFile(oren_nayar_gouraud.vert + basic_gouraud.frag)
          └─ LoadShaderFile(basic_phong.vert + oren_nayar_phong.frag)
  └─ ui.InitializeUI() — подготовка GUI (no-op для текстового UI)
  └─ utils.InitScene(projectRoot)
      ├─ textures.LoadTexture(1.jpg)     — текстура номера 1
      ├─ textures.LoadTexture(2.png)     — текстура номера 2
      ├─ textures.LoadTexture(3.png)     — текстура номера 3
      ├─ textures.LoadTexture(metal.jpg) — материал для номера 1
      ├─ textures.LoadTexture(marble.jpg)— материал для номера 2
      ├─ textures.LoadTexture(wood.jpg)  — материал для номера 3
      ├─ objects.CreateCube(0.8)         — геометрия куба (36 вершин, 1 VAO)
      ├─ scene.NewPodium(...)            — 4 кубика + параметры
      ├─ objects.LoadOBJ(heart.obj)      — модель сердца
      └─ textures.LoadTexture(onyx.jpg)  — текстура сердца
  └─ gl.Enable(DEPTH_TEST)
  └─ mgl32.Perspective(45°, 1280/720, 0.1, 100.0) — матрица проекции
```

### Цикл рендера (каждый кадр)

```
ui.BeginFrame()
utils.DrawScene(window, projection)
  ├─ input.ProcessInput() — обработка клавиатуры (камера, свет, шейдеры, веса)
  ├─ shaders.GetCurrentLightingProgram() — выбор одной из 10 программ
  └─ scene.DrawPodium(program, podium, camera, projection, light, material, weights)
      ├─ gl.Clear(COLOR | DEPTH) — очистка буферов
      ├─ gl.UseProgram(program) — активация шейдера
      ├─ UniformCache.Refresh(program) — получение uniform-локаций
      ├─ установка View, Projection, ViewPos (общие для всех кубиков)
      ├─ для каждого из 4 кубиков:
      │   ├─ setTransformUniforms(ModelMatrix) — матрица куба
      │   ├─ setMaterialUniforms(MaterialConfig) — ambient/diffuse/specular/sheen/roughness
      │   ├─ setLightUniforms(LightConfig) — позиция/цвет/затухание света
      │   ├─ setBlendUniforms(MatTex, NumTex, Color, weights)
      │   │   ├─ gl.ActiveTexture(GL_TEXTURE0) → BindTexture(matTex) → MaterialMap=0
      │   │   ├─ gl.ActiveTexture(GL_TEXTURE1) → BindTexture(numTex) → NumberMap=1
      │   │   ├─ glUniform3f(u_cubeColor, color)
      │   │   ├─ glUniform1f(u_materialWeight, matWeight)
      │   │   └─ glUniform1f(u_numberWeight, numWeight)
      │   └─ glDrawArrays(GL_TRIANGLES, 0, 36) — отрисовка куба
      └─ если Heart != nil:
          ├─ setTransformUniforms(модель сердца) — сдвиг и масштаб
          ├─ setMaterialUniforms — те же material/light uniform
          ├─ setBlendUniforms(Onyx, Onyx, Red, 1.0, 0.0) — только материал
          └─ glDrawArrays — отрисовка сердца
  └─ glfw.PollEvents()
  └─ window.SwapBuffers()
└─ ui.EndFrame()
└─ обновление заголовка окна — текущие параметры сцены
```

### Завершение

```
defer glfw.Terminate()
defer shaders.CleanupLightingVariants() — удаление 10 шейдерных программ
defer ui.Cleanup()
defer utils.Cleanup() — удаление VAO кубов и сердца
```

## 4. Геометрия сцены

### Создание куба (objects/model.go: CreateCube)

- **Размер:** 0.8 единиц (сторона), параметризуемый.
- **Вершины:** 36 вершин (6 граней × 2 треугольника × 3 вершины).
- **Формат вершины (interleaved VBO):**
  - Позиция (location 0): `vec3` — 3 × float32 = 12 байт
  - Нормаль (location 1): `vec3` — 3 × float32 = 12 байт
  - UV (location 3): `vec2` — 2 × float32 = 8 байт
  - **Итого на вершину:** 32 байта (8 float)
- **Страйд:** 32 байта (8 × 4).
- **Каждая грань** имеет уникальную нормаль (плоское затенение) и UV-координаты [0,1]².
- **Порядок треугольников:** CCW (counter-clockwise) для всех граней.
- **Отрисовка:** `glDrawArrays(GL_TRIANGLES, 0, 36)` — без индексного буфера.

### Компоновка пьедестала (scene/podium.go: NewPodium)

Все четыре кубика разделяют **один VAO/VBO** — геометрия создаётся один раз через `CreateCube(cubeSize)`. Каждый кубик имеет свою матрицу Model (трансляция + масштаб 1.0):

| Индекс | Позиция (x, y, z) | Номер | Материал | Цвет |
|--------|-------------------|-------|----------|------|
| 0 (нижний центр) | (0, 0, 0) | 1 — 1.jpg | metal | жёлтый (1,1,0) |
| 1 (левый) | (−0.8, 0, 0) | 2 — 2.png | marble | серый (0.5,0.5,0.5) |
| 2 (правый) | (0.8, 0, 0) | 3 — 3.png | wood | оранжевый (1,0.5,0) |
| 3 (верхний центр) | (0, 0.8, 0) | 1 — 1.jpg | metal | жёлтый (1,1,0) |

Кубики стоят **вплотную**: при size=0.8 и spacing=0.8 грань левого кубика касается грани центрального.

### Сердце (heart.obj)

Расположено над нижним центральным кубиком (индекс 0) на высоте y = 0.4 + 0.8 = 1.2. Масштаб 0.2. Использует отдельную модель (OBJ), загружаемую через `objects.LoadOBJ`.

## 5. Мультитекстурирование

### Загружаемые текстуры

| Файл | Формат | Назначение |
|------|--------|------------|
| textures/imgs/1.jpg | JPEG | Номер 1 (центральные кубики) |
| textures/imgs/2.png | PNG | Номер 2 (левый кубик) |
| textures/imgs/3.png | PNG | Номер 3 (правый кубик) |
| textures/materials/metal.jpg | JPEG | Металл (центральные кубики) |
| textures/materials/marble.jpg | JPEG | Мрамор (левый кубик) |
| textures/materials/wood.jpg | JPEG | Дерево (правый кубик) |
| textures/materials/onyx.jpg | JPEG | Оникс (сердце) |

### Привязка к текстурным блокам

- **GL_TEXTURE0** (unit 0) — текстура материала (`u_materialTexture`)
- **GL_TEXTURE1** (unit 1) — текстура номера (`u_numberTexture`)

Устанавливается в `scene/scene.go: setBlendUniforms()`.

### Формула смешивания (во всех 10 фрагментных и 5 вершинных Gouraud-шейдерах)

```glsl
// Выборка текстур
vec3 matColor = texture(u_materialTexture, texCoord).rgb;
vec3 numColor = texture(u_numberTexture, texCoord).rgb;

// Нормализация весов
float totalWeight = u_materialWeight + u_numberWeight;

// Смешивание: цвет кубика × (материал × вес + номер × вес) / сумма весов
vec3 surfaceColor = u_cubeColor * (u_materialWeight * matColor
                                   + u_numberWeight * numColor)
                   / max(totalWeight, 0.001);
```

Далее `surfaceColor` — это базовый цвет поверхности, который участвует в расчётах освещения (умножается на ambient/diffuse/specular).

### Uniform-переменные смешивания

| Имя в шейдере | Имя в Go (UniformCache) | Клавиши | Диапазон |
|---------------|------------------------|---------|----------|
| `u_materialTexture` | `MaterialMap` | — | texture unit 0 |
| `u_numberTexture` | `NumberMap` | — | texture unit 1 |
| `u_cubeColor` | `CubeColor` | — | vec3 [0,1] |
| `u_materialWeight` | `MaterialWeight` | 5 (-) / 6 (+) | [0.0, 1.0] |
| `u_numberWeight` | `NumberWeight` | 7 (-) / 8 (+) | [0.0, 1.0] |

### Поведение весов

- **materialWeight = 1.0, numberWeight = 0.0:** видна только текстура материала, умноженная на цвет кубика.
- **materialWeight = 0.0, numberWeight = 1.0:** видна только текстура номера, умноженная на цвет кубика.
- **materialWeight = 0.5, numberWeight = 0.5:** обе текстуры смешиваются 50/50.

### Сердце

Сердце рендерится с `materialWeight = 1.0, numberWeight = 0.0` — только текстура материала (onyx) × красный цвет. Номер не накладывается.

## 6. Модели освещения

### Ламберт (Lambert)

Только диффузная компонента: `L_d = k_d * I * max(0, N·L)`, где:
- `k_d` — коэффициент диффузного отражения из `material.diffuse`
- `I` — интенсивность света из `light.diffuse`
- `N·L` — косинус угла между нормалью и направлением на свет

Зеркальная компонента отсутствует. Шейдеры: `lambert_gouraud.vert`, `lambert_phong.frag`.

### Фонг (Phong)

Полная модель: `I = I_a + I_d + I_s`, где:
- `I_a = light.ambient * base_color * material.ambient * light.ambient_strength * attenuation`
- `I_d = light.diffuse * base_color * material.diffuse * max(0, N·L) * attenuation`
- `I_s = light.specular * material.specular * max(0, R·V)^sheen * attenuation`
- `R = reflect(-L, N)` — отражённый вектор
- `V` — направление от поверхности к камере

Шейдеры: `phong_gouraud.vert` (specular в вершинах), `phong_phong.frag` (specular в пикселях).

### Блинн-Фонг (Blinn-Phong)

Отличается от Фонга только вычислением зеркальной компоненты:
- `I_s = light.specular * material.specular * max(0, N·H)^sheen * attenuation`
- `H = normalize(L + V)` — half-vector (вектор половины угла)

Быстрее и точнее классического Фонга. Шейдеры: `blinn_phong_gouraud.vert`, `blinn_phong_phong.frag`.

### Тун (Toon/ Cel-shading)

Ступенчатое освещение с силуэтной обводкой:
- `toon_diff = floor(N·L * levels) / levels` — квантование диффузной компоненты на `levels = 3` уровня
- `toon_spec = floor(spec * levels) / levels` — квантование зеркальной
- `edge_factor = smoothstep(0.0, 0.25, max(N·V, 0.0))` — затемнение на ребрах (обводка)

Итоговый цвет: `(ambient + diffuse + specular) * edge_factor`.

Шейдеры: `toon_gouraud.vert`, `toon_phong.frag`.

### Орен-Наяр (Oren-Nayar)

Диффузная модель для шероховатых поверхностей:
- Расширение Ламберта: учитывает шероховатость `roughness`
- `rough2 = roughness²`
- `A = 1 - 0.5 * rough² / (rough² + 0.33)`
- `B = 0.45 * rough² / (rough² + 0.09)`
- `OrenNayar = N·L * (A + B * max(0, cos(Δφ)) * sin(α) * tan(β)) * attenuation`

Зеркальная компонента отсутствует. Шейдеры: `oren_nayar_gouraud.vert`, `oren_nayar_phong.frag`.

### Взаимодействие с мультитекстурированием

Во всех моделях последовательность одинакова:
1. Выборка текстур и вычисление `surfaceColor` (смешивание материала, номера и цвета)
2. Расчёт освещения с использованием `surfaceColor` вместо `base_color`

## 7. Шейдинг: Гуро против Фонга

### Гуро (Gouraud)

- **Вершинный шейдер** (например `lambert_gouraud.vert`): вычисляет `surfaceColor`, затем освещение, выдаёт `vert_color`.
- **Фрагментный шейдер** (`basic_gouraud.frag`): выборка текстур (материал + номер), смешивание цветов, умножение на интерполированный `vert_color`:
  ```glsl
  vec3 surfaceColor = ...; // смешивание текстур
  frag_color = vec4(surfaceColor * vert_color, 1.0);
  ```
- **Артефакты:** на крупных гранях куба зеркальные блики могут потеряться из-за интерполяции цвета между вершинами.

### Фонг (Phong)

- **Вершинный шейдер** (`basic_phong.vert`): передаёт во фрагментный нормаль, направление света, направление взгляда, UV-координаты через `out Vertex { ... } vert;`.
- **Фрагментный шейдер** (например `lambert_phong.frag`): принимает интерполированные данные, вычисляет `surfaceColor` (смешивание текстур), затем освещение попиксельно.
- **Преимущества:** корректные блики, плавные переходы освещения.

### Механизм переключения

Переключение происходит через **полную смену шейдерной программы** (`glUseProgram`). 10 программ закэшированы в `lightingVariants`. Клавиши:
- **T/G** — циклический перебор моделей освещения (Lambert → Phong → Blinn-Phong → Toon → Oren-Nayar)
- **Y** — переключение Gouraud ↔ Phong для текущей модели

## 8. Система освещения

### Тип источника

Точечный (PointLight). Свет распространяется от одной точки во все стороны.

### Структура LightConfig (lighting/light.go)

| Поле | Тип | Описание | Uniform |
|------|-----|----------|---------|
| Position | mgl32.Vec3 | Положение в world space | light.position |
| Ambient | mgl32.Vec3 | Цвет фонового освещения (RGB) | light.ambient |
| Diffuse | mgl32.Vec3 | Цвет диффузного освещения (RGB) | light.diffuse |
| Specular | mgl32.Vec3 | Цвет зеркального освещения (RGB) | light.specular |
| Constant | float32 | Постоянный коэфф. затухания (c) | light.constant |
| Linear | float32 | Линейный коэфф. затухания (k_l) | light.linear |
| Quadratic | float32 | Квадратичный коэфф. затухания (k_q) | light.quadratic |
| LinearCoef | float32 | Множитель линейного затухания | linear_coef |
| QuadraticCoef | float32 | Множитель квадратичного затухания | quadratic_coef |
| AmbientStrength | float32 | Мощность фонового света [0, 1] | light.ambient_strength |
| Mode | AttenuationMode | Тип затухания | attenuation_mode |

### Затухание

Формула затухания в зависимости от `attenuation_mode`:

| Mode | Формула |
|------|---------|
| Both (0) | `1.0 / (c + k_l * d + k_q * d²)` |
| Linear (1) | `1.0 / (c + k_l * d)` |
| Quadratic (2) | `1.0 / (c + k_q * d²)` |

Переключение: клавиша **M**.

### Фоновое освещение (Ambient)

`ambient = light.ambient * surfaceColor * light.ambient_strength * attenuation`

Регулировка: клавиши **B** (уменьшить) / **N** (увеличить), диапазон [0.0, 1.0].

## 9. GUI и управление

GUI реализован как текстовый вывод в консоль (пакет `ui/ui.go`). Основная информация выводится в заголовок окна GLFW.

### Элементы управления

| Название | Тип | Переменная/Uniform | Клавиши | Диапазон |
|----------|-----|-------------------|---------|----------|
| Выбор модели освещения | циклический | LightingVariant.Model | T (вперёд) / G (назад) | 5 моделей |
| Переключение шейдинга | тумблер | LightingVariant.Mode → u_shadingModel | Y | Gouraud/Phong |
| Тип затухания | циклический | AttenuationMode → attenuation_mode | M | Both/Linear/Quad |
| Коэфф. линейного затухания | инкремент | LightConfig.LinearCoef → linear_coef | Z (−) / X (+) | [0, ∞) |
| Коэфф. квадратичного затухания | инкремент | LightConfig.QuadraticCoef → quadratic_coef | C (−) / V (+) | [0, ∞) |
| Ambient strength | инкремент | LightConfig.AmbientStrength → light.ambient_strength | B (−) / N (+) | [0, 1] |
| Вес текстуры материала | инкремент | LightConfig.MaterialWeight → u_materialWeight | 5 (−) / 6 (+) | [0, 1] |
| Вес текстуры номера | инкремент | LightConfig.NumberWeight → u_numberWeight | 7 (−) / 8 (+) | [0, 1] |
| Вращение камеры | непрерывный | Camera.Yaw/Pitch | ← ↑ → ↓ | — |
| Панорама | непрерывный | Camera.Target | WASD | — |
| Панорама вверх/вниз | непрерывный | Camera.Target.Y | Space/Shift | — |
| Зум | инкремент | Camera.Distance | +/- | [1, 50] |
| Перемещение света | непрерывный | LightConfig.Position | Alt+IJKLUO | — |
| Сброс | кнопка | камера + объект | Ctrl+R | — |

## 10. Система координат и матрицы

### Библиотека

Матричные вычисления: `github.com/go-gl/mathgl` (mgl32).

### Матрицы

- **Model** (mat4): локальная система куба → мировая. Для каждого кубика: `Translate(Position) × Scale(size)`.
- **View** (mat4): мировая → система камеры. `LookAt(eye, target, up)`.
- **Projection** (mat4): камера → экран. `Perspective(45°, 1280/720, 0.1, 100.0)`.

### Порядок умножения в вершинном шейдере

```glsl
gl_Position = transform.projection * transform.view * transform.model * vec4(position, 1.0);
```

### Положение камеры по умолчанию

- Target: (0, 0.5, 0)
- Distance: 6.0
- Yaw: 0.7 рад (~40°) — камера спереди-справа
- Pitch: 0.35 рад (~20°) — камера чуть сверху

### Источник света

Начальная позиция: (4, 3, 5) — спереди-справа-сверху.

## 11. Известные ограничения и возможные улучшения

### Ограничения

1. **Нет теней.** Точечный источник света не отбрасывает тени. Все объекты освещаются одинаково.
2. **Один набор весов на все кубики.** Параметры `materialWeight` и `numberWeight` глобальные — регулируются для всех четырёх кубиков одновременно.
3. **Нет нормал-маппинга.** Детализация поверхности только через текстуры материалов.
4. **Артефакты шейдинга Гуро.** На крупных гранях куба зеркальные блики могут теряться.
5. **Нет back-face culling.** Отрисовываются обе стороны каждого треугольника, хотя задние грани скрыты depth test.
6. **Текстовый UI.** Вместо графического интерфейса — вывод в консоль и заголовок окна.

### Возможные улучшения

1. **Specular map.** Чёрно-белая текстура для управления зеркальным отражением на разных участках.
2. **Normal mapping.** Карта нормалей для имитации микрорельефа.
3. **Cubemap-отражения.** Окружение для отражений на металлических поверхностях.
4. **Раздельные веса для каждого кубика.** Возможность настраивать смешивание текстур индивидуально.
5. **ImGui.** Полноценный графический интерфейс вместо текстового вывода.
6. **Загрузка произвольных моделей.** OBJ-загрузчик уже есть, но можно добавить поддержку других форматов.
7. **Анимация.** Вращение пьедестала или пульсация сердца.
