package shaders

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
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
	rl.SetShaderValue(Prism, rl.GetShaderLocation(kaleidoscope, "resolution"), []float32{float32(windowWidth), float32(windowHeight)}, rl.ShaderUniformVec2)
	rl.SetShaderValue(kaleidoscope, rl.GetShaderLocation(kaleidoscope, "time"), []float32{float32(rl.GetTime())}, rl.ShaderUniformFloat)
	rl.SetShaderValueTexture(kaleidoscope, rl.GetShaderLocation(kaleidoscope, "texture0"), renderTexture.Texture)

}
func (k *PrismShader) EndShader() {
	rl.EndShaderMode()
}
func (k *PrismShader) Setup(parameters map[string]interface{}){
	// Pass the pointer to the shader to the RepackageShaderParams function
	if err := RepackageShaderParams(k, parameters); err != nil {
		log.Printf("Error repackaging shader parameters: %v", err)
		return
	}
}