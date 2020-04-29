package interactor

import "math"

var (
	serverLimit   = 0xf4240
	nodePerMin    = 10
	nodePerLimit  = 0xf4240
	contractLimit = 0xf4240
)

func BasicSanitizer(raw map[string]interface{}) map[string]interface{} {
	clean := make(map[string]interface{})
	for key, val := range raw {
		if key == ServerCount {
			sCount := val.(int)
			clean[key] = basicMinMaxIntSanitize(0, serverLimit, sCount)
		} else if key == NodePerCount {
			nCount := val.(int)
			clean[key] = basicMinMaxIntSanitize(nodePerMin, nodePerLimit, nCount)
		} else if key == ContractCount {
			cCount := val.(int)
			clean[key] = basicMinMaxIntSanitize(0, contractLimit, cCount)
		} else if key == TransLimit {
			tLimit := val.(int)
			clean[key] = basicMinMaxIntSanitize(0, math.MaxInt32, tLimit)
		}
	}

	return clean
}

func basicMinMaxIntSanitize(min int, max int, val int) int {
	cleanCount := val
	if cleanCount < min {
		cleanCount = min
	}
	if cleanCount > max {
		cleanCount = max
	}
	return cleanCount
}
