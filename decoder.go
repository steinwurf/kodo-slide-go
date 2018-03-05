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

// Decoder is used for encoding data
type Decoder struct {
	mDecoder *C.kslide_decoder_t
}

// freeDecoder deallocates and release the memory consumed by an decoder
// @param decoder The decoder which should be deallocated
func freeDecoder(decoder *Decoder) {
	C.kslide_delete_decoder(decoder.mDecoder)
}

// @param decoder The decoder to query
// @return the symbol size in bytes
func (decoder *Decoder) symbolSize() uint32 {
	return uint32(C.kslide_decoder_symbol_size(decoder.mDecoder))
}

// @param decoder The decoder to query
// @return The number of symbols in the stream
func (decoder *Decoder) streamSymbols() uint32 {
	return uint32(C.kslide_decoder_stream_symbols(decoder.mDecoder))
}

// @param decoder The decoder to query
// @return The number of symbols in the stream
func (decoder *Decoder) streamLowerBound() uint32 {
	return uint32(C.kslide_decoder_stream_lower_bound(decoder.mDecoder))
}

// Adds a symbol to the front of the decoder. Increments the stream front
// index.
// @param decoder The decoder to push to.
// @param data The data pointer to push.
// @return The stream index of the symbol being added.
func (decoder *Decoder) pushFrontSymbol(data *[]uint8) uint64 {
	return uint64(C.kslide_decoder_push_front_symbol(
		decoder.mDecoder, (*C.uint8_t)(&(*data)[0])))
}

// @param decoder The decoder to query
// @return The number of symbols decoded
func (decoder *Decoder) symbolsDecoded() uint32 {
	return uint32(C.kslide_decoder_symbols_decoded(decoder.mDecoder))
}

// @param decoder The decoder to query
// @param lower_bound The index of the "oldest" symbol in the window
// @param symbols The number of symbols in the window
func (decoder *Decoder) setWindow(lowerBound uint32, symbols uint32) {
	C.kslide_decoder_set_window(
		decoder.mDecoder,
		C.uint32_t(lowerBound),
		C.uint32_t(symbols))
}

func (decoder *Decoder) readSymbol(symbol *[]uint8, coefficients *[]uint8) {
	C.kslide_decoder_read_symbol(
		decoder.mDecoder,
		(*C.uint8_t)(&(*symbol)[0]),
		(*C.uint8_t)(&(*coefficients)[0]))
}
