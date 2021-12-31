//go:build !curve1174_no_precompute && curve1174_precompute_big

package curve1174

var precomputedBase [32][256]Point

func init() {
	var p Point
	p.Set(UBase)
	sp := &p
	for i := 0; i < 32; i++ {
		el := &precomputedBase[i]
		el[0].Set(UE)
		el[1].Set(sp)
		for i := 2; i < 255; i += 2 {
			el[i].Double(&el[i-1]).ToAffine(&el[i])
			el[i+1].AddZ1(&el[i], sp).ToAffine(&el[i+1])
		}
		sp.Double(&el[128]).ToAffine(sp)
	}
}

//ScalarBaseMult multiplies base point UBase by scalar b (b<2^251-9) and stores result in p. Execution time doesn't depend on b.
func (p *Point) ScalarBaseMult(b *FieldElement) *Point {
	index := b[0] & 0xFF
	p.Set(&precomputedBase[0][index])

	for i := 1; i < 32; i++ {
		index = (b[i/256] >> ((i % 256) * 8)) & 0xFF
		p.AddZ1(p, &precomputedBase[i][index])
	}

	return p
}
