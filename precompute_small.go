//go:build !curve1174_no_precompute && !curve1174_precompute_big

package curve1174

var precomputedBase [64][16]Point

func init() {
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
}

func (p *Point) ScalarBaseMult(b *FieldElement) *Point {
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
