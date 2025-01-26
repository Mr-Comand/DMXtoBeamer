#version 330

in vec2 fragTexCoord;
out vec4 finalColor;

uniform sampler2D texture0;  // The input texture (previous result or the scene)
uniform float time;          // Time-based animation for the ripple effect

void main() {
    // Applying a ripple effect using the sine function for distortion
    vec2 uv = fragTexCoord;
    uv.x += sin(uv.y * 10.0 + time * 2.0) * 0.05; // Apply distortion based on Y
    finalColor = texture(texture0, uv); // Sample from the texture with the new distorted coordinates
}
