package snapx

import (
	"fmt"
	"io/fs"

	"go.yaml.in/yaml/v3"

	"github.com/auvn/go-snap/internal/errlog"
)

func loadYamlFile[T any](
	fs fs.FS,
	f string,
	dest *T,
) error {
	cfgFile, err := fs.Open(f)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	defer func() {
		errlog.SwallowError(cfgFile.Close(), "close config file")
	}()

	dec := yaml.NewDecoder(cfgFile)
	if err := dec.Decode(dest); err != nil {
		return fmt.Errorf("decode yaml: %w", err)
	}
	return nil
}
