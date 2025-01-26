#version 330

// Input from vertex shader
in vec2 fragTexCoord;
out vec4 finalColor;

uniform sampler2D texture0;  // The input texture (previous result or the scene)
uniform float hueShift; // Hue shift value (-1.0 to 1.0)

// Function to rotate hue
vec3 hueShiftFunc(vec3 color, float shift) {
    float angle = shift * 6.283185; // Convert shift to radians (2 * PI)
    float c = cos(angle);
    float s = sin(angle);
    mat3 hueRotation = mat3(
        vec3(0.299, 0.587, 0.114) + vec3(0.701, -0.587, -0.114) * c + vec3(0.168, -0.329, 1.111) * s,
        vec3(0.299, 0.587, 0.114) + vec3(-0.299, 0.413, -0.114) * c + vec3(0.328, 0.035, -0.292) * s,
        vec3(0.299, 0.587, 0.114) + vec3(-0.3, -0.588, 0.886) * c + vec3(-1.079, 1.057, 0.021) * s
    );
    return clamp(color * hueRotation, 0.0, 1.0);
}

void main() {
    vec4 color = texture(texture0, fragTexCoord);              // Use the input color
    finalColor = vec4(hueShiftFunc(color.rgb, hueShift), color.a); // Apply hue shift
}
