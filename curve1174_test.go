package curve1174

import (
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestAddingNeutralElement(t *testing.T) {
	var out Point
	out.Add(Base, E)
	out.ToAffine(&out)
	if !Base.Equals(&out) {
		t.Error("not equal", Base, out)
	}
	out.Add(E, E).Add(&out, E)
	out.ToAffine(&out)
	if !E.Equals(&out) {
		t.Error("not equal", Base, out)
	}
}

func TestSpecificDoubling(t *testing.T) {
	var d Point
	var f FieldElement
	f.Sub(&f, UOne)
	d.X.Set(UOne)
	d.Z.Set(UOne)
	d.Double(&d).ToAffine(&d)
	if !d.X.Equals(&UZero) || !d.Y.Equals(&f) {
		t.Error("not equal", d)
	}
}

func TestDoublingNeutralElement(t *testing.T) {
	var d1, d2 Point
	d1.Double(E).Double(&d1)
	d1.ToAffine(&d1)
	if !d1.Equals(E) {
		t.Error("not equal", d1, "\n", d2)
	}
}

func TestDoubling(t *testing.T) {
	var d1, d2 Point
	d1.Add(Base, Base).ToAffine(&d1)
	d2.Double(Base).ToAffine(&d2)
	if !d1.Equals(&d2) {
		t.Error("not equal", d1, "\n", d2)
	}
}

func TestSpecificScalarBaseMult(t *testing.T) {
	vectors := []string{"0", "1", "2", "3", "16", "17", "256", "257", "1050", "1000000006",
		"45073459087905739457391209340239423224234234457"}
	for _, vector := range vectors {
		b2, _ := new(big.Int).SetString(vector, 10)
		b := FromBigInt(b2)
		var p1, p2 Point
		p1.ScalarBaseMult(b).ToAffine(&p1)
		p2.ScalarMult(Base, b).ToAffine(&p2)
		if !p1.Equals(&p2) {
			t.Error(b, p1, p2)
		}
	}
}

func TestScalarBaseMult(t *testing.T) {
	randomTestOp(t, func(res *FieldElement, x *FieldElement, y *FieldElement) {
		var p Point
		p.ScalarBaseMult(x).ToAffine(&p)
		res.Set(&p.X)
	}, func(res *big.Int, x *big.Int, y *big.Int) {
		a := FromBigInt(x)
		var p Point
		p.ScalarMult(Base, a).ToAffine(&p)
		res.Set(p.X.ToBigInt())
	}, 10000)
}

func TestScalarMult(t *testing.T) {
	var b, mult Point
	b.Set(E)
	for i := 0; i < 5; i++ {
		b.Add(&b, Base)
	}
	mult.ScalarMult(Base, &FieldElement{5})
	mult.ToAffine(&mult)
	b.ToAffine(&b)
	if !mult.Equals(&b) {
		t.Error("not equal", mult, "\n", b)
	}
}

func TestAdd(t *testing.T) {
	randomTest(t, func(res *FieldElement, x *FieldElement, y *FieldElement) {
		add(res, x, y)
	}, func(res *big.Int, x *big.Int, y *big.Int) {
		res.Add(x, y)
	})
}

func TestSub(t *testing.T) {
	randomTest(t, func(res, x, y *FieldElement) {
		sub(res, x, y)
	}, func(res, x, y *big.Int) {
		res.Sub(x, y)
	})
}

func TestMul(t *testing.T) {
	randomTest(t, func(res *FieldElement, x *FieldElement, y *FieldElement) {
		mul(res, x, y)
	}, func(res *big.Int, x *big.Int, y *big.Int) {
		res.Mul(x, y)
	})
}

func TestMulD(t *testing.T) {
	randomTest(t, func(res *FieldElement, x *FieldElement, y *FieldElement) {
		mulD(res, x)
	}, func(res *big.Int, x *big.Int, y *big.Int) {
		res.Mul(x, big.NewInt(-1174))
	})
}

func TestSqr(t *testing.T) {
	randomTest(t, func(res *FieldElement, x *FieldElement, y *FieldElement) {
		sqr(res, x)
	}, func(res *big.Int, x *big.Int, y *big.Int) {
		res.Mul(x, x)
	})
}

func TestMul2(t *testing.T) {
	randomTest(t, func(res *FieldElement, x *FieldElement, y *FieldElement) {
		mul2(res, x)
	}, func(res *big.Int, x *big.Int, y *big.Int) {
		res.Mul(x, big.NewInt(2))
	})
}

func TestSpecificMulD(t *testing.T) {
	vectors := []string{
		"6ba08004fbbffffc023f61fffdc90000071872dfffcc9ae91bb81db9400170f",
		"7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF6",
	}
	for _, vector := range vectors {
		b2, _ := new(big.Int).SetString(vector, 16)
		b3 := new(big.Int).Mul(big.NewInt(-1174), b2)
		b3.Mod(b3, P)
		p2 := FromBigInt(b2)
		var res FieldElement
		mulD(&res, p2)
		b4 := res.ToBigInt()
		if b4.Cmp(b3) != 0 {
			t.Errorf("\n%x\n%x\n%x", b2, b3, b4)
		}
	}
}

func TestSpecificMul(t *testing.T) {
	vectors := [][]string{
		{"145f7ffb04400003fdc09e000236fffff8e78d2000336506e447e246bffe8e8",
			"6ba08004fbbffffc023f61fffdc90000071872dfffcc9ae91bb81db9400170f"},
		{"145f7ffb04400003fdc09e000236fffff8e78d2000336506e447e246bffe8e8",
			"6ba08004fbbffffc023f61fffdc90000071872dfffcc9af91bb81db9400170f"},
		{"7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF6",
			"7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF6"}}
	for _, vector := range vectors {
		b1, _ := new(big.Int).SetString(vector[0], 16)
		b2, _ := new(big.Int).SetString(vector[1], 16)
		b3 := new(big.Int).Mul(b1, b2)
		b3.Mod(b3, P)
		p := FromBigInt(b1)
		p2 := FromBigInt(b2)
		var res FieldElement
		mul(&res, p2, p)
		b4 := res.ToBigInt()
		if b4.Cmp(b3) != 0 {
			t.Errorf("\n%x\n%x\n%x\n%x", b1, b2, b3, b4)
		}
	}
}

func TestMod(t *testing.T) {
	randomTest(t, func(res *FieldElement, x *FieldElement, y *FieldElement) {
		mod(res, x)
	}, func(res *big.Int, x *big.Int, y *big.Int) {
		res.Mod(x, P)
	})
}

func TestSelect(t *testing.T) {
	var points [16]Point
	var res Point
	for i := 0; i < 16; i++ {
		u := uint64(i + 1)
		u = u + (u << 32)
		points[i].X = FieldElement{u, u, u, u}
		points[i].Y = FieldElement{0xCAFEBABE}
		points[i].Z = FieldElement{0xDEADBEEF}
		points[i].T = FieldElement{0xCAFEBABE}
	}
	for i := 0; i < 16; i++ {
		selectPoint(&res, &points, uint64(i))
		if !res.Equals(&points[i]) {
			t.Errorf("\n%x\n%x", res, points[i])
		}
	}
}

const TestsCount = 1000000

func randomTest(t *testing.T,
	fieldFunc func(res *FieldElement, x *FieldElement, y *FieldElement),
	bigIntFunc func(res *big.Int, x *big.Int, y *big.Int)) {
	randomTestOp(t, fieldFunc, bigIntFunc, TestsCount)
}

func randomTestOp(t *testing.T,
	fieldFunc func(res *FieldElement, x *FieldElement, y *FieldElement),
	bigIntFunc func(res *big.Int, x *big.Int, y *big.Int),
	testsCount int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fails := 0
	for i := 0; i < testsCount; i++ {
		b1 := new(big.Int).Rand(r, P)
		b2 := new(big.Int).Rand(r, P)
		p := FromBigInt(b1)
		p2 := FromBigInt(b2)
		var res FieldElement

		fieldFunc(&res, p, p2)
		b3 := new(big.Int)
		bigIntFunc(b3, b1, b2)
		two := big.NewInt(1)
		two.Lsh(two, 256)
		b3.Mod(b3, P)
		b4 := res.ToBigInt()
		if b4.Cmp(b3) != 0 {
			t.Errorf("\n%x\n%x\n%x\n%x", b1, b2, b3, b4)
			fails++
			if fails > 10 {
				break
			}
		}
	}
}

func TestInverse(t *testing.T) {
	randomTestOp(t, func(res *FieldElement, x *FieldElement, y *FieldElement) {
		res.Inverse(x)
	}, func(res *big.Int, x *big.Int, y *big.Int) {
		res.ModInverse(x, P)
	}, 10000)
}

func TestSqrMul(t *testing.T) {
	p := FieldElement{0x02f9052f8017dbde, 0xc2d6b27b4453a8ad, 0x51b0fcf0fa3f3df8, 0x2ca9df7680ba163b}

	var res1, res2 FieldElement

	mul(&res1, &p, &p)
	sqr(&res2, &p)

	if !res1.Equals(&res2) {
		t.Errorf("\n%x\n%x", res1, res2)
	}
}
