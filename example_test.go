package irutil_test

import (
	"fmt"
	"log"

	"github.com/llir/irutil"
	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/pkg/errors"
)

func Example() {
	// Parse LLVM IR module.
	const module = `
define void @g() {
	; potential side effects
	ret void
}

define i32 @f(i32 %x, i32 %y) {
; <label>:0
	%1 = add i32 %x, %y
	%2 = xor i32 %1, 10
	%3 = sub i32 %x, %y
	call void @g()
	%4 = mul i32 3, 5
	ret i32 %3
}`
	m, err := asm.ParseString("foo.ll", module)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	// Perform dead code elimination (dce).
	for _, f := range m.Funcs {
		dce(f)
	}
	// Print LLVM IR module after dce.
	fmt.Println(m)

	// Output:
	//
	// define void @g() {
	// ; <label>:0
	// 	ret void
	// }
	//
	// define i32 @f(i32 %x, i32 %y) {
	// ; <label>:0
	// 	%1 = sub i32 %x, %y
	// 	call void @g()
	// 	ret i32 %1
	// }
}

// dce performs dead code elimination on f, pruning value instructions with
// unused results in f. Call instructions are never pruned as they may have side
// effects.
func dce(f *ir.Func) error {
	// Prune until we reach fixed-point. If we prune an unused value instruction,
	// this may result in another local variable becoming unused.
	for {
		// Get names of local variables used in f.
		usedIdents, err := getUsedIdents(f)
		if err != nil {
			return errors.WithStack(err)
		}
		// Prune value instructions with unused results in f.
		if !pruneUnusedInsts(f, usedIdents) {
			// Fix-point reached, no more pruning possible.
			break
		}
	}
	return nil
}

// pruneUnusedInsts prunes value instructions with unused results in f. Call
// instructions are never pruned as they may have side effects. The boolean
// return value indicates whether an instruction was pruned.
func pruneUnusedInsts(f *ir.Func, usedIdents map[string]bool) bool {
	pruned := false
	// Prune value instructions with unused results in f.
	for _, block := range f.Blocks {
		insts := block.Insts[:0] // filter instructions in-place.
		for _, inst := range block.Insts {
			if isUnused(inst, usedIdents) {
				// Prune unused value instruction, skip append to insts.
				pruned = true
				continue
			}
			insts = append(insts, inst)
		}
		block.Insts = insts
	}
	if pruned {
		// Reset IDs of unnamed local variables, as one or more instructions have
		// been pruned.
		irutil.ResetNames(f)
	}
	return pruned
}

// ValueInstruction is an instruction producing a local variable result.
//
// Note: call instructions are value instructions, but only produce a result if
// the return type of their callee is non-void.
type ValueInstruction interface {
	ir.Instruction
	value.Named
}

// isUnused reports whether the given instruction is a value instruction with
// unused results. Call instructions are never considered unused, as they may
// have side effects.
func isUnused(inst ir.Instruction, usedIdents map[string]bool) bool {
	valueInst, ok := inst.(ValueInstruction)
	if !ok {
		// Non-value instructions are never considered unused.
		return false
	}
	if _, ok := inst.(*ir.InstCall); ok {
		// Call instructions are never considered unused as they may have side
		// effects.
		return false
	}
	return !usedIdents[valueInst.Name()]
}

// getUsedIdents returns the names of local variables used in f.
func getUsedIdents(f *ir.Func) (map[string]bool, error) {
	// Assign IDs to unnamed local variables (as weneed to record the names of
	// local variables used in f).
	if err := f.AssignIDs(); err != nil {
		return nil, errors.WithStack(err)
	}
	// Pre-allocate backing array on stack for uses slice.
	var backing [100]irutil.Use
	// Record names of local variables used used in f.
	uses := irutil.FuncUses(backing[:0], f) // the backing array is optional, pass nil if not needed.
	usedIdents := make(map[string]bool)
	for _, use := range uses {
		valueInst, ok := (*use.Val).(ValueInstruction)
		if !ok {
			// Only record the names of local variables produced by value
			// instructions.
			continue
		}
		usedIdents[valueInst.Name()] = true
	}
	return usedIdents, nil
}
