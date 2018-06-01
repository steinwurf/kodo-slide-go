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

// freeDecoder deallocates and release the memory consumed by an decoder
// @param decoder The decoder which should be deallocated
func deleteDecoder(decoder *Decoder) {
	C.kslide_delete_decoder(decoder.mDecoder)
}

// SymbolSize returns the symbol size in bytes
// @return the symbol size in bytes
func (decoder *Decoder) SymbolSize() uint32 {
	return uint32(C.kslide_decoder_symbol_size(decoder.mDecoder))
}

// StreamSymbols returns the number of symbols in the stream
// @return The number of symbols in the stream
func (decoder *Decoder) StreamSymbols() uint32 {
	return uint32(C.kslide_decoder_stream_symbols(decoder.mDecoder))
}

// StreamLowerBound returns the number of symbols in the stream
// @return The number of symbols in the stream
func (decoder *Decoder) StreamLowerBound() uint32 {
	return uint32(C.kslide_decoder_stream_lower_bound(decoder.mDecoder))
}

// PushFrontSymbol adds a symbol to the front of the decoder,
// Increments the stream front index, and returns the stream index of the symbol
// being added.
// @param data The data pointer to push.
// @return The stream index of the symbol being added.
func (decoder *Decoder) PushFrontSymbol(data *[]uint8) uint64 {
	return uint64(C.kslide_decoder_push_front_symbol(
		decoder.mDecoder, (*C.uint8_t)(&(*data)[0])))
}

// SymbolsDecoded returns the number of symbols decoded
// @return The number of symbols decoded
func (decoder *Decoder) SymbolsDecoded() uint32 {
	return uint32(C.kslide_decoder_symbols_decoded(decoder.mDecoder))
}

// SetWindow changes the index of the "oldest" symbol in the window and the
// number of symbols in the window.
// @param lower_bound The index of the "oldest" symbol in the window
// @param symbols The number of symbols in the window
func (decoder *Decoder) SetWindow(lowerBound uint32, symbols uint32) {
	C.kslide_decoder_set_window(
		decoder.mDecoder,
		C.uint32_t(lowerBound),
		C.uint32_t(symbols))
}

// ReadSymbol reads a symbol to the decoder.
// @param symbol A pointer to the symbol
// @param coefficients A pointer to the coefficients associated with the symbol
func (decoder *Decoder) ReadSymbol(symbol *[]uint8, coefficients *[]uint8) {
	C.kslide_decoder_read_symbol(
		decoder.mDecoder,
		(*C.uint8_t)(&(*symbol)[0]),
		(*C.uint8_t)(&(*coefficients)[0]))
}
