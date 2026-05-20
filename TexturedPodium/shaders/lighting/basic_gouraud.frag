#version 330 core

in vec3 vert_color;
in vec2 uv_coords;

uniform sampler2D u_materialTexture;
uniform sampler2D u_numberTexture;
uniform vec3 u_cubeColor;
uniform float u_materialWeight;
uniform float u_numberWeight;

out vec4 frag_color;

void main() {
    // Multi-texture blending per fragment (standard Gouraud:
    // lighting computed per-vertex, texture applied per-fragment)
    vec3 matColor = texture(u_materialTexture, uv_coords).rgb;
    vec3 numColor = texture(u_numberTexture, uv_coords).rgb;
    float totalWeight = u_materialWeight + u_numberWeight;
    vec3 surfaceColor = u_cubeColor * (u_materialWeight * matColor + u_numberWeight * numColor) / max(totalWeight, 0.001);
    frag_color = vec4(surfaceColor * vert_color, 1.0);
}