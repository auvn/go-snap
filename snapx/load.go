package snapx

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/mitchellh/mapstructure"
	_ "github.com/spf13/viper"
)

type options struct {
	ConfigName  string
	EnvPrefix   string
	TagName     string
	FileName    string
	DisableEnv  bool
	ErrorUnused bool
	FS          func() fs.FS
}

func (o *options) withDefaults() {
	if o.ConfigName == "" {
		o.ConfigName = "config.yaml"
	}

	if o.TagName == "" {
		o.TagName = "snapx"
	}

	if o.FS == nil {
		o.FS = func() fs.FS { return os.DirFS(".") }
	}
}

type Option func(*options)

func WithConfigName(cfg string) Option {
	return func(o *options) {
		o.ConfigName = cfg
	}
}

func WithFS(f fs.FS) Option {
	return func(o *options) {
		o.FS = func() fs.FS { return f }
	}
}

func WithEnvPrefix(prefix string) Option {
	return func(o *options) {
		o.EnvPrefix = prefix
	}
}

func DisableEnv() Option {
	return func(o *options) {
		o.DisableEnv = true
	}
}

func WithTagName(tag string) Option {
	return func(o *options) {
		o.TagName = tag
	}
}

func WithFileName(f string) Option {
	return func(o *options) {
		o.FileName = f
	}
}

func Load[T any](
	dest *T,
	oo ...Option,
) error {
	var opts options

	for _, o := range oo {
		o(&opts)
	}

	opts.withDefaults()

	var raw map[string]any
	if err := loadYamlFile(opts.FS(), opts.ConfigName, &raw); err != nil {
		return fmt.Errorf("load raw config: %w", err)
	}

	if err := Decode(
		raw,
		dest,
		func(dc *mapstructure.DecoderConfig) {
			dc.TagName = opts.TagName
		},
	); err != nil {
		return fmt.Errorf("decode raw config: %w", err)
	}

	if !opts.DisableEnv {
		if err := loadEnv(opts.EnvPrefix, dest); err != nil {
			return fmt.Errorf("load env: %w", err)
		}
	}

	if err := validate(dest); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	return nil
}

func MustLoad[T any](
	dest *T,
	cfgFile string,
	opts ...Option,
) {
	if err := Load(dest, opts...); err != nil {
		panic(err)
	}
}
