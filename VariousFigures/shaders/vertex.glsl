#version 410 core

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;
layout (location = 2) in vec2 aTexCoord;

out vec3 ourColor;
out vec2 TexCoord;
out vec3 vPosition; // передаем позицию во фрагментный шейдер для паттернов

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
    gl_Position = projection * view * model * vec4(aPos, 1.0);
    
    // Сохраняем позицию в мировых координатах для паттернов
    vPosition = vec3(model * vec4(aPos, 1.0));
    
    ourColor = aColor;
    TexCoord = aTexCoord;
}