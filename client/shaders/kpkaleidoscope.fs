#version 330

in vec2 fragTexCoord;    // Fragment's texture coordinate (screen space)
out vec4 finalColor;     // Output color
in vec4 fragColor;
uniform vec2 resolution; // Resolution of the screen
uniform float segments;  // Number of mirrored segments (e.g. 6 for a hexagon)

void main() {
    // Center the coordinates around the middle of the screen
    vec2 centeredCoord = fragTexCoord * resolution - resolution / 2.0;

    // Convert the coordinates to polar (angle and radius)
    float angle = atan(centeredCoord.y, centeredCoord.x);
    float radius = length(centeredCoord);

    // Map the angle to one of the segments to create the kaleidoscope effect
    float segmentAngle = 3.14159 / segments; // Angle for each segment
    angle = mod(angle, segmentAngle); // Wrap the angle into the segment
    angle += floor(abs(angle) / segmentAngle) * segmentAngle;

    // Convert back to Cartesian coordinates
    vec2 finalCoord = vec2(cos(angle), sin(angle)) * radius;

    // Map the final coordinates back to texture space
    finalCoord += resolution / 2.0;
    finalCoord /= resolution;

    // Check if the point is within the circle radius
    float circleRadius = 0.2; // Circle radius in normalized space (0 to 1)
    if (length(finalCoord - vec2(0.5, 0.5)) < circleRadius) {
        finalColor = fragColor; // Red color for the circle
    } else {
        finalColor = vec4(0.0, 0.0, 0.0, 0.0); // Transparent outside the circle
    }
}
