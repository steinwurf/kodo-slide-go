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

// DecoderFactory builds Decoders
type DecoderFactory struct {
	mFactory *C.kslide_decoder_factory_t
}

// NewDecoderFactory builds a new decoder factory
func NewDecoderFactory() *DecoderFactory {
	factory := new(DecoderFactory)
	factory.mFactory = C.kslide_new_decoder_factory()
	runtime.SetFinalizer(factory, deleteDecoderFactory)
	return factory
}

// deleteDecoderFactory deallocates the memory consumed by a factory
func deleteDecoderFactory(factory *DecoderFactory) {
	C.kslide_delete_decoder_factory(factory.mFactory)
}

// SymbolSize returns the symbol size in bytes
func (factory *DecoderFactory) SymbolSize() uint64 {
	return uint64(C.kslide_decoder_factory_symbol_size(factory.mFactory))
}

// SetSymbolSize sets the symbol size in bytes
func (factory *DecoderFactory) SetSymbolSize(symbolSize uint64) {
	C.kslide_decoder_factory_set_symbol_size(
		factory.mFactory, C.uint64_t(symbolSize))
}

// Field returns the finite field
func (factory *DecoderFactory) Field() int32 {
	return int32(C.kslide_decoder_factory_field(factory.mFactory))
}

// SetField sets the finite field
func (factory *DecoderFactory) SetField(field int32) {
	C.kslide_decoder_factory_set_field(factory.mFactory, C.int32_t(field))
}

// Build builds a decoder
func (factory *DecoderFactory) Build() *Decoder {
	decoder := new(Decoder)
	decoder.mDecoder = C.kslide_decoder_factory_build(factory.mFactory)
	runtime.SetFinalizer(decoder, deleteDecoder)
	return decoder
}

// Initialize re-initializes the given decoder so that it can be reused.
func (factory *DecoderFactory) Initialize(decoder *Decoder) {
	C.kslide_decoder_factory_initialize(factory.mFactory, decoder.mDecoder)
}
