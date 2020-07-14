package bases

import (
	"strconv"
	"strings"

	lru "github.com/hashicorp/golang-lru"
)

type Base int64

const (
	max = 36
)

func (e Base) ConvertDigit(digit int64) string {
	base := int(e)
	if base <= 10 {
		return strconv.FormatInt(digit, base)
	}

	var simple string
	if e > max {
		return strings.ToUpper(strconv.FormatInt(digit, 10))
	} else {
		simple = strings.ToUpper(strconv.FormatInt(digit, base))
	}

	switch {
	case base <= 12:
		return strings.ReplaceAll(strings.ReplaceAll(simple, "A", "X"), "B", "E")
	case base == 13:
		return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(simple, "A", "X"), "B", "Y"), "C", "Z")
	case base <= 36:
		return strconv.FormatInt(digit, base)
	}

	return simple
}

func (b Base) DigitsFor(denominator int) int64 {
	ret := int64(0)
	base := int64(b)

	num := int64(1)
	den := int64(denominator)

	for den/num > base {
		num *= base
		num %= den
	}

	remainder := num % den * base

	data := map[int64]int64{}

	for remainder != 0 {
		if _, ok := data[remainder]; ok {
			return ret
		}
		data[remainder] = ret
		ret++
		remainder = remainder % den * base
	}

	return ret
}

func (b Base) DecimalFor(denominator int) (string, string) {
	num := int64(1)
	den := int64(denominator)

	intResult := num / den
	remainder := num % den * int64(b)

	data := map[int64]int64{}

	var result strings.Builder
	result.WriteString("0.")

	for remainder != 0 {
		if val, ok := data[remainder]; ok {
			pre := result.String()[0:val]
			post := result.String()[val:]

			return pre, post
		}

		data[remainder] = int64(len(result.String()))
		intResult = remainder / den

		if int64(b) > 36 && result.String() != "0." {
			result.WriteString(":")
		}

		result.WriteString(b.ConvertDigit(intResult))
		remainder = remainder % den * int64(b)
	}
	return result.String(), ""
}

var baseCache, _ = lru.New(1024)
var tidier = strings.NewReplacer(
	"aa", "a",
	"ae", "e",
	"ai", "i",
	"ao", "o",
	"au", "u",
	"oa", "a",
	"oe", "e",
	"oi", "i",
	"oo", "o",
	"ou", "u",
	"ii", "i",
	"iu", "u",
)

func (b Base) Name() string {
	if val, ok := names[b]; ok {
		return val
	}
	if val, ok := baseCache.Get(b); ok {
		return val.(string)
	}
	ret := tidier.Replace(strings.ReplaceAll(NameBase(b), "-", ""))
	baseCache.Add(b, ret)
	return ret
}
