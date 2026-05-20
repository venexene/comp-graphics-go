#version 330 core

// Vertex attributes (see objects/model.go for layout)
layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal_in;
layout (location = 2) in vec2 uv_coords_in;
layout (location = 3) in vec3 tangent_in;

uniform struct Transform {
	mat4 model;
	mat4 view;
	mat4 projection;
	mat3 normal_mat;
	vec3 view_pos;
} transform;

uniform struct PointLight {
	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
    vec3 position;

    float constant;
    float linear;
    float quadratic;

	float ambient_strength;
} light;

// Outputs to fragment shader — everything in tangent space for normal mapping
out TangentSpace {
	vec3 light_dir;   // direction to light in tangent space
	vec3 view_dir;    // direction to viewer in tangent space
	vec2 uv_coords;
	float distance;   // distance to light (for attenuation)
} vert;

void main() {
	// World-space position of this vertex
	vec4 world_pos = transform.model * vec4(position, 1.0);

	// Transform normal and tangent to world space, then build TBN
	vec3 N = normalize(transform.normal_mat * normal_in);
	vec3 T = normalize(transform.normal_mat * tangent_in);

	// Gram–Schmidt re-orthogonalisation: T' = T - N * dot(N, T)
	T = normalize(T - dot(T, N) * N);
	vec3 B = cross(N, T);

	// TBN matrix: transforms from tangent space to world space.
	// Its transpose (which equals its inverse for orthonormal matrices)
	// transforms from world space to tangent space.
	mat3 TBN = transpose(mat3(T, B, N));

	// Directions in world space
	vec3 world_light_dir = light.position - world_pos.xyz;
	vec3 world_view_dir  = transform.view_pos - world_pos.xyz;

	// Pass directions in tangent space (so the fragment shader can sample
	// the normal map directly)
	vert.light_dir = TBN * world_light_dir;
	vert.view_dir  = TBN * world_view_dir;
	vert.uv_coords = vec2(uv_coords_in.x, 1.0 - uv_coords_in.y);
	vert.distance  = length(world_light_dir);

	gl_Position = transform.projection * transform.view * world_pos;
}
