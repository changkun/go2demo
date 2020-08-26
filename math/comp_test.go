// Code generated by go2go; DO NOT EDIT.


//line comp_test.go2:5
package math

//line comp_test.go2:5
import "testing"

//line comp_test.go2:9
func TestMax(t *testing.T) {
	got := instantiate୦୦Max୦int(1, 2, 3)
	if got != 3 {
		t.Fatalf("want 3")
	}
	got = instantiate୦୦Max୦int(1, 3, 2)
	if got != 3 {
		t.Fatalf("want 3")
	}
	got = instantiate୦୦Max୦int(2, 1, 3)
	if got != 3 {
		t.Fatalf("want 3")
	}
	got = instantiate୦୦Max୦int(2, 3, 1)
	if got != 3 {
		t.Fatalf("want 3")
	}
	got = instantiate୦୦Max୦int(3, 1, 2)
	if got != 3 {
		t.Fatalf("want 3")
	}
	got = instantiate୦୦Max୦int(3, 2, 1)
	if got != 3 {
		t.Fatalf("want 3")
	}
}

func TestMin(t *testing.T) {
	m := instantiate୦୦Min୦int(1, 2, 3)
	if m != 1 {
		t.Fatalf("want 3, got %v", m)
	}
}
//line comp.go2:33
func instantiate୦୦Max୦int(v0 int, vn ...int,) int {
	switch l := len(vn); {
	case l == 0:
		return v0
	case l == 1:
		if v0 > vn[0] {
//line comp.go2:38
   return v0
//line comp.go2:38
  }
				return vn[0]
	default:
		vv := instantiate୦୦Max୦int(vn[0], vn[1:]...)
		if v0 > vv {
//line comp.go2:42
   return v0
//line comp.go2:42
  }
				return vv
	}
}
//line comp.go2:19
func instantiate୦୦Min୦int(s ...int,) int {
	if len(s) == 0 {
		panic("Min of no elements")
	}
	r := s[0]
	for _, v := range s[1:] {
		if v < r {
			r = v
		}
	}
	return r
}

//line comp.go2:30
var _ = testing.AllocsPerRun
