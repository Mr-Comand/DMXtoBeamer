package shaders

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Wavy rl.Shader

type WavyShader struct {
	DistortionAmount float32 `parameter:"DistortionAmount,default=0.1"`
	RotatingSpeed    float32 `parameter:"RotatingSpeed,default=1"`
	BaseRotation     float32 `parameter:"baseRotation,default=0"`
}

func (k *WavyShader) Load() {
	Wavy = rl.LoadShader("", "shaders/wavy.fs") // The fragment shader we created earlier

}
func (k *WavyShader) StartShader(renderTexture rl.RenderTexture2D) {
	rl.BeginShaderMode(Wavy)
	rl.SetShaderValue(Wavy, rl.GetShaderLocation(Wavy, "time"), []float32{float32(rl.GetTime())}, rl.ShaderUniformFloat)
	rl.SetShaderValue(Wavy, rl.GetShaderLocation(Wavy, "distortionAmount"), []float32{k.DistortionAmount / 10}, rl.ShaderUniformFloat)
	rl.SetShaderValue(Wavy, rl.GetShaderLocation(Wavy, "rotatingSpeed"), []float32{k.RotatingSpeed * 6.28}, rl.ShaderUniformFloat)
	rl.SetShaderValue(Wavy, rl.GetShaderLocation(Wavy, "baseRotation"), []float32{k.BaseRotation * 6.28}, rl.ShaderUniformFloat)
	rl.SetShaderValueTexture(Wavy, rl.GetShaderLocation(Wavy, "texture0"), renderTexture.Texture)

}
func (k *WavyShader) EndShader() {
	rl.EndShaderMode()
}
func (k *WavyShader) Setup(parameters map[string]interface{}) {
	Parse(parameters, k)
}
