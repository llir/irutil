package irutil

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

// Ident is the global or local identifier of a named value.
type Ident interface {
	value.Named
	// ID returns the ID of the identifier.
	ID() int64
	// SetID sets the ID of the identifier.
	SetID(id int64)
	// IsUnnamed reports whether the identifier is unnamed.
	IsUnnamed() bool
}

// ResetNames resets the IDs of unnamed local variables in the given function.
func ResetNames(f *ir.Func) {
	for _, block := range f.Blocks {
		// clear ID of unnamed basic block.
		if block.IsUnnamed() {
			block.SetName("")
		}
		for _, inst := range block.Insts {
			if inst, ok := inst.(Ident); ok {
				if inst.IsUnnamed() {
					// clear ID of unnamed variable.
					inst.SetName("")
				}
			}
		}
	}
}
