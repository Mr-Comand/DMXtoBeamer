#version 330

// Input from the vertex shader
in vec2 fragTexCoord;   // Texture coordinates for the fragment
in vec4 fragColor;      // Color from vertex shader (not used, but could be)

// Output color of the fragment
out vec4 finalColor;

// Uniforms (input values)
uniform sampler2D texture0; // Main texture of the image

void main() {
    finalColor = texture(texture0, fragTexCoord); // Sample the texture at the given UV coordinates
}
