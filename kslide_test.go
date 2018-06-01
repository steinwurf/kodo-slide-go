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
	var symbolSize uint32 = 1300
	decoderFactory := NewDecoderFactory()

	c.Assert(decoderFactory.SymbolSize(), Equals, symbolSize)

	var newSymbolSize uint32 = 300
	decoderFactory.SetSymbolSize(newSymbolSize)
	c.Assert(decoderFactory.SymbolSize(), Equals, newSymbolSize)
	c.Assert(decoderFactory.Field(), Equals, Binary8)

	for _, field := range [4]int32{Binary, Binary4, Binary8, Binary16} {
		decoderFactory.SetField(field)
		c.Assert(decoderFactory.Field(), Equals, field)
	}
}

func (s *MySuite) TestDecoder(c *C) {
	var symbolSize uint32 = 1300
	decoderFactory := NewDecoderFactory()
	decoder := decoderFactory.Build()
	c.Assert(decoder.SymbolSize(), Equals, symbolSize)
	c.Assert(decoder.StreamSymbols(), Equals, uint32(0))
	c.Assert(decoder.StreamLowerBound(), Equals, uint32(0))
}

func (s *MySuite) TestEncoderFactory(c *C) {
	var symbolSize uint32 = 1300
	encoderFactory := NewEncoderFactory()

	c.Assert(encoderFactory.SymbolSize(), Equals, symbolSize)

	var newSymbolSize uint32 = 300
	encoderFactory.SetSymbolSize(newSymbolSize)
	c.Assert(newSymbolSize, Equals, encoderFactory.SymbolSize())
	c.Assert(encoderFactory.Field(), Equals, Binary8)

	for _, field := range [4]int32{Binary, Binary4, Binary8, Binary16} {
		encoderFactory.SetField(field)
		c.Assert(encoderFactory.Field(), Equals, field)
	}
}

func (s *MySuite) TestEncoder(c *C) {
	var symbolSize uint32 = 1300
	encoderFactory := NewEncoderFactory()
	encoder := encoderFactory.Build()
	c.Assert(encoder.SymbolSize(), Equals, symbolSize)
}
