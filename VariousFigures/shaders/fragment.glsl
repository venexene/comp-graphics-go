#version 410 core
out vec4 FragColor; // Итоговый цвет фргамента или пикселя

in vec3 ourColor; // Интерполированный цвет из фрагментного шейдера
in vec2 TexCoord; // Интерполированные текстурные координаты

uniform sampler2D ourTexture; // Текстура
uniform bool useTexture; // Использовать текстуру или цвет

void main()
{
    if (useTexture) {
        FragColor = texture(ourTexture, TexCoord); // Определение итогового цвета пикселя на основе текстуры
    } else {
        // Иначе используем цвет вершины
        FragColor = vec4(ourColor, 1.0); // Определение итогового цвета пикселя с добавлением альфа-канала
    }
}