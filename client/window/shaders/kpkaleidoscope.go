package shaders

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var kpKaleidoscope rl.Shader

type KpKaleidoscopeShader struct {
	Segments int8
}

func (k *KpKaleidoscopeShader) Load() {
	kpKaleidoscope = rl.LoadShader("", "shaders/kpkaleidoscope.fs")
}
func (k *KpKaleidoscopeShader) StartShader() {
	windowWidth := rl.GetScreenWidth()
	windowHeight := rl.GetScreenHeight()
	rl.BeginShaderMode(kpKaleidoscope)
	rl.SetShaderValue(kpKaleidoscope, rl.GetShaderLocation(kpKaleidoscope, "resolution"), []float32{float32(windowWidth), float32(windowHeight)}, rl.ShaderUniformVec2)
	rl.SetShaderValue(kpKaleidoscope, rl.GetShaderLocation(kpKaleidoscope, "segments"), []float32{1.0}, rl.ShaderUniformFloat) // 6 segments for the kaleidoscope
}
func (k *KpKaleidoscopeShader) EndShader() {
	rl.EndShaderMode()
}
func (k *KpKaleidoscopeShader) Setup(parameters map[string]interface{}) {
	Parse(parameters, k)
}
