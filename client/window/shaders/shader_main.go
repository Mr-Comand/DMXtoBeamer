package shaders

import (
	"fmt"
	"log"
	"reflect"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var ElementShaders map[string]ElementShaderInterface
var TextureShaders map[string]TextureShaderInterface
var TextureShaderRenderTexture map[string]rl.RenderTexture2D

type ElementShaderInterface interface {
	StartShader()
	EndShader()
	Load()
	Setup(parameters map[string]interface{})
}
type TextureShaderInterface interface {
	StartShader(renderTexture rl.RenderTexture2D)
	EndShader()
	Load()
	Setup(parameters map[string]interface{})
}

func InitShaders() {
	windowWidth := int32(rl.GetScreenWidth())
	windowHeight := int32(rl.GetScreenHeight())

	ElementShaders = make(map[string]ElementShaderInterface)
	ElementShaders["kaleidoscope"] = &KaleidoscopeShader{Segments: 6}

	for _, v := range ElementShaders {
		v.Load()
	}

	TextureShaders = make(map[string]TextureShaderInterface)
	TextureShaderRenderTexture = make(map[string]rl.RenderTexture2D)
	TextureShaders["hueShift"] = &HueShiftShader{HueShift: 0.5}
	TextureShaders["prism"] = &PrismShader{}
	for name, v := range TextureShaders {
		v.Load()
		TextureShaderRenderTexture[name] = rl.LoadRenderTexture(windowWidth, windowHeight)
	}

}
func StartElementShader(shaderName string) {
	if shader := ElementShaders[shaderName]; shader != nil {
		shader.StartShader()
	} else {
		log.Println("Shader not found!")
	}
}
func EndElementShader(shaderName string) {
	if shader := ElementShaders[shaderName]; shader != nil {
		shader.EndShader()
	} else {
		log.Println("Shader not found!")
	}
}
func SetupElementShader(shaderName string, parameters map[string]interface{}) {
	if shader := ElementShaders[shaderName]; shader != nil {
		shader.Setup(parameters)
	} else {
		log.Println("Shader not found!")
	}
}

func StartTextureShader(shaderName string) {
	if shader := TextureShaders[shaderName]; shader != nil {
		rl.BeginTextureMode(TextureShaderRenderTexture[shaderName])
		rl.ClearBackground(rl.Black)
	} else {
		log.Println("Shader not found!")
	}
}
func EndTextureShader(shaderName string) {
	if shader := TextureShaders[shaderName]; shader != nil {
		rl.EndTextureMode()
		rl.ClearBackground(rl.Black)
		shader.StartShader(TextureShaderRenderTexture[shaderName])
		rl.DrawTexture(TextureShaderRenderTexture[shaderName].Texture, 0, 0, rl.White)
		shader.EndShader()
	} else {
		log.Println("Shader not found!")
	}
}
func SetupTextureShader(shaderName string, parameters map[string]interface{}) {
	if shader := TextureShaders[shaderName]; shader != nil {
		shader.Setup(parameters)
		
	} else {
		log.Println("Shader not found!")
	}
}
func AnotherTextureShader(shaderName string, shaderName2 string) {
	if shader2 := TextureShaders[shaderName2]; shader2 != nil {
		rl.EndTextureMode()
		rl.BeginTextureMode(TextureShaderRenderTexture[shaderName2])
		if shader := TextureShaders[shaderName]; shader != nil {
			rl.ClearBackground(rl.Black)
			shader.StartShader(TextureShaderRenderTexture[shaderName])
			rl.DrawTexture(TextureShaderRenderTexture[shaderName].Texture, 0, 0, rl.White)
			shader.EndShader()
		} else {
			log.Println("Shader not found!")
		}
	} else {
		log.Println("Shader2 not found!")
	}

}

func RepackageShaderParams(shader interface{}, parameters map[string]interface{}) error {
	// Ensure that shader is a pointer to a struct or a pointer to a pointer to a struct
	val := reflect.ValueOf(shader)
	if val.Kind() == reflect.Ptr {
		// If it's a pointer to a pointer (e.g. **Shader)
		if val.Elem().Kind() == reflect.Ptr {
			return fmt.Errorf("expected a pointer to a struct, got pointer to pointer")
		}
		// If it's a pointer to a struct
		if val.Elem().Kind() == reflect.Struct {
			val = val.Elem() // Dereference to get the struct value
		} else {
			return fmt.Errorf("expected a pointer to a struct, got %s", val.Elem().Kind())
		}
	} else {
		return fmt.Errorf("expected a pointer to a struct, got %s", val.Kind())
	}

	// Iterate through the parameters map and assign to the corresponding struct fields
	for key, paramValue := range parameters {
		// Check if the field exists and is valid
		field := val.FieldByName(key)
		if !field.IsValid() {
			log.Printf("Field %s does not exist in shader struct", key)
			continue
		}

		// Ensure the field is settable
		if !field.CanSet() {
			log.Printf("Field %s is not settable (likely unexported)", key)
			continue
		}

		// Check the type of the field and set the value accordingly
		switch field.Kind() {
		case reflect.Float32:
			if v, ok := paramValue.(float32); ok {
				field.SetFloat(float64(v))
			} else {
				log.Printf("Type mismatch for field %s: expected float32, got %T", key, paramValue)
				return fmt.Errorf("type mismatch for field %s", key)
			}
		case reflect.Int32:
			if v, ok := paramValue.(int32); ok {
				field.SetInt(int64(v))
			} else {
				log.Printf("Type mismatch for field %s: expected int32, got %T", key, paramValue)
				return fmt.Errorf("type mismatch for field %s", key)
			}
		case reflect.String:
			if v, ok := paramValue.(string); ok {
				field.SetString(v)
			} else {
				log.Printf("Type mismatch for field %s: expected string, got %T", key, paramValue)
				return fmt.Errorf("type mismatch for field %s", key)
			}
		default:
			log.Printf("Unsupported field type for field %s: %s", key, field.Kind())
			return fmt.Errorf("unsupported field type for field %s", key)
		}
	}
	return nil
}