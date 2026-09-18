#version 330

// Input from vertex shader
in vec2 fragTexCoord;
out vec4 finalColor;

uniform float time;
uniform sampler2D texture0;  // The input texture (previous result or the scene)
uniform float rotatingSpeed=0;
uniform float speed;
uniform float elements;
uniform float baseRotation;
uniform float shape;  // Select animation (0-9)

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
vec4 linear(vec2 uv, float elements){

    vec4 color = texture(texture0, vec2(uv.x, 1.0 - uv.y));              // Use the input color
    float angle = time * rotatingSpeed + baseRotation;
    float hueShift = mod((uv.x*sin(angle)+uv.y*cos(angle)+time*speed)*elements, 1.0);
    return vec4(hueShiftFunc(color.rgb, hueShift), color.a); // Apply hue shift
}
vec4 uniformShift(vec2 uv){
    vec4 color = texture(texture0, vec2(uv.x, 1.0 - uv.y));
    float hueShift = mod(time*speed,1.0);
    return vec4(hueShiftFunc(color.rgb, hueShift), color.a); // Apply hue shift

}
vec4 radial(vec2 uv, float elements){
    vec4 color = texture(texture0, vec2(uv.x, 1.0 - uv.y));

    // Convert (x, y) to centered coordinates (-0.5 to 0.5)
    vec2 uvShift = uv - 0.5;

    // Compute the angle from the center outward
    float angle = atan(uvShift.y, uvShift.x); 

    // Create evenly spaced radial rays using mod()
    float hueShift = mod(angle * elements / 6.283185 + time * speed + baseRotation, 1.0);

    return vec4(hueShiftFunc(color.rgb, hueShift), color.a);

}
vec4 grid(vec2 uv, float elements){
    vec4 color = texture(texture0, vec2(uv.x, 1.0 - uv.y));

    float shift = mod(sin(uv.x*6.24*elements)*sin(uv.y*6.24*elements)+(time*speed), 1.0);
    return vec4(hueShiftFunc(color.rgb, shift), color.a); // Apply hue shift
}
vec4 circle(vec2 uv, float elements){
    vec4 color = texture(texture0, vec2(uv.x, 1.0 - uv.y));
    uv -= 0.5;
    uv *= 2.0;
    float shift = mod(sqrt(uv.x*uv.x+uv.y*uv.y)*elements+time*speed,1.0);
    return vec4(hueShiftFunc(color.rgb, shift), color.a); // Apply hue shift
}
void main() {
    vec2 uv = fragTexCoord;
    // vec4 finalColor = vec4(0.0);
int mode = int(shape + 0.5);  // Convert float to nearest int

    switch(mode) {
        case 0: finalColor = uniformShift(uv); break;
        case 1: finalColor = linear(uv, elements); break;
        case 2: finalColor = radial(uv, elements); break;
        case 3: finalColor = grid(uv, elements); break;
        case 4: finalColor = circle(uv, elements); break;
        default: finalColor = linear(uv, elements); break;  // Fallback if invalid
    }
}