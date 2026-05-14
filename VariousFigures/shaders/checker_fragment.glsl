#version 410 core
out vec4 FragColor;

in vec3 vPosition;

void main()
{
    float k = 5.0;
    int sum = int(vPosition.x * k) + int(vPosition.y * k) + int(vPosition.z * k);
    
    if (mod(sum, 2) == 0) {
        FragColor = vec4(0.8, 0.8, 0, 1.0);
    }
    else {
        FragColor = vec4(0.5, 0.0, 0, 1.0);
    }
}