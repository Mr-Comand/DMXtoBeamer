#version 330

in vec2 fragTexCoord;
out vec4 finalColor;

uniform sampler2D texture0;  
uniform float time;

// Default values
uniform float segments  = 8.0;     // Number of reflections
uniform float levels    = 4.0;     // Recursive depth
uniform float rotation  = 0.0;     // Base rotation
uniform float zoom      = 1.2;     // Adjusted zoom
uniform float distortion = 0.02;   // Reduced distortion to keep visibility
uniform float fade      = 0.05;    // Softer fade effect only at edges
uniform float speed     = 1.0;     // Speed of animation
vec3 getPos(vec2 uv, float rotation){
    float angle = atan(uv.y, uv.x) + rotation + time * 0.1 * speed;
    float radius = length(uv) * zoom;
    
    // Calculate alpha: start at 1.0, which will be faded near edges
    float alpha = 1.0;

    // Apply recursive kaleidoscope effect
    for (float i = 1.0; i <= levels; i++) {
        float segmentAngle = 3.14159 / (segments * i);
        
        // Add subtle distortion to segment boundaries
        float distortionOffset = sin(angle * segments * 2.0 + time * speed) * distortion / i;
        angle = mod(angle + distortionOffset, segmentAngle * 2.0) - segmentAngle;
        
        radius = mix(radius, radius * 0.9, 0.5); // Prevent excessive shrinking
        
        // Apply the alpha fade only near the edges (not the center)
        float edgeFactor = (cos(angle * segments))*0.5+0.5;
        alpha *= (pow(edgeFactor, 1.0/fade));
       
    }

    // Convert back to Cartesian coordinates
    return vec3(cos(angle+time * speed) * radius + 0.5, sin(angle+time * speed) * radius + 0.5, alpha);
}

void main() {
    vec2 uv = fragTexCoord; 
    uv.y = 1 - uv.y;
    uv -= 0.5;
    vec3 kaleidoUV = getPos(uv, rotation);
    vec3 kaleido2UV = getPos(uv, rotation+3.14/segments);
    vec4 texColor = texture(texture0, kaleidoUV.xy)*kaleidoUV.z+texture(texture0, kaleido2UV.xy)*kaleido2UV.z;
    
    // Keep a strong center and smooth fade on the edges
    finalColor = vec4(texColor.rgb, texColor.a);
}
