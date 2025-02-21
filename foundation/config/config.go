package config

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

func Init(path string) {
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		log.Fatal(fmt.Errorf("error loading config: %w", err))
	}
}

func Load[T any](path string) (T, error) {
	var out T
	err := k.UnmarshalWithConf(path, &out, koanf.UnmarshalConf{Tag: "yaml"})
	if err != nil {
		return out, fmt.Errorf("error unmarshaling config: %w", err)
	}

	if err := validator.New().Struct(out); err != nil {
		return out, fmt.Errorf("error validating config: %w", err)
	}

	return out, nil
}
