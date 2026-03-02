#version 410 core
out vec4 FragColor;

in vec3 vPosition;

void main()
{
    float k = 10.0;
    float sum = floor(vPosition.y * k);
    
    if (mod(sum, 2) == 0) {
        FragColor = vec4(0.3, 0.8, 0.3, 1.0);
    }
    else {
        FragColor = vec4(0.0, 0.4, 0.0, 1.0);
    }
}