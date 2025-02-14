package preset_animation

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"technikflg.com/dmxToProjector/animations/animation_helpers"
)

type AnimationParameters map[string]interface{}
type AnimationInterface interface {
	Render(audioData *[]float64, dt float64)
	Unload()
	Reset()
	Configure(config AnimationParameters)
}

type Animation struct {
	FrameCount float32
}

func (a *Animation) Render(audioData *[]float64) {
}
func (a *Animation) Reset() {
	a.FrameCount = 0
}
func (a *Animation) Unload() {
}
func (a *Animation) Configure(config AnimationParameters) {
}
func Parse[T any](params AnimationParameters, out *T) error {
	if out == nil {
		return errors.New("output struct cannot be nil")
	}

	outValue := reflect.ValueOf(out).Elem()
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

		value, exists := params[key]
		if !exists && defaultValue != "" {
			// Use default value if not provided in params
			var err error
			value, err = convertDefaultValue(defaultValue, field.Type)
			if err != nil {
				fmt.Printf("failed to convert default value for field %s: %v", field.Name, err)
				continue
				// return fmt.Errorf("failed to convert default value for field %s: %v", field.Name, err)
			}
		}

		if exists || defaultValue != "" {
			val := reflect.ValueOf(value)
			if fieldValue.CanSet() {
				if val.Type().ConvertibleTo(fieldValue.Type()) {
					fieldValue.Set(val.Convert(fieldValue.Type()))
				} else if field.Type == reflect.TypeOf(animation_helpers.Color{}) {
					c, err := animation_helpers.FromString(value.(string))
					if err != nil {
						return fmt.Errorf("invalid color format for field %s", field.Name)
					}
					fieldValue.Set(reflect.ValueOf(c))
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
	case reflect.Struct:
		if targetType == reflect.TypeOf(animation_helpers.Color{}) { // Correctly compare with Color{}
			return animation_helpers.FromString(defaultValue)
		}
		return nil, fmt.Errorf("unsupported default type: %s", targetType.Kind())
	default:
		return nil, fmt.Errorf("unsupported default type: %s", targetType.Kind())
	}
}
