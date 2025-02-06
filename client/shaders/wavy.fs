#version 330

in vec2 fragTexCoord;
out vec4 finalColor;

uniform sampler2D texture0;  // The input texture (previous result or the scene)
uniform float time;          // Time-based animation for the prism effect
uniform float distortionAmount;
void main() {
    // Applying a prism-like effect using sine and cosine for distortion
    vec2 uv = fragTexCoord;
    
    
    float angle = time * 0.5;  // Angle of distortion increases with time
    
    // Distorting the texture coordinates by applying a sine/cosine function
    uv.x += sin(uv.y * 20.0 + angle) * distortionAmount;
    uv.y += cos(uv.x * 20.0 + angle) * distortionAmount;
    
    // Sampling the texture with the distorted coordinates
    finalColor = texture(texture0, vec2(uv.x, 1.0 - uv.y)); // Invert Y for correct sampling
}
