package shaders

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var RainbowRlShader rl.Shader

type RainbowShader struct {
	Shape         float32 `parameter:"Shape,default=0"`
	Speed         float32 `parameter:"Speed,default=0.5"`
	RotatingSpeed float32 `parameter:"RotatingSpeed,default=0"`
	Tiles         float32 `parameter:"Tiles,default=5"`
	BaseRotation  float32 `parameter:"BaseRotation,default=0"`
}

func (k *RainbowShader) Load() {
	RainbowRlShader = rl.LoadShader("", "shaders/rainbow.fs") // The fragment shader we created earlier

}
func (k *RainbowShader) StartShader(renderTexture rl.RenderTexture2D) {
	rl.BeginShaderMode(RainbowRlShader)
	rl.SetShaderValue(RainbowRlShader, rl.GetShaderLocation(RainbowRlShader, "time"), []float32{float32(rl.GetTime())}, rl.ShaderUniformFloat)
	rl.SetShaderValue(RainbowRlShader, rl.GetShaderLocation(RainbowRlShader, "speed"), []float32{k.Speed}, rl.ShaderUniformFloat)
	rl.SetShaderValue(RainbowRlShader, rl.GetShaderLocation(RainbowRlShader, "elements"), []float32{k.Tiles}, rl.ShaderUniformFloat)
	rl.SetShaderValue(RainbowRlShader, rl.GetShaderLocation(RainbowRlShader, "baseRotation"), []float32{k.BaseRotation / 360}, rl.ShaderUniformFloat)
	rl.SetShaderValue(RainbowRlShader, rl.GetShaderLocation(RainbowRlShader, "rotatingSpeed"), []float32{k.RotatingSpeed}, rl.ShaderUniformFloat)
	rl.SetShaderValue(RainbowRlShader, rl.GetShaderLocation(RainbowRlShader, "shape"), []float32{k.Shape}, rl.ShaderUniformFloat)

	rl.SetShaderValueTexture(RainbowRlShader, rl.GetShaderLocation(RainbowRlShader, "texture0"), renderTexture.Texture)

}
func (k *RainbowShader) EndShader() {
	rl.EndShaderMode()

}
func (k *RainbowShader) Setup(parameters map[string]interface{}) {
	Parse(parameters, k)
}
