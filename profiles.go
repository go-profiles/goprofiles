package goprofiles

import (
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type WithOptions func(*Options)

type Options struct {
	file     string
	profiles []string
}

func defaultOptions() Options {
	return Options{
		file: "profiles.yaml",
	}
}

func WithFile(file string) WithOptions {
	return func(o *Options) {
		o.file = file
	}
}

func WithProfile(profile ...string) WithOptions {
	return func(o *Options) {
		o.profiles = append(o.profiles, profile...)
	}
}

type Profiles struct {
	Options
	values map[interface{}]interface{}
}

func New(opts ...WithOptions) *Profiles {
	p := &Profiles{
		Options: defaultOptions(),
	}
	for _, opt := range opts {
		opt(&p.Options)
	}

	values, err := getValues(p.Options.file, p.Options.profiles...)
	if err != nil {
		panic(err)
	}

	p.values = values
	return p
}

func (p *Profiles) GetString(query string) string {
	v, _ := p.get(query)
	return v.(string)
}

func (p *Profiles) GetInt(query string) int {
	v, _ := p.get(query)
	return v.(int)
}

func (p *Profiles) GetBool(query string) bool {
	v, _ := p.get(query)
	return v.(bool)
}

func (p *Profiles) GetFloat(query string) float64 {
	v, _ := p.get(query)
	return v.(float64)
}

func (p *Profiles) GetIntSlice(query string) []int {
	_, slice := p.get(query)
	var ints []int
	for _, v := range slice {
		ints = append(ints, v.(int))
	}
	return ints
}

func (p *Profiles) GetStringSlice(query string) []string {
	_, slice := p.get(query)
	var strings []string
	for _, v := range slice {
		strings = append(strings, v.(string))
	}
	return strings
}

func (p *Profiles) GetBoolSlice(query string) []bool {
	_, slice := p.get(query)
	var bools []bool
	for _, v := range slice {
		bools = append(bools, v.(bool))
	}
	return bools
}

func (p *Profiles) GetFloatSlice(query string) []float64 {
	_, slice := p.get(query)
	var floats []float64
	for _, v := range slice {
		floats = append(floats, v.(float64))
	}
	return floats
}

func (p *Profiles) get(query string) (interface{}, []interface{}) {
	keys := strings.Split(query, ".")
	values := p.values

	for _, key := range keys {
		val, ok := values[key]
		if !ok {
			panic("🔍: no value found for key: " + query)
		}

		if nestedMap, isMap := val.(map[interface{}]interface{}); isMap {
			values = nestedMap
		} else {
			if sliceVal, ok := val.([]interface{}); ok {
				return nil, sliceVal
			}
			return val, nil
		}
	}
	panic("🔍: no value found for key: " + query)
}

func getValues(file string, profiles ...string) (map[interface{}]interface{}, error) {
	if file == "" {
		return nil, errors.New("no file specified")
	}

	switch {
	case strings.HasSuffix(file, ".yaml"):
		return getYamlValues(file, profiles...)
	}

	return nil, errors.New("unsupported file type")
}

func getYamlValues(file string, profiles ...string) (map[interface{}]interface{}, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.New("💾: file not found")
	}
	// Define a map to unmarshal the YAML content.
	var ym map[string]interface{}
	// Unmarshal the YAML data into the map.

	err = yaml.Unmarshal(bytes, &ym)
	if err != nil {
		return nil, err
	}

	var values = make(map[interface{}]interface{})
	for _, profile := range profiles {
		pm, ok := ym["goprofiles"].(map[interface{}]interface{})[profile]
		if !ok {
			return nil, errors.New("📝: profile not found")
		}

		for k, v := range pm.(map[interface{}]interface{}) {
			if _, ok := values[k]; ok {
				return nil, errors.New("💥 value conflict found for key: " + k.(string))
			}
			values[k] = v
		}
	}
	return values, nil
}
