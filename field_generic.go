//go:build !amd64 || curve1174_purego

package curve1174

import (
	"fmt"
	"math/bits"
)

var _ = fmt.Sprintf

func add(res, p1, p2 *FieldElement) {
	r0, carry := bits.Add64(p1[0], p2[0], 0)
	r1, carry := bits.Add64(p1[1], p2[1], carry)
	r2, carry := bits.Add64(p1[2], p2[2], carry)
	r3, carry := bits.Add64(p1[3], p2[3], carry)

	r0, carry = bits.Add64(r0, ^(carry-1)&288, 0)
	res[1], carry = bits.Add64(r1, 0, carry)
	res[2], carry = bits.Add64(r2, 0, carry)
	res[3], carry = bits.Add64(r3, 0, carry)

	res[0] = r0 + ^(carry-1)&288
}

func sub(res, p1, p2 *FieldElement) {
	r0, borrow := bits.Sub64(p1[0], p2[0], 0)
	r1, borrow := bits.Sub64(p1[1], p2[1], borrow)
	r2, borrow := bits.Sub64(p1[2], p2[2], borrow)
	r3, borrow := bits.Sub64(p1[3], p2[3], borrow)

	r0, borrow = bits.Sub64(r0, ^(borrow-1)&288, 0)
	res[1], borrow = bits.Sub64(r1, 0, borrow)
	res[2], borrow = bits.Sub64(r2, 0, borrow)
	res[3], borrow = bits.Sub64(r3, 0, borrow)

	res[0] = r0 - ^(borrow-1)&288
}

func mulAdd(a, b, r0, r1, c uint64) (o0, o1, carry uint64) {
	a, b = bits.Mul64(a, b)
	o0, carry = bits.Add64(r0, b, c)
	o1, carry = bits.Add64(r1, a, carry)
	return
}

func sqr(res, p *FieldElement) {
	//var r0, r1, r2, r3, r4, r5, r6, r7, carry uint64
	//
	//r2, r1 = bits.Mul64(p[0], p[1])
	//r4, r3 = bits.Mul64(p[0], p[3])
	//r6, r5 = bits.Mul64(p[2], p[3])
	//
	//r2, r3, carry = mulAdd(p[0], p[2], r2, r3, 0)
	//r4, r5, carry = mulAdd(p[1], p[3], r4, r5, carry)
	//r6, carry = bits.Add64(0, r6, carry)
	//r7, _ = bits.Add64(0, r7, carry)
	//
	//r3, r4, carry = mulAdd(p[1], p[2], r3, r4, 0)
	//r5, _ = bits.Add64(0, r5, carry)
	//r5, carry = bits.Add64(0, r5, carry)
	//r6, carry = bits.Add64(0, r6, carry)
	//r7, _ = bits.Add64(0, r7, carry)
	//
	//r7 = r7<<1 | r6>>63
	//r6 = r6<<1 | r5>>63
	//r5 = r5<<1 | r4>>63
	//r4 = r4<<1 | r3>>63
	//r3 = r3<<1 | r2>>63
	//r2 = r2<<1 | r1>>63
	//r1 = r1 << 1
	//
	//carry, r0 = bits.Mul64(p[0], p[0])
	//r1, carry = bits.Add64(r1, carry, 0)
	//r2, r3, carry = mulAdd(p[1], p[1], r2, r3, carry)
	//r4, r5, carry = mulAdd(p[2], p[2], r4, r5, carry)
	//r6, r7, _ = mulAdd(p[3], p[3], r6, r7, carry)
	//
	//extendedMod(res, r0, r1, r2, r3, r4, r5, r6, r7)
	mul(res, p, p)
}

func mul(res, p1, p2 *FieldElement) {
	_, _ = p1[3], p2[3]
	var r0, r1, r2, r3, r4, r5, r6, r7, carry uint64

	r1, r0 = bits.Mul64(p1[0], p2[0])
	r3, r2 = bits.Mul64(p1[0], p2[2])
	r5, r4 = bits.Mul64(p1[3], p2[1])
	r7, r6 = bits.Mul64(p1[3], p2[3])

	r1, r2, carry = mulAdd(p1[1], p2[0], r1, r2, 0)
	r3, r4, carry = mulAdd(p1[0], p2[3], r3, r4, carry)
	r5, r6, carry = mulAdd(p1[2], p2[3], r5, r6, carry)
	r7, _ = bits.Add64(0, r7, carry)

	r1, r2, carry = mulAdd(p1[0], p2[1], r1, r2, 0)
	r3, r4, carry = mulAdd(p1[1], p2[2], r3, r4, carry)
	r5, r6, carry = mulAdd(p1[3], p2[2], r5, r6, carry)
	r7, _ = bits.Add64(0, r7, carry)

	r2, r3, carry = mulAdd(p1[1], p2[1], r2, r3, 0)
	r4, r5, carry = mulAdd(p1[1], p2[3], r4, r5, carry)
	r6, carry = bits.Add64(0, r6, carry)
	r7, _ = bits.Add64(0, r7, carry)

	r2, r3, carry = mulAdd(p1[2], p2[0], r2, r3, 0)
	r4, r5, carry = mulAdd(p1[2], p2[2], r4, r5, carry)
	r6, carry = bits.Add64(0, r6, carry)
	r7, _ = bits.Add64(0, r7, carry)

	r3, r4, carry = mulAdd(p1[2], p2[1], r3, r4, 0)
	r5, carry = bits.Add64(0, r5, carry)
	r6, carry = bits.Add64(0, r6, carry)
	r7, _ = bits.Add64(0, r7, carry)

	r3, r4, carry = mulAdd(p1[3], p2[0], r3, r4, 0)
	r5, carry = bits.Add64(0, r5, carry)
	r6, carry = bits.Add64(0, r6, carry)
	r7, _ = bits.Add64(0, r7, carry)

	r8 := r7 >> 59
	r7 = r7<<5 | r6>>59
	r6 = r6<<5 | r5>>59
	r5 = r5<<5 | r4>>59
	r4 = r4<<5 | r3>>59

	r3 &= P3

	r4a, carry := bits.Add64(r4, r4<<3, 0)
	r5a, carry := bits.Add64(r5, r4>>61|r5<<3, carry)
	r6a, carry := bits.Add64(r6, r5>>61|r6<<3, carry)
	r7a, carry := bits.Add64(r7, r6>>61|r7<<3, carry)
	r8a, _ := bits.Add64(r8, r7>>61|r8<<3, carry)

	r0, carry = bits.Add64(r0, r4a, 0)
	r1, carry = bits.Add64(r1, r5a, carry)
	r2, carry = bits.Add64(r2, r6a, carry)
	r3, carry = bits.Add64(r3, r7a, carry)
	r4, _ = bits.Add64(r8a, 0, carry)

	r4 = (r4<<5 | r3>>59) * 9

	r3 &= P3

	res[0], carry = bits.Add64(r0, r4, 0)
	res[1], carry = bits.Add64(r1, 0, carry)
	res[2], carry = bits.Add64(r2, 0, carry)
	res[3], _ = bits.Add64(r3, 0, carry)

}

func mulD(res, p2 *FieldElement) *FieldElement {
	r0h, r0 := bits.Mul64(1174, p2[0])
	r1h, r1l := bits.Mul64(1174, p2[1])
	r2h, r2l := bits.Mul64(1174, p2[2])
	r3h, r3l := bits.Mul64(1174, p2[3])
	r1, carry := bits.Add64(r0h, r1l, 0)
	r2, carry := bits.Add64(r1h, r2l, carry)
	r3, carry := bits.Add64(r2h, r3l, carry)
	r4, _ := bits.Add64(r3h, 0, carry)

	r0, carry = bits.Add64(r0, (r4<<5|r3>>59)*9, 0)
	r1, carry = bits.Add64(r1, 0, carry)
	r2, carry = bits.Add64(r2, 0, carry)
	r3, _ = bits.Add64(r3&P3, 0, carry)

	r0, borrow := bits.Sub64(P0, r0, 0)
	r1, borrow = bits.Sub64(P1, r1, borrow)
	r2, borrow = bits.Sub64(P2, r2, borrow)
	r3, borrow = bits.Sub64(P3, r3, borrow)

	r0, borrow = bits.Sub64(r0, ^(borrow-1)&288, 0)
	res[1], borrow = bits.Sub64(r1, 0, borrow)
	res[2], borrow = bits.Sub64(r2, 0, borrow)
	res[3], borrow = bits.Sub64(r3, 0, borrow)

	res[0] = r0 - ^(borrow-1)&288

	return res
}

func mul2(res, p2 *FieldElement) {
	add(res, p2, p2)
}

func mod(res, p *FieldElement) {
	top := (p[3] >> 59) * 9
	r3 := p[3] & P3
	r0, carry := bits.Add64(p[0], top, 0)
	r1, carry := bits.Add64(p[1], 0, carry)
	r2, carry := bits.Add64(p[2], 0, carry)
	r3, _ = bits.Add64(r3, 0, carry)

	rr0, borrow := bits.Sub64(r0, P0, 0)
	rr1, borrow := bits.Sub64(r1, P1, borrow)
	rr2, borrow := bits.Sub64(r2, P2, borrow)
	rr3, borrow := bits.Sub64(r3, P3, borrow)

	b := ^(borrow - 1)

	res[0], carry = bits.Add64(rr0, P0&b, 0)
	res[1], carry = bits.Add64(rr1, b, carry)
	res[2], carry = bits.Add64(rr2, b, carry)
	res[3], _ = bits.Add64(rr3, P3&b, carry)
}

func selectPoint(res *Point, table *[16]Point, index uint64) {
	res.Set(&Point{X: UZero, Y: UZero, Z: UZero, T: UZero})
	for i := 0; i < 16; i++ {
		_, b1 := bits.Sub64(index, uint64(i), 0)
		_, b2 := bits.Sub64(uint64(i), index, 0)
		b1 = (b1 | b2) - 1
		for j := 0; j < 4; j++ {
			res.X[j] |= table[i].X[j] & b1
			res.Y[j] |= table[i].Y[j] & b1
			res.Z[j] |= table[i].Z[j] & b1
			res.T[j] |= table[i].T[j] & b1
		}
	}
}
