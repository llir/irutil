package irutil

import (
	"testing"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/stretchr/testify/assert"
)

func TestIntAdd(t *testing.T) {
	testCases := []struct {
		name     string
		expected constant.Constant
		from     constant.Constant
	}{
		{"IntAdd", constant.NewInt(types.I32, 3),
			constant.NewAdd(constant.NewInt(types.I32, 1), constant.NewInt(types.I32, 2))},
		{"IntSub", constant.NewInt(types.I32, -1),
			constant.NewSub(constant.NewInt(types.I32, 1), constant.NewInt(types.I32, 2))},
		{"IntMul", constant.NewInt(types.I32, 2),
			constant.NewMul(constant.NewInt(types.I32, 1), constant.NewInt(types.I32, 2))},
		{"IntDiv", constant.NewInt(types.I32, 1),
			constant.NewSDiv(constant.NewInt(types.I32, 2), constant.NewInt(types.I32, 2))},
		{"IntDiv", constant.NewInt(types.I32, 1),
			constant.NewSDiv(constant.NewInt(types.I32, 3), constant.NewInt(types.I32, 2))},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, Simplify(testCase.from))
		})
	}
}
