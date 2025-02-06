package shaders

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Prism rl.Shader

type PrismShader struct {
}

func (k *PrismShader) Load() {
	Prism = rl.LoadShader("", "shaders/prism.fs") // The fragment shader we created earlier

}
func (k *PrismShader) StartShader(renderTexture rl.RenderTexture2D) {
	windowWidth := rl.GetScreenWidth()
	windowHeight := rl.GetScreenHeight()
	rl.BeginShaderMode(Prism)
	rl.SetShaderValue(Prism, rl.GetShaderLocation(Prism, "resolution"), []float32{float32(windowWidth), float32(windowHeight)}, rl.ShaderUniformVec2)
	rl.SetShaderValue(Prism, rl.GetShaderLocation(Prism, "time"), []float32{float32(rl.GetTime())}, rl.ShaderUniformFloat)
	rl.SetShaderValueTexture(Prism, rl.GetShaderLocation(Prism, "texture0"), renderTexture.Texture)

}
func (k *PrismShader) EndShader() {
	rl.EndShaderMode()
}
func (k *PrismShader) Setup(parameters map[string]interface{}) {
	// Pass the pointer to the shader to the RepackageShaderParams function
	if err := RepackageShaderParams(k, parameters); err != nil {
		log.Printf("Error repackaging shader parameters: %v", err)
		return
	}
}
