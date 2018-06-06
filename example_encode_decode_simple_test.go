package kslide_test

// Copyright Steinwurf ApS 2018.
// Distributed under the "STEINWURF RESEARCH LICENSE 1.0".
// See accompanying file LICENSE.rst or
// http://www.steinwurf.com/licensing

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/steinwurf/kodo-slide-go"
)

type SymbolStorage struct {
	symbolSize uint64
	symbols    uint64
	data       []uint8
}

func allocateStorage(symbolSize uint64, symbols uint64) *SymbolStorage {
	symbolStorage := new(SymbolStorage)

	symbolStorage.symbolSize = symbolSize
	symbolStorage.symbols = symbols
	symbolStorage.data = make([]uint8, symbols*symbolSize)
	return symbolStorage
}

func randomizeStorage(symbolStorage *SymbolStorage) {
	size := symbolStorage.symbolSize * symbolStorage.symbols
	for i := uint64(0); i < size; i++ {
		symbolStorage.data[i] = uint8(rand.Uint32())
	}
}

func storageSymbol(symbolStorage *SymbolStorage, index uint64) []uint8 {
	return symbolStorage.data[index*symbolStorage.symbolSize : (index+1)*symbolStorage.symbolSize]
}

func Example_encodeDecodeSimple() {
	// Seed random number generator to produce different results every time
	rand.Seed(time.Now().UTC().UnixNano())
	symbols := uint64(100)
	symbolSize := uint64(750)

	// Initialization of encoder and decoder
	encoderFactory := NewEncoderFactory()
	decoderFactory := NewDecoderFactory()

	encoderFactory.SetSymbolSize(symbolSize)
	decoderFactory.SetSymbolSize(symbolSize)

	encoder := encoderFactory.Build()
	decoder := decoderFactory.Build()

	// Allocate memory for the encoder and decoder
	decoderStorage := allocateStorage(symbolSize, symbols)
	encoderStorage := allocateStorage(symbolSize, symbols)

	// Fill the encoder storage with random data
	randomizeStorage(encoderStorage)

	// Provide the decoder with storage
	for i := uint64(0); i < symbols; i++ {
		symbol := storageSymbol(decoderStorage, i)
		decoder.PushFrontSymbol(&symbol)
	}

	iterations := uint32(0)
	maxIterations := uint32(1000)
	symbolsDecoded := uint64(0)

	for symbolsDecoded < symbols && iterations < maxIterations {

		if encoder.StreamSymbols() < symbols && rand.Uint32()%2 == 0 {
			symbol := storageSymbol(encoderStorage, encoder.StreamSymbols())
			encoder.PushFrontSymbol(&symbol)
		}

		if encoder.StreamSymbols() == 0 {
			continue
		}

		encoder.SetWindow(encoder.StreamLowerBound(), encoder.StreamSymbols())
		decoder.SetWindow(encoder.StreamLowerBound(), encoder.StreamSymbols())

		coefficients := make([]uint8, encoder.CoefficientVectorSize())

		symbol := make([]uint8, encoder.SymbolSize())

		encoder.SetSeed(rand.Uint64())
		encoder.Generate(&coefficients)

		encoder.WriteSymbol(&symbol, &coefficients)
		decoder.ReadSymbol(&symbol, &coefficients)

		symbolsDecoded = decoder.SymbolsDecoded()
		iterations++
	}

	// Check if we properly decoded the data
	for i, v := range encoderStorage.data {
		if v != decoderStorage.data[i] {
			fmt.Println("Unexpected failure to decode")
			fmt.Println("Please file a bug report :)")
			return
		}
	}
	fmt.Println("Data decoded correctly")
	// Output: Data decoded correctly
}
