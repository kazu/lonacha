package lonacha

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/thoas/go-funk"

	"github.com/stretchr/testify/assert"
)

type Element struct {
	ID   int
	Name string
}

const (
	CREATE_SLICE_MAX int = 10000
)

func MakeSliceSample() (slice []Element) {
	slice = make([]Element, 0, CREATE_SLICE_MAX*2)

	for i := 0; i < CREATE_SLICE_MAX; i++ {
		slice = append(slice,
			Element{
				ID:   int(math.Abs(float64(rand.Intn(CREATE_SLICE_MAX)))),
				Name: fmt.Sprintf("aaa%d", i),
			})
	}
	return
}

type Elements []Element
type PtrElements []*Element

type EqInt map[string]int
type EqString map[string]string

type Eq map[string]interface{}

func (eq EqInt) Func(slice Elements) (funcs []CondFunc) {
	funcs = make([]CondFunc, 0, len(eq))
	for key, _ := range eq {
		switch key {
		case "ID":
			fn := func(i int) bool {
				return slice[i].ID == eq[key]
			}
			funcs = append(funcs, fn)
		}
	}
	return
}

func (eq EqString) Func(slice Elements) (funcs []CondFunc) {
	funcs = make([]CondFunc, 0, len(eq))
	for key, _ := range eq {
		switch key {
		case "Name":
			fn := func(i int) bool {
				return slice[i].Name == eq[key]
			}
			funcs = append(funcs, fn)
		}
	}
	return
}

func (slice Elements) Where(q Eq) Elements {
	eqInt := make(EqInt)
	eqString := make(EqString)
	funcs := make([]CondFunc, 0, len(q))
	for key, value := range q {
		switch key {
		case "ID":
			eqInt[key] = value.(int)
			funcs = append(funcs, eqInt.Func(slice)[0])
		case "Name":
			eqString[key] = value.(string)
			funcs = append(funcs, eqString.Func(slice)[0])
		}
	}
	oslice := &slice
	Filter(oslice, funcs...)
	return *oslice
}

func MakePtrSliceSample() (slice []*Element) {
	slice = make([]*Element, 0, CREATE_SLICE_MAX*2)

	for i := 0; i < CREATE_SLICE_MAX; i++ {
		slice = append(slice,
			&Element{
				ID:   int(math.Abs(float64(rand.Intn(CREATE_SLICE_MAX)))),
				Name: fmt.Sprintf("aaa%d", i),
			})
	}
	return
}

func TestWhere(t *testing.T) {
	slice := Elements(MakeSliceSample())

	nSlice := slice.Where(Eq{"ID": 555})

	assert.True(t, nSlice[0].ID == 555, nSlice)
	assert.True(t, len(nSlice) < 100, len(nSlice))
	t.Logf("nSlice.len=%d cap=%d\n", len(nSlice), cap(nSlice))
}

func TestFind(t *testing.T) {
	nSlice := Elements(MakeSliceSample())
	id := nSlice[50].ID
	data, err := Find(&nSlice, func(i int) bool {
		return nSlice[i].ID == id
	})

	assert.NoError(t, err)
	elm := data.(Element)
	assert.True(t, elm.ID == id, elm)

}

func TestFilter(t *testing.T) {
	nSlice := Elements(MakeSliceSample())
	id := nSlice[50].ID
	Filter(&nSlice, func(i int) bool {
		return nSlice[i].ID == id
	})

	assert.True(t, nSlice[0].ID == id, nSlice)
	assert.True(t, len(nSlice) < CREATE_SLICE_MAX, len(nSlice))
	t.Logf("nSlice.len=%d cap=%d\n", len(nSlice), cap(nSlice))
}

func TestDelete(t *testing.T) {
	nSlice := Elements(MakeSliceSample())
	size := len(nSlice)
	Delete(&nSlice, func(i int) bool {
		return nSlice[i].ID == 555
	})

	assert.True(t, nSlice[0].ID != 555, nSlice)
	assert.True(t, len(nSlice) < size, len(nSlice))
	t.Logf("nSlice.len=%d cap=%d\n", len(nSlice), cap(nSlice))
}

func TestUniq(t *testing.T) {
	nSlice := Elements(MakeSliceSample())
	nSlice = append(nSlice, Element{ID: nSlice[0].ID})
	size := len(nSlice)

	fn := func(i int) interface{} { return i }
	assert.NotEqual(t, fn(1), fn(2))
	Uniq(&nSlice, func(i int) interface{} {
		return nSlice[i].ID
	})

	assert.True(t, len(nSlice) < size, len(nSlice))
	t.Logf("nSlice.len=%d cap=%d\n", len(nSlice), cap(nSlice))
}

func TestSelect(t *testing.T) {
	slice := MakeSliceSample()

	ret, err := Select(&slice, func(i int) bool {
		return slice[i].ID < 50
	})
	nSlice, ok := ret.([]Element)

	assert.NoError(t, err)
	assert.True(t, nSlice[0].ID < 50, nSlice)
	assert.True(t, ok)
	assert.True(t, len(nSlice) < 100, len(nSlice))
	t.Logf("nSlice.len=%d cap=%d\n", len(nSlice), cap(nSlice))
}

func TestPtrSelect(t *testing.T) {
	slice := MakePtrSliceSample()

	ret, err := Select(&slice, func(i int) bool {
		return slice[i].ID < 50
	})
	nSlice, ok := ret.([]*Element)

	assert.NoError(t, err)
	assert.True(t, nSlice[0].ID < 50, nSlice)
	assert.True(t, ok)
	assert.True(t, len(nSlice) < 100, len(nSlice))
	t.Logf("nSlice.len=%d cap=%d\n", len(nSlice), cap(nSlice))
}

func BenchmarkFilter(b *testing.B) {

	orig := MakeSliceSample()

	b.ResetTimer()
	b.Run("lonacha.Filter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			objs := make([]Element, len(orig))
			copy(objs, orig)
			b.StartTimer()
			Filter(&objs, func(i int) bool {
				return objs[i].ID == 555
			})
		}
	})

	pObjs := MakePtrSliceSample()
	b.ResetTimer()
	b.Run("lonacha.Filter pointer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			objs := make([]*Element, 0, len(pObjs))
			copy(objs, pObjs)
			b.StartTimer()
			Filter(&objs, func(i int) bool {
				return objs[i].ID == 555
			})
		}
	})

	b.ResetTimer()
	b.Run("hand Filter pointer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			objs := make([]*Element, len(orig))
			copy(objs, pObjs)
			b.StartTimer()
			result := make([]*Element, 0, len(orig))
			for idx, _ := range objs {
				if objs[idx].ID == 555 {
					result = append(result, objs[idx])
				}
			}
		}
	})

}

// BenchmarkUniq/lonacha.Uniq-16         	    			1000	    997543 ns/op	  548480 B/op	   16324 allocs/op
// BenchmarkUniq/lonacha.UniqWithSort-16 	    			1000	   2237924 ns/op	     256 B/op	       7 allocs/op
// BenchmarkUniq/lonacha.UniqWithSort(sort)-16         	    1000	    260283 ns/op	     144 B/op	       4 allocs/op
// BenchmarkUniq/hand_Uniq-16                          	    1000	    427765 ns/op	  442642 B/op	       8 allocs/op
// BenchmarkUniq/hand_Uniq_iface-16                    	    1000	    808895 ns/op	  632225 B/op	    6322 allocs/op
// BenchmarkUniq/go-funk.Uniq-16                       	    1000	   1708396 ns/op	  655968 B/op	   10004 allocs/op
func BenchmarkUniq(b *testing.B) {

	orig := MakeSliceSample()

	b.ResetTimer()
	b.Run("lonacha.Uniq", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			objs := make([]Element, len(orig))
			copy(objs, orig)
			b.StartTimer()
			Uniq(&objs, func(i int) interface{} {
				return objs[i].ID
			})
		}
	})

	b.ResetTimer()
	b.Run("lonacha.UniqWithSort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			objs := make([]Element, len(orig))
			copy(objs, orig)
			b.StartTimer()
			UniqWithSort(&objs, func(i, j int) bool {
				return objs[i].ID < objs[j].ID
			})
		}
	})

	b.ResetTimer()
	b.Run("lonacha.UniqWithSort(sort)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			objs := make([]Element, len(orig))
			copy(objs, orig)
			sort.Slice(objs, func(i, j int) bool {
				return objs[i].ID < objs[j].ID
			})
			b.StartTimer()
			UniqWithSort(&objs, func(i, j int) bool {
				return objs[i].ID < objs[j].ID
			})
		}
	})

	b.ResetTimer()
	b.Run("hand Uniq", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			objs := make([]Element, len(orig))
			copy(objs, orig)
			b.StartTimer()
			exists := make(map[int]bool, len(objs))
			result := make([]Element, 0, len(orig))
			for idx, _ := range objs {
				if !exists[objs[idx].ID] {
					exists[objs[idx].ID] = true
					result = append(result, orig[idx])
				}
			}
		}
	})

	b.ResetTimer()
	b.Run("hand Uniq iface", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			objs := make([]Element, len(orig))
			copy(objs, orig)
			b.StartTimer()
			exists := make(map[interface{}]bool, len(objs))
			result := make([]Element, 0, len(orig))
			for idx, _ := range objs {
				if !exists[objs[idx].ID] {
					exists[objs[idx].ID] = true
					result = append(result, orig[idx])
				}
			}
		}
	})

	b.ResetTimer()
	b.Run("go-funk.Uniq", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			objs := make([]Element, len(orig))
			copy(objs, orig)
			b.StartTimer()
			funk.Uniq(objs)
		}
	})
}

func BenchmarkSelect(b *testing.B) {
	orig := MakeSliceSample()

	b.ResetTimer()
	b.Run("lonacha.Select", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			objs := make([]Element, len(orig))
			copy(objs, orig)
			b.StartTimer()
			Select(&objs, func(i int) bool {
				return objs[i].ID == 555
			})
		}
	})

	b.ResetTimer()
	b.Run("lonacha.FilterAndCopy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			objs := make([]Element, len(orig))
			copy(objs, orig)
			b.StartTimer()
			Filter(&objs, func(i int) bool {
				return objs[i].ID == 555
			})
			newObjs := make([]Element, len(orig))
			copy(newObjs, objs)
		}
	})

	b.ResetTimer()
	b.Run("hand Select", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			objs := make([]Element, len(orig))
			copy(objs, orig)
			b.StartTimer()
			result := make([]Element, len(orig))
			for idx, _ := range objs {
				if objs[idx].ID == 555 {
					result = append(result, objs[idx])
				}
			}
		}
	})

}

type TestInterface interface {
	Inc() int
	Name() string
}

type TestObject struct {
	Cnt  int
	name string
}

func (o TestObject) Inc() int {
	o.Cnt++
	return o.Cnt
}

func (o TestObject) Name() string {
	return o.name
}

func BenchmarkCall(b *testing.B) {

	b.ResetTimer()
	b.Run("struct call", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			object := TestObject{Cnt: 0, name: "Test"}
			b.StartTimer()
			for j := 0; j < 100000; j++ {
				object.Inc()
			}
		}
	})

	b.ResetTimer()
	b.Run("interface call", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			object := TestInterface(TestObject{Cnt: 0, name: "Test"})
			b.StartTimer()
			for j := 0; j < 100000; j++ {
				object.Inc()
			}
		}
	})
}

func (list PtrElements) Len() int           { return len(list) }
func (list PtrElements) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }
func (list PtrElements) Less(i, j int) bool { return list[i].ID < list[j].ID }

// BenchmarkSortPtr/sort.Sort-16         	    1000	   1712284 ns/op	      32 B/op	       1 allocs/op
// BenchmarkSortPtr/sort.Slice-16        	    2000	   1170132 ns/op	      64 B/op	       2 allocs/op

func BenchmarkSortPtr(b *testing.B) {
	orig := MakePtrSliceSample()

	b.ResetTimer()
	b.Run("sort.Sort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			data := make(PtrElements, len(orig))
			copy(data, orig)
			b.StartTimer()
			sort.Sort(data)
		}
	})

	b.ResetTimer()
	b.Run("sort.Slice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			data := make([]*Element, len(orig))
			copy(data, orig)
			b.StartTimer()
			sort.Slice(data, func(i, j int) bool { return data[i].ID < data[j].ID })
		}
	})
}

func (list Elements) Len() int           { return len(list) }
func (list Elements) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }
func (list Elements) Less(i, j int) bool { return list[i].ID < list[j].ID }

// BenchmarkSort/sort.Sort-16         	    1000	   1648947 ns/op	      34 B/op	       1 allocs/op
// BenchmarkSort/sort.Slice-16        	    1000	   1973036 ns/op	     112 B/op	       3 allocs/op

func BenchmarkSort(b *testing.B) {
	orig := MakeSliceSample()

	b.ResetTimer()
	b.Run("sort.Sort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			data := make(Elements, len(orig))
			copy(data, orig)
			b.StartTimer()
			sort.Sort(data)
		}
	})

	b.ResetTimer()
	b.Run("sort.Slice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			data := make([]Element, len(orig))
			copy(data, orig)
			b.StartTimer()
			sort.Slice(data, func(i, j int) bool { return data[i].ID < data[j].ID })
		}
	})
}
