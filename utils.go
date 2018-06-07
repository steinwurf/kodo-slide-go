package kslide

func GetSymbol(data []uint8, symbolSize uint64, index uint64) []uint8 {
	return data[index*symbolSize : (index+1)*symbolSize]
}

func GetSymbolWraparound(
	data []uint8, symbolSize uint64, index uint64, capacity uint64) []uint8 {
	mappedIndex := index % capacity
	return data[mappedIndex*symbolSize : (mappedIndex+1)*symbolSize]
}
