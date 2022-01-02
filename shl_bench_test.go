//go:build !curve1174_purego || !amd64

package curve1174

import "testing"

//go:noescape
func shl(res, x *[8]uint64)

//go:noescape
func shl2(res, x *[8]uint64)

var r [8]uint64

func BenchmarkShl(b *testing.B) {
	x1, x2 := [8]uint64{}, [8]uint64{}
	for i := 0; i < b.N; i++ {
		shl(&x1, &x2)
		shl(&x2, &x1)
	}
	b.StopTimer()
	r = x1
}

func BenchmarkShl2(b *testing.B) {
	x1, x2 := [8]uint64{}, [8]uint64{}
	for i := 0; i < b.N; i++ {
		shl2(&x1, &x2)
		shl2(&x2, &x1)
	}
	b.StopTimer()
	r = x1
}
