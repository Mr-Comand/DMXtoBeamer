#version 330

in vec2 fragTexCoord;
out vec4 finalColor;

uniform sampler2D texture0;
uniform float time;
uniform float speed; // Rotation speed

void main() {
    vec2 uv = fragTexCoord; // Centering
    uv.y = 1.0-uv.y;
    uv -= 0.5;
    float angle = time * speed;
    mat2 rotation = mat2(cos(angle), -sin(angle), sin(angle), cos(angle));
    uv = rotation * uv;
    finalColor = texture(texture0, uv + 0.5);
}
