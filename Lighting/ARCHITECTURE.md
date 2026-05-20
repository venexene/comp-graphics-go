# Архитектура проекта: Лабораторная работа по моделям освещения в OpenGL

## 1. Обзор проекта

**Название:** Lighting — Интерактивная демонстрация моделей освещения в OpenGL.

**Сцена:** Программа рендерит трёхмерную сцену, состоящую из трёх объектов:
- **Снеговик** (`models/snowman.obj`) — главный объект, расположен в центре сцены. Создан из трёх сфер в Blender, экспортирован в OBJ.
- **Сердце** (`models/heart.obj`) — дополнительный объект, расположен справа от снеговика (по +X).
- **Default** (`models/default.obj`) — дополнительный объект, расположен слева (по -X).

**Демонстрируемые технологии:**
- OpenGL 4.1 Core Profile через биндинги Go (go-gl/gl, go-gl/glfw, go-gl/mathgl)
- Четыре модели освещения: Lambert, Phong, Blinn-Phong, Toon
- Два режима шейдинга: Gouraud (вершинный) и Phong (попиксельный)
- Затухание света с тремя режимами: Both, Linear, Quadratic
- Интерактивное управление камерой (орбита), объектами и светом с клавиатуры
- Загрузка 3D-моделей из Wavefront OBJ-файлов

---

## 2. Дерево проекта

```
Lighting/
├── go.mod                    # Модуль Go, зависимости: gl, glfw, mathgl
├── ARCHITECTURE.md           # Этот файл
│
├── cmd/
│   └── main.go               # Точка входа: инициализация окна, загрузка моделей,
│                              #   цикл рендеринга
│
├── lighting/
│   ├── light.go              # Структура LightConfig (точечный источник света)
│   │                         #   и AttenuationMode (режимы затухания)
│   ├── material.go           # Структура MaterialConfig (параметры материала)
│   └── uniforms.go           # Кэш UniformCache — расположения uniform-переменных
│                              #   шейдерных программ
│
├── shaders/
│   ├── shaders.go            # Загрузка и компиляция GLSL-шейдеров
│   ├── lighting.go           # 8 вариантов освещения: компиляция, переключение
│   └── lighting/             # GLSL-файлы шейдеров (12 файлов)
│       ├── basic_gouraud.frag        # Общий фрагментный шейдер для Gouraud
│       ├── basic_phong.vert          # Общий вершинный шейдер для Phong
│       ├── lambert_gouraud.vert      # Lambert Gouraud (только diffuse)
│       ├── lambert_phong.frag        # Lambert Phong
│       ├── phong_gouraud.vert        # Phong Gouraud (с отражённым вектором R)
│       ├── phong_phong.frag          # Phong Phong
│       ├── blinn_phong_gouraud.vert  # Blinn-Phong Gouraud (с half-vector H)
│       ├── blinn_phong_phong.frag    # Blinn-Phong Phong
│       ├── toon_gouraud.vert         # Toon Gouraud (дискретные уровни + контур)
│       ├── toon_phong.frag           # Toon Phong
│       ├── oren_nayar_gouraud.vert   # Oren-Nayar Gouraud (шероховатые пов-ти)
│       └── oren_nayar_phong.frag     # Oren-Nayar Phong
│
├── scene/
│   ├── scene.go              # Отрисовка сцены: очистка буферов, uniform-переменные
│   ├── camera.go             # Орбитальная камера (сферические координаты)
│   └── object.go             # Трансформации объектов, переключение выбора
│
├── objects/
│   └── model.go              # Парсинг OBJ, создание VAO/VBO, отрисовка
│
├── input/
│   └── input.go              # Клавиатурный ввод: камера, объекты, свет, шейдеры
│
├── ui/
│   └── ui.go                 # Текстовый UI (ASCII-панель в терминале)
│
├── utils/
│   └── utils.go              # Слой совместимости, глобальное состояние сцены
│
└── models/                   # OBJ-файлы 3D-моделей
    ├── snowman.obj            # Снеговик (главный объект)
    ├── heart.obj              # Сердце (дополнительный)
    └── default.obj            # Default-модель (дополнительная)
```

---

## 3. Графический пайплайн

```
┌──────────────────────────────────────────────────────────────────┐
│                     ИНИЦИАЛИЗАЦИЯ (один раз)                      │
├──────────────────────────────────────────────────────────────────┤
│                                                                   │
│  1. findProjectRoot() → поиск go.mod вверх по дереву              │
│  2. initGlfw() → создание окна 1280×720, Core Profile 4.1        │
│  3. initOpenGL() → gl.Init(), InitLightingVariants()              │
│     ├── CompileShader(lambert_gouraud.vert, VERTEX)               │
│     ├── CompileShader(basic_gouraud.frag, FRAGMENT)               │
│     ├── LinkProgram() → программа #0 (Lambert Gouraud)            │
│     ├── ... (ещё 7 вариантов)                                     │
│     └── Все 8 программ готовы                                     │
│  4. CreateWhiteTexture() → белая текстура 1×1                     │
│  5. LoadOBJ(snowman.obj) → VAO/VBO снеговика                     │
│  6. LoadOBJ(heart.obj), LoadOBJ(default.obj)                     │
│  7. glEnable(DEPTH_TEST), LookAt, Perspective                     │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│                      ЦИКЛ РЕНДЕРИНГА (каждый кадр)                │
├──────────────────────────────────────────────────────────────────┤
│                                                                   │
│  1. ProcessInput(window, &cam, &mainState, &lightCfg, sel)       │
│     ├── Стрелки → вращение камеры                                 │
│     ├── WASD → панорамирование камеры                             │
│     ├── Q/E → масштаб объекта                                     │
│     ├── IJKL/UO → перемещение объекта (Alt — свет)               │
│     ├── T/G → следующий/предыдущий вариант освещения              │
│     ├── Y → переключение Gouraud↔Phong                            │
│     ├── M → режим затухания                                       │
│     ├── Z/X → линейный коэффициент                                │
│     ├── C/V → квадратичный коэффициент                            │
│     └── B/N → мощность фонового света                             │
│                                                                   │
│  2. GetCurrentLightingProgram() → ID текущей шейдерной программы  │
│                                                                   │
│  3. DrawScene(program, model, mainState, extras, cam, ...)       │
│     ├── glClearColor(0.2, 0.3, 0.3)                              │
│     ├── glClear(COLOR | DEPTH)                                    │
│     ├── glUseProgram(program)                                     │
│     ├── UniformCache.Refresh(program)                              │
│     ├── cam.ViewMatrix() → view                                   │
│     ├── glUniform3f(view_pos, eye)                                │
│     ├── glBindTexture(diffuse_map, defaultTex)                    │
│     │                                                             │
│     ├── Для снеговика:                                            │
│     │   modelMat = mainState.ModelMatrix()                        │
│     │   → setTransformUniforms(modelMat, view, proj)              │
│     │   → setMaterialUniforms(matCfg)                             │
│     │   → setLightUniforms(lightCfg)                              │
│     │   → model.Draw() → glDrawArrays(GL_TRIANGLES)               │
│     │                                                             │
│     ├── Для сердца (heartObj):                                    │
│     │   modelMat = heartObj.ModelMatrix()                         │
│     │   → ...те же uniform... → model.Draw()                     │
│     │                                                             │
│     └── Для default (defaultObj):                                 │
│         modelMat = defaultObj.ModelMatrix()                       │
│         → ...те же uniform... → model.Draw()                     │
│                                                                   │
│  4. glfw.PollEvents()                                             │
│  5. window.SwapBuffers()                                          │
│  6. window.SetTitle(...) → обновление заголовка                   │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│                      ЗАВЕРШЕНИЕ                                   │
├──────────────────────────────────────────────────────────────────┤
│  1. model.Delete() → glDeleteVertexArrays, glDeleteBuffers       │
│  2. CleanupLightingVariants() → glDeleteProgram для всех 8       │
│  3. glfw.Terminate()                                              │
└──────────────────────────────────────────────────────────────────┘
```

---

## 4. Модели освещения

### 4.1. Lambert (Ламберт)

**Формула:** I = I_a · k_a + I_d · k_d · max(N·L, 0)

Где:
- I_a — цвет фонового освещения (light.ambient)
- k_a — коэффициент фонового отражения материала (material.ambient)
- I_d — цвет диффузного освещения (light.diffuse)
- k_d — коэффициент диффузного отражения (material.diffuse)
- N — нормаль поверхности в world space
- L — направление от поверхности к источнику света

**Где реализована:**
- Вершинный шейдер: `shaders/lighting/lambert_gouraud.vert` (Gouraud)
- Фрагментный шейдер: `shaders/lighting/lambert_phong.frag` (Phong)

**Выбор варианта:**
- Нажатие T/G переключает между всеми 8 вариантами по циклу
- Lambert Gouraud — `lightingVariants[0]`
- Lambert Phong — `lightingVariants[1]`

**Особенности:** Только диффузная составляющая. Зеркальные блики отсутствуют. Самый быстрый вариант.

---

### 4.2. Phong (Фонг)

**Формула:** I = I_a · k_a + I_d · k_d · max(N·L, 0) + I_s · k_s · (max(V·R, 0))^sh

Где:
- I_s — цвет зеркального освещения (light.specular)
- k_s — коэффициент зеркального отражения (material.specular)
- V — направление от поверхности к камере
- R — отражённый вектор: reflect(-L, N)
- sh — экспонента блеска (material.sheen_coef)

**Где реализована:**
- Вершинный шейдер: `shaders/lighting/phong_gouraud.vert` (Gouraud)
- Фрагментный шейдер: `shaders/lighting/phong_phong.frag` (Phong)

**Выбор варианта:**
- Phong Gouraud — `lightingVariants[2]`
- Phong Phong — `lightingVariants[3]`

**Особенности:** Классическая модель Фонга с отражённым лучом. Более реалистичные блики, чем у Блинн-Фонга, но дороже в вычислении (нужен reflect).

---

### 4.3. Blinn-Phong (Блинн-Фонг)

**Формула:** I = I_a · k_a + I_d · k_d · max(N·L, 0) + I_s · k_s · (max(N·H, 0))^sh

Где:
- H — половинный вектор: normalize(L + V)
- Вместо R используется H, что дешевле

**Где реализована:**
- Вершинный шейдер: `shaders/lighting/blinn_phong_gouraud.vert` (Gouraud)
- Фрагментный шейдер: `shaders/lighting/blinn_phong_phong.frag` (Phong)

**Выбор варианта:**
- Blinn-Phong Gouraud — `lightingVariants[4]`
- Blinn-Phong Phong — `lightingVariants[5]`

**Особенности:** Стандарт де-факто в современной компьютерной графике. Даёт slightly different блики по сравнению с Phong, но быстрее.

---

### 4.4. Toon (Мультяшный)

**Формула:** I = (I_a · k_a + I_d · k_d · toon_diff + I_s · k_s · toon_spec) · edge_factor

Где:
- toon_diff = floor(max(N·L, 0) · levels) / levels — квантование на `levels=3` ступени
- toon_spec = floor((max(N·H, 0))^sh · levels) / levels — квантование зеркального
- edge_factor = smoothstep(0.15, 0.05, N·V) — затемнение контурных граней

**Где реализована:**
- Вершинный шейдер: `shaders/lighting/toon_gouraud.vert` (Gouraud)
- Фрагментный шейдер: `shaders/lighting/toon_phong.frag` (Phong)

**Выбор варианта:**
- Toon Gouraud — `lightingVariants[6]`
- Toon Phong — `lightingVariants[7]`

**Особенности:** Создаёт эффект мультипликационной графики (cel-shading). Тени имеют резкие границы. Контурные грани затемняются для имитации силуэтных линий.

### 4.5. Oren-Nayar

**Формула:** I = I_a + I_d · ρ · max(N·L, 0) · (A + B · max(0, cos(Δφ)) · sin(α) · tan(β))

Где:
- ρ — коэффициент шероховатости (roughness)
- A = 1 - 0.5 · σ² / (σ² + 0.33)
- B = 0.45 · σ² / (σ² + 0.09)
- σ = roughness
- α = max(θ_i, θ_r), β = min(θ_i, θ_r)
- θ_i — угол между N и L, θ_r — угол между N и V
- cos(Δφ) — азимутальная разница между проекциями L и V

**Где реализована:**
- Вершинный шейдер: `shaders/lighting/oren_nayar_gouraud.vert` (Gouraud)
- Фрагментный шейдер: `shaders/lighting/oren_nayar_phong.frag` (Phong)

**Особенности:** Улучшенная модель для матовых/шероховатых поверхностей. Не имеет зеркальной составляющей. Использует дополнительную uniform `roughness`.

---

## 5. Система освещения

### 5.1. Тип источника — точечный (Point Light)

Структура `LightConfig` (`lighting/light.go`):

| Поле | Тип | Описание | uniform в шейдере |
|------|-----|----------|-------------------|
| Position | mgl32.Vec3 | Позиция в world space | light.position |
| Ambient | mgl32.Vec3 | Цвет фонового света (RGB) | light.ambient |
| Diffuse | mgl32.Vec3 | Цвет диффузного света (RGB) | light.diffuse |
| Specular | mgl32.Vec3 | Цвет зеркального света (RGB) | light.specular |
| Constant | float32 | Постоянный коэффициент c | light.constant |
| Linear | float32 | Линейный коэффициент k_l | light.linear |
| Quadratic | float32 | Квадратичный коэффициент k_q | light.quadratic |
| LinearCoef | float32 | Множитель линейного затухания (0..∞) | linear_coef |
| QuadraticCoef | float32 | Множитель квадратичного затухания (0..∞) | quadratic_coef |
| AmbientStrength | float32 | Мощность фонового света [0, 1] | light.ambient_strength |
| Mode | AttenuationMode | Режим затухания (0=Both,1=Linear,2=Quad) | attenuation_mode |

Начальное положение света: `(2.0, 4.0, 3.0)` — справа-сверху-спереди.

### 5.2. Затухание (Attenuation)

Формула затухания: **attenuation = 1 / max(c + k_l · d + k_q · d², ε)**

Где ε = 0.0001 (защита от деления на ноль).

Три режима переключения (клавиша M):
0. **Both** — используются все три коэффициента
1. **Linear** — только линейный: `c + k_l · d`
2. **Quadratic** — только квадратичный: `c + k_q · d²`

Пользователь управляет множителями `LinearCoef` и `QuadraticCoef` клавишами Z/X и C/V соответственно. Базовые коэффициенты:
- Linear: 0.09 (физически корректное значение)
- Quadratic: 0.032 (физически корректное значение)

### 5.3. Фоновое освещение (Ambient)

Мощность фонового света `AmbientStrength` регулируется клавишами B/N в диапазоне [0, 1]. Начальное значение: 0.6. Учитывается во всех шейдерах путём умножения на `light.ambient_strength`.

---

## 6. Шейдинг: Гуро против Фонга

### 6.1. Шейдинг Гуро (Gouraud)

- Освещение вычисляется **в вершинном шейдере** для каждой вершины треугольника.
- Результирующий цвет вершины (`vert_color`) интерполируется по площади треугольника.
- **Плюсы:** Быстрее (меньше вычислений на кадр).
- **Минусы:** Артефакты на низкополигональных моделях — блики могут быть угловатыми или пропадать между вершинами.
- **Файлы:** Все `*_gouraud.vert` + `basic_gouraud.frag`.

### 6.2. Шейдинг Фонга (Phong)

- Нормаль (`normal`), направление на свет (`light_dir`) и на камеру (`view_dir`) интерполируются через интерфейсный блок `Vertex`.
- Освещение вычисляется **во фрагментном шейдере** для каждого пикселя.
- **Плюсы:** Более качественные блики, корректное освещение на низкополигональных моделях.
- **Минусы:** Дороже (вычисления на каждый пиксель).
- **Файлы:** `basic_phong.vert` + все `*_phong.frag`.

### 6.3. Переключение в коде

Переключение между Гуро и Фонг для одной модели освещения — клавиша **Y**. В коде (`shaders/lighting.go:ToggleShadingMode()`):
- Находит вариант с тем же `Model` и противоположным `Mode`.
- Индекс `currentLightingIndex` переключается на найденный вариант.
- Следующий кадр использует новую шейдерную программу.

---

## 7. Сцена и модели

### 7.1. Снеговик

- Файл: `models/snowman.obj`
- Создан в Blender из трёх сфер (нижняя, средняя, голова).
- Экспортирован в формат Wavefront OBJ.
- Расположен в центре сцены (0, 0, 0).

### 7.2. Сердце и Default

- Файлы: `models/heart.obj`, `models/default.obj`
- Расположены на оси X на расстоянии ±4 единицы от центра.
- Масштабированы 0.6 от исходного размера.
- Без текстур — используют белую текстуру-заглушку.

### 7.3. Загрузка и отрисовка

Загрузка: `objects.LoadOBJ(path)` — парсинг OBJ-файла:
1. Чтение вершин (v), текстурных координат (vt), нормалей (vn).
2. Триангуляция граней (f).
3. Усреднение нормалей для сглаженного шейдинга.
4. Создание `Model` (VAO + VBO) через `CreateModelFromVertices()`.

Отрисовка: `model.Draw()` → `glDrawArrays(GL_TRIANGLES)`.

---

## 8. GUI и управление

### 8.1. Тип UI

Текстовый UI — ASCII-панель, выводится в терминал при старте. В реальном времени информация отображается в заголовке окна GLFW. Графический UI (Dear ImGui, Nuklear) не используется.

### 8.2. Элементы управления

| Клавиша | Действие | Переменная в коде | Uniform в шейдере | Диапазон |
|---------|----------|-------------------|-------------------|----------|
| T/G | Следующий/предыдущий вариант освещения | `currentLightingIndex` | — | 0..7 |
| Y | Переключение Гуро↔Фонг | `shaders.ToggleShadingMode()` | — | — |
| M | Режим затухания (Both→Linear→Quad) | `lightCfg.Mode` | `attenuation_mode` | 0,1,2 |
| Z/X | Линейный коэффициент | `lightCfg.LinearCoef` | `linear_coef` | [0, ∞) |
| C/V | Квадратичный коэффициент | `lightCfg.QuadraticCoef` | `quadratic_coef` | [0, ∞) |
| B/N | Мощность фонового света | `lightCfg.AmbientStrength` | `light.ambient_strength` | [0, 1] |
| Стрелки | Вращение камеры | `cam.Rotate()` | — | — |
| WASD | Панорамирование камеры | `cam.PanForward/Right()` | — | — |
| +/- | Зум | `cam.Zoom()` | — | [1, 50] |
| Q/E | Масштаб объекта | `mainState.Scale` | — | [0.1, 3.0] |
| R/F | Вращение вокруг Z | `mainState.RotationZ` | — | — |
| 1/2 | Вращение вокруг X | `mainState.RotationX` | — | — |
| 3/4 | Вращение вокруг Y | `mainState.RotationY` | — | — |
| IJKL | Перемещение объекта XY | `Position` | — | — |
| U/O | Перемещение объекта Z | `Position` | — | — |
| Alt+IJKL | Перемещение света XY | `lightCfg.Position` | — | — |
| Alt+U/O | Перемещение света Z | `lightCfg.Position` | — | — |
| Tab | Выбор следующего объекта | `sel.CycleForward()` | — | — |
| Ctrl+R | Сброс камеры и объекта | `DefaultObjectState()`, `DefaultCamera()` | — | — |

### 8.3. Отображение

- **Заголовок окна** обновляется каждый кадр: содержит выбранный объект, позицию света, модель освещения, режим шейдинга, параметры затухания.
- **ASCII-панель** выводится в терминал один раз при запуске.

---

## 9. Система координат и матрицы

### 9.1. Библиотека

Для работы с матрицами используется `github.com/go-gl/mathgl/mgl32`.

### 9.2. Матрицы

| Матрица | Пространство | uniform в шейдере | Вычисление |
|---------|--------------|-------------------|------------|
| Model | Object → World | `transform.model` | `ObjectState.ModelMatrix()`: T · Rz · Ry · Rx · S |
| View | World → View | `transform.view` | `Camera.ViewMatrix()`: `mgl32.LookAt(eye, target, up)` |
| Projection | View → Clip | `transform.projection` | `mgl32.Perspective(fov, aspect, near, far)` |
| Normal | Object → World (нормали) | `transform.normal_mat` | `(M⁻¹)ᵀ`, подматрица 3×3 от Model |

### 9.3. Параметры проекции

- FOV: 45°
- Aspect ratio: 16:9 (1280 / 720)
- Near plane: 0.1
- Far plane: 100.0

### 9.4. Камера

Начальное положение: `(5, 2, 5)`, смотрит в `(0, 0, 0)`.
Орбитальная камера — вращается вокруг целевой точки по сферическим координатам:
- Distance (расстояние до цели): начальное 5.0
- Yaw (рыскание): вращение вокруг оси Y
- Pitch (тангаж): вертикальный угол, ограничен ±π/2

### 9.5. Источник света

Начальное положение в world space: `(2.0, 4.0, 3.0)`.

---

## 10. Известные ограничения и возможные улучшения

### Ограничения

1. **Отсутствие теней.** Источник света точечный, но тени не отбрасываются. Все объекты освещаются одинаково, даже если находятся за другими объектами с точки зрения света.
2. **Один источник света.** Сцена использует только один точечный источник. Нет поддержки направленного, прожекторного или множественных источников.
3. **Нет PBR.** Модели освещения — классические (Lambert, Phong, Blinn-Phong). Физически корректный рендеринг (PBR) не реализован.
4. **Текстуры — только заглушка.** Все модели используют белую текстуру 1×1. Настоящие текстурные карты не загружаются.
5. **Фиксированный материал.** Параметры материала заданы глобально для всех объектов. Нет per-объектных материалов.
6. **Нет нормальных карт.** Детализация поверхностей только за счёт геометрии.
7. **Текстовый UI.** Вместо полноценного графического интерфейса — информация в заголовке окна.

### Возможные артефакты

- На Gouraud-шейдинге при низкой детализации модели могут появляться угловатые блики.
- При сильном затухании свет может становиться почти невидимым на больших расстояниях.
- Toon-шейдинг может давать резкие границы на гладких поверхностях.

### Идеи для расширения

1. **Карты теней** — реализация shadow mapping для точечного источника (omnidirectional shadow maps).
2. **Множественные источники** — поддержка нескольких точечных, направленных и прожекторных источников.
3. **PBR** — переход на физически корректную модель (GGX, Fresnel, BRDF).
4. **Normal mapping** — загрузка карт нормалей для повышения детализации.
5. **HDR и Tonemapping** — расширение динамического диапазона рендеринга.
6. **Post-processing** — bloom, SSAO, размытие.
7. **ImGui** — графический интерфейс вместо текстового.
8. **Анимация** — вращение объектов по времени, анимированный свет.
9. **Skybox** — кубическая карта окружения.
10. **Мультисемплинг (MSAA)** — сглаживание краёв.
