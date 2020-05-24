// generated by github.com/kazu/loncha/structer

// ContainerListWriter is a base of http://golang.org/pkg/container/list/
// this is tuning performancem, reduce heap usage.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package example

import (
    "github.com/kazu/loncha/list_head"
    "unsafe"
    "errors"
)

func (d *Sample) Init() {
	d.ListHead.Init()
}

// Next returns the next list element or nil.
func (d *Sample) Next() *Sample {
	if d.ListHead.Next() == nil {
		panic(errors.New("d.next is nil"))
	}
	return (*Sample)(unsafe.Pointer(uintptr(unsafe.Pointer(d.ListHead.Next())) - unsafe.Offsetof(d.ListHead)))
}
// Prev returns the previous list element or nil.
func (d *Sample) Prev() *Sample {
	if d.ListHead.Next() == nil {
		panic(errors.New("d.prev is nil"))
	}
	return (*Sample)(unsafe.Pointer(uintptr(unsafe.Pointer(d.ListHead.Prev())) - unsafe.Offsetof(d.ListHead)))
}

// New returns an initialized list.
func NewSampleList(h *Sample) *Sample {
	h.Init()
	return h
}

func (d *Sample) Len() int {
	return d.ListHead.Len()
}

func (d *Sample) Add(n *Sample)  *Sample {
	d.ListHead.Add(&n.ListHead)
	return n
}

func (d *Sample) Delete() *Sample {
	ptr := d.ListHead.Delete()
	return (*Sample)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) - unsafe.Offsetof(d.ListHead)))
}

func (d *Sample) Remove() *Sample {
	return d.Delete()
}

func (d *Sample) ContainOf(ptr *list_head.ListHead) *Sample {
	return (*Sample)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) - unsafe.Offsetof(d.ListHead)))
}

func (d *Sample) Front() *Sample {
	return d.ContainOf(d.ListHead.Front())
}


func (d *Sample) Back() *Sample {
	return d.ContainOf(d.ListHead.Back())
}

// PushFront inserts a new value v at the front of list l and returns e.
func (d *Sample) PushFront(v *Sample) *Sample {
	front := d.Front()
	v.Add(front)
	return v
}


// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *Sample) PushBack(v *Sample) *Sample {
	last := l.Back()
	last.Add(v)
	return v
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
func (l *Sample) InsertBefore(v *Sample) *Sample {
	l.Prev().Add(v)
	return v
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
func (l *Sample) InsertAfter(v *Sample) *Sample {
	l.Next().Add(v)
	return v
}


// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
func (l *Sample) MoveToFront(v *Sample) *Sample {
	v.Remove()
	return l.PushFront(v)
}

func (l *Sample) MoveToBack(v *Sample) *Sample {
	v.Remove()
	return l.PushBack(v)
}


// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
func (l *Sample) MoveBefore(v *Sample) *Sample {
	v.Remove()
	l.Prev().Add(v)
	return v
}

// MoveAfter moves element e to its new position after mark.
// If e is not an element of l, or e == mark, the list is not modified.
func (l *Sample) MoveAfter(v *Sample) *Sample {
	v.Remove()
	l.Add(v)
	return v
}

func (l *Sample) PushBackList(other *Sample) {
	l.Back().Add(other)
	return 
}

func (l *Sample) PushFrontList(other *Sample) {
	other.PushBackList(l)
	return
}


func (l *Sample) Each(fn func(e *Sample)) {

	cur := l.Cursor()

	for cur.Next() {
		fn(l.ContainOf(cur.Pos))
	}

}
