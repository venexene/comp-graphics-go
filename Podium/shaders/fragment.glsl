#version 410 core
out vec4 FragColor;

in vec3 ourColor;
uniform vec3 cubeColor;

void main()
{
    FragColor = vec4(cubeColor, 1.0);
}