package shaders

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Prism rl.Shader

type PrismShader struct {
	XSwing     float32 `parameter:"XSwing,default=0.2"`
	YSwing     float32 `parameter:"YSwing,default=0.2"`
	PhaseShift float32 `parameter:"PhaseShift,default=0.5"`
}

func (s *PrismShader) Load() {
	Prism = rl.LoadShader("", "shaders/prism.fs") // The fragment shader we created earlier

}
func (s *PrismShader) StartShader(renderTexture rl.RenderTexture2D) {
	windowWidth := rl.GetScreenWidth()
	windowHeight := rl.GetScreenHeight()
	rl.BeginShaderMode(Prism)
	rl.SetShaderValue(Prism, rl.GetShaderLocation(Prism, "resolution"), []float32{float32(windowWidth), float32(windowHeight)}, rl.ShaderUniformVec2)
	rl.SetShaderValue(Prism, rl.GetShaderLocation(Prism, "time"), []float32{float32(rl.GetTime())}, rl.ShaderUniformFloat)
	rl.SetShaderValue(Prism, rl.GetShaderLocation(Prism, "xSwing"), []float32{s.XSwing}, rl.ShaderUniformFloat)
	rl.SetShaderValue(Prism, rl.GetShaderLocation(Prism, "ySwing"), []float32{s.YSwing}, rl.ShaderUniformFloat)
	rl.SetShaderValue(Prism, rl.GetShaderLocation(Prism, "phaseShift"), []float32{s.PhaseShift}, rl.ShaderUniformFloat)
	rl.SetShaderValueTexture(Prism, rl.GetShaderLocation(Prism, "texture0"), renderTexture.Texture)

}
func (s *PrismShader) EndShader() {
	rl.EndShaderMode()
}
func (s *PrismShader) Setup(parameters map[string]interface{}) {
	// Pass the pointer to the shader to the RepackageShaderParams function
	Parse(parameters, s)

}
