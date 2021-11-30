//+build !purego

package curve1174

import "github.com/klauspost/cpuid"

var cpuSupported = cpuid.CPU.BMI2() && cpuid.CPU.ADX()

// res=x * y % 2^251-9
//go:noescape
func mul(res uintptr, x uintptr, y uintptr)

// res=x * y % 2^251-9
//go:noescape
func mulNoAdx(res *FieldElement, x *FieldElement, y *FieldElement)

// res=x % 2^251-9
//go:noescape
func mod(res *FieldElement, x *FieldElement)

//go:noescape
func sqr(res *FieldElement, x *FieldElement)

// res=x * 2 % 2^251-9
//go:noescape
func mul2(res *FieldElement, x *FieldElement)

// res=x >> 1
//go:noescape
func div2(res *FieldElement, x *FieldElement)

// res=x - y % 2^251-9
//go:noescape
func sub(res *FieldElement, x *FieldElement, y *FieldElement)

// res=x + y % 2^251-9
//go:noescape
func add(res *FieldElement, x *FieldElement, y *FieldElement)

// res=x * -1174 % 2^251-9
//go:noescape
func mulD(res *FieldElement, x *FieldElement)

//go:noescape
func selectPoint(res *Point, table *[16]Point, index uint64)

//go:noescape
func fastInverse(res, x *FieldElement)