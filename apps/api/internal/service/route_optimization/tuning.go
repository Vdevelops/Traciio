package route_optimization

import (
	"os"
	"strconv"
	"sync"
	"time"
)

type improvementTuning struct {
	budget          time.Duration
	budgetThreshold int
	twoOptMaxIter   int
	threeOptMaxIter int
}

var (
	loadImprovementTuningOnce sync.Once
	cachedImprovementTuning   improvementTuning
)

func getEnvInt(name string, defaultValue int) int {
	v := os.Getenv(name)
	if v == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func loadTuning() {
	// Defaults are conservative and only affect large routes.
	// Set ROUTE_OPTIMIZATION_IMPROVEMENT_BUDGET_MS=0 to disable.
	budgetMs := getEnvInt("ROUTE_OPTIMIZATION_IMPROVEMENT_BUDGET_MS", 75)
	threshold := getEnvInt("ROUTE_OPTIMIZATION_IMPROVEMENT_BUDGET_THRESHOLD", 20)
	if threshold < 0 {
		threshold = 0
	}

	max2 := getEnvInt("ROUTE_OPTIMIZATION_TWO_OPT_MAX_ITERATIONS", 100)
	if max2 < 1 {
		max2 = 1
	}
	max3 := getEnvInt("ROUTE_OPTIMIZATION_THREE_OPT_MAX_ITERATIONS", 50)
	if max3 < 1 {
		max3 = 1
	}

	if budgetMs < 0 {
		budgetMs = 0
	}

	cachedImprovementTuning = improvementTuning{
		budget:          time.Duration(budgetMs) * time.Millisecond,
		budgetThreshold: threshold,
		twoOptMaxIter:   max2,
		threeOptMaxIter: max3,
	}
}

func getImprovementTuning() improvementTuning {
	loadImprovementTuningOnce.Do(loadTuning)
	return cachedImprovementTuning
}

func improvementDeadline(waypointCount int) (time.Time, bool) {
	tuning := getImprovementTuning()
	if tuning.budget <= 0 {
		return time.Time{}, false
	}
	if waypointCount < tuning.budgetThreshold {
		return time.Time{}, false
	}
	return time.Now().Add(tuning.budget), true
}
