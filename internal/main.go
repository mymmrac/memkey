package main

import (
	"fmt"

	"github.com/mymmrac/memkey"
)

func main() {
	s := &memkey.Store[int]{}

	fmt.Println(memkey.Has[string](s, 1), memkey.HasRaw(s, 1))

	text, ok := memkey.Get[string](s, 1)
	fmt.Println(text, ok)

	memkey.Set(s, 1, "hmm")

	fmt.Println(memkey.HasRaw(s, 1), memkey.Has[string](s, 1), memkey.Has[int](s, 1))

	text, ok = memkey.Get[string](s, 1)
	fmt.Println(text, ok)

	number, ok := memkey.Get[float64](s, 1)
	fmt.Println(number, ok)

	memkey.Set(s, 1, 5.2)

	number, ok = memkey.Get[float64](s, 1)
	fmt.Println(number, ok)

	fmt.Println(memkey.Has[float64](s, 1))

	raw, ok := memkey.GetRaw(s, 1)
	fmt.Println(raw, ok)

	fmt.Println("====")

	fmt.Println(memkey.KeysRaw(s), memkey.ValuesRaw(s))
	fmt.Println(memkey.Keys[int](s), memkey.Values[int](s))
}
