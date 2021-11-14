package basic_functions

import (
	"CoolProp"
	"math"
)

func DefaultEnv() map[string]interface{} {
	output := map[string]interface{}{
		// Constants
		"e":   math.E,
		"pi":  math.Pi,
		"phi": math.Phi,
		// Functions
		"abs":          math.Abs,
		"acos":         math.Acos,
		"acosh":        math.Acosh,
		"asin":         math.Asin,
		"asinh":        math.Asinh,
		"atan":         math.Atan,
		"atan2":        math.Atan2,
		"atanh":        math.Atanh,
		"cbrt":         math.Cbrt,
		"cos":          math.Cos,
		"cosh":         math.Cosh,
		"ceil":         math.Ceil,
		"dim":          math.Dim,
		"copysign":     math.Copysign,
		"erf":          math.Erf,
		"erfc":         math.Erfc,
		"erfcinv":      math.Erfcinv,
		"erfinv":       math.Erfcinv,
		"exp":          math.Exp,
		"exp2":         math.Exp2,
		"expm1":        math.Expm1,
		"fma":          math.FMA,
		"floor":        math.Floor,
		"frexp":        math.Frexp,
		"gamma":        math.Gamma,
		"hypot":        math.Hypot,
		"ilogb":        math.Ilogb,
		"zerobessel":   math.J0,
		"onebessel":    math.J1,
		"nbessel":      math.Jn,
		"ldexp":        math.Ldexp,
		"lgamma":       math.Lgamma,
		"log":          math.Log,
		"log10":        math.Log10,
		"log2":         math.Log2,
		"logb":         math.Logb,
		"max":          math.Max,
		"min":          math.Min,
		"mod":          math.Mod,
		"pow":          math.Pow,
		"pow10":        math.Pow10,
		"reminder":     math.Remainder,
		"round":        math.Round,
		"roundtoeven":  math.RoundToEven,
		"sin":          math.Sin,
		"sincos":       math.Sincos,
		"sinh":         math.Sinh,
		"sqrt":         math.Sqrt,
		"tan":          math.Tan,
		"tanh":         math.Tanh,
		"trunc":        math.Trunc,
		"zerobessel2":  math.Y0,
		"onebessel2":   math.Y1,
		"nbessel2":     math.Yn,
		"prop1SI":      CoolProp.Props1SI,
		"propsSI":      CoolProp.PropsSI,
		"propsSImulti": CoolProp.PropsSImulti,
	}
	return output
}
