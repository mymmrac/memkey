package memkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_zero(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Zero(t, zero[int]())
	})

	t.Run("string", func(t *testing.T) {
		assert.Zero(t, zero[string]())
	})

	t.Run("interface", func(t *testing.T) {
		assert.Zero(t, zero[testInterface]())
	})
}
