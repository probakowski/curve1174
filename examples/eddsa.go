package main

import (
	"unsafe"
)

var a func(x, y, z uintptr)

func main() {
	data := [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}
	ab(&data, &data, &data)
}

func ab(data1, data2 ,data3 *[8]uint64) {
	a(uintptr(unsafe.Pointer(data1)), uintptr(unsafe.Pointer(data2)), uintptr(unsafe.Pointer(data3)))
}