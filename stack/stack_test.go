// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stack

import (
	"testing"
)

func TestStack(t *testing.T) {
	var s Stack[string]
	s.Push("hi")
	s.Push("bye")
	if s.IsEmpty() {
		t.Fatalf("unexpected IsEmpty")
	}
	if s.Len() != 2 {
		panic("bad Len")
	}
	if v, ok := s.Pop(); !ok || v != "bye" {
		t.Fatalf("bad Pop 1")
	}
	if v, ok := s.Pop(); !ok || v != "hi" {
		t.Fatalf("bad Pop 2")
	}
	if !s.IsEmpty() {
		t.Fatalf("expected IsEmpty")
	}
	if s.Len() != 0 {
		t.Fatalf("bad Len 2")
	}
}