#version 410 core
out vec4 FragColor;

in vec3 ourColor; // не используется, но оставляем для совместимости
in vec2 TexCoord; // не используется
in vec3 vPosition; // передаем позицию из вершинного шейдера

uniform bool useTexture; // не используется, но оставляем

void main()
{
    float stripeWidth = 0.1; // ширина полосок
    
    // Используем координату X для создания вертикальных полосок
    // Можно заменить на vPosition.y для горизонтальных полосок
    float stripe = mod(vPosition.x, stripeWidth * 2.0);
    
    if (stripe < stripeWidth) {
        // Ярко-голубой цвет для четных полосок
        FragColor = vec4(0.0, 0.8, 1.0, 1.0); // ярко-голубой (cyan)
    } else {
        // Белый для нечетных полосок
        FragColor = vec4(1.0, 1.0, 1.0, 1.0); // белый
    }
}