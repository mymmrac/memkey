package memkey

import (
	"sync"
	"testing"
)

type testInterface interface {
	test()
}

type testInterfaceImpl struct{}

func (t testInterfaceImpl) test() {
	panic("implement me")
}

var (
	testIntKey     int
	testIntKeyLock sync.Mutex
)

func testKey(t *testing.T) int {
	t.Helper()

	testIntKeyLock.Lock()
	defer testIntKeyLock.Unlock()

	testIntKey++
	return testIntKey
}
