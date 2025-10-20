package snapx

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/auvn/go-snap/internal/errlog"
	"github.com/mitchellh/mapstructure"
)

type options struct {
	EnvPrefix       string
	TagName         string
	DisableEnv      bool
	ErrorUnused     bool
	RawConfigFormat RawConfigFormat
	RawConfigReader RawConfigReader
}

func (o *options) withDefaults() {
	if o.TagName == "" {
		o.TagName = "snapx"
	}

	if o.RawConfigReader == nil {
		o.RawConfigReader = func() (io.Reader, func(), error) {
			return nil, nil, errors.New("no config file specified")
		}
	}

	if o.RawConfigFormat <= 0 {
		o.RawConfigFormat = RawConfigFormatYAML
	}
}

type Option func(*options)

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
		o.RawConfigReader = func() (io.Reader, func(), error) {
			f, err := os.Open(f)
			if err != nil {
				return nil, nil, err
			}

			return f, func() { errlog.SwallowError(f.Close(), "close file") }, nil
		}
	}
}

func WithRawConfigFormat(format RawConfigFormat) Option {
	return func(o *options) {
		o.RawConfigFormat = format
	}
}

func WithRawConfigReader(r RawConfigReader) Option {
	return func(o *options) {
		o.RawConfigReader = r
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

	if err := loadRawConfig(
		opts.RawConfigFormat,
		opts.RawConfigReader,
		&raw,
	); err != nil {
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
	opts ...Option,
) {
	if err := Load(dest, opts...); err != nil {
		panic(err)
	}
}
