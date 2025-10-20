package snapx_test

import (
	"embed"
	_ "embed"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/auvn/go-snap/snapx"
)

//go:embed testdata
var _testdata embed.FS

type testConfig struct {
	Key        string
	AnotherKey string
	Map        any
	MapOfMaps  map[string]map[string]any
}

func TestLoad(t *testing.T) {
	t.Run("yaml", func(t *testing.T) {
		var cfg testConfig
		err := snapx.Load(
			&cfg,
			snapx.WithRawConfigReader(func() (io.Reader, func(), error) {
				f, err := _testdata.Open("testdata/test-config.yaml")
				if err != nil {
					return nil, nil, err
				}

				return f, func() { f.Close() }, nil
			}),
		)
		require.NoError(t, err)

		want := testConfig{
			Key:        "value",
			AnotherKey: "value",
			Map: map[string]any{
				"key": "value",
			},

			MapOfMaps: map[string]map[string]any{
				"super.Puper.Key": {
					"A": "b",
					"C": "d",
				},
			},
		}
		assert.Equal(t, want, cfg)
	})

	t.Run("json", func(t *testing.T) {
		var cfg testConfig
		err := snapx.Load(
			&cfg,
			snapx.WithRawConfigFormat(snapx.RawConfigFormatJSON),
			snapx.WithRawConfigReader(func() (io.Reader, func(), error) {
				f, err := _testdata.Open("testdata/test-config.json")
				if err != nil {
					return nil, nil, err
				}
				return f, func() { f.Close() }, nil
			}),
		)
		require.NoError(t, err)

		want := testConfig{
			Key:        "value",
			AnotherKey: "value",
			Map: map[string]any{
				"key": "value",
			},
			MapOfMaps: map[string]map[string]any{
				"super.Puper.Key": {
					"A": "b",
					"C": "d",
				},
			},
		}
		assert.Equal(t, want, cfg)
	})
}
