package camera

func SmoothApproach(current, target, factor float64) float64 {
	return current + (target-current)*factor
}

