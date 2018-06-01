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
// @return A new factory capable of building decoders using the
//         selected parameters.
func NewDecoderFactory() *DecoderFactory {
	factory := new(DecoderFactory)
	factory.mFactory = C.kslide_new_decoder_factory()
	runtime.SetFinalizer(factory, deleteDecoderFactory)
	return factory
}

// deleteDecoderFactory deallocates the memory consumed by a factory
// @param factory The factory which should be deallocated
func deleteDecoderFactory(factory *DecoderFactory) {
	C.kslide_delete_decoder_factory(factory.mFactory)
}

// Field returns the finite field
// @param factory The factory to query
// @return the the finite field
func (factory *DecoderFactory) Field() int32 {
	return int32(C.kslide_decoder_factory_field(factory.mFactory))
}

// SetField sets the finite field
// @param factory The factory which should be configured
// @param field the finite field
func (factory *DecoderFactory) SetField(field int32) {
	C.kslide_decoder_factory_set_field(factory.mFactory, C.int32_t(field))
}

// SymbolSize returns the symbol size in bytes
// @param factory The factory to query
// @return the symbol size in bytes
func (factory *DecoderFactory) SymbolSize() uint32 {
	return uint32(C.kslide_decoder_factory_symbol_size(factory.mFactory))
}

// SetSymbolSize sets the symbol size
// @param factory The factory which should be configured
// @param symbolSize the symbol size in bytes
func (factory *DecoderFactory) SetSymbolSize(symbolSize uint32) {
	C.kslide_decoder_factory_set_symbol_size(
		factory.mFactory, C.uint32_t(symbolSize))
}

// Build builds the actual decoder
// @param factory The decoder factory which should be used to build the decoder
// @return pointer to an instantiation of an decoder
func (factory *DecoderFactory) Build() *Decoder {
	decoder := new(Decoder)
	decoder.mDecoder = C.kslide_decoder_factory_build(factory.mFactory)
	runtime.SetFinalizer(decoder, deleteDecoder)
	return decoder
}
