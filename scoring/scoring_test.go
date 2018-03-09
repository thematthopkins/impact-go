package scoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnclampedPerformance_StraightPercentage(t *testing.T) {
	standard := Standard{
		ScoringMethod: StraightPercentage,
		LowThreshold:  0,
		HighThreshold: 0,
	}

	response := Response{
		PercentResponse: 1,
	}

	result := unclampedPerformance(response, standard)
	assert.Equal(t, 0.01, result)
}

func TestUnclampedPerformance_InversePercentage(t *testing.T) {
	standard := Standard{
		ScoringMethod: InversePercentage,
		LowThreshold:  0,
		HighThreshold: 100,
	}

	response := Response{
		PercentResponse: 25,
	}

	result := unclampedPerformance(response, standard)
	assert.Equal(t, 0.75, result)
}

func TestUnclampedPerformance_LowHighThreshold(t *testing.T) {
	standard := Standard{
		ScoringMethod: LowHighThreshold,
		LowThreshold:  50,
		HighThreshold: 100,
	}

	response := Response{
		PercentResponse: 51,
	}

	result := unclampedPerformance(response, standard)
	assert.Equal(t, 0.02, result)
}

func TestUnclampedPerformance_SumOfAnswerValues(t *testing.T) {
	standard := Standard{
		ScoringMethod: SumOfAnswerValues,
	}

	response := Response{
		PercentResponse: 50,
	}

	result := unclampedPerformance(response, standard)
	assert.Equal(t, 0.5, result)
}

func TestUnclampedPerformance_UnknownScoringMethod(t *testing.T) {
	standard := Standard{
		ScoringMethod: "Unknown Scoring Method",
	}

	response := Response{
		PercentResponse: 50,
	}

	result := unclampedPerformance(response, standard)
	assert.Equal(t, 0.0, result)
}

func TestScore(t *testing.T) {
	assert.Equal(t, 25.0, score(500, 25))
	assert.Equal(t, 0.0, score(-100, 25))
	assert.Equal(t, 10.0, score(0.1, 100))
}
