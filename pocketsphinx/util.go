package pocketsphinx

/*
#cgo pkg-config: pocketsphinx
#include "pocketsphinx.h"
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"

func Args() []Arg {
	__ret := C.ps_args()
	var __v = make([]Arg, 2048)
	packSArg(__v, __ret)
	return __v
}
