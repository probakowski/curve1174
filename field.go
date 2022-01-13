package curve1174

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"math/bits"
)

//P3 is 4th (highest) digit (base 2^64) of P=2^251-9
const P3 uint64 = 0x07ffffffffffffff

//P2 is 3rd digit (base 2^64) of P=2^251-9
const P2 uint64 = 0xffffffffffffffff

//P1 is 2nd digit (base 2^64) of P=2^251-9
const P1 uint64 = 0xffffffffffffffff

//P0 is 1st (lowest) digit (base 2^64) of P=2^251-9
const P0 uint64 = 0xfffffffffffffff7

//UOne represents 1
var UOne = &FieldElement{1}

//UP represents 2^251-9
var UP = &FieldElement{P0, P1, P2, P3}

//UZero represents 0
var UZero FieldElement

//FieldElement is element of finite field F_p, p=2^251-9
type FieldElement [4]uint64

//MulD multiplies field element by 1174 mod 2^251-9. Execution time doesn't depend on values
func (out *FieldElement) MulD(p *FieldElement) *FieldElement {
	mulD(out, p)
	return out
}

//Mul2 multiplies field element by 2 mod 2^251-9. Execution time doesn't depend on values
func (out *FieldElement) Mul2(p *FieldElement) *FieldElement {
	mul2(out, p)
	return out
}

//Mul multiplies two field elements mod 2^251-9. Execution time doesn't depend on values
func (out *FieldElement) Mul(p, p2 *FieldElement) *FieldElement {
	mul(out, p, p2)
	return out
}

//Sqr squares field element mod 2^251-9. Execution time doesn't depend on values
func (out *FieldElement) Sqr(p *FieldElement) *FieldElement {
	sqr(out, p)
	return out
}

//Add adds two field elements mod 2^251-9. Execution time doesn't depend on values
func (out *FieldElement) Add(p, p2 *FieldElement) *FieldElement {
	add(out, p, p2)
	return out
}

//Sub subtracts two field elements mod 2^251-9. Execution time doesn't depend on values
func (out *FieldElement) Sub(p, p2 *FieldElement) *FieldElement {
	sub(out, p, p2)
	return out
}

//Mod returns number mod 2^251-9. Execution time doesn't depend on values
func (out *FieldElement) Mod(p *FieldElement) *FieldElement {
	mod(out, p)
	return out
}

//Set sets one field element to be equal to the other. Execution time doesn't depend on values
func (out *FieldElement) Set(p2 *FieldElement) *FieldElement {
	out[0] = p2[0]
	out[1] = p2[1]
	out[2] = p2[2]
	out[3] = p2[3]
	return out
}

//Inverse sets out to be inverse of p2 mod 2^251-9 (out * p2 == 1 | 2^251-9). It uses Euler's theorem and computes
//inverse by raising p2 to power 2^251-11 (m=2^251-9 is prime, a^-1 == a^(m-2) | m). Execution time doesn't depend on value.
//Addition chain from https://github.com/mmcloughlin/addchain/blob/master/doc/results.md#curve1174-field-inversion (250sqr+13mul)
func (out *FieldElement) Inverse(p2 *FieldElement) *FieldElement {
	var x2, x3, x6, x7, x14, x15, x30, x60, x120, x240, x247 FieldElement
	x2.Sqr(p2).Mul(&x2, p2)
	x3.Sqr(&x2).Mul(&x3, p2)
	x6.sqrTimes(&x3, 3).Mul(&x6, &x3)
	x7.Sqr(&x6).Mul(&x7, p2)
	x14.sqrTimes(&x7, 7).Mul(&x14, &x7)
	x15.Sqr(&x14).Mul(&x15, p2)
	x30.sqrTimes(&x15, 15).Mul(&x30, &x15)
	x60.sqrTimes(&x30, 30).Mul(&x60, &x30)
	x120.sqrTimes(&x60, 60).Mul(&x120, &x60)
	x240.sqrTimes(&x120, 120).Mul(&x240, &x120)
	x247.sqrTimes(&x240, 7).Mul(&x247, &x7)
	return out.sqrTimes(&x247, 2).Mul(out, p2).sqrTimes(out, 2).Mul(out, p2)
}

func (out *FieldElement) sqrTimes(p *FieldElement, n int) *FieldElement {
	out.Set(p)
	for i := 0; i < n; i++ {
		out.Sqr(out)
	}
	return out
}

//Equals checks if 2 field elements has the same value
func (out *FieldElement) Equals(p2 *FieldElement) bool {
	var f FieldElement
	return f.Sub(out, p2).IsZero()
}

//ToBigInt returns element value as *big.Int
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

//SetBigInt sets field elements to value from *big.Int.
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

//FromBigInt returns field element with the same value as provided big.Int
func FromBigInt(b1 *big.Int) *FieldElement {
	var p FieldElement
	return p.SetBigInt(b1)
}

func (out *FieldElement) String() string {
	return fmt.Sprintf("%x", out)
}

//Format to implement fmt.Formatter interface
func (out *FieldElement) Format(s fmt.State, c rune) {
	out.ToBigInt().Format(s, c)
}

//IsOne checks if out==1
func (out *FieldElement) IsOne() bool {
	return out[0] == 1 && out[1]|out[2]|out[3] == 0
}

//IsZero checks if out==0
func (out *FieldElement) IsZero() bool {
	var o FieldElement
	o.Mod(out)
	return o[0]|o[1]|o[2]|o[3] == 0
}

//Cmp compares 2 field elements. Can be used for sorting
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
