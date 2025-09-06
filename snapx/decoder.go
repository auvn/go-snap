package snapx

import (
	"encoding"
	"fmt"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

type DecoderConfigOption func(*mapstructure.DecoderConfig)

func Decode[T any](
	input any,
	dest *T,
	oo ...DecoderConfigOption,
) error {
	cfg := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			decodeHook,
		),
		WeaklyTypedInput: true,
		ErrorUnused:      true,
	}

	for _, o := range oo {
		o(&cfg)
	}

	cfg.Result = dest

	dec, err := mapstructure.NewDecoder(&cfg)
	if err != nil {
		return fmt.Errorf("create decoder: %w", err)
	}

	if err := dec.Decode(input); err != nil {
		return fmt.Errorf("decode input: %w", err)
	}

	return nil
}

type ValueDecoder interface {
	Decode(data any) error
}

func decodeHook(from reflect.Type, to reflect.Type, data any) (any, error) {
	var (
		typeOfTimeDuration = reflect.TypeOf(time.Duration(0))
		typeOfTime         = reflect.TypeOf(time.Time{})
	)

	if from.Kind() == reflect.String {
		switch to {
		case typeOfTimeDuration:
			converted, err := time.ParseDuration(data.(string))
			if err != nil {
				return nil, fmt.Errorf("parse duration: %w", err)
			}

			return converted, nil
		case typeOfTime:
			tm, err := time.Parse(time.RFC3339, data.(string))
			if err != nil {
				return nil, fmt.Errorf("parse time: %w", err)
			}
			return tm, nil
		}
	}

	toValPtr := reflect.New(to).Elem().Addr()
	if v, ok := toValPtr.Interface().(ValueDecoder); ok {
		if err := v.Decode(data); err != nil {
			return nil, fmt.Errorf("decode value: %w", err)
		}
		return v, nil
	}

	// byte compatible unmarshallers

	var b []byte
	switch v := data.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return data, nil
	}

	switch v := toValPtr.Interface().(type) {
	case encoding.TextUnmarshaler:
		if err := v.UnmarshalText(b); err != nil {
			return nil, fmt.Errorf("unmarshal text: %w", err)
		}
		return v, nil
	case encoding.BinaryUnmarshaler:
		if err := v.UnmarshalBinary(b); err != nil {
			return nil, fmt.Errorf("unmarshal binary: %w", err)
		}
		return v, nil
	}

	return data, nil
}
