#version 410 core
out vec4 FragColor;

in vec3 ourColor;
in vec2 TexCoord;

uniform sampler2D ourTexture;
uniform bool useTexture;
uniform vec3 cubeColor;

void main()
{
    if (useTexture) {
        FragColor = texture(ourTexture, TexCoord);
    } else {
        FragColor = vec4(cubeColor, 1.0);
    }
}