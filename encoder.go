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

// Encoder is used for encoding data
type Encoder struct {
	mEncoder *C.kslide_encoder_t
}

// freeEncoder deallocates and release the memory consumed by an encoder
// @param encoder The encoder which should be deallocated
func freeEncoder(encoder *Encoder) {
	C.kslide_delete_encoder(encoder.mEncoder)
}

/// @param encoder The encoder to query
/// @return the symbol size in bytes
func (encoder *Encoder) SymbolSize() uint32 {
	return uint32(C.kslide_encoder_symbol_size(encoder.mEncoder))
}

/// @param decoder The encoder to query
/// @return The number of symbols in the stream
func (encoder *Encoder) StreamSymbols() uint32 {
	return uint32(C.kslide_encoder_stream_symbols(encoder.mEncoder))
}

/// @param encoder The encoder to query
/// @return The number of symbols in the stream
func (encoder *Encoder) StreamLowerBound() uint32 {
	return uint32(C.kslide_encoder_stream_lower_bound(encoder.mEncoder))
}

/// @param encoder The encoder to query
/// @return The size of a coefficient vector in bytes
func (encoder *Encoder) CoefficientsVectorSize() uint32 {
	return uint32(C.kslide_encoder_coefficients_vector_size(encoder.mEncoder))
}

/// Adds a symbol to the front of the encoder. Increments the stream front
/// index.
/// @param data The data pointer to push.
/// @return The stream index of the symbol being added.
func (encoder *Encoder) PushFrontSymbol(data *[]uint8) uint64 {

	return uint64(C.kslide_encoder_push_front_symbol(
		encoder.mEncoder, (*C.uint8_t)(&(*data)[0])))
}

func (encoder *Encoder) Generate(data *[]uint8) {
	C.kslide_encoder_generate(encoder.mEncoder, (*C.uint8_t)(&(*data)[0]))
}

func (encoder *Encoder) SetSeed(seed uint32) {
	C.kslide_encoder_set_seed(encoder.mEncoder, C.uint32_t(seed))
}

func (encoder *Encoder) WriteSymbol(symbol *[]uint8, coefficients *[]uint8) {
	C.kslide_encoder_write_symbol(
		encoder.mEncoder,
		(*C.uint8_t)(&(*symbol)[0]),
		(*C.uint8_t)(&(*coefficients)[0]))
}

/// @param encoder The encoder to query
/// @param lower_bound The index of the "oldest" symbol in the window
/// @param symbols The number of symbols in the window
func (encoder *Encoder) SetWindow(lowerBound uint32, symbols uint32) {
	C.kslide_encoder_set_window(
		encoder.mEncoder,
		C.uint32_t(lowerBound),
		C.uint32_t(symbols))
}
