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

// Decoder is used for decoding data
type Decoder struct {
	mDecoder *C.kslide_decoder_t
}

// freeDecoder deallocates and release the memory consumed by the given decoder
func deleteDecoder(decoder *Decoder) {
	C.kslide_delete_decoder(decoder.mDecoder)
}

// SymbolSize returns the size of a symbol in the stream in bytes.
func (decoder *Decoder) SymbolSize() uint64 {
	return uint64(C.kslide_decoder_symbol_size(decoder.mDecoder))
}

// StreamSymbols returns the total number of symbols known at the decoder.
// The number of symbols in the decoding window MUST be less than or equal to
// this number.
func (decoder *Decoder) StreamSymbols() uint64 {
	return uint64(C.kslide_decoder_stream_symbols(decoder.mDecoder))
}

// StreamLowerBound returns the index of the oldest symbol known by the decoder.
// This symbol may not be inside the window but can be included in the window
// if needed.
func (decoder *Decoder) StreamLowerBound() uint64 {
	return uint64(C.kslide_decoder_stream_lower_bound(decoder.mDecoder))
}

// StreamUpperBound returns the upper bound of the stream.
// The range of valid symbol indices goes from [Decoder::StreamLowerBound(),
// Decoder::StreamUpperBound()).
// Note the stream is a half-open interval. Going from
// Decoder::StreamLowerBound() to Decoder::StreamUpperBound() - 1.
func (decoder *Decoder) StreamUpperBound() uint64 {
	return uint64(C.kslide_decoder_stream_upper_bound(decoder.mDecoder))
}

// PushFrontSymbol adds a new symbol to the front of the decoder,
// increments the number of symbols in the stream and increases the stream
// upper bound, and returns the stream index of the symbol being added.
// Note, the caller must ensure that the memory of the given symbol remains
// valid as long as the symbol is included in the stream.
// The caller is responsible for freeing the memory if needed.
// Once the symbol is popped from the stream.
func (decoder *Decoder) PushFrontSymbol(symbol *[]uint8) uint64 {
	return uint64(C.kslide_decoder_push_front_symbol(
		decoder.mDecoder,
		(*C.uint8_t)(&(*symbol)[0])))
}

// PopBackSymbol removes the "oldest" symbol from the stream. Increments the
// Decoder::StreamLowerBound(), and returns the index of the symbol being
// removed.
func (decoder *Decoder) PopBackSymbol() uint64 {
	return uint64(C.kslide_decoder_pop_back_symbol(decoder.mDecoder))
}

// WindowSymbols returns the number of symbols currently in the coding window.
// The window must be within the bounds of the stream.
func (decoder *Decoder) WindowSymbols() uint64 {
	return uint64(C.kslide_decoder_window_symbols(decoder.mDecoder))
}

// WindowLowerBound returns the index of the "oldest" symbol in the coding
// window.
func (decoder *Decoder) WindowLowerBound() uint64 {
	return uint64(C.kslide_decoder_window_lower_bound(decoder.mDecoder))
}

// WindowUpperBound returns the upper bound of the window.
// The range of valid symbol indices goes from
// [Decoder::WindowLowerBound(), Decoder::WindowUpperBound()).
// Note the window is a half-open interval.
// Going from Decoder::WindowLowerBound() to Decoder::WindowUpperBound() - 1.
func (decoder *Decoder) WindowUpperBound() uint64 {
	return uint64(C.kslide_decoder_window_upper_bound(decoder.mDecoder))
}

// SetWindow sets the window.
// The window represents the symbols which will be included in the next
// decoding.
// The window cannot exceed the bounds of the stream.
//
// Example: If lowerBound=4 and symbols=3 the following
//          symbol indices will be included 4,5,6
func (decoder *Decoder) SetWindow(lowerBound uint64, symbols uint64) {
	C.kslide_decoder_set_window(
		decoder.mDecoder,
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
// See Decoder::SetSeed(...) and Decoder::Generate(...).
func (decoder *Decoder) CoefficientVectorSize() uint64 {
	return uint64(C.kslide_decoder_coefficient_vector_size(decoder.mDecoder))
}

// SetSeed sets the seed of the internal random generator function.
// If using the same seed on the decoder and encoder the exact same set of
// coefficients will be generated.
func (decoder *Decoder) SetSeed(seedValue uint64) {
	C.kslide_decoder_set_seed(decoder.mDecoder, C.uint64_t(seedValue))
}

// Generate generates coding coefficients for the symbols in the coding window
// according to the specified seed (see Decoder::SetSeed(...)).
// The given buffer must be Decoder::CoefficientVectorSize() large in bytes.
func (decoder *Decoder) Generate(coefficients *[]uint8) {
	C.kslide_decoder_generate(
		decoder.mDecoder,
		(*C.uint8_t)(&(*coefficients)[0]))
}

// ReadSymbol decodes a coded symbol according to the coding coefficients.
//
// Both buffers may be modified during this call. The reason for this
// is that the decoder will directly operate on the provided memory
// for performance reasons.
func (decoder *Decoder) ReadSymbol(symbol *[]uint8, coefficients *[]uint8) {
	C.kslide_decoder_read_symbol(
		decoder.mDecoder,
		(*C.uint8_t)(&(*symbol)[0]),
		(*C.uint8_t)(&(*coefficients)[0]))
}

// ReadSourceSymbol adds a source symbol at the decoder.
func (decoder *Decoder) ReadSourceSymbol(symbol *[]uint8, index uint64) {
	C.kslide_decoder_read_source_symbol(
		decoder.mDecoder,
		(*C.uint8_t)(&(*symbol)[0]),
		C.uint64_t(index))
}

// Rank returns the rank of a decoder indicates how many symbols have been
// partially or fully decoded. This number is also equivalent to the
// number of pivot elements we have in the stream.
func (decoder *Decoder) Rank() uint64 {
	return uint64(C.kslide_decoder_rank(decoder.mDecoder))
}

// SymbolsMissing returns the number of missing symbols at the decoder
func (decoder *Decoder) SymbolsMissing() uint64 {
	return uint64(C.kslide_decoder_symbols_missing(decoder.mDecoder))
}

// SymbolsPartiallyDecoded returns the number of partially decoded symbols at
// the decoder
func (decoder *Decoder) SymbolsPartiallyDecoded() uint64 {
	return uint64(C.kslide_decoder_symbols_partially_decoded(decoder.mDecoder))
}

// SymbolsDecoded returns the number of decoded symbols at the decoder
func (decoder *Decoder) SymbolsDecoded() uint64 {
	return uint64(C.kslide_decoder_symbols_decoded(decoder.mDecoder))
}

// IsSymbolDecoded return true if the symbol at the given index is is decoded
// (i.e. it corresponds to a source symbol), and otherwise false.
func (decoder *Decoder) IsSymbolDecoded(index uint64) bool {
	return C.kslide_decoder_is_symbol_decoded(decoder.mDecoder, C.uint64_t(index)) == 1
}
