package memkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	s := &Store[int]{}

	t.Run("int_not_found", func(t *testing.T) {
		value, ok := Get[int](s, testKey(t))
		assert.False(t, ok)
		assert.Zero(t, value)
		assert.IsType(t, 0, value)
	})

	t.Run("float_not_found", func(t *testing.T) {
		value, ok := Get[float64](s, testKey(t))
		assert.False(t, ok)
		assert.Zero(t, value)
		assert.IsType(t, 0.0, value)
	})

	t.Run("interface_not_found", func(t *testing.T) {
		value, ok := Get[testInterface](s, testKey(t))
		assert.False(t, ok)
		assert.Zero(t, value)
		assert.IsType(t, testInterface(nil), value)
	})

	t.Run("int_found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		value, ok := Get[int](s, k)
		assert.True(t, ok)
		assert.Equal(t, 1, value)
	})

	t.Run("float_found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1.0)
		value, ok := Get[float64](s, k)
		assert.True(t, ok)
		assert.Equal(t, 1.0, value)
	})

	t.Run("interface_found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, testInterfaceImpl{})
		value, ok := Get[testInterface](s, k)
		assert.True(t, ok)
		assert.Equal(t, testInterfaceImpl{}, value)
	})
}
