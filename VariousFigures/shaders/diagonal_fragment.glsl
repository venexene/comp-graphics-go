#version 410 core
out vec4 FragColor;

in vec3 vPosition;

void main()
{
    float k = 8.0;

    float sum = floor((vPosition.x + vPosition.y) * k);

    if (mod(sum, 2.0) == 0.0)
        FragColor = vec4(0.2, 0.6, 0.8, 1.0);
    else
        FragColor = vec4(0.1, 0.1, 0.4, 1.0);
}