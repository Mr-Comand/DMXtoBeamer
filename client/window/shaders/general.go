package shaders

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var GeneralRlShader rl.Shader

type GeneralShader struct {
	HueShift float32 `parameter:"HueShift,default=0"`
	Dimmer   float32 `parameter:"Dimmer,default=1"`
}

func (k *GeneralShader) Load() {
	GeneralRlShader = rl.LoadShader("", "shaders/general.fs") // The fragment shader we created earlier

}
func (k *GeneralShader) StartShader(renderTexture rl.RenderTexture2D) {
	rl.BeginShaderMode(GeneralRlShader)
	rl.SetShaderValue(GeneralRlShader, rl.GetShaderLocation(GeneralRlShader, "hueShift"), []float32{k.HueShift}, rl.ShaderUniformFloat)
	rl.SetShaderValue(GeneralRlShader, rl.GetShaderLocation(GeneralRlShader, "dimmer"), []float32{k.Dimmer}, rl.ShaderUniformFloat)
	rl.SetShaderValueTexture(GeneralRlShader, rl.GetShaderLocation(GeneralRlShader, "texture0"), renderTexture.Texture)

}
func (k *GeneralShader) EndShader() {
	rl.EndShaderMode()

}
func (k *GeneralShader) Setup(parameters map[string]interface{}) {
	Parse(parameters, k)
}
