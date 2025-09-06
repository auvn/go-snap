package snapx

import (
	"fmt"
	"io"

	"go.yaml.in/yaml/v3"
)

type RawConfigReader func() (io.Reader, func(), error)

func loadYamlFile[T any](
	reader RawConfigReader,
	dest *T,
) error {
	f, closeF, err := reader()
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	defer closeF()

	dec := yaml.NewDecoder(f)
	if err := dec.Decode(dest); err != nil {
		return fmt.Errorf("decode yaml: %w", err)
	}
	return nil
}
