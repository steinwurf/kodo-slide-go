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

// deleteEncoder deallocates and release the memory consumed by an encoder
// @param encoder The encoder which should be deallocated
func deleteEncoder(encoder *Encoder) {
	C.kslide_delete_encoder(encoder.mEncoder)
}

// SymbolSize returns the symbol size in bytes
// @return the symbol size in bytes
func (encoder *Encoder) SymbolSize() uint32 {
	return uint32(C.kslide_encoder_symbol_size(encoder.mEncoder))
}

// StreamSymbols returns the number of symbols in the stream
// @param decoder The encoder to query
// @return The number of symbols in the stream
func (encoder *Encoder) StreamSymbols() uint32 {
	return uint32(C.kslide_encoder_stream_symbols(encoder.mEncoder))
}

// StreamLowerBound returns the number of symbols in the stream
// @return The number of symbols in the stream
func (encoder *Encoder) StreamLowerBound() uint32 {
	return uint32(C.kslide_encoder_stream_lower_bound(encoder.mEncoder))
}

// CoefficientVectorSize returns the size of a coefficient vector in bytes
// @return The size of a coefficient vector in bytes
func (encoder *Encoder) CoefficientVectorSize() uint32 {
	return uint32(C.kslide_encoder_coefficient_vector_size(encoder.mEncoder))
}

// PushFrontSymbol adds a symbol to the front of the encoder, increments the
// stream front index, and returns the stream index of the symbol being added.
// @param data The data pointer to push.
// @return The stream index of the symbol being added.
func (encoder *Encoder) PushFrontSymbol(data *[]uint8) uint64 {

	return uint64(C.kslide_encoder_push_front_symbol(
		encoder.mEncoder, (*C.uint8_t)(&(*data)[0])))
}

// Generate generates coding coefficients for the symbols in the coding window
// according to the specified seed
// @param coefficients Buffer where the coding coefficients should be
//        stored. This buffer must be encoder::coefficient_vector_size()
func (encoder *Encoder) Generate(data *[]uint8) {
	C.kslide_encoder_generate(encoder.mEncoder, (*C.uint8_t)(&(*data)[0]))
}

// SetSeed sets the seed of the internal random generator function. If using the
// same seed on the encoder and decoder the exact same set of coefficients will
// be generated.
// @param seed_value A value for the seed.
func (encoder *Encoder) SetSeed(seed uint32) {
	C.kslide_encoder_set_seed(encoder.mEncoder, C.uint32_t(seed))
}

// WriteSymbol writes an encoded symbol according to the coding coefficients.
// @param symbol The buffer where the encoded symbol will be stored.
//        The symbol buffer must be Encoder::SymbolSize.
// @param coefficients The coding coefficients. A compatible format can
//        be created using Encoder::Generate.
func (encoder *Encoder) WriteSymbol(symbol *[]uint8, coefficients *[]uint8) {
	C.kslide_encoder_write_symbol(
		encoder.mEncoder,
		(*C.uint8_t)(&(*symbol)[0]),
		(*C.uint8_t)(&(*coefficients)[0]))
}

// SetWindow changes the index of the "oldest" symbol in the window and the
// number of symbols in the window.
// @param lower_bound The index of the "oldest" symbol in the window
// @param symbols The number of symbols in the window
func (encoder *Encoder) SetWindow(lowerBound uint32, symbols uint32) {
	C.kslide_encoder_set_window(
		encoder.mEncoder,
		C.uint32_t(lowerBound),
		C.uint32_t(symbols))
}
