package skiplist

import (
	"fmt"
	"testing"
)

func TestSkipList(t *testing.T) {
	sl := New(5)
	sl.Put([]byte("1"), []byte("1"))
	sl.Put([]byte("2"), []byte("2"))
	b, err := sl.Get([]byte("3"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	sl.Range(func(key, value []byte) bool {
		return true	
	})
}