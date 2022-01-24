//go:build !amd64 || curve1174_purego

package curve1174

import (
	"crypto/subtle"
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

func sqr(res, p1 *FieldElement) {
	//_ = p1[3]
	//var r0, r1, r2, r3, r4, r5, r6, r7, carry uint64
	//
	//p01h, r1 := bits.Mul64(p1[0], p1[1])
	//p02h, p02l := bits.Mul64(p1[0], p1[2])
	//p03h, p03l := bits.Mul64(p1[0], p1[3])
	//p13h, p13l := bits.Mul64(p1[1], p1[3])
	//p23h, p23l := bits.Mul64(p1[2], p1[3])
	//
	//r2, carry = bits.Add64(p01h, p02l, 0)
	//r3, carry = bits.Add64(p02h, p03l, carry)
	//r4, carry = bits.Add64(p03h, p13l, carry)
	//r5, carry = bits.Add64(p13h, p23l, carry)
	//r6, carry = bits.Add64(p23h, 0, carry)
	//r7, _ = bits.Add64(0, 0, carry)
	//
	//p12h, p12l := bits.Mul64(p1[1], p1[2])
	//r3, carry = bits.Add64(r3, p12l, 0)
	//r4, carry = bits.Add64(r4, p12h, 0)
	//r5, carry = bits.Add64(r5, 0, carry)
	//r6, carry = bits.Add64(r6, 0, carry)
	//r7, _ = bits.Add64(r7, 0, carry)
	//
	//r7 = r7<<1 + r6>>63
	//r6 = r6<<1 + r5>>63
	//r5 = r5<<1 + r4>>63
	//r4 = r4<<1 + r3>>63
	//r3 = r3<<1 + r2>>63
	//r2 = r2<<1 + r1>>63
	//r1 = r1 << 1
	//
	//p00h, r0 := bits.Mul64(p1[0], p1[0])
	//p22h, p22l := bits.Mul64(p1[2], p1[2])
	//p33h, p33l := bits.Mul64(p1[3], p1[3])
	//p11h, p11l := bits.Mul64(p1[1], p1[1])
	//r1, carry = bits.Add64(r1, p00h, 0)
	//r2, carry = bits.Add64(r2, p11l, carry)
	//r3, carry = bits.Add64(r3, p11h, carry)
	//r4, carry = bits.Add64(r4, p22l, carry)
	//r5, carry = bits.Add64(r5, p22h, carry)
	//r6, carry = bits.Add64(r6, p33l, carry)
	//r7, _ = bits.Add64(r7, p33h, carry)
	//
	//r8 := r7 >> 59
	//r7 = r7<<5 | r6>>59
	//r6 = r6<<5 | r5>>59
	//r5 = r5<<5 | r4>>59
	//r4 = r4<<5 | r3>>59
	//
	//r3 &= P3
	//
	//r4a, carry := bits.Add64(r4, r4<<3, 0)
	//r5a, carry := bits.Add64(r5, r4>>61|r5<<3, carry)
	//r6a, carry := bits.Add64(r6, r5>>61|r6<<3, carry)
	//r7a, carry := bits.Add64(r7, r6>>61|r7<<3, carry)
	//r8a, _ := bits.Add64(r8, r7>>61|r8<<3, carry)
	//
	//r0, carry = bits.Add64(r0, r4a, 0)
	//r1, carry = bits.Add64(r1, r5a, carry)
	//r2, carry = bits.Add64(r2, r6a, carry)
	//r3, carry = bits.Add64(r3, r7a, carry)
	//r4, _ = bits.Add64(r8a, 0, carry)
	//
	//r4 = (r4<<5 + r3>>59) * 9
	//
	//r3 &= P3
	//
	//res[0], carry = bits.Add64(r0, r4, 0)
	//res[1], carry = bits.Add64(r1, 0, carry)
	//res[2], carry = bits.Add64(r2, 0, carry)
	//res[3], _ = bits.Add64(r3, 0, carry)
	mul(res, p1, p1)
}

func mul(res, p1, p2 *FieldElement) {
	_, _ = p1[3], p2[3]
	var r0, r1, r2, r3, r4, r5, r6, r7, carry uint64

	p00h, r0 := bits.Mul64(p1[0], p2[0])
	p01h, p01l := bits.Mul64(p1[0], p2[1])
	p02h, p02l := bits.Mul64(p1[0], p2[2])
	p03h, p03l := bits.Mul64(p1[0], p2[3])
	p10h, p10l := bits.Mul64(p1[1], p2[0])
	p11h, p11l := bits.Mul64(p1[1], p2[1])
	p12h, p12l := bits.Mul64(p1[1], p2[2])
	p13h, p13l := bits.Mul64(p1[1], p2[3])
	p20h, p20l := bits.Mul64(p1[2], p2[0])
	p21h, p21l := bits.Mul64(p1[2], p2[1])
	p22h, p22l := bits.Mul64(p1[2], p2[2])
	p23h, p23l := bits.Mul64(p1[2], p2[3])
	p30h, p30l := bits.Mul64(p1[3], p2[0])
	p31h, p31l := bits.Mul64(p1[3], p2[1])
	p32h, p32l := bits.Mul64(p1[3], p2[2])
	p33h, p33l := bits.Mul64(p1[3], p2[3])

	r1, carry = bits.Add64(p00h, p01l, 0)
	r2, carry = bits.Add64(p01h, p02l, carry)
	r3, carry = bits.Add64(p02h, p03l, carry)
	r4, carry = bits.Add64(p03h, p13l, carry)
	r5, carry = bits.Add64(p13h, p23l, carry)
	r6, carry = bits.Add64(p23h, p33l, carry)
	r7, _ = bits.Add64(p33h, 0, carry)

	r1, carry = bits.Add64(r1, p10l, 0)
	r2, carry = bits.Add64(r2, p10h, carry)
	r3, carry = bits.Add64(r3, p21l, carry)
	r4, carry = bits.Add64(r4, p21h, carry)
	r5, carry = bits.Add64(r5, p22h, carry)
	r6, carry = bits.Add64(r6, p32h, carry)
	r7, _ = bits.Add64(r7, 0, carry)

	r2, carry = bits.Add64(r2, p11l, 0)
	r3, carry = bits.Add64(r3, p11h, carry)
	r4, carry = bits.Add64(r4, p22l, carry)
	r5, carry = bits.Add64(r5, p32l, carry)
	r6, carry = bits.Add64(r6, 0, carry)
	r7, _ = bits.Add64(r7, 0, carry)

	r2, carry = bits.Add64(r2, p20l, 0)
	r3, carry = bits.Add64(r3, p12l, carry)
	r4, carry = bits.Add64(r4, p12h, carry)
	r5, carry = bits.Add64(r5, p31h, carry)
	r6, carry = bits.Add64(r6, 0, carry)
	r7, _ = bits.Add64(r7, 0, carry)

	r3, carry = bits.Add64(r3, p20h, 0)
	r4, carry = bits.Add64(r4, p30h, carry)
	r5, carry = bits.Add64(r5, 0, carry)
	r6, carry = bits.Add64(r6, 0, carry)
	r7, _ = bits.Add64(r7, 0, carry)

	r3, carry = bits.Add64(r3, p30l, 0)
	r4, carry = bits.Add64(r4, p31l, carry)
	r5, carry = bits.Add64(r5, 0, carry)
	r6, carry = bits.Add64(r6, 0, carry)
	r7, _ = bits.Add64(r7, 0, carry)

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

	r4 = (r4<<5 + r3>>59) * 9

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
	res.Set(&Point{})
	for i := 0; i < 16; i++ {
		b1 := ^(uint64(subtle.ConstantTimeEq(int32(index), int32(i))) - 1)
		res.X[0] |= table[i].X[0] & b1
		res.X[1] |= table[i].X[1] & b1
		res.X[2] |= table[i].X[2] & b1
		res.X[3] |= table[i].X[3] & b1
		res.Y[0] |= table[i].Y[0] & b1
		res.Y[1] |= table[i].Y[1] & b1
		res.Y[2] |= table[i].Y[2] & b1
		res.Y[3] |= table[i].Y[3] & b1
		res.T[0] |= table[i].T[0] & b1
		res.T[1] |= table[i].T[1] & b1
		res.T[2] |= table[i].T[2] & b1
		res.T[3] |= table[i].T[3] & b1
		res.Z[0] |= table[i].Z[0] & b1
		res.Z[1] |= table[i].Z[1] & b1
		res.Z[2] |= table[i].Z[2] & b1
		res.Z[3] |= table[i].Z[3] & b1
	}
}
