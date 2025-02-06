package shaders

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var hueShift rl.Shader

type HueShiftShader struct {
	HueShift float32
}

func (k *HueShiftShader) Load() {
	hueShift = rl.LoadShader("", "shaders/hue_shift.fs") // The fragment shader we created earlier

}
func (k *HueShiftShader) StartShader(renderTexture rl.RenderTexture2D) {
	rl.BeginShaderMode(hueShift)
	rl.SetShaderValue(hueShift, rl.GetShaderLocation(hueShift, "hueShift"), []float32{k.HueShift}, rl.ShaderUniformFloat)
	rl.SetShaderValueTexture(hueShift, rl.GetShaderLocation(hueShift, "texture0"), renderTexture.Texture)

}
func (k *HueShiftShader) EndShader() {
	rl.EndShaderMode()

}
func (k *HueShiftShader) Setup(parameters map[string]interface{}) {
	Parse(parameters, k)
}
