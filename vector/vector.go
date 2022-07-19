package vector

import (
	"errors"
	"sync"
)

type Iterable interface {
	Iter() *Iterator
	ValueAt(int) (any, error)
}

type Vector[T any] struct {
	buf  []T
	lock sync.Mutex
}

func New[T any]() *Vector[T] {
	return new(Vector[T]).init()
}

func (v *Vector[T]) init() *Vector[T] {
	v.buf = make([]T, 0)
	return v
}

func (v *Vector[T]) PushBack(value T) {
	v.lock.Lock()
	v.buf = append(v.buf, value)
	v.lock.Unlock()
}

func (v *Vector[T]) PushFront(value T) {
	v.lock.Lock()
	newBuf := make([]T, 1)
	newBuf[0] = value
	newBuf = append(newBuf, v.buf...)
	v.buf = newBuf
	v.lock.Unlock()
}

func (v *Vector[T]) Front() T {
	return v.buf[0]
}

func (v *Vector[T]) Back() T {
	return v.buf[len(v.buf)-1]
}

func (v *Vector[T]) PopBack() (ret T) {
	v.lock.Lock()
	ret = v.buf[len(v.buf)-1]
	v.buf = v.buf[:len(v.buf)-1]
	v.lock.Unlock()
	return
}

func (v *Vector[T]) PopFront() (ret T) {
	v.lock.Lock()
	ret = v.buf[0]
	v.buf = v.buf[1:]
	v.lock.Unlock()
	return
}

func (v *Vector[T]) Size() int {
	return len(v.buf)
}

func (v *Vector[T]) Resize(size int) {
	v.lock.Lock()
	v.buf = v.buf[:size]
	v.lock.Unlock()
}

func (v *Vector[T]) Remove(index int) (ret any, err error) {
	if index > len(v.buf)-1 || index < 0 {
		return nil, errors.New("index out of range")
	} else {
		if index == len(v.buf)-1 {
			ret = v.PopBack()
			return
		} else if index == 0 {
			ret = v.PopFront()
			return
		} else {
			v.lock.Lock()
			ret = v.buf[index]
			v.buf = append(v.buf[:index], v.buf[index+1:]...)
			v.lock.Unlock()
			return
		}
	}
}

func (v *Vector[T]) ValueAt(index int) (any, error) {
	if index >= len(v.buf) {
		return nil, errors.New("index out of range")
	} else {
		return v.buf[index], nil
	}
}

func (v *Vector[T]) Iter() (ret *Iterator) {
	ret = new(Iterator)
	ret.container = v
	return
}

type Iterator struct {
	container Iterable
	index     int
}

func (i *Iterator) Index() int {
	return i.index
}

func (i *Iterator) Begin() (ret any) {
	i.index = 0
	ret, err := i.container.ValueAt(0)
	if err != nil {
		return nil
	}
	return ret
}

func (i *Iterator) Next() (ret any) {
	i.index++
	ret, err := i.container.ValueAt(i.index)
	if err != nil {
		return nil
	}
	return ret
}
