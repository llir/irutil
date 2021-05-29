package irutil

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
)

// Layout provides configuration about metadata of compiler
type Layout interface {
	SizeOf(typ types.Type) int
}

// DefaultLayout provides a default implementation for size of type, you should create your owned implementation by embedding this structure, and call SizeOf by default
type DefaultLayout struct{}

// SizeOf returns size of types.Type, value is how many bits
func (l DefaultLayout) SizeOf(typ types.Type) int {
	switch typ := typ.(type) {
	case *types.IntType:
		return int(typ.BitSize)
	case *types.FloatType:
		switch typ.Kind {
		case types.FloatKindHalf:
			return 16
		case types.FloatKindFloat:
			return 32
		case types.FloatKindDouble:
			return 64
		case types.FloatKindFP128:
			return 128
		case types.FloatKindX86_FP80:
			return 80
		case types.FloatKindPPC_FP128:
			return 128
		}
	}
	panic(fmt.Sprintf("unimplemented size of this type, %v", typ))
}
