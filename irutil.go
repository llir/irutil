// Package irutil implements LLVM IR utility functions.
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
	for _, b := range f.Blocks {
		for _, i := range b.Insts {
			if i, ok := i.(Ident); ok {
				if i.IsUnnamed() {
					// clear ID of unnamed variable.
					i.SetName("")
				}
			}
		}
	}
}
