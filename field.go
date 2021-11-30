package curve1174

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"math/bits"
	"unsafe"
)

const P3 uint64 = 0x07ffffffffffffff
const P2 uint64 = 0xffffffffffffffff
const P1 uint64 = 0xffffffffffffffff
const P0 uint64 = 0xfffffffffffffff7

var UOne = &FieldElement{1}
var UD = &FieldElement{0xfffffffffffffb61, P1, P2, P3}
var UP = &FieldElement{P0, P1, P2, P3}
var UZero FieldElement

type FieldElement [4]uint64

func (out *FieldElement) MulD(p *FieldElement) *FieldElement {
	mulD(out, p)
	return out
}

func (out *FieldElement) Mul2(p *FieldElement) *FieldElement {
	mul2(out, p)
	return out
}

func (out *FieldElement) Mul(p, p2 *FieldElement) *FieldElement {
	mul(uintptr(unsafe.Pointer(out)), uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(p2)))
	return out
}

func (out *FieldElement) Sqr(p *FieldElement) *FieldElement {
	sqr(out, p)
	return out
}

func (out *FieldElement) Add(p, p2 *FieldElement) *FieldElement {
	add(out, p, p2)
	return out
}

func (out *FieldElement) Sub(p, p2 *FieldElement) *FieldElement {
	sub(out, p, p2)
	return out
}

func (out *FieldElement) Mod(p *FieldElement) *FieldElement {
	mod(out, p)
	return out
}

func (out *FieldElement) Set(p2 *FieldElement) *FieldElement {
	out[0] = p2[0]
	out[1] = p2[1]
	out[2] = p2[2]
	out[3] = p2[3]
	return out
}

func (out *FieldElement) IsEven() bool {
	return out[0]&1 == 0
}

func (out *FieldElement) subNoMod(p, p2 *FieldElement) *FieldElement {
	var borrow uint64
	out[0], borrow = bits.Sub64(p[0], p2[0], 0)
	out[1], borrow = bits.Sub64(p[1], p2[1], borrow)
	out[2], borrow = bits.Sub64(p[2], p2[2], borrow)
	out[3], _ = bits.Sub64(p[3], p2[3], borrow)

	return out
}

func (out *FieldElement) subP(p *FieldElement) *FieldElement {
	var borrow uint64
	out[0], borrow = bits.Sub64(p[0], P0, 0)
	out[1], borrow = bits.Sub64(p[1], P1, borrow)
	out[2], borrow = bits.Sub64(p[2], P2, borrow)
	out[3], _ = bits.Sub64(p[3], P3, borrow)

	return out
}

func (out *FieldElement) addP(p *FieldElement) *FieldElement {
	var carry uint64
	out[0], carry = bits.Add64(p[0], P0, 0)
	out[1], carry = bits.Add64(p[1], P1, carry)
	out[2], carry = bits.Add64(p[2], P2, carry)
	out[3], _ = bits.Add64(p[3], P3, carry)
	return out
}

func (out *FieldElement) FastInverse(p2 *FieldElement) *FieldElement {
	var u, v, b FieldElement
	d := FieldElement{1}
	u.Set(UP)
	v.Set(p2)
	for !u.IsZero() {
		for u.IsEven() {
			if !b.IsEven() {
				b.subP(&b)

			}
			div2(&b, &b)
			div2(&u, &u)
		}
		for v.IsEven() {
			if !d.IsEven() {
				d.subP(&d)
			}
			div2(&d, &d)
			div2(&v, &v)
		}
		if u.Cmp(&v) >= 0 {
			u.subNoMod(&u, &v)
			b.subNoMod(&b, &d)
		} else {
			v.subNoMod(&v, &u)
			d.subNoMod(&d, &b)
		}
	}
	if d[3]&(1<<63) == 0 {
		out.Set(&d)
	} else {
		out.addP(&d)
	}
	//fastInverse(out, p2)
	return out
}

func (out *FieldElement) Inverse(p2 *FieldElement) *FieldElement {
	var p3, pF, p2F, p4F, p8F, p16F, p32F FieldElement
	p3.Sqr(p2).Mul(&p3, p2)
	pF.sqrTimes(&p3, 2).Mul(&pF, &p3)
	p2F.sqrTimes(&pF, 4).Mul(&p2F, &pF)
	p4F.sqrTimes(&p2F, 8).Mul(&p4F, &p2F)
	p8F.sqrTimes(&p4F, 16).Mul(&p8F, &p4F)
	p16F.sqrTimes(&p8F, 32).Mul(&p16F, &p8F)
	p32F.sqrTimes(&p16F, 64).Mul(&p32F, &p16F)
	out.sqrTimes(&p32F, 64).Mul(out, &p16F)
	out.sqrTimes(out, 32).Mul(out, &p8F)
	out.sqrTimes(out, 16).Mul(out, &p4F)
	out.sqrTimes(out, 4).Mul(out, &pF)
	out.sqrTimes(out, 2).Mul(out, &p3)
	out.Sqr(out).Mul(out, p2)
	out.sqrTimes(out, 2).Mul(out, p2)
	out.sqrTimes(out, 2)
	out.Mul(out, p2)

	return out
}

func (out *FieldElement) sqrTimes(p *FieldElement, n int) *FieldElement {
	out.Set(p)
	for i := 0; i < n; i++ {
		out.Sqr(out)
	}
	return out
}

func (out *FieldElement) Equals(p2 *FieldElement) bool {
	return out[0] == p2[0] && out[1] == p2[1] && out[2] == p2[2] && out[3] == p2[3]
}

func (out *FieldElement) ToBigInt() *big.Int {
	l := len(out)
	var b [64]byte
	bytes := b[:l*8]
	for i := 0; i < l; i++ {
		binary.BigEndian.PutUint64(bytes[i*8:], out[l-i-1])
	}
	b4 := new(big.Int).SetBytes(bytes)
	return b4
}

func (out *FieldElement) SetBigInt(b1 *big.Int) *FieldElement {
	b := b1.Bits()
	for i := 0; i < len(out); i++ {
		out[i] = 0
	}
	if bits.UintSize == 64 {
		for i := 0; i < len(b) && i < len(out); i++ {
			out[i] = uint64(b[i])
		}
	} else {
		for i := 0; i < len(b) && i < len(out)*2; i++ {
			out[i/2] |= uint64(b[i]) << ((i / 2) * 32)
		}
	}
	return out
}

func FromBigInt(b1 *big.Int) *FieldElement {
	var p FieldElement
	return p.SetBigInt(b1)
}

func (out FieldElement) String() string {
	return fmt.Sprintf("%x%016x%016x%016x", out[3], out[2], out[1], out[0])
}

func (out FieldElement) Format(s fmt.State, c rune) {
	out.ToBigInt().Format(s, c)
}

func (out *FieldElement) IsOne() bool {
	return out[0] == 1 && out[1]|out[2]|out[3] == 0
}

func (out *FieldElement) IsZero() bool {
	return out[0]|out[1]|out[2]|out[3] == 0
}

func (out *FieldElement) Cmp(p2 *FieldElement) int {
	for i := 3; i >= 0; i-- {
		if out[i] > p2[i] {
			return 1
		} else if out[i] < p2[i] {
			return -1
		}
	}
	return 0
}