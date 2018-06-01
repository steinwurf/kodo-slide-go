package kslide_test

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/steinwurf/kodo-slide-go"
)

type payload struct {
	coefficients []uint8
	symbol       []uint8
}

func BenchmarkEncoder(b *testing.B) {
	// Seed random number generator to produce different results every time
	rand.Seed(time.Now().UTC().UnixNano())

	// Set the capacity of the decoder (this is the number of encoded symbols
	// that are used in the decoding process).
	symbols := uint32(10)

	// The size of a symbol in bytes
	symbolSize := uint32(1600)
	fields := map[string]int32{
		"Binary":   Binary,
		"Binary4":  Binary4,
		"Binary8":  Binary8,
		"Binary16": Binary16,
	}

	for name, field := range fields {
		encoderFactory := NewEncoderFactory()

		encoderFactory.SetField(field)
		encoderFactory.SetSymbolSize(symbolSize)

		// Allocate some data to encode.
		encoderStorage := make([][]uint8, symbols)
		for i := uint32(0); i < symbols; i++ {
			symbolData := make([]uint8, encoderFactory.SymbolSize())
			// Just for fun - fill the data with random data
			for j := range symbolData {
				symbolData[j] = uint8(rand.Uint32())
			}
			encoderStorage[i] = symbolData
		}

		var payloads []payload
		b.Run(name+"Encode", func(b *testing.B) { payloads = encodeData(b, encoderFactory, &encoderStorage) })

		decoderFactory := NewDecoderFactory()

		decoderFactory.SetField(field)
		decoderFactory.SetSymbolSize(symbolSize)

		decoderStorage := make([][]uint8, symbols)
		for i := uint32(0); i < symbols; i++ {
			decoderStorage[i] = make([]uint8, decoderFactory.SymbolSize())
		}
		b.Run(name+"Decode", func(b *testing.B) { decodeData(b, decoderFactory, &decoderStorage, &payloads) })

		// var success bool = true
		// // Check if we properly decoded the data
		// for i, v := range dataIn {
		// 	if v != dataOut[i] {
		// 		success = false
		// 		break
		// 	}
		// }

		// if success == true {
		// 	fmt.Println("Data decoded correctly")
		// } else {
		// 	fmt.Println("Decoding failed")
		// 	b.FailNow()
		// }
	}
}

func encodeData(
	b *testing.B, encoderFactory *EncoderFactory, encoderStorage *[][]uint8) []payload {

	var encoder *Encoder
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder = encoderFactory.Build()

		for encoder.StreamSymbols() < uint32(len(*encoderStorage)) {
			encoder.PushFrontSymbol(&(*encoderStorage)[encoder.StreamSymbols()])

			encoder.SetWindow(encoder.StreamLowerBound(),
				encoder.StreamSymbols())

			coefficients := make([]uint8, encoder.CoefficientVectorSize())
			symbol := make([]uint8, encoder.SymbolSize())

			encoder.SetSeed(rand.Uint32())
			encoder.Generate(&coefficients)
			encoder.WriteSymbol(&symbol, &coefficients)
		}
	}

	payloads := make([]payload, uint32(len(*encoderStorage)))
	for encoder.StreamSymbols() < uint32(len(*encoderStorage)) {

		data := (*encoderStorage)[encoder.StreamSymbols()]
		encoder.PushFrontSymbol(&data)

		encoder.SetWindow(encoder.StreamLowerBound(),
			encoder.StreamSymbols())

		p := payload{
			coefficients: make([]uint8, encoder.CoefficientVectorSize()),
			symbol:       make([]uint8, encoder.SymbolSize())}

		encoder.SetSeed(rand.Uint32())
		encoder.Generate(&p.coefficients)
		encoder.WriteSymbol(&p.symbol, &p.coefficients)
		payloads[encoder.StreamSymbols()-1] = p
	}

	return payloads
}

func decodeData(
	b *testing.B,
	decoderFactory *DecoderFactory,
	decoderStorage *[][]uint8,
	payloads *[]payload) {

	// f := func(decoder *Decoder, payloads []payload) {

	// }

	var decoder *Decoder
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoder = decoderFactory.Build()
		for decoder.StreamSymbols() < uint32(len(*decoderStorage)) {
			index := decoder.StreamSymbols()
			decoder.PushFrontSymbol(&(*decoderStorage)[index])
		}
		for _, payload := range *payloads {
			// 	if decoder.IsComplete() {
			// 		break
			//	}
			decoder.ReadSymbol(&payload.symbol, &payload.coefficients)
		}
		// f(decoder, *payloads)
	}
}

// 	kodo_slide::encoder::factory encoder_factory;
// 	kodo_slide::decoder::factory decoder_factory;

// 	encoder_factory.set_symbol_size(symbol_size);
// 	decoder_factory.set_symbol_size(symbol_size);

// 	auto encoder = encoder_factory.build();
// 	auto decoder = decoder_factory.build();

// 	// encoder.set_trace_stdout();
// 	// encoder.set_zone_prefix("encoder");
// 	// decoder.set_trace_stdout();

// 	// Allocate our finite memory for encoding and decoding
// 	symbol_storage decoder_storage(symbols, symbol_size);
// 	symbol_storage encoder_storage(symbols, symbol_size);

// 	// Provide the decoder with the storage
// 	for (uint32_t i = 0; i < symbols; ++i)
// 	{
// 		uint8_t* symbol = decoder_storage.symbol(i);
// 		decoder.push_front_symbol(symbol);
// 	}

// 	// Fill the encoder symbols storage
// 	for (uint32_t i = 0; i < symbols; ++i)
// 	{
// 		uint8_t* symbol = encoder_storage.symbol(i);
// 		std::generate_n(symbol, symbol_size, rand);
// 	}

// 	std::vector<uint8_t> coefficients;
// 	std::vector<uint8_t> symbol;

// 	// Make sure we will not hang on bugs that cause infinite loops
// 	uint32_t iterations = 0;

// 	while (decoder.symbols_decoded() < symbols && iterations < 1000U)
// 	{
// 		if (encoder.stream_symbols() < symbols && rand() % 2)
// 		{
// 			uint8_t* s = encoder_storage.symbol(encoder.stream_symbols());
// 			encoder.push_front_symbol(s);
// 		}

// 		if (encoder.stream_symbols() == 0)
// 		{
// 			continue;
// 		}

// 		encoder.set_window(encoder.stream_lower_bound(),
// 							encoder.stream_symbols());
// 		decoder.set_window(encoder.window_lower_bound(),
// 							encoder.window_symbols());

// 		coefficients.resize(encoder.coefficient_vector_size());
// 		symbol.resize(encoder.symbol_size());

// 		uint32_t seed = rand();

// 		encoder.set_seed(seed);
// 		encoder.generate(coefficients.data());

// 		encoder.write_symbol(symbol.data(), coefficients.data());
// 		decoder.read_symbol(symbol.data(), coefficients.data());

// 		++iterations;
// 	}

// 	EXPECT_EQ(encoder_storage, decoder_storage);
// }
