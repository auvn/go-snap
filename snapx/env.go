package snapx

import (
	"fmt"

	"github.com/hexdigest/envconfig"
)

func loadEnv[T any](
	prefix string,
	dest *T,
) error {
	if err := envconfig.Process(prefix, dest); err != nil {
		return fmt.Errorf("process env: %w", err)
	}
	return nil
}
