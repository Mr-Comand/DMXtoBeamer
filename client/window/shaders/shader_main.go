package shaders

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"reflect"
	"strconv"
	"strings"

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

func InitShaders(windowWidth, windowHeight int32) {

	ElementShaders = make(map[string]ElementShaderInterface)
	ElementShaders["kaleidoscope"] = &KaleidoscopeShader{Segments: 6}

	for _, v := range ElementShaders {
		v.Load()
	}

	TextureShaders = make(map[string]TextureShaderInterface)
	TextureShaderRenderTexture = make(map[string]rl.RenderTexture2D)
	TextureShaders["general"] = &GeneralShader{}
	TextureShaders["hueShift"] = &HueShiftShader{HueShift: 0.5}
	TextureShaders["wavy"] = &WavyShader{}
	TextureShaders["wavyColors"] = &WavyColorsShader{}
	TextureShaders["spin"] = &SpinShader{}
	TextureShaders["pixelation"] = &PixelationShader{}
	TextureShaders["prism"] = &PrismShader{}
	TextureShaders["repeat"] = &RepeatShader{}
	TextureShaders["rainbowwave"] = &RainbowShader{}
	for name, v := range TextureShaders {
		v.Load()
		TextureShaderRenderTexture[name] = rl.LoadRenderTexture(windowWidth+300, windowHeight+300)
	}

}
func OnWindowResize(windowWidth, windowHeight int32) {
	for name := range TextureShaders {
		rl.UnloadRenderTexture(TextureShaderRenderTexture[name])
		TextureShaderRenderTexture[name] = rl.LoadRenderTexture(windowWidth+300, windowHeight+300)
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
		rl.ClearBackground(color.RGBA{A: 0})
	} else {
		log.Println("Shader not found!")
	}
}
func EndTextureShader(shaderName string) {
	if shader := TextureShaders[shaderName]; shader != nil {
		rl.EndTextureMode()
		shader.StartShader(TextureShaderRenderTexture[shaderName])
		rl.DrawTexture(TextureShaderRenderTexture[shaderName].Texture, -150, -150, rl.White)
		shader.EndShader()
	} else {
		log.Println("Shader not found!")
	}
}
func SetupTextureShader(shaderName string, parameters map[string]interface{}) bool {
	if shader := TextureShaders[shaderName]; shader != nil {
		shader.Setup(parameters)
		return true
	} else {
		log.Println("Shader not found!")
		return false
	}
}
func AnotherTextureShader(shaderName string, shaderName2 string) {
	if shader2 := TextureShaders[shaderName2]; shader2 != nil {
		rl.EndTextureMode()
		rl.BeginTextureMode(TextureShaderRenderTexture[shaderName2])
		if shader := TextureShaders[shaderName]; shader != nil {
			rl.ClearBackground(color.RGBA{A: 0})
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
func Parse[T any](parameters map[string]interface{}, shader *T) error {
	if shader == nil {
		return errors.New("output struct cannot be nil")
	}

	outValue := reflect.ValueOf(shader).Elem()
	outType := outValue.Type()

	for i := 0; i < outType.NumField(); i++ {
		field := outType.Field(i)
		fieldValue := outValue.Field(i)

		// Get the corresponding key in the map
		tag := field.Tag.Get("parameter")
		parts := strings.Split(tag, ",")
		key := parts[0] // Get the field name from the JSON tag
		defaultValue := ""

		// Check if a default value is provided
		for _, part := range parts[1:] {
			if strings.HasPrefix(part, "default=") {
				defaultValue = strings.TrimPrefix(part, "default=")
			}
		}

		if key == "" {
			key = field.Name // Fallback to struct field name if no JSON tag
		}

		value, exists := parameters[key]
		if !exists && defaultValue != "" {
			// Use default value if not provided in params
			var err error
			value, err = convertDefaultValue(defaultValue, field.Type)
			if err != nil {
				return fmt.Errorf("failed to convert default value for field %s: %v", field.Name, err)
			}
		}

		if exists || defaultValue != "" {
			val := reflect.ValueOf(value)

			// Ensure the field is settable
			if fieldValue.CanSet() {
				// Try to convert the value to the correct type
				if val.Type().ConvertibleTo(fieldValue.Type()) {
					fieldValue.Set(val.Convert(fieldValue.Type()))
				} else {
					return fmt.Errorf("type mismatch for field %s", field.Name)
				}
			}
		}
	}

	return nil
}

// Convert default string value to the appropriate type
func convertDefaultValue(defaultValue string, targetType reflect.Type) (interface{}, error) {
	switch targetType.Kind() {
	case reflect.String:
		return defaultValue, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.ParseInt(defaultValue, 10, targetType.Bits())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.ParseUint(defaultValue, 10, targetType.Bits())
	case reflect.Float32, reflect.Float64:
		return strconv.ParseFloat(defaultValue, targetType.Bits())
	case reflect.Bool:
		return strconv.ParseBool(defaultValue)
	default:
		return nil, fmt.Errorf("unsupported default type: %s", targetType.Kind())
	}
}
