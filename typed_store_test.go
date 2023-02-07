package memkey

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTypedStore_SetWithTTL(t *testing.T) {
	s := &TypedStore[int, bool]{}

	done := make(chan struct{})

	go s.ExpireTTL(time.Millisecond, func(key int, value bool) {
		assert.Equal(t, 2, key)
		assert.Equal(t, true, value)
		done <- struct{}{}
	})

	s.Set(1, true)
	s.SetWithTTL(2, true, time.Millisecond*2)
	s.Set(3, true)

	assert.Equal(t, 3, s.Len())

	select {
	case <-time.After(time.Millisecond * 5):
		assert.FailNow(t, "timeout")
		return
	case <-done:
		assert.Equal(t, 2, s.Len())
	}
}
