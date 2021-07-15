// Copyright 2021 Mark Mandriota. All right reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package stack

import (
	r "reflect"
	"sync"
)

// flag
const (
	ATOMIC = 1 << iota // enables multi-thread safety.
)

type Buffer [1 << 12]interface{}

type Frame struct {
	sync.Mutex

	Flag uint32

	cp uint32
	bf Buffer

	Prev, Next *Frame
}

func NewStack(f uint32) *Frame {
	return &Frame{Flag: f}
}

// Len - returns len of current frame + len of prev frame.
func (f *Frame) Len() int {
	if f.Prev != nil {
		return int(f.cp+f.Prev.cp)
	}
	return int(f.cp)
}

// Ret - returns to n frames.
// This can be used to clear and allocate memory.
func (f *Frame) Ret(n uint64) *Frame {
	for i := n; i > 0; i-- {
		if f.Prev == nil {
			f.Prev = &Frame{Flag: f.Flag, Next: f}
		}
		f.Prev.cp = 0
		f = f.Prev
	}

	return f
}

// Add - Pushes src onto the stack.
// If Frame is not full returns it, else returns next Frame with calling Add.
//
// Warning: Return value must be assigned!
func (f *Frame) Add(src interface{}) *Frame {
	if f.Flag&ATOMIC == ATOMIC {
		f.Lock()
		defer f.Unlock()
	}

	if f.cp < uint32(len(f.bf)) {
		f.bf[f.cp] = src
		f.cp++
	} else if f.Next == nil {
		f.Next = &Frame{Flag: f.Flag, cp: 1, bf: Buffer{src}, Prev: f}
	} else {
		return f.Next.Add(src)
	}

	return f
}

// Sub - Popes value from stack to dst. Panics if cannot assign value to dst.
// If current Frame is not empty returns it, else returns prev Frame with calling Sub if is not nil.
//
// Warning: Return value must be assigned!
func (f *Frame) Sub(dst interface{}) *Frame {
	if f.Flag&ATOMIC == ATOMIC {
		f.Lock()
		defer f.Unlock()
	}

	if f.cp > 0 {
		f.cp--
		r.ValueOf(dst).Elem().Set(r.ValueOf(f.bf[f.cp]))
		return f
	}

	if f.Prev != nil {
		return f.Prev.Sub(dst)
	}
	return nil
}