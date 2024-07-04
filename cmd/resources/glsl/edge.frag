#version 440 core

uniform vec4 edgeColor;
out vec4  outColor;

void main() {
    outColor = edgeColor;

    if (outColor.a < 1e-6) {
        discard;
    }
}