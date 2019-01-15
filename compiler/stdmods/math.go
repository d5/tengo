package stdmods

import (
	"math"

	"github.com/d5/tengo/objects"
)

// Math is a math module.
var Math = Module{
	Name: "math",
	Globals: map[string]objects.Object{
		"e":        &objects.Float{Value: math.E},
		"pi":       &objects.Float{Value: math.Pi},
		"phi":      &objects.Float{Value: math.Phi},
		"sqrt2":    &objects.Float{Value: math.Sqrt2},
		"sqrtE":    &objects.Float{Value: math.SqrtE},
		"sqrtPi":   &objects.Float{Value: math.SqrtPi},
		"sqrtPhi":  &objects.Float{Value: math.SqrtPhi},
		"ln2":      &objects.Float{Value: math.Ln2},
		"log2E":    &objects.Float{Value: math.Log2E},
		"ln10":     &objects.Float{Value: math.Ln10},
		"log10E":   &objects.Float{Value: math.Log10E},
		"abs":      funcAFRF(math.Abs),
		"acos":     funcAFRF(math.Acos),
		"acosh":    funcAFRF(math.Acosh),
		"asin":     funcAFRF(math.Asin),
		"asinh":    funcAFRF(math.Asinh),
		"atan":     funcAFRF(math.Atan),
		"atan2":    funcAFFRF(math.Atan2),
		"atanh":    funcAFRF(math.Atanh),
		"cbrt":     funcAFRF(math.Cbrt),
		"ceil":     funcAFRF(math.Ceil),
		"copysign": funcAFFRF(math.Copysign),
		"cos":      funcAFRF(math.Cos),
		"cosh":     funcAFRF(math.Cosh),
		"dim":      funcAFFRF(math.Dim),
		"erf":      funcAFRF(math.Erf),
		"erfc":     funcAFRF(math.Erfc),
		"erfcinv":  funcAFRF(math.Erfcinv),
		"erfinv":   funcAFRF(math.Erfinv),
		"exp":      funcAFRF(math.Exp),
		"exp2":     funcAFRF(math.Exp2),
		"expm1":    funcAFRF(math.Expm1),
		"floor":    funcAFRF(math.Floor),
		// TODO: continue adding more functions
	},
}
