package irutil

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

// ResetTypes resets the (cached) types of instructions in the given function.
func ResetTypes(f *ir.Func) {
	for _, b := range f.Blocks {
		for _, inst := range b.Insts {
			valueInst, ok := inst.(value.Value)
			if !ok {
				continue
			}
			resetType(valueInst)
		}
	}
}

// resetType resets the (cached) type of the given instruction.
func resetType(inst value.Value) {
	switch inst := inst.(type) {
	// Unary instructions
	case *ir.InstFNeg:
		inst.Typ = nil
	// Binary instructions
	case *ir.InstAdd:
		inst.Typ = nil
	case *ir.InstFAdd:
		inst.Typ = nil
	case *ir.InstSub:
		inst.Typ = nil
	case *ir.InstFSub:
		inst.Typ = nil
	case *ir.InstMul:
		inst.Typ = nil
	case *ir.InstFMul:
		inst.Typ = nil
	case *ir.InstUDiv:
		inst.Typ = nil
	case *ir.InstSDiv:
		inst.Typ = nil
	case *ir.InstFDiv:
		inst.Typ = nil
	case *ir.InstURem:
		inst.Typ = nil
	case *ir.InstSRem:
		inst.Typ = nil
	case *ir.InstFRem:
		inst.Typ = nil
	// Bitwise instructions
	case *ir.InstShl:
		inst.Typ = nil
	case *ir.InstLShr:
		inst.Typ = nil
	case *ir.InstAShr:
		inst.Typ = nil
	case *ir.InstAnd:
		inst.Typ = nil
	case *ir.InstOr:
		inst.Typ = nil
	case *ir.InstXor:
		inst.Typ = nil
	// Vector instructions
	case *ir.InstExtractElement:
		inst.Typ = nil
	case *ir.InstInsertElement:
		inst.Typ = nil
	case *ir.InstShuffleVector:
		inst.Typ = nil
	// Aggregate instructions
	case *ir.InstExtractValue:
		inst.Typ = nil
	case *ir.InstInsertValue:
		inst.Typ = nil
	// Memory instructions
	case *ir.InstAlloca:
		inst.Typ = nil
	case *ir.InstLoad:
		// type not cached.
	case *ir.InstCmpXchg:
		inst.Typ = nil
	case *ir.InstAtomicRMW:
		inst.Typ = nil
	case *ir.InstGetElementPtr:
		inst.Typ = nil
	// Conversion instructions
	case *ir.InstTrunc:
		// type not cached.
	case *ir.InstZExt:
		// type not cached.
	case *ir.InstSExt:
		// type not cached.
	case *ir.InstFPTrunc:
		// type not cached.
	case *ir.InstFPExt:
		// type not cached.
	case *ir.InstFPToUI:
		// type not cached.
	case *ir.InstFPToSI:
		// type not cached.
	case *ir.InstUIToFP:
		// type not cached.
	case *ir.InstSIToFP:
		// type not cached.
	case *ir.InstPtrToInt:
		// type not cached.
	case *ir.InstIntToPtr:
		// type not cached.
	case *ir.InstBitCast:
		// type not cached.
	case *ir.InstAddrSpaceCast:
		// type not cached.
	// Other instructions
	case *ir.InstICmp:
		inst.Typ = nil
	case *ir.InstFCmp:
		inst.Typ = nil
	case *ir.InstPhi:
		inst.Typ = nil
	case *ir.InstSelect:
		inst.Typ = nil
	case *ir.InstCall:
		inst.Typ = nil
	case *ir.InstVAArg:
		// type not cached.
	case *ir.InstLandingPad:
		// type not cached.
	case *ir.InstCatchPad:
		// type not cached.
	case *ir.InstCleanupPad:
		// type not cached.
	default:
		panic(fmt.Errorf("support for instruction type %T not yet implemented", inst))
	}
}
