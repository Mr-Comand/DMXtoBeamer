#version 330

in vec2 fragTexCoord;
out vec4 finalColor;

uniform sampler2D texture0;
uniform float pixelSize=200; // Controls pixelation level

void main() {
    vec2 uv = floor(fragTexCoord * pixelSize) / pixelSize; // Pixelate the texture
    uv.y = 1.0-uv.y;
    finalColor = texture(texture0, uv);
}
