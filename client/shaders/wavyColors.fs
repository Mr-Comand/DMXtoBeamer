#version 330

in vec2 fragTexCoord;
out vec4 finalColor;

uniform float time;
uniform sampler2D texture0;
uniform float distortionAmount; 
uniform float rotatingSpeed; 
uniform float baseRotation; 

void main() {
    vec2 uv = fragTexCoord;
    uv.y = 1.0-uv.y;
    float angle = time * 0.5 * rotatingSpeed + baseRotation;  // Angle of distortion increases with time

    // Prism distortion effect based on sine waves
    float distortionX = sin(uv.y * 10.0 + time *2.0)*sin(angle) * distortionAmount;
    float distortionY = sin(uv.x * 10.0  + time *2.0) *cos(angle) * distortionAmount;
    
    // Chromatic aberration
    vec4 textureR = texture(texture0, uv + vec2( distortionX, distortionY));
    vec4 textureG = texture(texture0, uv);
    vec4 textureB = texture(texture0, uv - vec2( distortionX, distortionY));

    float r = textureR.r;
    float g = textureG.g;
    float b = textureB.b;
    float a = min((textureR.a + textureG.a + textureB.a),1.0) ; // Averaging alpha values

    
    finalColor = vec4(r, g, b, a);


    
}
