package pocketsphinx

/*
#cgo pkg-config: pocketsphinx
#include "pocketsphinx.h"
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"
import (
	"unsafe"
)

func Args() []Arg {
	__ret := C.ps_args()
	if __ret == nil {
		return nil
	}

	// The array returned by ps_args() is terminated by CMDLN_EMPTY_OPTION.
	// This matches the arg-counting logic in arg_strlen() in sphinxbase's
	// src/libsphinxbase/util/cmd_ln.c.
	var nargs uintptr
	base := uintptr(unsafe.Pointer(__ret))
	for {
		arg := (*C.arg_t)(unsafe.Pointer(base + nargs*sizeOfArgValue))
		if arg.name == nil {
			break
		}
		nargs++
	}

	var __v = make([]Arg, nargs)
	packSArg(__v, __ret)
	return __v
}
