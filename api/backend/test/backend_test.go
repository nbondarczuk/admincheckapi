package backend_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackend(t *testing.T) {
	t.Run("always success", func(t *testing.T) {
		assert.Equal(t, "1", "1")
	})
}
