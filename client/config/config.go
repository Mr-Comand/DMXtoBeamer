package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Artnet struct {
		Interface string `yaml:"interface" default:"10.8.8.1/24"`
	} `yaml:"Artnet"`

	Window struct {
		MaxFPS int32 `yaml:"max_fps" default:"100"`
		Vsync  bool  `yaml:"vsync" default:"true"`
	} `yaml:"window"`
}

// Global config variable
var Cfg *Config

// Load reads the YAML file or creates it with defaults if missing
func Load(path string) *Config {
	cfg := &Config{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// file missing -> apply defaults and create
		applyDefaults(cfg)
		if err := cfg.Save(path); err != nil {
			log.Fatalf("failed to create default config: %v", err)
		}
		fmt.Println("Created default config at", path)
		Cfg = cfg
		return cfg
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatalf("failed to parse YAML: %v", err)
	}

	applyDefaults(cfg)
	Cfg = cfg
	return cfg
}

// Save writes the Config to a YAML file
func (cfg *Config) Save(path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// applyDefaults sets default values from struct tags
func applyDefaults(s interface{}) {
	val := reflect.ValueOf(s).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		if field.Kind() == reflect.Struct {
			applyDefaults(field.Addr().Interface())
			continue
		}

		defaultTag := structField.Tag.Get("default")
		if defaultTag == "" {
			continue
		}

		zero := reflect.Zero(field.Type()).Interface()
		if reflect.DeepEqual(field.Interface(), zero) {
			switch field.Kind() {
			case reflect.String:
				field.SetString(defaultTag)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if v, err := strconv.ParseInt(defaultTag, 10, 64); err == nil {
					field.SetInt(v)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if v, err := strconv.ParseUint(defaultTag, 10, 64); err == nil {
					field.SetUint(v)
				}
			case reflect.Float32, reflect.Float64:
				if v, err := strconv.ParseFloat(defaultTag, 64); err == nil {
					field.SetFloat(v)
				}
			case reflect.Bool:
				if v, err := strconv.ParseBool(defaultTag); err == nil {
					field.SetBool(v)
				}
			}
		}
	}
}
