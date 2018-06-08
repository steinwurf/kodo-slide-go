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
func NewEncoderFactory() *EncoderFactory {
	factory := new(EncoderFactory)
	factory.mFactory = C.kslide_new_encoder_factory()
	runtime.SetFinalizer(factory, deleteEncoderFactory)
	return factory
}

// deleteEncoderFactory deallocates the memory consumed by a factory
func deleteEncoderFactory(factory *EncoderFactory) {
	C.kslide_delete_encoder_factory(factory.mFactory)
}

// SymbolSize returns the symbol size in bytes
func (factory *EncoderFactory) SymbolSize() uint64 {
	return uint64(C.kslide_encoder_factory_symbol_size(factory.mFactory))
}

// SetSymbolSize sets the symbol size in bytes.
func (factory *EncoderFactory) SetSymbolSize(symbolSize uint64) {
	C.kslide_encoder_factory_set_symbol_size(
		factory.mFactory, C.uint64_t(symbolSize))
}

// Field returns the finite field.
func (factory *EncoderFactory) Field() int32 {
	return int32(C.kslide_encoder_factory_field(factory.mFactory))
}

// SetField sets the finite field.
func (factory *EncoderFactory) SetField(field int32) {
	C.kslide_encoder_factory_set_field(factory.mFactory, C.int32_t(field))
}

// Build builds an encoder
func (factory *EncoderFactory) Build() *Encoder {
	encoder := new(Encoder)
	encoder.mEncoder = C.kslide_encoder_factory_build(factory.mFactory)
	runtime.SetFinalizer(encoder, deleteEncoder)
	return encoder
}

// Initialize re-initializes the given encoder so that it can be reused.
func (factory *EncoderFactory) Initialize(encoder *Encoder) {
	C.kslide_encoder_factory_initialize(factory.mFactory, encoder.mEncoder)
}
