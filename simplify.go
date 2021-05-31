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
			z := constant.NewInt(x.Typ, 0)
			z.X = (&big.Int{}).Add(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprSub:
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			z := constant.NewInt(x.Typ, 0)
			z.X = (&big.Int{}).Sub(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprMul:
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			z := constant.NewInt(x.Typ, 0)
			z.X = (&big.Int{}).Mul(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprSDiv:
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			z := constant.NewInt(x.Typ, 0)
			z.X = (&big.Int{}).Div(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprUDiv:
		x, ok := Simplify(c.X).(*constant.Int)
		y, ok2 := Simplify(c.Y).(*constant.Int)
		if ok && ok2 {
			z := constant.NewInt(x.Typ, 0)
			z.X = (&big.Int{}).Div(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprFAdd:
		x, ok := Simplify(c.X).(*constant.Float)
		y, ok2 := Simplify(c.Y).(*constant.Float)
		if ok && ok2 {
			z := constant.NewFloat(x.Typ, 0)
			z.X = (&big.Float{}).Add(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprFSub:
		x, ok := Simplify(c.X).(*constant.Float)
		y, ok2 := Simplify(c.Y).(*constant.Float)
		if ok && ok2 {
			z := constant.NewFloat(x.Typ, 0)
			z.X = (&big.Float{}).Sub(x.X, y.X)
			return z
		}
		return c
	case *constant.ExprFMul:
		x, ok := Simplify(c.X).(*constant.Float)
		y, ok2 := Simplify(c.Y).(*constant.Float)
		if ok && ok2 {
			z := constant.NewFloat(x.Typ, 0)
			z.X = (&big.Float{}).Mul(x.X, y.X)
			return z
		}
		return c
	default:
		return c
	}
}
