package processor

import (
	"context"
	"errors"
)

// Process contains service-private business logic for latest-state-writer.
func Process(ctx context.Context, payload map[string]any) error {
	if ctx == nil {
		return errors.New("context is required")
	}
	if payload == nil {
		return errors.New("payload is required")
	}

	return nil
}
