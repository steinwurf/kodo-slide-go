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
func deleteEncoder(encoder *Encoder) {
	C.kslide_delete_encoder(encoder.mEncoder)
}

// SymbolSize returns the size of a symbol in the stream in bytes.
// @return The size of a symbol in the stream in bytes.
func (encoder *Encoder) SymbolSize() uint64 {
	return uint64(C.kslide_encoder_symbol_size(encoder.mEncoder))
}

// StreamSymbols returns the total number of symbols available in memory at the
// decoder.
// The number of symbols in the coding window MUST be less than or equal to
// this number.
func (encoder *Encoder) StreamSymbols() uint64 {
	return uint64(C.kslide_encoder_stream_symbols(encoder.mEncoder))
}

// StreamLowerBound returns the index of the oldest symbol known by the encoder.
// This symbol may not be inside the window but can be included in the window
// if needed.
func (encoder *Encoder) StreamLowerBound() uint64 {
	return uint64(C.kslide_encoder_stream_lower_bound(encoder.mEncoder))
}

// StreamUpperBound returns the upper bound of the stream.
// The range of valid symbol indices goes from [Encoder::StreamLowerBound(),
// Encoder::StreamUpperBound()).
// Note the stream is a half-open interval. Going from
// Encoder::StreamLowerBound() to Encoder::StreamUpperBound() - 1.
func (encoder *Encoder) StreamUpperBound() uint64 {
	return uint64(C.kslide_encoder_stream_upper_bound(encoder.mEncoder))
}

// PushFrontSymbol adds a new symbol to the front of the encoder,
// increments the number of symbols in the stream and increases the stream
// upper bound, and returns the stream index of the symbol being added.
// Note, the caller must ensure that the memory of the given symbol remains
// valid as long as the symbol is included in the stream.
// The caller is responsible for freeing the memory if needed.
// Once the symbol is popped from the stream.
func (encoder *Encoder) PushFrontSymbol(data *[]uint8) uint64 {
	return uint64(C.kslide_encoder_push_front_symbol(
		encoder.mEncoder, (*C.uint8_t)(&(*data)[0])))
}

// PopBackSymbol removes the "oldest" symbol from the stream. Increments the
// Encoder::StreamLowerBound(), and returns the index of the symbol being
// removed.
func (encoder *Encoder) PopBackSymbol() uint64 {
	return uint64(C.kslide_encoder_pop_back_symbol(encoder.mEncoder))
}

// WindowSymbols returns the number of symbols currently in the coding window.
// The window must be within the bounds of the stream.
func (encoder *Encoder) WindowSymbols() uint64 {
	return uint64(C.kslide_encoder_window_symbols(encoder.mEncoder))
}

// WindowLowerBound returns the index of the "oldest" symbol in the coding
// window.
func (encoder *Encoder) WindowLowerBound() uint64 {
	return uint64(C.kslide_encoder_window_lower_bound(encoder.mEncoder))
}

// WindowUpperBound returns the upper bound of the window.
// The range of valid symbol indices goes from
// [Encoder::WindowLowerBound(), Encoder::WindowUpperBound()).
// Note the window is a half-open interval.
// Going from Encoder::WindowLowerBound() to Encoder::WindowUpperBound() - 1.
func (encoder *Encoder) WindowUpperBound() uint64 {
	return uint64(C.kslide_encoder_window_upper_bound(encoder.mEncoder))
}

// SetWindow sets the window.
// The window represents the symbols which will be included in the next
// encoding.
// The window cannot exceed the bounds of the stream.
//
// Example: If lowerBound=4 and symbols=3 the following
//          symbol indices will be included 4,5,6
func (encoder *Encoder) SetWindow(lowerBound uint64, symbols uint64) {
	C.kslide_encoder_set_window(
		encoder.mEncoder,
		C.uint64_t(lowerBound),
		C.uint64_t(symbols))
}

// CoefficientVectorSize returns the size of the coefficient vector in the
// current window in bytes.
// The number of coefficients is equal to the number of symbols in the window.
// The size in bits of each coefficients depends on the finite field chosen.
// A custom coding scheme can be implemented by generating the coding vector
// manually.
// Alternatively the built-in generator can be used.
// See Encoder::SetSeed(...) and Encoder::Generate(...).
func (encoder *Encoder) CoefficientVectorSize() uint64 {
	return uint64(C.kslide_encoder_coefficient_vector_size(encoder.mEncoder))
}

// SetSeed sets the seed of the internal random generator function.
// If using the same seed on the encoder and decoder the exact same set of
// coefficients will be generated.
func (encoder *Encoder) SetSeed(seedValue uint64) {
	C.kslide_encoder_set_seed(encoder.mEncoder, C.uint64_t(seedValue))
}

// Generate generates coding coefficients for the symbols in the coding window
// according to the specified seed (see Encoder::SetSeed(...)).
// The given buffer must be Encoder::CoefficientVectorSize() large in bytes.
func (encoder *Encoder) Generate(coefficients *[]uint8) {
	C.kslide_encoder_generate(
		encoder.mEncoder, (*C.uint8_t)(&(*coefficients)[0]))
}

// WriteSymbol Writes an encoded symbol according to the coding coefficients.
// The symbol buffer must be Encoder::SymbolSize() large.
// The  coefficients must have the memory layout required.
// A compatible format can be created using Encoder::Generate(...)
func (encoder *Encoder) WriteSymbol(symbol *[]uint8, coefficients *[]uint8) {
	C.kslide_encoder_write_symbol(
		encoder.mEncoder,
		(*C.uint8_t)(&(*symbol)[0]),
		(*C.uint8_t)(&(*coefficients)[0]))
}

// WriteSourceSymbol Writes a source symbol to the symbol buffer at the given
// index.
// The symbol buffer must be Encoder::SymbolSize() large.
func (encoder *Encoder) WriteSourceSymbol(symbol *[]uint8, index uint64) {
	C.kslide_encoder_write_source_symbol(
		encoder.mEncoder,
		(*C.uint8_t)(&(*symbol)[0]),
		C.uint64_t(index))
}
