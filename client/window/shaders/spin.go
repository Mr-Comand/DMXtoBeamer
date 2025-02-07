package shaders

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var SpinRlShader rl.Shader

type SpinShader struct {
	Speed float32 `parameter:"Speed,default=1"`
}

func (k *SpinShader) Load() {
	SpinRlShader = rl.LoadShader("", "shaders/spin.fs") // The fragment shader we created earlier

}
func (k *SpinShader) StartShader(renderTexture rl.RenderTexture2D) {
	rl.BeginShaderMode(SpinRlShader)
	rl.SetShaderValue(SpinRlShader, rl.GetShaderLocation(SpinRlShader, "time"), []float32{float32(rl.GetTime())}, rl.ShaderUniformFloat)
	rl.SetShaderValue(SpinRlShader, rl.GetShaderLocation(SpinRlShader, "speed"), []float32{k.Speed * 6.28}, rl.ShaderUniformFloat)
	rl.SetShaderValueTexture(SpinRlShader, rl.GetShaderLocation(SpinRlShader, "texture0"), renderTexture.Texture)

}
func (k *SpinShader) EndShader() {
	rl.EndShaderMode()

}
func (k *SpinShader) Setup(parameters map[string]interface{}) {
	Parse(parameters, k)
}
