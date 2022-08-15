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

func TestStore_GetAndSet(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		s.Set(k, 1)
		value, ok := s.Get(k)
		assert.True(t, ok)
		assert.Equal(t, 1, value)
	})

	t.Run("not_found", func(t *testing.T) {
		k := testKey(t)
		s.Set(k, 1)
		value, ok := s.Get(-1)
		assert.False(t, ok)
		assert.Nil(t, value)
	})
}

func TestMustGet(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.Equal(t, 1, MustGet[int](s, k))
	})

	t.Run("not_found", func(t *testing.T) {
		Set(s, testKey(t), 1)
		assert.Equal(t, 0, MustGet[int](s, -1))
	})
}

func TestStore_MustGet(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		s.Set(k, 1)
		assert.Equal(t, 1, s.MustGet(k))
	})

	t.Run("not_found", func(t *testing.T) {
		s.Set(testKey(t), 1)
		assert.Equal(t, nil, s.MustGet(-1))
	})
}

func TestType(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		kind, ok := Type(s, k)
		assert.True(t, ok)
		assert.Equal(t, "int", kind)
	})

	t.Run("not_found", func(t *testing.T) {
		Set(s, testKey(t), 1)
		kind, ok := Type(s, -1)
		assert.False(t, ok)
		assert.Equal(t, "", kind)
	})
}

func TestStore_Type(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		typ, ok := s.Type(k)
		assert.True(t, ok)
		assert.Equal(t, "int", typ)
	})

	t.Run("not_found", func(t *testing.T) {
		Set(s, testKey(t), 1)
		typ, ok := s.Type(-1)
		assert.False(t, ok)
		assert.Equal(t, "", typ)
	})
}

func TestMustType(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.Equal(t, "int", MustType(s, k))
	})

	t.Run("not_found", func(t *testing.T) {
		Set(s, testKey(t), 1)
		assert.Equal(t, "", MustType(s, -1))
	})
}

func TestStore_MustType(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		s.Set(k, 1)
		assert.Equal(t, "int", s.MustType(k))
	})

	t.Run("not_found", func(t *testing.T) {
		s.Set(testKey(t), 1)
		assert.Equal(t, "", s.MustType(-1))
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

func TestStore_Has(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.True(t, s.Has(k))
	})

	t.Run("not_found", func(t *testing.T) {
		Set(s, testKey(t), 1)
		assert.False(t, s.Has(testKey(t)))
	})
}

func TestDeleteOk(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.True(t, Delete[int](s, k))
		assert.False(t, Has[int](s, k))
	})

	t.Run("not_found_by_type", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.False(t, Delete[float64](s, k))
		assert.True(t, Has[int](s, k))
	})

	t.Run("not_found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.False(t, Delete[float64](s, -1))
		assert.True(t, Has[int](s, k))
	})
}

func TestStore_Delete(t *testing.T) {
	s := &Store[int]{}

	t.Run("found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.True(t, s.Delete(k))
		assert.False(t, s.Has(k))
	})

	t.Run("not_found", func(t *testing.T) {
		k := testKey(t)
		Set(s, k, 1)
		assert.False(t, s.Delete(-1))
		assert.True(t, s.Has(k))
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

func TestStore_Len(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, 0, s.Len())

	Set(s, testKey(t), 1)
	Set(s, testKey(t), 2.0)
	assert.Equal(t, 2, s.Len())
}

func TestKeys(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, []int{}, Keys[int](s))

	k := testKey(t)
	Set(s, k, 1)
	Set(s, testKey(t), 2.0)
	assert.Equal(t, []int{k}, Keys[int](s))
}

func TestStore_Keys(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, []int{}, s.Keys())

	k1 := testKey(t)
	Set(s, k1, 1)
	k2 := testKey(t)
	Set(s, k2, 2.0)
	assert.ElementsMatch(t, []int{k1, k2}, s.Keys())
}

func TestValues(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, []int{}, Values[int](s))

	Set(s, testKey(t), 1)
	Set(s, testKey(t), 2.0)
	assert.Equal(t, []int{1}, Values[int](s))
}

func TestStore_Values(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, []any{}, s.Values())

	Set(s, testKey(t), 1)
	Set(s, testKey(t), 2.0)
	assert.ElementsMatch(t, []any{1, 2.0}, s.Values())
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

func TestStore_Entries(t *testing.T) {
	s := &Store[int]{}

	assert.Equal(t, []Entry[int, any]{}, s.Entries())

	k1 := testKey(t)
	Set(s, k1, 1)
	k2 := testKey(t)
	Set(s, k2, 2.0)
	assert.ElementsMatch(t, []Entry[int, any]{{k1, 1}, {k2, 2.0}}, s.Entries())
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

func TestStore_ForEach(t *testing.T) {
	s := &Store[int]{}

	count := 0
	ForEach(s, func(_ int, _ any) {
		count++
	})
	assert.Equal(t, 0, count)

	Set(s, testKey(t), 1)
	Set(s, testKey(t), 2.0)
	s.ForEach(func(_ int, _ any) {
		count++
	})
	assert.Equal(t, 2, count)
}
