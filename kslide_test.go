package kslide_test

import (
	"math/rand"
	"testing"

	. "github.com/steinwurf/kodo-slide-go"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestDecoderFactory(c *C) {
	var symbolSize uint64 = 1300
	decoderFactory := NewDecoderFactory()

	c.Assert(decoderFactory.SymbolSize(), Equals, symbolSize)

	var newSymbolSize uint64 = 300
	decoderFactory.SetSymbolSize(newSymbolSize)
	c.Assert(decoderFactory.SymbolSize(), Equals, newSymbolSize)
	c.Assert(decoderFactory.Field(), Equals, Binary8)

	for _, field := range [4]int32{Binary, Binary4, Binary8, Binary16} {
		decoderFactory.SetField(field)
		c.Assert(decoderFactory.Field(), Equals, field)
	}
}

func (s *MySuite) TestDecoder(c *C) {
	var symbolSize uint64 = 1300
	decoderFactory := NewDecoderFactory()
	decoder := decoderFactory.Build()
	c.Assert(decoder.SymbolSize(), Equals, symbolSize)
	c.Assert(decoder.StreamSymbols(), Equals, uint64(0))
	c.Assert(decoder.StreamLowerBound(), Equals, uint64(0))

	c.Assert(decoder.SymbolSize(), Equals, symbolSize)
	c.Assert(decoder.StreamSymbols(), Equals, uint64(0))
	c.Assert(decoder.StreamLowerBound(), Equals, uint64(0))
	c.Assert(decoder.StreamUpperBound(), Equals, uint64(0))
	c.Assert(decoder.WindowSymbols(), Equals, uint64(0))
	c.Assert(decoder.WindowLowerBound(), Equals, uint64(0))
	c.Assert(decoder.WindowUpperBound(), Equals, uint64(0))
	c.Assert(decoder.CoefficientVectorSize(), Equals, uint64(0))
	c.Assert(decoder.Rank(), Equals, uint64(0))
	c.Assert(decoder.SymbolsMissing(), Equals, uint64(0))
	c.Assert(decoder.SymbolsPartiallyDecoded(), Equals, uint64(0))
	c.Assert(decoder.SymbolsDecoded(), Equals, uint64(0))

	var newSymbolSize uint64 = 300
	decoderFactory.SetSymbolSize(newSymbolSize)

	c.Assert(decoder.SymbolSize(), Equals, symbolSize)
	decoderFactory.Initialize(decoder)
	c.Assert(decoder.SymbolSize(), Equals, newSymbolSize)
}

func (s *MySuite) TestEncoderFactory(c *C) {
	var symbolSize uint64 = 1300
	encoderFactory := NewEncoderFactory()

	c.Assert(encoderFactory.SymbolSize(), Equals, symbolSize)

	var newSymbolSize uint64 = 300
	encoderFactory.SetSymbolSize(newSymbolSize)
	c.Assert(newSymbolSize, Equals, encoderFactory.SymbolSize())
	c.Assert(encoderFactory.Field(), Equals, Binary8)

	for _, field := range [4]int32{Binary, Binary4, Binary8, Binary16} {
		encoderFactory.SetField(field)
		c.Assert(encoderFactory.Field(), Equals, field)
	}
}

func (s *MySuite) TestEncoder(c *C) {
	var symbolSize uint64 = 1300
	encoderFactory := NewEncoderFactory()
	encoder := encoderFactory.Build()

	c.Assert(encoder.SymbolSize(), Equals, symbolSize)
	c.Assert(encoder.StreamSymbols(), Equals, uint64(0))
	c.Assert(encoder.StreamLowerBound(), Equals, uint64(0))
	c.Assert(encoder.StreamUpperBound(), Equals, uint64(0))
	c.Assert(encoder.WindowSymbols(), Equals, uint64(0))
	c.Assert(encoder.WindowLowerBound(), Equals, uint64(0))
	c.Assert(encoder.WindowUpperBound(), Equals, uint64(0))
	c.Assert(encoder.CoefficientVectorSize(), Equals, uint64(0))

	var newSymbolSize uint64 = 300
	encoderFactory.SetSymbolSize(newSymbolSize)

	c.Assert(encoder.SymbolSize(), Equals, symbolSize)
	encoderFactory.Initialize(encoder)
	c.Assert(encoder.SymbolSize(), Equals, newSymbolSize)
}

func (s *MySuite) TestCodec(c *C) {

	for _, field := range [4]int32{Binary, Binary4, Binary8, Binary16} {
		mixCodedUncoded(c, field)
	}
}

type RateController struct {
	mN        uint32
	mK        uint32
	mPosition uint32
}

func NewRateController(n uint32, k uint32) *RateController {
	controller := new(RateController)
	controller.mN = n
	controller.mK = k
	controller.mPosition = 0

	return controller
}

func (controller *RateController) Advance() {
	controller.mPosition = (controller.mPosition + 1) % controller.mN
}

func (controller *RateController) GenerateData() bool {
	return controller.mPosition < controller.mK
}

func mixCodedUncoded(c *C, field int32) {
	// Set the capacity of the decoder (this is the number of encoded symbols
	// that are used in the decoding process).
	const capacity uint64 = 15

	// Set the window size (this is the number of symbols included in an
	// encoded symbol).
	const windowSymbols uint64 = 5

	// The size of a symbol in bytes
	const symbolSize uint64 = 16

	// Maximum number of interations
	const maxIterations uint32 = 1000

	encoderFactory := NewEncoderFactory()
	decoderFactory := NewDecoderFactory()

	encoderFactory.SetSymbolSize(symbolSize)
	decoderFactory.SetSymbolSize(symbolSize)

	encoderFactory.SetField(field)
	decoderFactory.SetField(field)

	encoder := encoderFactory.Build()
	decoder := decoderFactory.Build()

	// Cache for all original source symbols added to the encoder - such that
	// we can check that they are decoded correctly.
	var sourceSymbolIndex = 0
	sourceSymbols := make([][]uint8, maxIterations)

	// Allocate our finite memory for decoding
	decoderStorage := make([]uint8, capacity*symbolSize)

	// Provide the decoder with storage
	for i := uint64(0); i < capacity; i++ {
		symbol := GetSymbol(decoderStorage, symbolSize, i)
		decoder.PushFrontSymbol(&symbol)
	}

	// // Initialize our rate controller
	control := NewRateController(8, 3)

	// Make sure we will not hang on bugs that cause infinite loops
	var iterations uint32 = 0

	// Counter for keeping track of the number of decoded symbols
	var decoded uint32 = 0

	for decoded < 100 && iterations < maxIterations {
		// Update loop state
		iterations++
		// Manage the encoder's window
		if control.GenerateData() {
			// Create a new source symbol
			var sourceSymbol = make([]uint8, symbolSize)
			for i := 0; i < len(sourceSymbol); i++ {
				sourceSymbol[i] = uint8(rand.Uint32())
			}

			sourceSymbols[sourceSymbolIndex] = sourceSymbol
			sourceSymbolIndex++

			if encoder.StreamSymbols() == windowSymbols {
				// If window is full - pop a symbol before pushing a new one
				encoder.PopBackSymbol()
			}
			encoder.PushFrontSymbol(&sourceSymbol)
		}
		control.Advance()

		// Uncoded or coded
		coded := rand.Uint32()%2 == 0
		var randomIndex uint64 = 0

		// Choose a seed for this encoding
		seed := rand.Uint64()

		// Encode a symbol
		encoder.SetWindow(encoder.StreamLowerBound(), encoder.StreamSymbols())

		coefficients := make([]uint8, encoder.CoefficientVectorSize())
		symbol := make([]uint8, encoder.SymbolSize())

		if coded {
			encoder.SetSeed(seed)
			encoder.Generate(&coefficients)
			encoder.WriteSymbol(&symbol, &coefficients)
		} else {
			// Warning: This approach is biased towards the lower end.
			maxIndex := encoder.WindowUpperBound() - uint64(1)
			minIndex := encoder.WindowLowerBound()
			randomIndex = rand.Uint64()%(maxIndex-minIndex+uint64(1)) + minIndex

			encoder.WriteSourceSymbol(&symbol, randomIndex)
		}

		if rand.Uint32()%2 == 0 {
			// Simulate 50% packet loss
			continue
		}

		// Move the decoders's window / stream if needed
		for decoder.StreamUpperBound() < encoder.StreamUpperBound() {
			lowerBound := decoder.StreamLowerBound()
			decoderSymbol := GetSymbolWraparound(
				decoderStorage, symbolSize, lowerBound, capacity)

			if decoder.IsSymbolDecoded(lowerBound) {
				decoded++

				// Compare with corresponding source symbol
				sourceSymbol := sourceSymbols[lowerBound]
				c.Assert(sourceSymbol, DeepEquals, decoderSymbol)
			}

			decoder.PopBackSymbol()

			// Moves the decoder's upper bound
			decoder.PushFrontSymbol(&decoderSymbol)
		}

		// Decode the symbol
		decoder.SetWindow(encoder.WindowLowerBound(), encoder.WindowSymbols())

		if coded {
			decoder.SetSeed(seed)
			decoder.Generate(&coefficients)
			decoder.ReadSymbol(&symbol, &coefficients)
		} else {
			decoder.ReadSourceSymbol(&symbol, randomIndex)
		}
	}

	c.Assert(iterations < maxIterations, Equals, true)
	c.Assert(decoded >= 100, Equals, true)
}
