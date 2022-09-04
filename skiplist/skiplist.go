package skiplist

import (
	"bytes"
	"errors"
	"math/rand"
)

var (
	ErrElementNotFound = errors.New("ElementNotFound")
)

type SkipList struct {
	header   *Element
	maxLevel int
}

type Element struct {
	record   *Record
	tag      float64
	level    int
	levelPtr []*Element
}

type Record struct {
	key   []byte
	value []byte
}

func New(maxLevel int) *SkipList {
	return new(SkipList).init(maxLevel)
}

func (s *SkipList) init(maxLevel int) *SkipList {
	s.maxLevel = maxLevel
	s.header = &Element{
		record:   nil,
		tag:      0,
		level:    0,
		levelPtr: make([]*Element, s.maxLevel),
	}

	return s
}

func (s *SkipList) Put(key []byte, value []byte) error {
	r := &Record{
		key:   key,
		value: value,
	}
	s.add(r)
	return nil
}

func (s *SkipList) Get(key []byte) ([]byte, error) {
	curNode := s.header
	i := s.maxLevel - 1
	for i >= 0 {
		for next := curNode.levelPtr[i]; next != nil; next = curNode.levelPtr[i] {
			cmp := bytes.Compare(next.record.key, key)
			if cmp <= 0 {
				if cmp == 0 {
					return next.record.value, nil
				}
				curNode = next
				continue
			} else {
				break
			}
		}
		i--
	}
	return nil, ErrElementNotFound
}

func (s *SkipList) add(rec *Record) {
	// tag := s.getTag(rec.key)
	curNode := s.header
	i := s.maxLevel - 1
	preNodes := make([]*Element, s.maxLevel)
	for i >= 0 {
		preNodes[i] = curNode
		for next := curNode.levelPtr[i]; next != nil; next = curNode.levelPtr[i] {
			cmp := bytes.Compare(next.record.key, rec.key)
			// next.key < rec.key
			if cmp <= 0 {
				if cmp == 0 {
					next.record = rec
					return
				}
				curNode = next
				continue
			// next.key > rec.key
			} else {
				break
			}
		}
		preNodes[i] = curNode
		i--
	}

	lev := s.randLevel()
	elem := &Element{
		record:   rec,
		tag:      s.getTag(rec.key),
		level:    lev,
		levelPtr: make([]*Element, s.maxLevel),
	}
	for i := 0; i < lev; i++ {
		elem.levelPtr[i] = preNodes[i].levelPtr[i]
		preNodes[i].levelPtr[i] = elem
	}

}

func (s *SkipList) getTag(key []byte) float64 {
	tag := uint64(0)
	l := len(key)
	if l > 8 {
		l = 8 
	}
	for i := 0; i < l; i++ {
		shift := uint(64 - (i+1)*8)
		tag |= uint64(key[i]) << shift
	}
	return float64(tag)
}

func (s *SkipList) randLevel() int {
	level := 1
	for ; level < s.maxLevel; level++ {
		if rand.Intn(2) == 0 {
			return level
		}
	}
	return level
}

type RangFunc func(key, value []byte) bool

func (s *SkipList) Range(rf RangFunc) {
	curNode := s.header
	for next := curNode.levelPtr[0]; next != nil; next = curNode.levelPtr[0] {
		if b := rf(next.record.key, next.record.value); b {
			continue
		} else {
			break
		}
	}
}
