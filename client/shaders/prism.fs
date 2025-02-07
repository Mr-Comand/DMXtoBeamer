#version 330

in vec2 fragTexCoord;
out vec4 finalColor;

uniform sampler2D texture0;  // The input texture
uniform float time;          // Time for animation
uniform float xSwing;
uniform float ySwing;
uniform float phaseShift;
void main() {
    vec2 uv = fragTexCoord;
uv.y = 1.0-uv.y;
    // Chromatic aberration offsets
    float xSwingShiftAmount = xSwing * sin(time); // Time-based movement of the prism effect
    float YSwingShiftAmount = ySwing * sin(time+(phaseShift*6.28)); // Time-based movement of the prism effect

    // Separate the texture sampling for RGB channels
    float r = texture(texture0, uv + vec2( xSwingShiftAmount, YSwingShiftAmount)).r;
    float g = texture(texture0, uv).g;
    float b = texture(texture0, uv - vec2( xSwingShiftAmount, YSwingShiftAmount)).b;

    // Combine the colors to form a prism effect
    finalColor = vec4(r, g, b, 1.0);
}
