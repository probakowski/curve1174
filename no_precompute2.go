//go:build curve1174_no_precompute

package curve1174

//ScalarBaseMult multiplies base point Base by scalar b (b<2^251-9) and stores result in p. Execution time doesn't depend on b.
//If there are precomputed tables (PrecomputeBase, PrecomputeBase2) they will be used for speedup.
func (p *Point) ScalarBaseMult(b *FieldElement) *Point {
	return p.ScalarMult(Base, b)
}
