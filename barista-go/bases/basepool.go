package bases

type BasePool map[Base]struct{}

func PoolForRange(from, to int64) BasePool {
	ret := BasePool{}
	for i := from; i <= to; i++ {
		ret[Base(i)] = struct{}{}
	}
	return ret
}

func (b BasePool) Smallest() Base {
	min := Base(99999)
	for val := range b {
		if val < min {
			min = val
		}
	}
	return min
}

func (b BasePool) LargestExpansionFor(denominator int) ([]Base, int64, bool) {
	length := int64(0)
	recurring := []Base{}
	currentLongest := []Base{}
	hasRecurring := false

	for base := range b {
		_, recurringNums := base.DecimalFor(denominator)
		if recurringNums != "" {
			recurring = append(recurring, base)
			hasRecurring = true
		}
		if hasRecurring {
			continue
		}
		baseLen := base.DigitsFor(denominator)
		if length < baseLen {
			length = baseLen
			currentLongest = []Base{}
		} else if length > baseLen {
			continue
		}
		currentLongest = append(currentLongest, base)
	}
	if hasRecurring {
		return recurring, -1, true
	}
	return currentLongest, length, false
}

func (b BasePool) RemoveLongestExpansionFor(denominator int) {
	toRemove, _, _ := b.LargestExpansionFor(denominator)
	for _, base := range toRemove {
		delete(b, base)
	}
}
