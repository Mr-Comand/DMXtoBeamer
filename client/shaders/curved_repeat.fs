    #version 330

in vec2 fragTexCoord;
out vec4 finalColor;

uniform sampler2D texture0;
uniform float time;
uniform float strength=10; // Strength of distortion

void main() {
    vec2 uv = fragTexCoord - 0.5; // Center the effect
    float dist = length(uv) * strength;
    vec2 warpedUV = uv + uv * dist; // Radial distortion
    finalColor = texture(texture0, warpedUV + 0.5);
}
