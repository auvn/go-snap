package snapx

import (
	"encoding/json"
	"fmt"
	"io"

	"go.yaml.in/yaml/v3"
)

type RawConfigFormat int

const (
	RawConfigFormatYAML RawConfigFormat = iota + 1
	RawConfigFormatJSON
)

type (
	RawConfigDecoder[T any] func(reader RawConfigReader, dest *T) error
	RawConfigReader         func() (io.Reader, func(), error)
)

func loadRawConfig[T any](
	format RawConfigFormat,
	reader RawConfigReader,
	dest *T,
) error {
	switch format {
	case RawConfigFormatYAML:
		return loadYamlFile(reader, dest)
	case RawConfigFormatJSON:
		return loadJSONFile(reader, dest)
	default:
		return fmt.Errorf("unexpected format: %d", format)
	}
}

func loadYamlFile[T any](
	reader RawConfigReader,
	dest *T,
) error {
	f, closeF, err := reader()
	if err != nil {
		return fmt.Errorf("reader: %w", err)
	}

	defer closeF()

	dec := yaml.NewDecoder(f)
	if err := dec.Decode(dest); err != nil {
		return fmt.Errorf("decode yaml: %w", err)
	}
	return nil
}

func loadJSONFile[T any](
	reader RawConfigReader,
	dest *T,
) error {
	f, closeF, err := reader()
	if err != nil {
		return fmt.Errorf("reader: %w", err)
	}

	defer closeF()

	dec := json.NewDecoder(f)
	if err := dec.Decode(dest); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	return nil
}
