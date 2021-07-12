// Copyright 2021 Mark Mandriota. All right reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package stack

import (
	"reflect"
	"testing"
)

var stack *Frame

func init() {
	_ = make([]byte, 1<<32-1)
	stack = NewStack(0).Ret(1<<14)
}

func TestStack(t *testing.T) {
	frameT := reflect.TypeOf(stack).Elem()
	t.Logf("FRAME:\tsize:%8X align:%8X", frameT.Size(), frameT.Align())

	for i := 0; i < frameT.NumField(); i++ {
		fieldT := frameT.Field(i)
		t.Logf("field:%d:%s)\toffset:%8X align:%8X size:%8X", i, fieldT.Name, fieldT.Offset, fieldT.Type.FieldAlign(), fieldT.Type.Size())
	}
}

func TestStack_Add(t *testing.T) {
	stack = stack.Add(42).Add("heh").Add(TestStack_Add)
	t.Logf("Current position:\t%d", stack.cp)
}

func TestStack_Sub(t *testing.T) {
	var (
		fun func(t *testing.T)
		str string
		num int
	)

	stack = stack.Sub(&fun).Sub(&str).Sub(&num)

	t.Logf("Result: %v %v %v", reflect.ValueOf(fun), str, num)
	t.Logf("Current position:\t%d", stack.cp)
}

func BenchmarkStack_Add(b *testing.B) {
	src := 1
	for i := 0; i < b.N; i++ {
		stack = stack.Add(src)
	}
}

func BenchmarkStack_Sub(b *testing.B) {
	dst := 0
	for i := 0; i < b.N; i++ {
		if stack = stack.Sub(&dst); stack == nil {
			panic("No elements onto a stack")
		}
	}
}