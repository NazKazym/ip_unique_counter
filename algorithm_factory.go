package main

func newIPSetFromConfig(cc CounterConfig, approxItems uint64) IPSet {
	threshold := cc.BitmapThreshold
	if threshold == 0 {
		threshold = defaultCounterBitmapThreshold
	}

	switch cc.Storage {
	case "map":
		// hint на capacity — просто подсказка, можно передать approxItems
		return newMapIPSet(int(approxItems))

	case "bitmap":
		return newBitmapIPSet()

	case "auto":
		fallthrough
	default:
		// если не смогли оценить — безопасно начинаем с map
		if approxItems == 0 {
			return newMapIPSet(0)
		}
		if approxItems > threshold {
			return newBitmapIPSet()
		}
		return newMapIPSet(int(approxItems))
	}
}
