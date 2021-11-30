package curve1174

import (
	"fmt"
	"math/big"
)

var P, _ = new(big.Int).SetString("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7", 16)

var UBase = &Point{
	X: FieldElement{0x16123f27bce29eda, 0xc021d96a492ecd65, 0x9343aee7c029a190, 0x37fbb0cea308c47},
	Y: FieldElement{0xa4ccb1bf9b46360e, 0x4fe2dee2af3f976b, 0x6656841169840e0c, 0x6b72f82d47fb7cc},
	Z: *UOne,
	T: FieldElement{0xfb1ebfece06620ec, 0x9c6c6daf574e84cb, 0x5083299c2d40b958, 0x18b74129cf1e5d9},
}

var UE = &Point{
	X: UZero,
	Y: *UOne,
	Z: *UOne,
	T: UZero,
}

type Point struct {
	X FieldElement
	Y FieldElement
	Z FieldElement
	T FieldElement
}

func (p *Point) Set(p2 *Point) *Point {
	p.X.Set(&p2.X)
	p.Y.Set(&p2.Y)
	p.Z.Set(&p2.Z)
	p.T.Set(&p2.T)
	return p
}

func (p Point) String() string {
	return fmt.Sprintf("Curve1174 Point\nX:%x\nY:%x\nZ:%x\nT:%x", p.X, p.Y, p.Z, p.T)
}

func (p *Point) ToAffine(pp *Point) *Point {
	//var inv FieldElement
	//inv.FastInverse(&pp.Z)
	//zInv := &inv
	//p.X.Mul(&pp.X, zInv)
	//p.Y.Mul(&pp.Y, zInv)
	//p.T.Mul(&pp.T, zInv)
	//p.Z.Set(UOne)
	z := pp.Z.ToBigInt()
	inverse := z.ModInverse(z, P)
	if inverse != nil {
		zInv := FromBigInt(inverse)
		p.X.Mul(&pp.X, zInv)
		p.Y.Mul(&pp.Y, zInv)
		p.T.Mul(&pp.T, zInv)
		p.Z.Set(UOne)
	}
	return p
}

func (p *Point) Equals(p2 *Point) bool {
	return p.Z.Equals(&p2.Z) && p.Y.Equals(&p2.Y) && p.X.Equals(&p2.X)
}

var precomputedBase [64][16]Point
var precomputed bool

func PrecomputeBase() {
	if precomputed {
		return
	}
	var p Point
	p.Set(UBase)
	sp := &p
	for i := 0; i < 64; i++ {
		el := &precomputedBase[i]
		el[0].Set(UE)
		el[1].Set(sp)
		el[2].Double(sp)
		el[3].Add(&el[2], sp)
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
		for j := 0; j < 16; j++ {
			el[j].ToAffine(&el[j])
		}
		sp.Double(&el[8]).ToAffine(sp)
	}

	precomputed = true
}

var precomputedBase2 [32][256]Point
var precomputed2 bool

func PrecomputeBase2() {
	if precomputed2 {
		return
	}
	var p Point
	p.Set(UBase)
	sp := &p
	for i := 0; i < 32; i++ {
		el := &precomputedBase2[i]
		el[0].Set(UE)
		el[1].Set(sp)
		for i := 2; i < 255; i += 2 {
			el[i].Double(&el[i-1]).ToAffine(&el[i])
			el[i+1].AddZ1(&el[i], sp).ToAffine(&el[i+1])
		}
		sp.Double(&el[128]).ToAffine(sp)
	}
}

func (p *Point) ScalarBaseMult2(b *FieldElement) *Point {
	index := b[0] & 0xFF
	p.Set(&precomputedBase2[0][index])

	for i := 1; i < 32; i++ {
		index = (b[i/256] >> ((i % 256) * 8)) & 0xFF
		p.AddZ1(p, &precomputedBase2[i][index])
	}

	return p
}

func (p *Point) ScalarBaseMult(b *FieldElement) *Point {
	if !precomputed {
		return p.ScalarMult(UBase, b)
	}

	index := b[0] & 0xF
	selectPoint(p, &precomputedBase[0], index)
	var pp Point

	for i := 1; i < 64; i++ {
		index = (b[i/16] >> ((i % 16) * 4)) & 0xF
		selectPoint(&pp, &precomputedBase[i], index)
		p.AddZ1(p, &pp)
	}

	return p
}

func (p *Point) ScalarMult(sp *Point, b *FieldElement) *Point {
	el := [16]Point{*UE, *sp}
	p.Set(UE)
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

func (p *Point) Add(p1, p2 *Point) *Point {
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

func (p *Point) addToProjective(p1, p2 *Point) *Point {
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
