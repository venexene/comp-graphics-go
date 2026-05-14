#version 410 core // Версия OpenGL

layout (location = 0) in vec3 aPos; // Позиция вершины
layout (location = 1) in vec3 aColor; // Цвет вершины

// out-переменные вершинного автоматически интерполируются для фрагментного
out vec3 ourColor; // Передача цвета во фрагментный шейдер

// uniform-переменные - глобальные константные переменные
// Одинаковы для всех вершин и фрагментов в данном draw call
// Устанавливаются из CPU перед отрисовкой
uniform mat4 model; // Матрица модели
uniform mat4 view; // Матрица вида 
uniform mat4 projection; // Матрица проекции

void main()
{
    gl_Position = projection * view * model * vec4(aPos, 1.0);
    ourColor = aColor; // Передача цвета дальше в фрагментный шейдер
}