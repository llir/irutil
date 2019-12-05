package irutil

import (
	"strings"

	"github.com/llir/llvm/ir"
)

// Comment is an LLVM IR comment represented as a pseudo-instruction. Comment
// implements ir.Instruction.
type Comment struct {
	// Comment text; may contain multiple lines.
	Text string

	// embed ir.Instruction to satisfy the ir.Instruction interface.
	ir.Instruction
}

// LLString returns the LLVM syntax representation of the value.
func (c *Comment) LLString() string {
	// handle multi-line comments.
	text := strings.ReplaceAll(c.Text, "\n", "; ")
	return "; " + text
}

// NewComment returns a new LLVM IR comment represented as a pseudo-instruction.
// Text may contain multiple lines.
func NewComment(text string) *Comment {
	return &Comment{
		Text: text,
	}
}
