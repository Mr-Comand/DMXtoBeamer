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
		minRaw := ""
		maxRaw := ""
		// Check if a default value is provided
		for _, part := range parts[1:] {
			if strings.HasPrefix(part, "default=") {
				defaultValue = strings.TrimPrefix(part, "default=")
			}
			if strings.HasPrefix(part, "min=") {
				minRaw = strings.TrimPrefix(part, "min=")
			}
			if strings.HasPrefix(part, "max=") {
				maxRaw = strings.TrimPrefix(part, "max=")
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
			}
		}
		if exists || defaultValue != "" {

			val := reflect.ValueOf(value)

			if val.Type().ConvertibleTo(fieldValue.Type()) {
				val = val.Convert(fieldValue.Type())
			} else if field.Type == reflect.TypeOf(animation_helpers.Color{}) {
				c, err := animation_helpers.FromString(value.(string))
				if err != nil {
					return fmt.Errorf("invalid color format for field %s", field.Name)
				}
				val = reflect.ValueOf(c)
			} else {
				fmt.Printf("Type mismatch for field %s: expected %s but got %s\n", field.Name, field.Type, val.Type())
				return fmt.Errorf("type mismatch for field %s", field.Name)
			}

			// Apply min/max after conversion
			if minRaw != "" || maxRaw != "" {
				val = applyMinMax(val, field.Type, minRaw, maxRaw)
			}
			// Assign final value
			if fieldValue.CanSet() {
				fieldValue.Set(val)
			} else {
				fmt.Printf("Cannot set field %s\n", field.Name)
			}

		} else {
			fmt.Printf("Parameter %s not found and no default value provided\n", key)
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

// applyMinMax clamps values to min/max. Missing min/max are ignored.
func applyMinMax(value reflect.Value, t reflect.Type, minRaw, maxRaw string) reflect.Value {
	kind := t.Kind()

	hasMin := minRaw != ""
	hasMax := maxRaw != ""

	// Parse min/max only if present
	parseInt := func(raw string, bits int) int64 {
		if raw == "" {
			return 0
		}
		v, _ := strconv.ParseInt(raw, 10, bits)
		return v
	}

	parseUint := func(raw string, bits int) uint64 {
		if raw == "" {
			return 0
		}
		v, _ := strconv.ParseUint(raw, 10, bits)
		return v
	}

	parseFloat := func(raw string, bits int) float64 {
		if raw == "" {
			return 0
		}
		v, _ := strconv.ParseFloat(raw, bits)
		return v
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := value.Int()
		minV := parseInt(minRaw, t.Bits())
		maxV := parseInt(maxRaw, t.Bits())

		if hasMin && v < minV {
			fmt.Printf("clamped %s: %d -> min %d\n", t.Name(), v, minV)
			v = minV
		}
		if hasMax && v > maxV {
			fmt.Printf("clamped %s: %d -> max %d\n", t.Name(), v, maxV)
			v = maxV
		}
		return reflect.ValueOf(v).Convert(t)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := value.Uint()
		minV := parseUint(minRaw, t.Bits())
		maxV := parseUint(maxRaw, t.Bits())

		if hasMin && v < minV {
			fmt.Printf("clamped %s: %d -> min %d\n", t.Name(), v, minV)
			v = minV
		}
		if hasMax && v > maxV {
			fmt.Printf("clamped %s: %d -> max %d\n", t.Name(), v, maxV)
			v = maxV
		}
		return reflect.ValueOf(v).Convert(t)

	case reflect.Float32, reflect.Float64:
		v := value.Float()
		minV := parseFloat(minRaw, t.Bits())
		maxV := parseFloat(maxRaw, t.Bits())

		if hasMin && v < minV {
			fmt.Printf("clamped %s: %f -> min %f\n", t.Name(), v, minV)
			v = minV
		}
		if hasMax && v > maxV {
			fmt.Printf("clamped %s: %f -> max %f\n", t.Name(), v, maxV)
			v = maxV
		}
		return reflect.ValueOf(v).Convert(t)

	default:
		return value
	}
}
