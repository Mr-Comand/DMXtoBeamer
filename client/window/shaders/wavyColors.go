package shaders

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var WavyColors rl.Shader

type WavyColorsShader struct {
	DistortionAmount float32 `parameter:"DistortionAmount,default=0.1"`
	RotatingSpeed    float32 `parameter:"RotatingSpeed,default=1"`
	BaseRotation     float32 `parameter:"baseRotation,default=0"`
}

func (k *WavyColorsShader) Load() {
	WavyColors = rl.LoadShader("", "shaders/wavyColors.fs") // The fragment shader we created earlier

}
func (k *WavyColorsShader) StartShader(renderTexture rl.RenderTexture2D) {
	rl.BeginShaderMode(WavyColors)
	rl.SetShaderValue(WavyColors, rl.GetShaderLocation(WavyColors, "time"), []float32{float32(rl.GetTime())}, rl.ShaderUniformFloat)
	rl.SetShaderValue(WavyColors, rl.GetShaderLocation(WavyColors, "distortionAmount"), []float32{k.DistortionAmount / 10}, rl.ShaderUniformFloat)
	rl.SetShaderValue(WavyColors, rl.GetShaderLocation(WavyColors, "rotatingSpeed"), []float32{k.RotatingSpeed * 6.28}, rl.ShaderUniformFloat)
	rl.SetShaderValue(WavyColors, rl.GetShaderLocation(WavyColors, "baseRotation"), []float32{k.BaseRotation * 6.28}, rl.ShaderUniformFloat)
	rl.SetShaderValueTexture(WavyColors, rl.GetShaderLocation(WavyColors, "texture0"), renderTexture.Texture)
}
func (k *WavyColorsShader) EndShader() {
	rl.EndShaderMode()
}
func (k *WavyColorsShader) Setup(parameters map[string]interface{}) {
	Parse(parameters, k)
}
