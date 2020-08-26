// Code generated by go2go; DO NOT EDIT.


//line map_test.go2:5
package sync

//line map_test.go2:5
import (
//line map_test.go2:5
 "math/rand"
//line map_test.go2:5
 "reflect"
//line map_test.go2:5
 "runtime"
//line map_test.go2:5
 "sync"
//line map_test.go2:5
 "sync/atomic"
//line map_test.go2:5
 "testing"
//line map_test.go2:5
 "testing/quick"
//line map_test.go2:5
 "unsafe"
//line map_test.go2:5
)

//line map_test.go2:16
type mapOp string

const (
	opLoad          = mapOp("Load")
	opStore         = mapOp("Store")
	opLoadOrStore   = mapOp("LoadOrStore")
	opLoadAndDelete = mapOp("LoadAndDelete")
	opDelete        = mapOp("Delete")
)

var mapOps = [...]mapOp{opLoad, opStore, opLoadOrStore, opLoadAndDelete, opDelete}

//line map_test.go2:29
type mapCall struct {
	op   mapOp
	k, v interface{}
}

func (c mapCall) apply(m mapInterface) (interface{}, bool) {
	switch c.op {
	case opLoad:
		return m.Load(c.k)
	case opStore:
		m.Store(c.k, c.v)
		return nil, false
	case opLoadOrStore:
		return m.LoadOrStore(c.k, c.v)
	case opLoadAndDelete:
		return m.LoadAndDelete(c.k)
	case opDelete:
		m.Delete(c.k)
		return nil, false
	default:
		panic("invalid mapOp")
	}
}

type mapResult struct {
	value interface{}
	ok    bool
}

func randValue(r *rand.Rand) interface{} {
	b := make([]byte, r.Intn(4))
	for i := range b {
		b[i] = 'a' + byte(rand.Intn(26))
	}
	return string(b)
}

func (mapCall) Generate(r *rand.Rand, size int) reflect.Value {
	c := mapCall{op: mapOps[rand.Intn(len(mapOps))], k: randValue(r)}
	switch c.op {
	case opStore, opLoadOrStore:
		c.v = randValue(r)
	}
	return reflect.ValueOf(c)
}

func applyCalls(m mapInterface, calls []mapCall) (results []mapResult, final map[interface{}]interface{}) {
	for _, c := range calls {
		v, ok := c.apply(m)
		results = append(results, mapResult{v, ok})
	}

	final = make(map[interface{}]interface{})
	m.Range(func(k, v interface{}) bool {
		final[k] = v
		return true
	})

	return results, final
}

func applyMap(calls []mapCall) ([]mapResult, map[interface{}]interface{}) {
	return applyCalls(new(instantiate୦୦Map୦interface୮4୮5୦interface୮4୮5), calls)
}

func applyRWMutexMap(calls []mapCall) ([]mapResult, map[interface{}]interface{}) {
	return applyCalls(new(RWMutexMap), calls)
}

func applyDeepCopyMap(calls []mapCall) ([]mapResult, map[interface{}]interface{}) {
	return applyCalls(new(DeepCopyMap), calls)
}

func TestMapMatchesRWMutex(t *testing.T) {
	if err := quick.CheckEqual(applyMap, applyRWMutexMap, nil); err != nil {
		t.Error(err)
	}
}

func TestMapMatchesDeepCopy(t *testing.T) {
	if err := quick.CheckEqual(applyMap, applyDeepCopyMap, nil); err != nil {
		t.Error(err)
	}
}

func TestConcurrentRange(t *testing.T) {
	const mapSize = 1 << 10

	m := new(instantiate୦୦Map୦interface୮4୮5୦interface୮4୮5)
	for n := int64(1); n <= mapSize; n++ {
		m.Store(n, int64(n))
	}

	done := make(chan struct{})
	var wg sync.WaitGroup
	defer func() {
		close(done)
		wg.Wait()
	}()
	for g := int64(runtime.GOMAXPROCS(0)); g > 0; g-- {
		r := rand.New(rand.NewSource(g))
		wg.Add(1)
		go func(g int64) {
			defer wg.Done()
			for i := int64(0); ; i++ {
				select {
				case <-done:
					return
				default:
				}
				for n := int64(1); n < mapSize; n++ {
					if r.Int63n(mapSize) == 0 {
						m.Store(n, n*i*g)
					} else {
						m.Load(n)
					}
				}
			}
		}(g)
	}

	iters := 1 << 10
	if testing.Short() {
		iters = 16
	}
	for n := iters; n > 0; n-- {
		seen := make(map[int64]bool, mapSize)

		m.Range(func(ki, vi interface{}) bool {
			k, v := ki.(int64), vi.(int64)
			if v%k != 0 {
				t.Fatalf("while Storing multiples of %v, Range saw value %v", k, v)
			}
			if seen[k] {
				t.Fatalf("Range visited key %v twice", k)
			}
			seen[k] = true
			return true
		})

		if len(seen) != mapSize {
			t.Fatalf("Range visited %v elements of %v-element Map", len(seen), mapSize)
		}
	}
}

//line map_test.go2:173
type instantiate୦୦Map୦interface୮4୮5୦interface୮4୮5 struct {
//line map.go2:15
 mu     sync.Mutex
			read   atomic.Value
			dirty  map[interface{}]*entry
			misses int
}

//line map.go2:36
func (m *instantiate୦୦Map୦interface୮4୮5୦interface୮4୮5,) Load(key interface{}) (value interface{}, ok bool) {
	read, _ := m.read.Load().(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5)
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return
	}
	v, loaded := e.load()
	ok = loaded
	if ok {
		value = v.(interface{})
	}
	return
}

//line map.go2:68
func (m *instantiate୦୦Map୦interface୮4୮5୦interface୮4୮5,) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *instantiate୦୦Map୦interface୮4୮5୦interface୮4୮5,) Delete(key interface{}) {
	m.LoadAndDelete(key)
}

func (m *instantiate୦୦Map୦interface୮4୮5୦interface୮4୮5,) LoadAndDelete(key interface{}) (value interface{}, loaded bool) {
	read, _ := m.read.Load().(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5)
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			delete(m.dirty, key)
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if ok {
		v, ok := e.delete()
		if ok {
			value = v.(interface{})
		}
		loaded = ok
		return
	}
	return
}

//line map.go2:119
func (m *instantiate୦୦Map୦interface୮4୮5୦interface୮4୮5,) LoadOrStore(key interface{},

//line map.go2:119
 value interface{},

//line map.go2:119
) (interface{},

//line map.go2:119
 bool) {
	read, _ := m.read.Load().(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5)
	if e, ok := read.m[key]; ok {
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual.(interface{}), loaded
		}
	}

	m.mu.Lock()
	read, _ = m.read.Load().(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5)
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		actual, loaded, _ := e.tryLoadOrStore(value)
		m.mu.Unlock()
		return actual.(interface{}), loaded
	} else if e, ok := m.dirty[key]; ok {
		actual, loaded, _ := e.tryLoadOrStore(value)
		m.missLocked()
		m.mu.Unlock()
		return actual.(interface{}), loaded
	} else {
		if !read.amended {
			m.dirtyLocked()
			m.read.Store(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry(value)
		actual, loaded := value, false
		m.mu.Unlock()
		return actual, loaded
	}
}

//line map.go2:182
func (m *instantiate୦୦Map୦interface୮4୮5୦interface୮4୮5,) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5)
	m.dirty = make(map[interface{}]*entry, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

//line map.go2:207
func (m *instantiate୦୦Map୦interface୮4୮5୦interface୮4୮5,) Range(f func(key interface{},

//line map.go2:207
 value interface{},

//line map.go2:207
) bool) {

//line map.go2:212
 read, _ := m.read.Load().(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5)
	if read.amended {

//line map.go2:218
  m.mu.Lock()
		read, _ = m.read.Load().(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5)
		if read.amended {
			read = instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5{m: m.dirty}
			m.read.Store(read)
			m.dirty = nil
			m.misses = 0
		}
		m.mu.Unlock()
	}

	for k, e := range read.m {
		v, ok := e.load()
		if !ok {
			continue
		}
		if !f(k, v.(interface{})) {
			break
		}
	}
}

func (m *instantiate୦୦Map୦interface୮4୮5୦interface୮4୮5,) Store(key interface{},

//line map.go2:240
 value interface{},

//line map.go2:240
) {
	read, _ := m.read.Load().(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5)
	vv := interface{}(value)
	if e, ok := read.m[key]; ok && e.tryStore(&vv) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5)
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		vv := interface{}(value)
		e.storeLocked(&vv)
	} else if e, ok := m.dirty[key]; ok {
		vv := interface{}(value)
		e.storeLocked(&vv)
	} else {
		if !read.amended {
			m.dirtyLocked()
			m.read.Store(instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry(value)
	}
	m.mu.Unlock()
}

//line map.go2:266
type instantiate୦୦readOnly୦interface୮4୮5୦interface୮4୮5 struct {
//line map.go2:22
 m       map[interface{}]*entry
			amended bool
}

//line map.go2:24
var _ = rand.ExpFloat64
//line map.go2:24
var _ = reflect.Append
//line map.go2:24
var _ = runtime.BlockProfile

//line map.go2:24
type _ sync.Cond

//line map.go2:24
var _ = atomic.AddInt32
//line map.go2:24
var _ = testing.AllocsPerRun
//line map.go2:24
var _ = quick.Check

//line map.go2:24
type _ unsafe.Pointer
