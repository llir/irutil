package irutil

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

// User is a user of a value.
//
// A User has one of the following underlying types.
//
//    ir.Instruction   // https://godoc.org/github.com/llir/llvm/ir#Instruction
//    ir.Terminator    // https://godoc.org/github.com/llir/llvm/ir#Terminator
type User interface {
	ir.LLStringer
}

// Use tracks the use of a value and its parent instruction or terminator.
type Use struct {
	// Use of value.
	Val *value.Value
	// Parent instruction or terminator which uses value as operand.
	User User // ir.Instruction or ir.Terminator
}

// Replace replaces the use of a value with the new value.
func (use Use) Replace(new value.Value) {
	*use.Val = new
}

// FuncUses returns value uses of instructions and terminators in the given
// function. To avoid memory allocation, an optional backing slice may be
// provided; if nil or insufficient space, a new backing slice will be
// allocated.
func FuncUses(backing []Use, f *ir.Func) []Use {
	for _, block := range f.Blocks {
		for _, inst := range block.Insts {
			backing = InstUses(backing, inst)
		}
		backing = TermUses(backing, block.Term)
	}
	return backing
}

// TODO: consider renaming InstUses to InstOperands.

// InstUses returns value uses of the given instruction. To avoid memory
// allocation, an optional backing slice may be provided; if nil or insufficient
// space, a new backing slice will be allocated.
func InstUses(backing []Use, inst ir.Instruction) []Use {
	switch inst := inst.(type) {
	// Unary instructions
	case *ir.InstFNeg:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		return backing
	// Binary instructions
	case *ir.InstAdd:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstFAdd:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstSub:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstFSub:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstMul:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstFMul:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstUDiv:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstSDiv:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstFDiv:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstURem:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstSRem:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstFRem:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	// Bitwise instructions
	case *ir.InstShl:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstLShr:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstAShr:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstAnd:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstOr:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstXor:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	// Vector instructions
	case *ir.InstExtractElement:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Index, User: inst})
		return backing
	case *ir.InstInsertElement:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Elem, User: inst})
		backing = append(backing, Use{Val: &inst.Index, User: inst})
		return backing
	case *ir.InstShuffleVector:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		backing = append(backing, Use{Val: &inst.Mask, User: inst})
		return backing
	// Aggregate instructions
	case *ir.InstExtractValue:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		return backing
	case *ir.InstInsertValue:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Elem, User: inst})
		return backing
	// Memory instructions
	case *ir.InstAlloca:
		if inst.NElems != nil {
			backing = append(backing, Use{Val: &inst.NElems, User: inst})
		}
		return backing
	case *ir.InstLoad:
		backing = append(backing, Use{Val: &inst.Src, User: inst})
		return backing
	case *ir.InstStore:
		backing = append(backing, Use{Val: &inst.Src, User: inst})
		backing = append(backing, Use{Val: &inst.Dst, User: inst})
		return backing
	case *ir.InstFence:
		// no value operands.
		return backing
	case *ir.InstCmpXchg:
		backing = append(backing, Use{Val: &inst.Ptr, User: inst})
		backing = append(backing, Use{Val: &inst.Cmp, User: inst})
		backing = append(backing, Use{Val: &inst.New, User: inst})
		return backing
	case *ir.InstAtomicRMW:
		backing = append(backing, Use{Val: &inst.Dst, User: inst})
		backing = append(backing, Use{Val: &inst.X, User: inst})
		return backing
	case *ir.InstGetElementPtr:
		backing = append(backing, Use{Val: &inst.Src, User: inst})
		for i := range inst.Indices {
			backing = append(backing, Use{Val: &inst.Indices[i], User: inst})
		}
		return backing
	// Conversion instructions
	case *ir.InstTrunc:
		backing = append(backing, Use{Val: &inst.From, User: inst})
		return backing
	case *ir.InstZExt:
		backing = append(backing, Use{Val: &inst.From, User: inst})
		return backing
	case *ir.InstSExt:
		backing = append(backing, Use{Val: &inst.From, User: inst})
		return backing
	case *ir.InstFPTrunc:
		backing = append(backing, Use{Val: &inst.From, User: inst})
		return backing
	case *ir.InstFPExt:
		backing = append(backing, Use{Val: &inst.From, User: inst})
		return backing
	case *ir.InstFPToUI:
		backing = append(backing, Use{Val: &inst.From, User: inst})
		return backing
	case *ir.InstFPToSI:
		backing = append(backing, Use{Val: &inst.From, User: inst})
		return backing
	case *ir.InstUIToFP:
		backing = append(backing, Use{Val: &inst.From, User: inst})
		return backing
	case *ir.InstSIToFP:
		backing = append(backing, Use{Val: &inst.From, User: inst})
		return backing
	case *ir.InstPtrToInt:
		backing = append(backing, Use{Val: &inst.From, User: inst})
		return backing
	case *ir.InstIntToPtr:
		backing = append(backing, Use{Val: &inst.From, User: inst})
		return backing
	case *ir.InstBitCast:
		backing = append(backing, Use{Val: &inst.From, User: inst})
		return backing
	case *ir.InstAddrSpaceCast:
		backing = append(backing, Use{Val: &inst.From, User: inst})
		return backing
	// Other instructions
	case *ir.InstICmp:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstFCmp:
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstPhi:
		for i := range inst.Incs {
			backing = append(backing, Use{Val: &inst.Incs[i].X, User: inst})
			backing = append(backing, Use{Val: &inst.Incs[i].Pred, User: inst})
		}
		return backing
	case *ir.InstSelect:
		backing = append(backing, Use{Val: &inst.Cond, User: inst})
		backing = append(backing, Use{Val: &inst.X, User: inst})
		backing = append(backing, Use{Val: &inst.Y, User: inst})
		return backing
	case *ir.InstCall:
		backing = append(backing, Use{Val: &inst.Callee, User: inst})
		for i := range inst.Args {
			backing = append(backing, Use{Val: &inst.Args[i], User: inst})
		}
		for i := range inst.OperandBundles {
			for j := range inst.OperandBundles[i].Inputs {
				backing = append(backing, Use{Val: &inst.OperandBundles[i].Inputs[j], User: inst})
			}
		}
		return backing
	case *ir.InstVAArg:
		backing = append(backing, Use{Val: &inst.ArgList, User: inst})
		return backing
	case *ir.InstLandingPad:
		for i := range inst.Clauses {
			backing = append(backing, Use{Val: &inst.Clauses[i].X, User: inst})
		}
		return backing
	case *ir.InstCatchPad:
		backing = append(backing, Use{Val: &inst.Scope, User: inst})
		for i := range inst.Args {
			backing = append(backing, Use{Val: &inst.Args[i], User: inst})
		}
		return backing
	case *ir.InstCleanupPad:
		backing = append(backing, Use{Val: &inst.Scope, User: inst})
		for i := range inst.Args {
			backing = append(backing, Use{Val: &inst.Args[i], User: inst})
		}
		return backing
	// TODO: figure out how to handle user defined instructions (e.g. Comment).
	// TODO: remove *Comment case and add support for more general way of
	// identifying user-defined instructions.
	case *Comment:
		// no value operands.
		return backing
	default:
		panic(fmt.Errorf("support for instruction %T not yet implemented", inst))
	}
}

// TermUses returns value uses of the given terminator. To avoid memory
// allocation, an optional backing slice may be provided; if nil or insufficient
// space, a new backing slice will be allocated.
func TermUses(backing []Use, term ir.Terminator) []Use {
	switch term := term.(type) {
	case *ir.TermRet:
		if term.X != nil {
			backing = append(backing, Use{Val: &term.X, User: term})
		}
		return backing
	case *ir.TermBr:
		backing = append(backing, Use{Val: &term.Target, User: term})
		return backing
	case *ir.TermCondBr:
		backing = append(backing, Use{Val: &term.Cond, User: term})
		backing = append(backing, Use{Val: &term.TargetTrue, User: term})
		backing = append(backing, Use{Val: &term.TargetFalse, User: term})
		return backing
	case *ir.TermSwitch:
		backing = append(backing, Use{Val: &term.X, User: term})
		backing = append(backing, Use{Val: &term.TargetDefault, User: term})
		for i := range term.Cases {
			backing = append(backing, Use{Val: &term.Cases[i].X, User: term})
			backing = append(backing, Use{Val: &term.Cases[i].Target, User: term})
		}
		return backing
	case *ir.TermIndirectBr:
		backing = append(backing, Use{Val: &term.Addr, User: term})
		for i := range term.ValidTargets {
			backing = append(backing, Use{Val: &term.ValidTargets[i], User: term})
		}
		return backing
	case *ir.TermInvoke:
		backing = append(backing, Use{Val: &term.Invokee, User: term})
		for i := range term.Args {
			backing = append(backing, Use{Val: &term.Args[i], User: term})
		}
		backing = append(backing, Use{Val: &term.Normal, User: term})
		backing = append(backing, Use{Val: &term.Exception, User: term})
		for i := range term.OperandBundles {
			for j := range term.OperandBundles[i].Inputs {
				backing = append(backing, Use{Val: &term.OperandBundles[i].Inputs[j], User: term})
			}
		}
		return backing
	case *ir.TermCallBr:
		backing = append(backing, Use{Val: &term.Callee, User: term})
		for i := range term.Args {
			backing = append(backing, Use{Val: &term.Args[i], User: term})
		}
		backing = append(backing, Use{Val: &term.Normal, User: term})
		for i := range term.Others {
			backing = append(backing, Use{Val: &term.Others[i], User: term})
		}
		for i := range term.OperandBundles {
			for j := range term.OperandBundles[i].Inputs {
				backing = append(backing, Use{Val: &term.OperandBundles[i].Inputs[j], User: term})
			}
		}
		return backing
	case *ir.TermResume:
		backing = append(backing, Use{Val: &term.X, User: term})
		return backing
	case *ir.TermCatchSwitch:
		backing = append(backing, Use{Val: &term.Scope, User: term})
		for i := range term.Handlers {
			backing = append(backing, Use{Val: &term.Handlers[i], User: term})
		}
		backing = append(backing, Use{Val: &term.UnwindTarget, User: term})
		return backing
	case *ir.TermCatchRet:
		backing = append(backing, Use{Val: &term.From, User: term})
		backing = append(backing, Use{Val: &term.To, User: term})
		return backing
	case *ir.TermCleanupRet:
		backing = append(backing, Use{Val: &term.From, User: term})
		backing = append(backing, Use{Val: &term.UnwindTarget, User: term})
		return backing
	case *ir.TermUnreachable:
		// no value operands.
		return backing
	default:
		panic(fmt.Errorf("support for terminator %T not yet implemented", term))
	}
}
