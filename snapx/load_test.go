package snapx_test

import (
	"embed"
	_ "embed"
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
	var cfg testConfig
	err := snapx.Load(
		&cfg,
		snapx.WithFS(_testdata),
		snapx.WithConfigName("testdata/test-config.yaml"),
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
}
