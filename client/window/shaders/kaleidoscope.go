package shaders

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
)

var kaleidoscope rl.Shader

type KaleidoscopeShader struct {
	Segments int8
}

func (k *KaleidoscopeShader) Load() {
	kaleidoscope = rl.LoadShader("", "shaders/kaleidoscope.fs")
}
func (k *KaleidoscopeShader) StartShader() {
	windowWidth := rl.GetScreenWidth()
	windowHeight := rl.GetScreenHeight()
	rl.BeginShaderMode(kaleidoscope)
	rl.SetShaderValue(kaleidoscope, rl.GetShaderLocation(kaleidoscope, "resolution"), []float32{float32(windowWidth), float32(windowHeight)}, rl.ShaderUniformVec2)
	rl.SetShaderValue(kaleidoscope, rl.GetShaderLocation(kaleidoscope, "segments"), []float32{1.0}, rl.ShaderUniformFloat) // 6 segments for the kaleidoscope
}
func (k *KaleidoscopeShader) EndShader() {
	rl.EndShaderMode()
}
func (k *KaleidoscopeShader) Setup(parameters map[string]interface{}){
	// Pass the pointer to the shader to the RepackageShaderParams function
	if err := RepackageShaderParams(k, parameters); err != nil {
		log.Printf("Error repackaging shader parameters: %v", err)
		return
	}
}