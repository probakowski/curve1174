package curve1174

import (
	"fmt"
	"math/big"
)

//P is order of F_p, 2^251-9
var P, _ = new(big.Int).SetString("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7", 16)

//Base is base point of curve in affine coordinates (Base.Z == 1)
var Base = &Point{
	X: FieldElement{0x16123f27bce29eda, 0xc021d96a492ecd65, 0x9343aee7c029a190, 0x37fbb0cea308c47},
	Y: FieldElement{0xa4ccb1bf9b46360e, 0x4fe2dee2af3f976b, 0x6656841169840e0c, 0x6b72f82d47fb7cc},
	Z: *UOne,
	T: FieldElement{0xfb1ebfece06620ec, 0x9c6c6daf574e84cb, 0x5083299c2d40b958, 0x18b74129cf1e5d9},
}

//E is identity element of curve's group (x:0, y:1)
var E = &Point{
	X: UZero,
	Y: *UOne,
	Z: *UOne,
	T: UZero,
}

//Point represents point on curve. It supports projective and extended coordinates
type Point struct {
	X FieldElement
	Y FieldElement
	Z FieldElement
	T FieldElement
}

//Set p to be exactly the same as p2
func (p *Point) Set(p2 *Point) *Point {
	p.X.Set(&p2.X)
	p.Y.Set(&p2.Y)
	p.Z.Set(&p2.Z)
	p.T.Set(&p2.T)
	return p
}

//Format to implement fmt.Formatter interface
func (p *Point) Format(s fmt.State, c rune) {
	switch c {
	case 'v':
		if s.Flag('#') {
			_, _ = fmt.Fprintf(s, "%T", p)
		}
		if s.Flag('#') || s.Flag('+') {
			_, _ = fmt.Fprintf(s, "{X:%x Y:%x Z:%x T:%x}", p.X, p.Y, p.Z, p.T)
		} else {
			_, _ = fmt.Fprintf(s, "{%x %x %x %x}", p.X, p.Y, p.Z, p.T)
		}
	case 'x', 'X', 'd', 'o', 'O', 'b':
		_, _ = fmt.Fprintf(s, "{")
		p.X.ToBigInt().Format(s, c)
		_, _ = fmt.Fprintf(s, " ")
		p.Y.ToBigInt().Format(s, c)
		_, _ = fmt.Fprintf(s, " ")
		p.Z.ToBigInt().Format(s, c)
		_, _ = fmt.Fprintf(s, " ")
		p.T.ToBigInt().Format(s, c)
		_, _ = fmt.Fprintf(s, "}")
	}
}

func (p *Point) String() string {
	return fmt.Sprintf("%#v", p)
}

//ToAffine transforms point pp to affine coordinates and store result in p (p.Z==1)
func (p *Point) ToAffine(pp *Point) *Point {
	var inv FieldElement
	inv.Inverse(&pp.Z)
	zInv := &inv
	p.X.Mul(&pp.X, zInv).Mod(&p.X)
	p.Y.Mul(&pp.Y, zInv).Mod(&p.Y)
	p.T.Mul(&pp.T, zInv).Mod(&p.T)
	p.Z.Set(UOne)
	return p
}

//Equals checks if two points have exactly the same representation (all components must be equal)
func (p *Point) Equals(p2 *Point) bool {
	return p.Z.Equals(&p2.Z) && p.Y.Equals(&p2.Y) && p.X.Equals(&p2.X)
}

//ScalarMult multiplies point on curve sp by scalar b (b<2^251-9) and stores result in p. Execution time doesn't depend on b.
func (p *Point) ScalarMult(sp *Point,
	b *FieldElement) *Point {
	el := [16]Point{*E, *sp}
	p.Set(E)
	el[2].Double(sp)
	el[3].Add(&el[2], &el[1])
	el[4].Double(&el[2])
	el[5].Add(&el[4], sp)
	el[6].Double(&el[3])
	el[7].Add(&el[6], sp)
	el[8].Double(&el[4])
	el[9].Add(&el[8], sp)
	el[10].Double(&el[5])
	el[11].Add(&el[10], sp)
	el[12].Double(&el[6])
	el[13].Add(&el[12], sp)
	el[14].Double(&el[7])
	el[15].Add(&el[14], sp)

	index := b[3] >> 60
	selectPoint(p, &el, index)
	var pp Point

	for j := 14; j >= 0; j-- {
		p.doubleProjective(p).doubleProjective(p).doubleProjective(p).Double(p)
		index = (b[3] >> (j * 4)) & 0xF
		selectPoint(&pp, &el, index)
		p.addToProjective(p, &pp)
	}

	for i := 2; i > 0; i-- {
		for j := 15; j >= 0; j-- {
			p.doubleProjective(p).doubleProjective(p).doubleProjective(p).Double(p)
			index = (b[i] >> (j * 4)) & 0xF
			selectPoint(&pp, &el, index)
			p.addToProjective(p, &pp)
		}
	}

	for j := 15; j > 0; j-- {
		p.doubleProjective(p).doubleProjective(p).doubleProjective(p).Double(p)
		index := (b[0] >> (j * 4)) & 0xF
		selectPoint(&pp, &el, index)
		p.addToProjective(p, &pp)
	}

	p.doubleProjective(p).doubleProjective(p).doubleProjective(p).Double(p)
	index = b[0] & 0xF
	selectPoint(&pp, &el, index)
	p.Add(p, &pp)

	return p
}

//AddZ1 adds two points on curve and store results in p. p2 has to be in affine coordinates (p2.Z == 1)
//Formula based on https://www.hyperelliptic.org/EFD/g1p/auto-twisted-extended.html#addition-madd-2008-hwcd
func (p *Point) AddZ1(p1, p2 *Point) *Point {
	var a, b, c, d, e, e1, f, g, h FieldElement
	a.Mul(&p1.X, &p2.X)
	b.Mul(&p1.Y, &p2.Y)
	c.Mul(&p1.T, &p2.T).MulD(&c)
	d.Set(&p1.Z)
	e1.Add(&p2.X, &p2.Y)
	e.Add(&p1.X, &p1.Y).Mul(&e, &e1).Sub(&e, &a).Sub(&e, &b)
	f.Sub(&d, &c)
	g.Add(&d, &c)
	h.Sub(&b, &a)
	p.X.Mul(&e, &f)
	p.Y.Mul(&g, &h)
	p.T.Mul(&e, &h)
	p.Z.Mul(&f, &g)
	return p
}

//Add adds any two points on curve and store results in p.
//Formula based on https://www.hyperelliptic.org/EFD/g1p/auto-twisted-extended.html#addition-add-2008-hwcd
func (p *Point) Add(p1,
	p2 *Point) *Point {
	var a, b, c, d, e, e1, f, g, h FieldElement
	a.Mul(&p1.X, &p2.X)
	b.Mul(&p1.Y, &p2.Y)
	c.Mul(&p1.T, &p2.T).MulD(&c)
	d.Mul(&p1.Z, &p2.Z)
	e1.Add(&p2.X, &p2.Y)
	e.Add(&p1.X, &p1.Y).Mul(&e, &e1).Sub(&e, &a).Sub(&e, &b)
	f.Sub(&d, &c)
	g.Add(&d, &c)
	h.Sub(&b, &a)
	p.X.Mul(&e, &f)
	p.Y.Mul(&g, &h)
	p.T.Mul(&e, &h)
	p.Z.Mul(&f, &g)
	return p
}

//addToProjective adds two points on curve store result in p. Result is in projective coordinates (p.T is not correct!)
func (p *Point) addToProjective(p1,
	p2 *Point) *Point {
	var a, b, c, d, e, e1, f, g, h FieldElement
	a.Mul(&p1.X, &p2.X)
	b.Mul(&p1.Y, &p2.Y)
	c.Mul(&p1.T, &p2.T).MulD(&c)
	d.Mul(&p1.Z, &p2.Z)
	e1.Add(&p2.X, &p2.Y)
	e.Add(&p1.X, &p1.Y).Mul(&e, &e1).Sub(&e, &a).Sub(&e, &b)
	f.Sub(&d, &c)
	g.Add(&d, &c)
	h.Sub(&b, &a)
	p.X.Mul(&e, &f)
	p.Y.Mul(&g, &h)
	p.Z.Mul(&f, &g)
	return p
}

//doubleProjective doubles point on curve using projective coordinates and store result in p (p.T is not correct!)
func (p *Point) doubleProjective(dp *Point) *Point {
	var b, c, d, f, h, j FieldElement
	b.Add(&dp.X, &dp.Y).Sqr(&b)
	c.Sqr(&dp.X)
	d.Sqr(&dp.Y)
	f.Add(&c, &d)
	h.Sqr(&dp.Z)
	j.Mul2(&h).Sub(&f, &j)
	p.X.Sub(&b, &c).Sub(&p.X, &d).Mul(&p.X, &j)
	p.Y.Sub(&c, &d).Mul(&p.Y, &f)
	p.Z.Mul(&f, &j)
	return p
}

//Double doubles point on curve and store result in p (p = dp+dp)
//Formula based on https://www.hyperelliptic.org/EFD/g1p/auto-twisted-extended.html#doubling-dbl-2008-hwcd
func (p *Point) Double(dp *Point) *Point {
	var a, b, c, e, f, g, h FieldElement
	a.Sqr(&dp.X)
	b.Sqr(&dp.Y)
	c.Sqr(&dp.Z).Mul2(&c)
	e.Add(&dp.X, &dp.Y).Sqr(&e).Sub(&e, &a).Sub(&e, &b)
	g.Add(&a, &b)
	f.Sub(&g, &c)
	h.Sub(&a, &b)
	p.X.Mul(&e, &f)
	p.Y.Mul(&g, &h)
	p.T.Mul(&e, &h)
	p.Z.Mul(&f, &g)
	return p
}

//DoubleZ1 doubles point on curve and store result in p (p = dp+dp). dp has to be in affine coordinates (dp.Z == 1)
//Formula based on https://www.hyperelliptic.org/EFD/g1p/auto-twisted-extended.html#doubling-mdbl-2008-hwcd
func (p *Point) DoubleZ1(dp *Point) *Point {
	var a, b, c, d, e, f, g, h FieldElement
	a.Sqr(&dp.X)
	b.Sqr(&dp.Y)
	d = a
	e.Add(&dp.X, &dp.Y).Sqr(&e).Sub(&e, &a).Sub(&e, &b)
	g.Add(&d, &b)
	f.Sub(&g, &FieldElement{2, 0, 0, 0})
	c.Sqr(&g)
	h.Sub(&d, &b)
	p.X.Mul(&e, &f)
	p.Y.Mul(&g, &h)
	p.T.Mul(&e, &h)
	h.Mul2(&g)
	p.Z.Sub(&c, &h)
	return p
}
