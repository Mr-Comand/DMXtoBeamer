package shaders

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var RepeatRlShader rl.Shader

type RepeatShader struct {
	Strength float32 `parameter:"Strength,default=10"`
}

func (k *RepeatShader) Load() {
	RepeatRlShader = rl.LoadShader("", "shaders/curved_repeat.fs") // The fragment shader we created earlier

}
func (k *RepeatShader) StartShader(renderTexture rl.RenderTexture2D) {
	rl.BeginShaderMode(RepeatRlShader)
	rl.SetShaderValue(RepeatRlShader, rl.GetShaderLocation(RepeatRlShader, "strength"), []float32{k.Strength}, rl.ShaderUniformFloat)
	rl.SetShaderValueTexture(RepeatRlShader, rl.GetShaderLocation(RepeatRlShader, "texture0"), renderTexture.Texture)
}
func (k *RepeatShader) EndShader() {
	rl.EndShaderMode()
}
func (k *RepeatShader) Setup(parameters map[string]interface{}) {
	Parse(parameters, k)
}
