package irutil

import (
	"math/big"

	"github.com/llir/llvm/ir/constant"
)

// Simplify returns an equivalent (and potentially simplified) constant to
// the constant expression.
func Simplify(c constant.Constant) constant.Constant {
	switch c := c.(type) {
	case *constant.ExprAdd:
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			x.X = (&big.Int{}).Add(x.X, y.X)
		}
		return x
	case *constant.ExprSub:
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			x.X = (&big.Int{}).Sub(x.X, y.X)
		}
		return x
	case *constant.ExprMul:
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			x.X = (&big.Int{}).Mul(x.X, y.X)
		}
		return x
	case *constant.ExprSDiv:
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			x.X = (&big.Int{}).Div(x.X, y.X)
		}
		return x
	case *constant.ExprUDiv:
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			x.X = (&big.Int{}).Div(x.X, y.X)
		}
		return x
	case *constant.ExprFAdd:
		x, ok := Simplify(c.X).(*constant.Float)
		y, ok2 := Simplify(c.Y).(*constant.Float)
		if ok && ok2 {
			x.X = (&big.Float{}).Add(x.X, y.X)
		}
		return x
	case *constant.ExprFSub:
		x, ok := Simplify(c.X).(*constant.Float)
		y, ok2 := Simplify(c.Y).(*constant.Float)
		if ok && ok2 {
			x.X = (&big.Float{}).Sub(x.X, y.X)
		}
		return x
	case *constant.ExprFMul:
		x, ok := Simplify(c.X).(*constant.Float)
		y, ok2 := Simplify(c.Y).(*constant.Float)
		if ok && ok2 {
			x.X = (&big.Float{}).Mul(x.X, y.X)
		}
		return x
	default:
		return c
	}
}
