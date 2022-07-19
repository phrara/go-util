package vector

import (
	"fmt"
	"testing"
)

func TestVector(t *testing.T) {
	vec := New[string]()
	vec.PushBack("elem_1")
	vec.PushFront("elem_0")
	vec.PushBack("elem_2")
	iter := vec.Iter()
	for i := iter.Begin(); i != nil; i = iter.Next() {
		fmt.Println(i)
	}

	fmt.Println("-----------------------------")

	fmt.Println(vec.PopFront())
	iter = vec.Iter()
	for i := iter.Begin(); i != nil; i = iter.Next() {
		fmt.Println(i, iter.Index())
	}

	fmt.Println("-----------------------------")

	vec.PushBack("elem_3")
	_, err := vec.Remove(1)
	if err != nil {
		return
	}
	iter = vec.Iter()
	for i := iter.Begin(); i != nil; i = iter.Next() {
		fmt.Println(i, iter.Index())
	}

}
