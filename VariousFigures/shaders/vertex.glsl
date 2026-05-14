#version 410 core

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out vec3 ourColor;
out vec3 vPosition; // передаем позицию во фрагментный шейдер для паттернов

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
    gl_Position = projection * view * model * vec4(aPos, 1.0);
    
    vPosition = aPos;   // локальные координаты!
    
    ourColor = aColor;
}