package kslide_test

import (
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
	var sourceSymbols = [maxIterations][symbolSize]uint8{}

	// // Allocate our finite memory for decoding
	// symbol_storage* decoder_storage = symbol_storage_alloc(capacity, symbol_size);

	// // Provide the decoder with storage
	// for (uint32_t i = 0; i < capacity; ++i)
	// {
	//     uint8_t* symbol = symbol_storage_symbol(decoder_storage, i);
	//     kslide_decoder_push_front_symbol(decoder, symbol);
	// }

	// // Initialize our rate controller
	// rate_controller control = rate_controller_init(8, 3);

	// // Make sure we will not hang on bugs that cause infinite loops
	// uint32_t iterations = 0;

	// // Counter for keeping track of the number of decoded symbols
	// uint32_t decoded = 0;

	// while (decoded < 100U && iterations < max_iterations)
	// {
	//     // Update loop state
	//     ++iterations;
	//     // Manage the encoder's window
	//     if (rate_controller_generate_data(control))
	//     {
	//         // Create a new source symbol
	//         uint8_t* source_symbol = (uint8_t*)malloc(symbol_size);
	//         source_symbols[source_symbols_index] = source_symbol;
	//         source_symbols_index++;
	//         assert(source_symbols_index < source_symbol_count);

	//         randomize_buffer(source_symbol, symbol_size);

	//         if (kslide_encoder_stream_symbols(encoder) == window_symbols)
	//         {
	//             // If window is full - pop a symbol before pushing a new one
	//             kslide_encoder_pop_back_symbol(encoder);
	//         }

	//         kslide_encoder_push_front_symbol(encoder, source_symbol);
	//     }

	//     // Uncoded or coded
	//     bool coded = rand() % 2;
	//     uint64_t random_index = 0;

	//     // Choose a seed for this encoding
	//     uint64_t seed = rand();

	//     // Encode a symbol
	//     kslide_encoder_set_window(
	//         encoder,
	//         kslide_encoder_stream_lower_bound(encoder),
	//         kslide_encoder_stream_symbols(encoder));

	//     uint8_t* coefficients = (uint8_t*) malloc(
	//         kslide_encoder_coefficient_vector_size(encoder));
	//     uint8_t* symbol = (uint8_t*)malloc(kslide_encoder_symbol_size(encoder));

	//     if (coded)
	//     {
	//         kslide_encoder_set_seed(encoder, seed);
	//         kslide_encoder_generate(encoder, coefficients);
	//         kslide_encoder_write_symbol(encoder, symbol, coefficients);
	//     }
	//     else
	//     {
	//         // Warning: This approach is biased towards the lower end.
	//         uint64_t max_index = kslide_encoder_window_upper_bound(encoder) - 1;
	//         uint64_t min_index = kslide_encoder_window_lower_bound(encoder);
	//         random_index = rand() % (max_index - min_index + 1) + min_index;

	//         kslide_encoder_write_source_symbol(encoder, symbol, random_index);
	//     }

	//     rate_controller_advance(control);

	//     if (rand() % 2)
	//     {
	//         free(coefficients);
	//         free(symbol);
	//         // Simulate 50% packet loss
	//         continue;
	//     }

	//     // Move the decoders's window / stream if needed
	//     while (kslide_decoder_stream_upper_bound(decoder) <
	//            kslide_encoder_stream_upper_bound(encoder))
	//     {
	//         uint64_t lower_bound = kslide_decoder_stream_lower_bound(decoder);
	//         uint8_t* decoder_symbol = symbol_storage_symbol(decoder_storage, lower_bound);

	//         if (kslide_decoder_is_symbol_decoded(decoder, lower_bound))
	//         {
	//             ++decoded;

	//             // Compare with corresponding source symbol
	//             uint8_t* source_symbol = source_symbols[lower_bound];
	//             EXPECT_EQ(0, memcmp(decoder_symbol, source_symbol, symbol_size));
	//         }

	//         uint64_t pop_index = kslide_decoder_pop_back_symbol(decoder);
	//         assert(pop_index == lower_bound);

	//         // Moves the decoder's upper bound
	//         kslide_decoder_push_front_symbol(decoder, decoder_symbol);
	//     }

	//     // Decode the symbol
	//     kslide_decoder_set_window(
	//         decoder,
	//         kslide_encoder_window_lower_bound(encoder),
	//         kslide_encoder_window_symbols(encoder));

	//     if (coded)
	//     {
	//         kslide_decoder_set_seed(decoder, seed);
	//         kslide_decoder_generate(decoder, coefficients);
	//         kslide_decoder_read_symbol(decoder, symbol, coefficients);
	//     }
	//     else
	//     {
	//         kslide_decoder_read_source_symbol(decoder, symbol, random_index);
	//     }

	//     free(coefficients);
	//     free(symbol);
	// }
	// SCOPED_TRACE(testing::Message() << "decoded = " << decoded);

	// EXPECT_LT(iterations, max_iterations);
	// EXPECT_GE(decoded, 100U);

	// symbol_storage_free(decoder_storage);

	// for (uint32_t i = 0; i < source_symbol_count; i++)
	// {
	//     free(source_symbols[i]);
	// }
	// free(source_symbols);
	// kslide_delete_encoder(encoder);
	// kslide_delete_encoder_factory(encoder_factory);
	// kslide_delete_decoder(decoder);
	// kslide_delete_decoder_factory(decoder_factory);
}
