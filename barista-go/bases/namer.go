package bases

import (
	"math/big"
	"strings"

	lru "github.com/hashicorp/golang-lru"
)

var names = map[Base]string{
	0:  "nullary",
	1:  "unary",
	2:  "binary",
	3:  "trinary",
	4:  "quaternary",
	5:  "quinary",
	6:  "seximal",
	7:  "septimal",
	8:  "octal",
	9:  "nonary",
	10: "decimal",
	11: "elevenary",
	12: "dozenal",
	13: "baker's dozenal",
	14: "biseptimal",
	15: "triquinary",
	16: "hex",
	17: "suboptimal",
	18: "triseximal",
	19: "untriseximal",
	20: "vigesimal",
	21: "triseptimal",
	22: "bielevenary",
	23: "unbielevenary",
	24: "tetraseximal",
	25: "pentaquinary",
	26: "biker's dozenal",
	27: "trinonary",
	28: "tetraseptimal",
	29: "untetraseptimal",
	30: "pentaseximal",
	31: "unpentaseximal",
	32: "tetroctal",
	33: "trielevenary",
	34: "bisuboptimal",
	35: "pentaseptimal",
	36: "niftimal",
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
var prefixList = []Base{100, 36, 20, 17, 16, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2}

var suffixes = map[Base]string{
	2:   "binary",
	3:   "trinary",
	4:   "tetra",
	5:   "quinary",
	6:   "seximal",
	7:   "septimal",
	8:   "octal",
	9:   "nonary",
	10:  "gesimal",
	11:  "elevenary",
	12:  "dozenal",
	13:  "ker's dozenal",
	16:  "hex",
	17:  "suboptimal",
	20:  "vigesimal",
	36:  "niftimal",
	100: "centesimal",
}
var suffixList = []Base{100, 36, 20, 17, 16, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2}

type factor struct {
	lower  Base
	larger Base
}

var factorCache, _ = lru.New(1 << 16)

func factors(base Base) (lower Base, larger Base) {
	if val, ok := factorCache.Get(base); ok {
		a := val.(factor)
		return a.lower, a.larger
	}
	lower = -1
	larger = -1
	for _, prime := range primes {
		if prime > base {
			break
		}
		if prime == base {
			lower = 1
			larger = base
			return
		}
		if base%prime == 0 {
			lower = prime
			larger = base / prime
			return
		}
	}
	if big.NewInt(int64(base)).ProbablyPrime(0) {
		lower = 1
		larger = base
		return
	}
	for i := Base(1); i < base; i++ {
		if base%i == 0 {
			if lower > base/i {
				factorCache.Add(base, factor{lower, larger})
				return
			}
			lower = i
			larger = base / i
		}
	}
	factorCache.Add(base, factor{lower, larger})
	return
}

func multPrefixForm(base Base) []string {
	if val, ok := prefixes[base]; ok {
		return []string{val}
	}
	retArr := []string{}
	if base > 17 && big.NewInt(int64(base)).ProbablyPrime(0) {
		return append(append([]string{"hen"}, multPrefixForm(base-1)...), "sna")
	}
	for base > 1 {
		pre := base
		for _, item := range prefixList {
			if base%item == 0 {
				retArr = append(retArr, prefixes[item])
				base = base / item
			}
		}
		if pre == base {
			smaller, greater := factors(base)
			return append(multPrefixForm(greater), multPrefixForm(smaller)...)
		}
	}
	return retArr
}

func suffixForm(base Base) []string {
	if val, ok := suffixes[base]; ok {
		return []string{val}
	}
	if base > 17 && big.NewInt(int64(base)).ProbablyPrime(0) {
		return append([]string{"un"}, suffixForm(base-1)...)
	}
	smaller, greater := factors(base)
	return append(multPrefixForm(greater), suffixForm(smaller)...)
}

func NameBase(base Base) (ret string) {
	if base < 0 {
		return "nega" + NameBase(-1*base)
	}
	return strings.Join(suffixForm(base), "-")
}
