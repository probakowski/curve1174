//go:generate go run asm.go -out ../field_amd64.s -pkg curve1174

package main

import (
	"fmt"
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

var regs [8]Register
var yPtr, xPtr, resPtr, zero, flag Register
var storedLo, storedHi Register
var storedX, storedY = -1, -1
var lastX int

func main() {
	Package("przemko-robakowski.pl/curve1174")
	ConstraintExpr("!purego")

	for i := 0; i < 8; i++ {
		regs[i] = GP64()
	}

	zero = GP64()

	mulFunc()
	mulNoAdxFunc()
	modFunc()
	mul2Func()
	div2Func()
	subFunc()
	addFunc()
	mulDFunc()
	selectFunc()
	fastInverse()
	sqrFunc()

	Generate()
}

func fastInverse() {
	TEXT("fastInverse", NOSPLIT, "func(res, x *FieldElement)")
	Pragma("noescape")
	v1 := []Register{R8, R9, R10, R11}
	v := []Op{v1[0], v1[1], v1[2], v1[3]}
	xPtr = Load(Param("x"), GP64())
	x := Param("x").Dereference(xPtr)
	for i := 0; i < 4; i++ {
		Load(x.Index(i), v1[i])
	}
	local := AllocLocal(8)
	MOVQ(RSP, X0)
	u := []Op{RAX, RBX, RCX, RDX}
	d1 := []Register{R12, R13, R14, R15}
	d := []Op{d1[0], d1[1], d1[2], d1[3]}
	b := []Op{local, RSI, RDI, RBP}

	MOVQ(Imm(0xfffffffffffffff7), u[0])
	MOVQ(Imm(0xffffffffffffffff), u[1])
	MOVQ(Imm(0xffffffffffffffff), u[2])
	MOVQ(Imm(0x07ffffffffffffff), u[3])

	//MOVQ(U32(7), u[0])
	//MOVQ(U32(0), u[1])
	//MOVQ(U32(0), u[2])
	//MOVQ(U32(0), u[3])

	for i := 0; i < 4; i++ {
		MOVQ(U32(0), b[i])
		XORQ(d[i], d[i])
	}

	MOVQ(U32(1), d[0])

	Label("mainloop")
	{
		for i := 0; i < 3; i++ {
			CMPQ(u[i], Imm(0))
			JNZ(LabelRef("notzero"))
		}
		CMPQ(u[3], Imm(0))
		JZ(LabelRef("aftermain"))

		Label("notzero")
		TESTQ(U32(1), u[0])
		JNZ(LabelRef("afteru"))
		Label("uloop")
		{
			TESTQ(U32(1), b[0])
			JZ(LabelRef("divu"))

			subNN(b)

			Label("divu")
			div2(b)
			div2(u)
		}
		TESTQ(U32(1), u[0])
		JZ(LabelRef("uloop"))
		Label("afteru")

		TESTQ(U32(1), v[0])
		JNZ(LabelRef("afterv"))
		Label("vloop")
		{
			TESTQ(U32(1), d[0])
			JZ(LabelRef("divv"))

			subNN(d)

			Label("divv")
			div2(d)
			div2(v)
		}
		TESTQ(U32(1), v[0])
		JZ(LabelRef("vloop"))
		Label("afterv")
		for i := 3; i >= 0; i-- {
			CMPQ(u[i], v[i])
			JB(LabelRef("greaterv"))
			JNE(LabelRef("greateru"))
		}
		Label("greateru")
		subNoMod(u, v)
		subNoMod(b, d)
		JMP(LabelRef("mainloop"))
		Label("greaterv")
		subNoMod(v, u)
		subNoMod(d, b)
		JMP(LabelRef("mainloop"))
	}
	Label("aftermain")

	MOVQ(X0, RSP)

	resPtr = Load(Param("res"), GP64())
	res := Param("res").Dereference(resPtr)
	CMPQ(d[3], Imm(0))
	JNS(LabelRef("store"))

	ADDQ(U8(0xf7), d[0])
	ADCQ(U8(0xff), d[1])
	ADCQ(U8(0xff), d[2])
	MOVQ(Imm(0x07ffffffffffffff), u[0])
	ADCQ(u[0], d[3])

	//ADDQ(U8(0x7), d[0])
	//ADCQ(U8(0x0), d[1])
	//ADCQ(U8(0x0), d[2])
	//ADCQ(U8(0x0), d[3])

	Label("store")
	for i := 0; i < 4; i++ {
		Store(d1[i], res.Index(i))
	}

	RET()
}

func subNoMod(u []Op, v []Op) {
	SUBQ(v[0], u[0])
	for i := 1; i < 4; i++ {
		SBBQ(v[i], u[i])
	}
}

func subNN(r []Op) {
	SUBQ(U8(0xf7), r[0])
	SBBQ(U8(0xff), r[1])
	SBBQ(U8(0xff), r[2])
	MOVQ(r[0], RSI)
	MOVQ(Imm(0x07ffffffffffffff), RSI)
	SBBQ(RSI, r[3])
	MOVQ(X1, RSI)
	//SUBQ(U8(0x7), r[0])
	//SBBQ(U8(0x0), r[1])
	//SBBQ(U8(0x0), r[2])
	//SBBQ(U8(0x0), r[3])
}

func div2(r []Op) {
	for i := 0; i < 3; i++ {
		SHRQ(Imm(1), r[i+1], r[i])
	}
	SARQ(Imm(1), r[3])
}

func selectFunc() {
	TEXT("selectPoint", NOSPLIT, "func(res *Point, table *[16]Point, index uint64)")
	Pragma("noescape")
	targetIndex := Load(Param("index"), XMM())
	xPtr = Load(Param("table"), GP64())
	resPtr = Load(Param("res"), GP64())
	var res [8]Register
	for i := 0; i < 8; i++ {
		res[i] = XMM()
	}
	currentIndex, k, one := XMM(), XMM(), XMM()
	index := GP64()
	PSHUFD(Imm(0), targetIndex, targetIndex)
	for i := 0; i < 8; i++ {
		PXOR(res[i], res[i])
	}
	MOVQ(U64(16), index)
	PCMPEQL(currentIndex, currentIndex)
	PXOR(one, one)
	PSUBL(currentIndex, one)
	PXOR(currentIndex, currentIndex)
	Label("loop")
	MOVO(currentIndex, k)
	PCMPEQL(targetIndex, k)

	for j := 0; j < 8; j++ {
		v := XMM()
		MOVOU(Mem{Base: xPtr, Disp: j * 16}, v)
		PAND(k, v)
		POR(v, res[j])
	}

	ADDQ(Imm(128), xPtr)
	PADDL(one, currentIndex)
	SUBQ(Imm(1), index)
	JNZ(LabelRef("loop"))
	for i := 0; i < 8; i++ {
		MOVOU(res[i], Mem{Base: resPtr, Disp: i * 16})
	}
	RET()
}

func addFunc() {
	TEXT("add", NOSPLIT, "func(res, x, y *FieldElement)")
	Pragma("noescape")
	Doc("res=x + y % 2^251-9")
	xPtr = Load(Param("x"), GP64())
	yPtr = Load(Param("y"), GP64())
	resPtr = Load(Param("res"), GP64())

	for i := 0; i < 4; i++ {
		MOVQ(mem(xPtr, i), regs[i])
	}

	ADDQ(mem(yPtr, 0), regs[0])
	for i := 1; i < 4; i++ {
		ADCQ(mem(yPtr, i), regs[i])
	}

	subN()

	storeResults()
	RET()
}

func subFunc() {
	TEXT("sub", NOSPLIT, "func(res, x, y *FieldElement)")
	Pragma("noescape")
	Doc("res=x - y % 2^251-9")
	xPtr = Load(Param("x"), GP64())
	yPtr = Load(Param("y"), GP64())

	for i := 0; i < 4; i++ {
		MOVQ(mem(xPtr, i), regs[i])
	}

	b := GP64()
	n0 := GP64()
	n3 := GP64()
	XORQ(b, b)

	SUBQ(mem(yPtr, 0), regs[0])
	for i := 1; i < 4; i++ {
		SBBQ(mem(yPtr, i), regs[i])
	}

	SBBQ(Imm(0), b)
	MOVQ(Imm(0xfffffffffffffff7), n0)
	ANDQ(b, n0)
	MOVQ(Imm(0x07ffffffffffffff), n3)
	ANDQ(b, n3)

	ADDQ(n0, regs[0])
	ADCQ(b, regs[1])
	ADCQ(b, regs[2])
	ADCQ(n3, regs[3])

	storeResults()

	RET()
}

func mulDFunc() {
	TEXT("mulD", NOSPLIT, "func(res, x *FieldElement)")
	Doc("res=x * -1174 % 2^251-9")
	Pragma("noescape")
	yPtr = Load(Param("x"), GP64())
	resPtr = Load(Param("res"), GP64())

	XORQ(zero, zero)
	for i := 2; i < 8; i++ {
		XORQ(regs[i], regs[i])
	}
	MOVQ(U64(1174), RDX)
	MULXQ(mem(yPtr, 0), regs[0], regs[1])
	lastX = 0
	mulAdd(0, 1)
	mulAdd(0, 2)
	mulAdd(0, 3)
	carry(4, 6)

	extendedMod()

	p0, p1, p2, p3 := GP64(), GP64(), GP64(), GP64()

	MOVQ(Imm(0xfffffffffffffff7), p0)
	MOVQ(Imm(0xffffffffffffffff), p1)
	MOVQ(Imm(0xffffffffffffffff), p2)
	MOVQ(Imm(0x07ffffffffffffff), p3)

	SUBQ(regs[0], p0)
	SBBQ(regs[1], p1)
	SBBQ(regs[2], p2)
	SBBQ(regs[3], p3)

	MOVQ(p0, mem(resPtr, 0))
	MOVQ(p1, mem(resPtr, 1))
	MOVQ(p2, mem(resPtr, 2))
	MOVQ(p3, mem(resPtr, 3))

	RET()
}

func div2Func() {
	TEXT("div2", NOSPLIT, "func(res, x *FieldElement)")
	Doc("res=x >> 1")
	Pragma("noescape")
	xPtr = Load(Param("x"), GP64())
	resPtr = Load(Param("res"), GP64())

	for i := 0; i < 4; i++ {
		MOVQ(mem(xPtr, i), regs[i])
	}

	for i := 0; i < 3; i++ {
		SHRQ(Imm(1), regs[i+1], regs[i])
	}
	SARQ(Imm(1), regs[3])

	storeResults()

	RET()
}

func mul2Func() {
	TEXT("mul2", NOSPLIT, "func(res, x *FieldElement)")
	Doc("res=x * 2 % 2^251-9")
	Pragma("noescape")
	xPtr = Load(Param("x"), GP64())

	for i := 0; i < 4; i++ {
		MOVQ(mem(xPtr, i), regs[i])
	}

	for i := 3; i >= 1; i-- {
		SHLQ(Imm(1), regs[i-1], regs[i])
	}
	SHLQ(Imm(1), regs[0])

	subN()

	storeResults()

	RET()
}

func subN() {
	p0, p12, p3 := GP64(), GP64(), GP64()

	MOVQ(Imm(0xffffffffffffffff), p12)
	MOVQ(Imm(0xfffffffffffffff7), p0)
	MOVQ(Imm(0x07ffffffffffffff), p3)

	for i := 0; i < 4; i++ {
		MOVQ(regs[i], regs[i+4])
	}

	SUBQ(p0, regs[4])
	SBBQ(p12, regs[5])
	SBBQ(p12, regs[6])
	SBBQ(p3, regs[7])

	for i := 0; i < 4; i++ {
		CMOVQCC(regs[i+4], regs[i])
	}
}

func storeResults() {
	Comment("Store results")
	resPtr = Load(Param("res"), GP64())
	for i := 0; i < 4; i++ {
		MOVQ(regs[i], mem(resPtr, i))
	}
}

func modFunc() {
	TEXT("mod", NOSPLIT, "func(res, x *FieldElement)")
	Doc("res=x % 2^251-9")
	Pragma("noescape")
	xPtr = Load(Param("x"), GP64())

	for i := 0; i < 4; i++ {
		MOVQ(mem(xPtr, i), regs[i])
	}
	mod()

	storeResults()

	RET()
}

func mulNoAdxFunc() {
	TEXT("mulNoAdx", NOSPLIT, "func(res, x, y *FieldElement)")
	Pragma("noescape")
	Doc("res=x * y % 2^251-9")

	mulNoAdx()
}

func mulNoAdx() {
	xPtr = Load(Param("x"), GP64())
	resPtr = Load(Param("res"), GP64())
	yPtr = Load(Param("y"), GP64())
	flag = GP8()

	mulq(0, 0)
	mulq(0, 2)
	mulq(3, 1)
	mulq(3, 3)

	MOVB(Imm(1), flag)
	mulqAdd(1, 0)
	mulqAdd(0, 3)
	mulqAdd(2, 3)
	ADCQ(Imm(0), regs[7])

	MOVB(Imm(1), flag)
	mulqAdd(0, 1)
	mulqAdd(1, 2)
	mulqAdd(3, 2)
	ADCQ(Imm(0), regs[7])

	MOVB(Imm(1), flag)
	mulqAdd(1, 1)
	mulqAdd(1, 3)
	ADCQ(Imm(0), regs[6])
	ADCQ(Imm(0), regs[7])

	MOVB(Imm(1), flag)
	mulqAdd(2, 0)
	mulqAdd(2, 2)
	ADCQ(Imm(0), regs[6])
	ADCQ(Imm(0), regs[7])

	MOVB(Imm(1), flag)
	mulqAdd(2, 1)
	ADCQ(Imm(0), regs[5])
	ADCQ(Imm(0), regs[6])
	ADCQ(Imm(0), regs[7])

	MOVB(Imm(1), flag)
	mulqAdd(3, 0)
	ADCQ(Imm(0), regs[5])
	ADCQ(Imm(0), regs[6])
	ADCQ(Imm(0), regs[7])

	extendedMod()

	storeResults()

	RET()
}

func mulqAdd(x, y int) {
	MOVQ(mem(xPtr, x), RAX)
	MULQ(mem(yPtr, y))
	CMPB(flag, Imm(1))
	ADCQ(RAX, regs[x+y])
	ADCQ(RDX, regs[x+y+1])
	SETCC(flag)
}

func mulq(x, y int) {
	MOVQ(mem(xPtr, x), RAX)
	MULQ(mem(yPtr, y))
	MOVQ(RAX, regs[x+y])
	MOVQ(RDX, regs[x+y+1])
}

func sqrFunc() {
	TEXT("sqr", NOSPLIT, "func(res, x *FieldElement)")
	Pragma("noescape")
	Doc("res=x * x % 2^251-9")

	//CMPB(Mem{Base: StaticBase, Symbol: Symbol{Name: "·cpuSupported"}}, U8(1))
	//JNE(LabelRef("mulNoAdx"))

	xPtr = Load(Param("x"), GP64())

	x0, x1, x2, x3 := GP64(), GP64(), GP64(), GP64()

	Comment("load x to registers")
	MOVQ(mem(xPtr, 0), x0)
	MOVQ(mem(xPtr, 1), x1)
	MOVQ(mem(xPtr, 2), x2)
	MOVQ(mem(xPtr, 3), x3)

	Comment("clear flags")
	XORQ(zero, zero)

	Comment("fill registers")

	Comment("x[3]*x[2]")
	mulS(3, x3, x2, regs[5], regs[6])
	Comment("x[0]*x[3]")
	mulS(0, x0, x3, regs[3], regs[4])
	Comment("x[0]*x[1]")
	mulS(0, x0, x1, regs[1], regs[2])

	Comment("2-4 pass")
	mulAddS(0, 2, x0, x2)
	mulAddS(1, 2, x1, x2)
	mulAddS(1, 3, x1, x3)
	carry(5, 8)

	Comment("clear 7")
	XORQ(regs[7], regs[7])

	Comment("multiply by 2 by shifting")
	SHLQ(Imm(1), regs[6], regs[7])
	SHLQ(Imm(1), regs[5], regs[6])
	SHLQ(Imm(1), regs[4], regs[5])
	SHLQ(Imm(1), regs[3], regs[4])
	SHLQ(Imm(1), regs[2], regs[3])
	SHLQ(Imm(1), regs[1], regs[2])
	SHLQ(Imm(1), regs[1])

	Comment("add all z*z")
	lo, hi := GP64(), GP64()
	mulS(0, x0, x0, regs[0], hi)
	ADDQ(hi, regs[1])
	mulS(1, x1, x1, lo, hi)
	ADCQ(lo, regs[2])
	ADCQ(hi, regs[3])
	mulS(2, x2, x2, lo, hi)
	ADCQ(lo, regs[4])
	ADCQ(hi, regs[5])
	mulS(3, x3, x3, lo, hi)
	ADCQ(lo, regs[6])
	ADCQ(hi, regs[7])

	extendedMod()

	storeResults()

	RET()

	//Label("mulNoAdx")
	//mulNoAdx()
}

func mulFunc() {
	TEXT("mul", NOSPLIT, "func(res, x, y *FieldElement)")
	Pragma("noescape")
	Doc("res=x * y % 2^251-9")

	CMPB(Mem{Base: StaticBase, Symbol: Symbol{Name: "·cpuSupported"}}, U8(1))
	JNE(LabelRef("mulNoAdx"))

	xPtr = Load(Param("x"), GP64())
	yPtr = Load(Param("y"), GP64())

	Comment("Fill all regs")
	mul(3, 1, regs[4], regs[5])
	mul(3, 3, regs[6], regs[7])
	mul(0, 0, regs[0], regs[1])
	mul(0, 2, regs[2], regs[3])

	XORQ(zero, zero)
	Comment("First 1-5 chain")
	mulAdd(0, 1)
	mulAdd(2, 0)
	mulAdd(2, 1)
	mulAdd(2, 2)
	mulAdd(2, 3)
	carry(6, 8)

	Comment("Second 1-5 chain")
	mulAdd(1, 0)
	mulAdd(1, 1)
	mulAdd(1, 2)
	mulAdd(1, 3)
	mulAdd(3, 2)
	carry(6, 8)

	mulAdd(3, 0)
	carry(4, 8)

	mulAdd(0, 3)
	carry(4, 8)

	extendedMod()

	storeResults()

	RET()

	Label("mulNoAdx")
	mulNoAdx()
}

func extendedMod() {
	Comment("Mod 1st stage")
	SHLQ(Imm(5), regs[6], regs[7])
	SHLQ(Imm(5), regs[5], regs[6])
	SHLQ(Imm(5), regs[4], regs[5])
	SHLQ(Imm(5), regs[3], regs[4])

	andReg := GP64()
	MOVQ(Imm(0x07ffffffffffffff), andReg)
	ANDQ(andReg, regs[3])

	ADDQ(regs[4], regs[0])
	ADCQ(regs[5], regs[1])
	ADCQ(regs[6], regs[2])
	ADCQ(regs[7], regs[3])

	SHLQ(Imm(3), regs[6], regs[7])
	SHLQ(Imm(3), regs[5], regs[6])
	SHLQ(Imm(3), regs[4], regs[5])
	SHLQ(Imm(3), regs[4])

	ADDQ(regs[4], regs[0])
	ADCQ(regs[5], regs[1])
	ADCQ(regs[6], regs[2])
	ADCQ(regs[7], regs[3])

	Comment("Mod 2nd stage")
	mod()
}

func carry(start, end int) {
	Comment(fmt.Sprintf("Carry %d-%d", start, end))
	ADCXQ(zero, regs[start])

	for i := start + 1; i < end; i++ {
		ADOXQ(zero, regs[i])
		ADCXQ(zero, regs[i])
	}
}

func mulAdd(x, y int) {
	Comment(fmt.Sprintf("x[%d]*y[%d]", x, y))
	hi := GP64()
	low := GP64()
	xy := x + y
	mul(x, y, low, hi)
	ADCXQ(low, regs[xy])
	ADOXQ(hi, regs[xy+1])
}

func mul(x, y int, low, hi Register) {
	if x != lastX {
		MOVQ(mem(xPtr, x), RDX)
		lastX = x
	}
	MULXQ(mem(yPtr, y), low, hi)
}

func mulAddS(xn, yn int, x, y Register) {
	Comment(fmt.Sprintf("x[%d]*y[%d]", xn, yn))
	hi := GP64()
	low := GP64()
	mulS(xn, x, y, low, hi)
	xy := xn + yn
	ADCXQ(low, regs[xy])
	ADOXQ(hi, regs[xy+1])
}

func mulS(xn int, x, y, low, hi Register) {
	if xn != lastX {
		MOVQ(x, RDX)
		lastX = xn
	}
	MULXQ(y, low, hi)
}

func mem(xPtr Register, i int) Mem {
	return Mem{Base: xPtr, Disp: 8 * i}
}

func mod() {
	MOVQ(regs[3], regs[7])
	SHRQ(Imm(59), regs[7])

	p3 := GP64()
	MOVQ(Imm(0x07ffffffffffffff), p3)
	ANDQ(p3, regs[3])

	Comment("regs[7] = regs[7]*9")
	LEAQ(Mem{Base: regs[7], Index: regs[7], Scale: 8}, regs[7])

	ADDQ(regs[7], regs[0])
	ADCQ(Imm(0), regs[1])
	ADCQ(Imm(0), regs[2])
	ADCQ(Imm(0), regs[3])

	p12 := GP64()
	p0 := GP64()
	MOVQ(Imm(0xffffffffffffffff), p12)
	MOVQ(Imm(0xfffffffffffffff7), p0)

	for i := 0; i < 4; i++ {
		MOVQ(regs[i], regs[i+4])
	}
	SUBQ(p0, regs[4])
	SBBQ(p12, regs[5])
	SBBQ(p12, regs[6])
	SBBQ(p3, regs[7])

	for i := 0; i < 4; i++ {
		CMOVQCC(regs[i+4], regs[i])
	}
}
