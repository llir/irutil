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
func (use *Use) Replace(new value.Value) {
	*use.Val = new
}

// FuncUses returns value uses of instructions and terminators in the given
// function.
func FuncUses(f *ir.Func) []*Use {
	var uses []*Use
	for _, block := range f.Blocks {
		for _, inst := range block.Insts {
			us := InstUses(inst)
			uses = append(uses, us...)
		}
		us := TermUses(block.Term)
		uses = append(uses, us...)
	}
	return uses
}

// TODO: consider renaming InstUses to InstOperands.

// InstUses returns value uses of the given instruction.
func InstUses(inst ir.Instruction) []*Use {
	switch inst := inst.(type) {
	// Unary instructions
	case *ir.InstFNeg:
		return []*Use{
			{Val: &inst.X, User: inst},
		}
	// Binary instructions
	case *ir.InstAdd:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstFAdd:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstSub:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstFSub:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstMul:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstFMul:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstUDiv:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstSDiv:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstFDiv:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstURem:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstSRem:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstFRem:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	// Bitwise instructions
	case *ir.InstShl:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstLShr:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstAShr:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstAnd:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstOr:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstXor:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	// Vector instructions
	case *ir.InstExtractElement:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Index, User: inst},
		}
	case *ir.InstInsertElement:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Elem, User: inst},
			{Val: &inst.Index, User: inst},
		}
	case *ir.InstShuffleVector:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
			{Val: &inst.Mask, User: inst},
		}
	// Aggregate instructions
	case *ir.InstExtractValue:
		return []*Use{
			{Val: &inst.X, User: inst},
		}
	case *ir.InstInsertValue:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Elem, User: inst},
		}
	// Memory instructions
	case *ir.InstAlloca:
		var uses []*Use
		if inst.NElems != nil {
			use := &Use{Val: &inst.NElems, User: inst}
			uses = append(uses, use)
		}
		return uses
	case *ir.InstLoad:
		return []*Use{
			{Val: &inst.Src, User: inst},
		}
	case *ir.InstStore:
		return []*Use{
			{Val: &inst.Src, User: inst},
			{Val: &inst.Dst, User: inst},
		}
	case *ir.InstFence:
		// no value operands.
		return nil
	case *ir.InstCmpXchg:
		return []*Use{
			{Val: &inst.Ptr, User: inst},
			{Val: &inst.Cmp, User: inst},
			{Val: &inst.New, User: inst},
		}
	case *ir.InstAtomicRMW:
		return []*Use{
			{Val: &inst.Dst, User: inst},
			{Val: &inst.X, User: inst},
		}
	case *ir.InstGetElementPtr:
		uses := []*Use{
			{Val: &inst.Src, User: inst},
		}
		for i := range inst.Indices {
			use := &Use{Val: &inst.Indices[i], User: inst}
			uses = append(uses, use)
		}
		return uses
	// Conversion instructions
	case *ir.InstTrunc:
		return []*Use{
			{Val: &inst.From, User: inst},
		}
	case *ir.InstZExt:
		return []*Use{
			{Val: &inst.From, User: inst},
		}
	case *ir.InstSExt:
		return []*Use{
			{Val: &inst.From, User: inst},
		}
	case *ir.InstFPTrunc:
		return []*Use{
			{Val: &inst.From, User: inst},
		}
	case *ir.InstFPExt:
		return []*Use{
			{Val: &inst.From, User: inst},
		}
	case *ir.InstFPToUI:
		return []*Use{
			{Val: &inst.From, User: inst},
		}
	case *ir.InstFPToSI:
		return []*Use{
			{Val: &inst.From, User: inst},
		}
	case *ir.InstUIToFP:
		return []*Use{
			{Val: &inst.From, User: inst},
		}
	case *ir.InstSIToFP:
		return []*Use{
			{Val: &inst.From, User: inst},
		}
	case *ir.InstPtrToInt:
		return []*Use{
			{Val: &inst.From, User: inst},
		}
	case *ir.InstIntToPtr:
		return []*Use{
			{Val: &inst.From, User: inst},
		}
	case *ir.InstBitCast:
		return []*Use{
			{Val: &inst.From, User: inst},
		}
	case *ir.InstAddrSpaceCast:
		return []*Use{
			{Val: &inst.From, User: inst},
		}
	// Other instructions
	case *ir.InstICmp:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstFCmp:
		return []*Use{
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstPhi:
		var uses []*Use
		for i := range inst.Incs {
			use := &Use{Val: &inst.Incs[i].X, User: inst}
			uses = append(uses, use)
			use = &Use{Val: &inst.Incs[i].Pred, User: inst}
			uses = append(uses, use)
		}
		return uses
	case *ir.InstSelect:
		return []*Use{
			{Val: &inst.Cond, User: inst},
			{Val: &inst.X, User: inst},
			{Val: &inst.Y, User: inst},
		}
	case *ir.InstCall:
		uses := []*Use{
			{Val: &inst.Callee, User: inst},
		}
		for i := range inst.Args {
			use := &Use{Val: &inst.Args[i], User: inst}
			uses = append(uses, use)
		}
		for i := range inst.OperandBundles {
			for j := range inst.OperandBundles[i].Inputs {
				use := &Use{Val: &inst.OperandBundles[i].Inputs[j], User: inst}
				uses = append(uses, use)
			}
		}
		return uses
	case *ir.InstVAArg:
		return []*Use{
			{Val: &inst.ArgList, User: inst},
		}
	case *ir.InstLandingPad:
		var uses []*Use
		for i := range inst.Clauses {
			use := &Use{Val: &inst.Clauses[i].X, User: inst}
			uses = append(uses, use)
		}
		return uses
	case *ir.InstCatchPad:
		uses := []*Use{
			{Val: &inst.Scope, User: inst},
		}
		for i := range inst.Args {
			use := &Use{Val: &inst.Args[i], User: inst}
			uses = append(uses, use)
		}
		return uses
	case *ir.InstCleanupPad:
		uses := []*Use{
			{Val: &inst.Scope, User: inst},
		}
		for i := range inst.Args {
			use := &Use{Val: &inst.Args[i], User: inst}
			uses = append(uses, use)
		}
		return uses
	// TODO: figure out how to handle user defined instructions (e.g. Comment).
	// TODO: remove *Comment case and add support for more general way of
	// identifying user-defined instructions.
	case *Comment:
		// no value operands.
		return nil
	default:
		panic(fmt.Errorf("support for instruction %T not yet implemented", inst))
	}
}

// TermUses returns value uses of the given terminator.
func TermUses(term ir.Terminator) []*Use {
	switch term := term.(type) {
	case *ir.TermRet:
		var uses []*Use
		if term.X != nil {
			use := &Use{Val: &term.X, User: term}
			uses = append(uses, use)
		}
		return uses
	case *ir.TermBr:
		return []*Use{
			{Val: &term.Target, User: term},
		}
	case *ir.TermCondBr:
		return []*Use{
			{Val: &term.Cond, User: term},
			{Val: &term.TargetTrue, User: term},
			{Val: &term.TargetFalse, User: term},
		}
	case *ir.TermSwitch:
		uses := []*Use{
			{Val: &term.X, User: term},
			{Val: &term.TargetDefault, User: term},
		}
		for i := range term.Cases {
			use := &Use{Val: &term.Cases[i].X, User: term}
			uses = append(uses, use)
			use = &Use{Val: &term.Cases[i].Target, User: term}
			uses = append(uses, use)
		}
		return uses
	case *ir.TermIndirectBr:
		uses := []*Use{
			{Val: &term.Addr, User: term},
		}
		for i := range term.ValidTargets {
			use := &Use{Val: &term.ValidTargets[i], User: term}
			uses = append(uses, use)
		}
		return uses
	case *ir.TermInvoke:
		uses := []*Use{
			{Val: &term.Invokee, User: term},
		}
		for i := range term.Args {
			use := &Use{Val: &term.Args[i], User: term}
			uses = append(uses, use)
		}
		use := &Use{Val: &term.Normal, User: term}
		uses = append(uses, use)
		use = &Use{Val: &term.Exception, User: term}
		uses = append(uses, use)
		for i := range term.OperandBundles {
			for j := range term.OperandBundles[i].Inputs {
				use := &Use{Val: &term.OperandBundles[i].Inputs[j], User: term}
				uses = append(uses, use)
			}
		}
		return uses
	case *ir.TermCallBr:
		uses := []*Use{
			{Val: &term.Callee, User: term},
		}
		for i := range term.Args {
			use := &Use{Val: &term.Args[i], User: term}
			uses = append(uses, use)
		}
		use := &Use{Val: &term.Normal, User: term}
		uses = append(uses, use)
		for i := range term.Others {
			use := &Use{Val: &term.Others[i], User: term}
			uses = append(uses, use)
		}
		for i := range term.OperandBundles {
			for j := range term.OperandBundles[i].Inputs {
				use := &Use{Val: &term.OperandBundles[i].Inputs[j], User: term}
				uses = append(uses, use)
			}
		}
		return uses
	case *ir.TermResume:
		return []*Use{
			{Val: &term.X, User: term},
		}
	case *ir.TermCatchSwitch:
		uses := []*Use{
			{Val: &term.Scope, User: term},
		}
		for i := range term.Handlers {
			use := &Use{Val: &term.Handlers[i], User: term}
			uses = append(uses, use)
		}
		use := &Use{Val: &term.UnwindTarget, User: term}
		uses = append(uses, use)
		return uses
	case *ir.TermCatchRet:
		return []*Use{
			{Val: &term.From, User: term},
			{Val: &term.To, User: term},
		}
	case *ir.TermCleanupRet:
		return []*Use{
			{Val: &term.From, User: term},
			{Val: &term.UnwindTarget, User: term},
		}
	case *ir.TermUnreachable:
		// no value operands.
		return nil
	default:
		panic(fmt.Errorf("support for terminator %T not yet implemented", term))
	}
}
