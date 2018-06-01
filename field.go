package kslide

// // Copyright Steinwurf ApS 2018.
// // Distributed under the "STEINWURF RESEARCH LICENSE 1.0".
// // See accompanying file LICENSE.rst or
// // http://www.steinwurf.com/licensing

/*
#cgo CFLAGS: -I../kodo-slide-c
#include <kodo_slide_c.h>
*/
import "C"

/// Enum specifying the available finite fields
const (
	Binary   = int32(C.kslide_binary)
	Binary4  = int32(C.kslide_binary4)
	Binary8  = int32(C.kslide_binary8)
	Binary16 = int32(C.kslide_binary16)
)
