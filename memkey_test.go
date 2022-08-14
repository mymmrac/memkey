package memkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAndSet(t *testing.T) {
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

	t.Run("int_not_found_by_type", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1.0)
		value, ok := Get[int](s, k)
		assert.False(t, ok)
		assert.Zero(t, value)
		assert.IsType(t, 0, value)
	})

	t.Run("redefine", func(t *testing.T) {
		k := testKey(t)

		Set(s, k, 1)
		value, ok := Get[int](s, k)
		assert.True(t, ok)
		assert.Equal(t, 1, value)

		Set(s, k, 2)
		value, ok = Get[int](s, k)
		assert.True(t, ok)
		assert.Equal(t, 2, value)
	})

	t.Run("redefine_type", func(t *testing.T) {
		k := testKey(t)

		Set(s, k, 1)
		valueInt, ok := Get[int](s, k)
		assert.True(t, ok)
		assert.Equal(t, 1, valueInt)

		Set(s, k, 2.0)
		valueFloat, ok := Get[float64](s, k)
		assert.True(t, ok)
		assert.Equal(t, 2.0, valueFloat)
	})
}

func TestGetRaw(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		value, ok := GetRaw(s, k)
		assert.True(t, ok)
		assert.Equal(t, 1, value)
	})

	t.Run("not_found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		value, ok := GetRaw(s, -1)
		assert.False(t, ok)
		assert.Nil(t, value)
	})
}

func TestHas(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.True(t, Has[int](s, k))
	})

	t.Run("not_found_by_type", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.False(t, Has[float64](s, k))
	})

	t.Run("not_found", func(t *testing.T) {
		Set(s, testKey(t), 1)
		assert.False(t, Has[int](s, testKey(t)))
	})
}

func TestHasRaw(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.True(t, HasRaw(s, k))
	})

	t.Run("not_found", func(t *testing.T) {
		Set(s, testKey(t), 1)
		assert.False(t, HasRaw(s, testKey(t)))
	})
}

func TestDelete(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		Delete[int](s, k)
		assert.False(t, Has[int](s, k))
	})

	t.Run("not_found_by_type", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		Delete[float64](s, k)
		assert.True(t, Has[int](s, k))
	})

	t.Run("not_found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		Delete[float64](s, -1)
		assert.True(t, Has[int](s, k))
	})
}

func TestDeleteOk(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.True(t, DeleteOk[int](s, k))
		assert.False(t, Has[int](s, k))
	})

	t.Run("not_found_by_type", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.False(t, DeleteOk[float64](s, k))
		assert.True(t, Has[int](s, k))
	})

	t.Run("not_found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.False(t, DeleteOk[float64](s, -1))
		assert.True(t, Has[int](s, k))
	})
}

func TestDeleteRaw(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		DeleteRaw(s, k)
		assert.False(t, HasRaw(s, k))
	})

	t.Run("not_found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		DeleteRaw(s, -1)
		assert.True(t, HasRaw(s, k))
	})
}

func TestDeleteRawOk(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.True(t, DeleteRawOk(s, k))
		assert.False(t, HasRaw(s, k))
	})

	t.Run("not_found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.False(t, DeleteRawOk(s, -1))
		assert.True(t, HasRaw(s, k))
	})
}

func TestLen(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, 0, Len[int](s))

	Set(s, testKey(t), 1)
	Set(s, testKey(t), 2.0)
	assert.Equal(t, 1, Len[int](s))
	assert.Equal(t, 1, Len[float64](s))
}

func TestLenRaw(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, 0, LenRaw(s))

	Set(s, testKey(t), 1)
	Set(s, testKey(t), 2.0)
	assert.Equal(t, 2, LenRaw(s))
}

func TestKeys(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, []int{}, Keys[int](s))

	k := testKey(t)
	Set(s, k, 1)
	Set(s, testKey(t), 2.0)
	assert.Equal(t, []int{k}, Keys[int](s))
}

func TestKeysRaw(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, []int{}, KeysRaw(s))

	k1 := testKey(t)
	Set(s, k1, 1)
	k2 := testKey(t)
	Set(s, k2, 2.0)
	assert.ElementsMatch(t, []int{k1, k2}, KeysRaw(s))
}

func TestValues(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, []int{}, Values[int](s))

	Set(s, testKey(t), 1)
	Set(s, testKey(t), 2.0)
	assert.Equal(t, []int{1}, Values[int](s))
}

func TestValuesRaw(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, []any{}, ValuesRaw(s))

	Set(s, testKey(t), 1)
	Set(s, testKey(t), 2.0)
	assert.ElementsMatch(t, []any{1, 2.0}, ValuesRaw(s))
}

func TestEntries(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, []Entry[int, int]{}, Entries[int](s))

	k1 := testKey(t)
	Set(s, k1, 1)
	k2 := testKey(t)
	Set(s, k2, 2.0)
	assert.Equal(t, []Entry[int, int]{{k1, 1}}, Entries[int](s))
}

func TestEntriesRaw(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, []Entry[int, any]{}, EntriesRaw(s))

	k1 := testKey(t)
	Set(s, k1, 1)
	k2 := testKey(t)
	Set(s, k2, 2.0)
	assert.ElementsMatch(t, []Entry[int, any]{{k1, 1}, {k2, 2.0}}, EntriesRaw(s))
}

func TestForEach(t *testing.T) {
	s := &Store[int]{}

	count := 0
	ForEach[int](s, func(_ int, _ int) {
		count++
	})
	assert.Equal(t, 0, count)

	Set(s, testKey(t), 1)
	Set(s, testKey(t), 2.0)
	ForEach[int](s, func(_ int, _ int) {
		count++
	})
	assert.Equal(t, 1, count)
}

func TestForEachRaw(t *testing.T) {
	s := &Store[int]{}

	count := 0
	ForEachRaw[int](s, func(_ int, _ any) {
		count++
	})
	assert.Equal(t, 0, count)

	Set(s, testKey(t), 1)
	Set(s, testKey(t), 2.0)
	ForEachRaw(s, func(_ int, _ any) {
		count++
	})
	assert.Equal(t, 2, count)
}
