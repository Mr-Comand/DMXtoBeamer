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
    uv.y = 1.0 - uv.y;

    // Chromatic aberration offsets
    float xSwingShiftAmount = xSwing * sin(time); 
    float ySwingShiftAmount = ySwing * sin(time + (phaseShift * 6.28));

    // Separate the texture sampling for RGB channels
    vec4 texColor = texture(texture0, uv);
    vec4 texColorR = texture(texture0, uv + vec2(xSwingShiftAmount, ySwingShiftAmount));
    vec4 texColorB = texture(texture0, uv - vec2(xSwingShiftAmount, ySwingShiftAmount));
    
    float r = texColorR.r;
    float g = texColor.g;
    float b = texColorB.b;
    float a = min((texColor.a + texColorR.a + texColorB.a),1.0) ; // Averaging alpha values

    // Preserve the alpha value while blending shifted positions
    finalColor = vec4(r, g, b, a);
}
