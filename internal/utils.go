package papersave

import (
	"math"
)

func computeChunks(dataSize, chunkSize int) int {
	ret := int(math.Ceil(float64(dataSize) / float64(chunkSize)))
	return ret
}
