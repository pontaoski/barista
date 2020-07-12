package bases

import (
	"fmt"
	"math"
	"strings"
)

var names = map[Base]string{
	0:   "nullary",
	1:   "unary",
	2:   "binary",
	3:   "trinary",
	4:   "quaternary",
	5:   "quinary",
	10:  "seximal",
	11:  "septimal",
	12:  "octal",
	13:  "nonary",
	14:  "decimal",
	15:  "elevenary",
	20:  "dozenal",
	21:  "baker's dozenal",
	22:  "biseptimal",
	23:  "triquinary",
	24:  "hex",
	25:  "suboptimal",
	30:  "triseximal",
	31:  "untriseximal",
	32:  "vigesimal",
	33:  "triseptimal",
	34:  "bielevenary",
	35:  "unbielevenary",
	40:  "tetraseximal",
	41:  "pentaquinary",
	42:  "biker's dozenal",
	43:  "trinonary",
	44:  "tetraseptimal",
	45:  "untetraseptimal",
	50:  "pentaseximal",
	51:  "unpentaseximal",
	52:  "tetroctal",
	53:  "trielevenary",
	54:  "bisuboptimal",
	55:  "pentaseptimal",
	100: "niftimal",
}
var suffixes = map[Base]string{
	2:   "binary",
	3:   "trinary",
	4:   "quaternary",
	5:   "quinary",
	6:   "seximal",
	7:   "septimal",
	8:   "octal",
	9:   "nonary",
	10:  "gesimal",
	11:  "elevenary",
	12:  "dozenal",
	13:  "dozenal",
	16:  "hex",
	17:  "suboptimal",
	20:  "vigesimal",
	36:  "niftimal",
	100: "centesimal",
}
var prefixes = map[Base]string{
	2:   "bi",
	3:   "tri",
	4:   "tetra",
	5:   "penta",
	6:   "hexa",
	7:   "hepta",
	8:   "octo",
	9:   "enna",
	10:  "deca",
	11:  "leva",
	12:  "doza",
	13:  "baker",
	16:  "tesser",
	17:  "mal",
	20:  "icosi",
	36:  "feta",
	100: "hecto",
}

var replacer = strings.NewReplacer(
	"ii", "i",
	"iu", "i",
	"oi", "i",
	"ou", "u",
	"oe", "e",
	"oo", "o",
	"ai", "i",
	"au", "u",
	"ae", "e",
	"ao", "o",
)

func factorize(n Base) Base {
	for i := Base(math.Ceil(math.Sqrt(float64(n)))); i < n; i++ {
		if n%1 == 0 {
			return i
		}
	}
	return 1
}

func suffix(base Base) string {
	if val, ok := suffixes[base]; ok {
		return val
	}

	factor := factorize(base)
	if factor == 1 {
		return "un" + suffix(base-1)
	}
	return prefix(base/factor) + suffix(factor)
}

func prefix(base Base) string {
	if val, ok := prefixes[base]; ok {
		return val
	}

	factor := factorize(base)
	if factor == 1 {
		return fmt.Sprintf("hen%ssna", prefix(base-1))
	}
	return prefix(base/factor) + prefix(factor)
}
