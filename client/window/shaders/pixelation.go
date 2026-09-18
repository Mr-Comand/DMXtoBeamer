package shaders

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var PixelationRlShader rl.Shader

type PixelationShader struct {
	Resolution float32 `parameter:"Resolution,default=150"`
}

func (k *PixelationShader) Load() {
	PixelationRlShader = rl.LoadShader("", "shaders/Pixelation.fs") // The fragment shader we created earlier
}
func (k *PixelationShader) StartShader(renderTexture rl.RenderTexture2D) {
	rl.BeginShaderMode(PixelationRlShader)
	rl.SetShaderValue(PixelationRlShader, rl.GetShaderLocation(PixelationRlShader, "pixelSize"), []float32{k.Resolution}, rl.ShaderUniformFloat)
	rl.SetShaderValueTexture(PixelationRlShader, rl.GetShaderLocation(PixelationRlShader, "texture0"), renderTexture.Texture)
}
func (k *PixelationShader) EndShader() {
	rl.EndShaderMode()

}
func (k *PixelationShader) Setup(parameters map[string]interface{}) {
	Parse(parameters, k)
}
