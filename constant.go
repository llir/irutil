package irutil

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"strings"
)

// NewZero returns a new zero value of the given type.
func NewZero(typ types.Type) value.Value {
	switch typ := typ.(type) {
	case *types.IntType:
		return constant.NewInt(typ, 0)
	case *types.FloatType:
		return constant.NewFloat(typ, 0)
	default:
		return constant.NewZeroInitializer(typ)
	}
}

// NewCString returns a new NULL-terminated character array constant based on
// the given UTF-8 string contents.
func NewCString(s string) *constant.CharArray {
	return constant.NewCharArrayFromString(s + "\x00")
}

// NewPString returns a pascal string
func NewPString(s string) *constant.CharArray {
	var sb strings.Builder
	sb.WriteRune(rune(len(s)))
	sb.WriteString(s)
	return constant.NewCharArrayFromString(sb.String())
}
