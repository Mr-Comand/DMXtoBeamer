package shaders

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Kaleidoscope rl.Shader

type KaleidoscopeShader struct {
	Segments   float32 `parameter:"Segments,default=8.0"`
	Levels     float32 `parameter:"Levels,default=4.0"`
	Rotation   float32 `parameter:"Rotation,default=0.0"`
	Zoom       float32 `parameter:"Zoom,default=1.2"`
	Distortion float32 `parameter:"Distortion,default=0.02"`
	Fade       float32 `parameter:"Fade,default=0.05"`
	Speed      float32 `parameter:"Speed,default=1"`
}

func (k *KaleidoscopeShader) Load() {
	Kaleidoscope = rl.LoadShader("", "shaders/kaleidoscope.fs")

}
func (k *KaleidoscopeShader) StartShader(renderTexture rl.RenderTexture2D) {
	rl.BeginShaderMode(Kaleidoscope)
	rl.SetShaderValue(Kaleidoscope, rl.GetShaderLocation(Kaleidoscope, "time"), []float32{float32(rl.GetTime())}, rl.ShaderUniformFloat)
	rl.SetShaderValueTexture(Kaleidoscope, rl.GetShaderLocation(Kaleidoscope, "texture0"), renderTexture.Texture)

	rl.SetShaderValue(Kaleidoscope, rl.GetShaderLocation(Kaleidoscope, "segments"), []float32{k.Segments}, rl.ShaderUniformFloat)
	rl.SetShaderValue(Kaleidoscope, rl.GetShaderLocation(Kaleidoscope, "levels"), []float32{k.Levels}, rl.ShaderUniformFloat)
	rl.SetShaderValue(Kaleidoscope, rl.GetShaderLocation(Kaleidoscope, "rotation"), []float32{k.Rotation}, rl.ShaderUniformFloat)
	rl.SetShaderValue(Kaleidoscope, rl.GetShaderLocation(Kaleidoscope, "zoom"), []float32{k.Zoom}, rl.ShaderUniformFloat)
	rl.SetShaderValue(Kaleidoscope, rl.GetShaderLocation(Kaleidoscope, "distortion"), []float32{k.Distortion}, rl.ShaderUniformFloat)
	rl.SetShaderValue(Kaleidoscope, rl.GetShaderLocation(Kaleidoscope, "fade"), []float32{k.Fade}, rl.ShaderUniformFloat)
	rl.SetShaderValue(Kaleidoscope, rl.GetShaderLocation(Kaleidoscope, "speed"), []float32{k.Speed}, rl.ShaderUniformFloat)

}
func (k *KaleidoscopeShader) EndShader() {
	rl.EndShaderMode()
}
func (k *KaleidoscopeShader) Setup(parameters map[string]interface{}) {
	Parse(parameters, k)
}
