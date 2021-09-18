package generators

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func ensureDirectoryExistence(dir string) error {
	if err := os.MkdirAll(dir, 0776); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create directory: %s", dir))
	}

	return nil
}
