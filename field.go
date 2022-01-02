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
//inverse by raising p2 to power 2^251-11 (m=2^251-9 is prime, a^-1 == a^(m-2) | m). Execution time doesn't depend on value
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

//Equals checks if 2 field elements has the same value
func (out *FieldElement) Equals(p2 *FieldElement) bool {
	return out[0] == p2[0] && out[1] == p2[1] && out[2] == p2[2] && out[3] == p2[3]
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
	return out[0]|out[1]|out[2]|out[3] == 0
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
