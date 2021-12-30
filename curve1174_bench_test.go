package curve1174

import (
	"crypto/elliptic"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

var scalar, _ = new(big.Int).SetString("5ffffff23ffffffffffff1758603537581ffffffffffffffffffffffffffff7", 16)

func BenchmarkCurve1174Add(b *testing.B) {
	var p1, p2 Point
	p1.Add(UBase, UBase)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p2.Add(UBase, &p1)
	}
}

func BenchmarkCurve1174Double(b *testing.B) {
	var p Point
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p.Double(UBase)
	}
}

func BenchmarkCurve1174ScalarMult(b *testing.B) {
	var p Point
	f := FromBigInt(scalar)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p.ScalarMult(UBase, f) //.ToAffine(&p)
	}
}

var pp Point

func BenchmarkCurve1174ScalarBaseMult(b *testing.B) {
	var p Point
	f := FromBigInt(scalar)
	PrecomputeBase()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p.ScalarBaseMult(f).ToAffine(&p)
	}
	b.StopTimer()
	pp = p
}

func BenchmarkCurve1174ScalarBaseMult2(b *testing.B) {
	var p Point
	f := FromBigInt(scalar)
	PrecomputeBase2()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p.ScalarBaseMult2(f).ToAffine(&p)
	}
}

func BenchmarkCurveP256Add(b *testing.B) {
	p256 := elliptic.P256()
	params := p256.Params()
	x, y := p256.Double(params.Gx, params.Gy)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p256.Add(params.Gx, params.Gy, x, y)
	}
}

func BenchmarkCurveP256Double(b *testing.B) {
	p256 := elliptic.P256()
	params := p256.Params()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p256.Double(params.Gx, params.Gy)
	}
}

func BenchmarkCurveP256ScalarMult(b *testing.B) {
	p256 := elliptic.P256()
	params := p256.Params()
	bytes := scalar.Bytes()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p256.ScalarMult(params.Gx, params.Gy, bytes)
	}
}

func BenchmarkCurveP256ScalarBaseMult(b *testing.B) {
	p256 := elliptic.P256()
	bytes := scalar.Bytes()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p256.ScalarBaseMult(bytes)
	}
}

func BenchmarkCurveP521ScalarMult(b *testing.B) {
	p521 := elliptic.P521()
	params := p521.Params()
	bytes := scalar.Bytes()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p521.ScalarMult(params.Gx, params.Gy, bytes)
	}
}

func BenchmarkCurveP521ScalarBaseMult(b *testing.B) {
	p521 := elliptic.P521()
	bytes := scalar.Bytes()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p521.ScalarBaseMult(bytes)
	}
}

func BenchmarkMod(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var p, p2 FieldElement
	p[0] = r.Uint64()
	p[1] = r.Uint64()
	p[2] = r.Uint64()
	p[3] = r.Uint64()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mod(&p, &p2)
	}
}

func BenchmarkInverseBigInt(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b2 := new(big.Int).Rand(r, P)
	b4 := new(big.Int)
	var p Point
	p.ScalarBaseMult(FromBigInt(b2))
	b.ReportAllocs()
	b.ResetTimer()
	bigInt := p.Z.ToBigInt()
	for i := 0; i < b.N; i++ {
		b4.ModInverse(bigInt, P)
	}
}

func BenchmarkFastInverse(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b2 := new(big.Int).Rand(r, P)
	e := FromBigInt(b2)
	var p Point
	p.ScalarBaseMult(e)
	b.ReportAllocs()
	var ee FieldElement
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ee.FastInverse(&p.Z)
	}
}

func BenchmarkInverse(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b2 := new(big.Int).Rand(r, P)
	e := FromBigInt(b2)
	var p Point
	p.ScalarBaseMult(e)
	b.ReportAllocs()
	var ee FieldElement
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ee.Inverse(&p.Z)
	}
}

func BenchmarkMul(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var p, p2, p3 FieldElement
	p[0] = r.Uint64()
	p[1] = r.Uint64()
	p[2] = r.Uint64()
	p[3] = r.Uint64()
	p2[0] = r.Uint64()
	p2[1] = r.Uint64()
	p2[2] = r.Uint64()
	p2[3] = r.Uint64()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mul(&p3, &p, &p2)
	}
}

func BenchmarkMulSqr(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var p, p2, p3 FieldElement
	p[0] = r.Uint64()
	p[1] = r.Uint64()
	p[2] = r.Uint64()
	p[3] = r.Uint64()
	p2[0] = r.Uint64()
	p2[1] = r.Uint64()
	p2[2] = r.Uint64()
	p2[3] = r.Uint64()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sqr(&p3, &p)
	}
}

func BenchmarkAdd(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var p, p2, p3 FieldElement
	p[0] = r.Uint64()
	p[1] = r.Uint64()
	p[2] = r.Uint64()
	p[3] = r.Uint64()
	p2[0] = r.Uint64()
	p2[1] = r.Uint64()
	p2[2] = r.Uint64()
	p2[3] = r.Uint64()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		add(&p3, &p, &p2)
	}
}

func BenchmarkSub(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var p, p2, p3 FieldElement
	p[0] = r.Uint64()
	p[1] = r.Uint64()
	p[2] = r.Uint64()
	p[3] = r.Uint64()
	p2[0] = r.Uint64()
	p2[1] = r.Uint64()
	p2[2] = r.Uint64()
	p2[3] = r.Uint64()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sub(&p3, &p, &p2)
	}
}

func BenchmarkMulD(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var p, p2 FieldElement
	p[0] = r.Uint64()
	p[1] = r.Uint64()
	p[2] = r.Uint64()
	p[3] = r.Uint64()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mulD(&p2, &p)
	}
}

func BenchmarkMul2(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var p, p2 FieldElement
	p[0] = r.Uint64()
	p[1] = r.Uint64()
	p[2] = r.Uint64()
	p[3] = r.Uint64()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mul2(&p2, &p)
	}
}

func BenchmarkPrecompute(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		precomputed = false
		PrecomputeBase()
	}
}

//go:noinline
func sumCopy(data [18]int) int {
	return data[0] + data[1] + data[2] + data[3] + data[4] + data[5] + data[6] + data[7]
}

//go:noinline
func sumPointer(data *[18]int) int {
	return data[0] + data[1] + data[2] + data[3] + data[4] + data[5] + data[6] + data[7]
}

func BenchmarkCallCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := [18]int{i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i}
		sumCopy(data)
	}
}

func BenchmarkCallPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := [18]int{i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i}
		sumPointer(&data)
	}
}
