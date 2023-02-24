package e

import "fmt"

// Wrap
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
