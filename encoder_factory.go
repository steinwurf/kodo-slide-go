package kslide

// Copyright Steinwurf ApS 2018.
// Distributed under the "STEINWURF RESEARCH LICENSE 1.0".
// See accompanying file LICENSE.rst or
// http://www.steinwurf.com/licensing

/*
#cgo CFLAGS: -I../kodo-slide-c
#cgo LDFLAGS: -L../kodo-slide-c -lkodo_slide_c_static -lkodo_slide -lfifi -lcpuid
#include <kodo_slide_c.h>
*/
import "C"
import "runtime"

// EncoderFactory builds Encoders
type EncoderFactory struct {
	mFactory *C.kslide_encoder_factory_t
}

// NewEncoderFactory builds a new encoder factory
// @return A new factory capable of building encoders using the
//         selected parameters.
func NewEncoderFactory() *EncoderFactory {
	factory := new(EncoderFactory)
	factory.mFactory = C.kslide_new_encoder_factory()
	runtime.SetFinalizer(factory, deleteEncoderFactory)
	return factory
}

// deleteEncoderFactory deallocates the memory consumed by a factory
// @param factory The factory which should be deallocated
func deleteEncoderFactory(factory *EncoderFactory) {
	C.kslide_delete_encoder_factory(factory.mFactory)
}

// Field returns the finite field
// @param factory The factory to query
// @return the the finite field
func (factory *EncoderFactory) Field() int32 {
	return int32(C.kslide_encoder_factory_field(factory.mFactory))
}

// SetField sets the finite field
// @param factory The factory which should be configured
// @param field the finite field
func (factory *EncoderFactory) SetField(field int32) {
	C.kslide_encoder_factory_set_field(factory.mFactory, C.int32_t(field))
}

// SymbolSize returns the symbol size in bytes
// @param factory The factory to query
// @return the symbol size in bytes
func (factory *EncoderFactory) SymbolSize() uint32 {
	return uint32(C.kslide_encoder_factory_symbol_size(factory.mFactory))
}

// SetSymbolSize sets the symbol size
// @param factory The factory which should be configured
// @param the symbol size in bytes
func (factory *EncoderFactory) SetSymbolSize(symbolSize uint32) {
	C.kslide_encoder_factory_set_symbol_size(
		factory.mFactory, C.uint32_t(symbolSize))
}

// Build builds the actual encoder
// @param factory The encoder factory which should be used to build the encoder
// @return pointer to an instantiation of an encoder
func (factory *EncoderFactory) Build() *Encoder {
	encoder := new(Encoder)
	encoder.mEncoder = C.kslide_encoder_factory_build(factory.mFactory)
	runtime.SetFinalizer(encoder, deleteEncoder)
	return encoder
}
