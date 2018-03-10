package scoring

// Response is user input
type Response struct {
	IsAnswered          bool
	TotalAnswer         float64
	Points              float64
	HiddenByContingency bool
	NumberResponse      float64
	CurrencyResponse    float64
	// takes place of percent response and "value percentage"
	PercentResponse float64
}

// Standard is metadata for a Response
type Standard struct {
	ScoringMethod ScoreType
	AnswerType    string
	LowThreshold  float64
	HighThreshold float64
	Worth         float64
}

// ScoreType scores
type ScoreType string

const (
	StraightPercentage ScoreType = "Straight Percentage"
	InversePercentage            = "Inverse Percentage"
	LowHighThreshold             = "Low/High Treshold"
	SumOfAnswerValues            = "Sum of Answer Values"
)

func clamp(
	input float64,
	low float64,
	high float64,
) float64 {
	if input < low {
		return low
	} else if input > high {
		return high
	} else {
		return input
	}
}

func inverseLerp(
	input float64,
	low float64,
	high float64,
) float64 {
	return (input - low) / (high - low)
}

func unclampedPerformance(
	response Response,
	standard Standard,
) float64 {

	switch standard.ScoringMethod {
	case StraightPercentage:
		return inverseLerp(response.PercentResponse, 0, 100)
	case InversePercentage:
		return 1 - inverseLerp(response.PercentResponse, standard.LowThreshold, standard.HighThreshold)
	case LowHighThreshold:
		return inverseLerp(response.PercentResponse, standard.LowThreshold, standard.HighThreshold)
	case SumOfAnswerValues:
		return inverseLerp(response.PercentResponse, 0, 100)
	default:
		return 0
	}
}

func score(
	unclampedPerformance float64,
	worth float64,
) float64 {
	return clamp(unclampedPerformance, 0, 1) * worth
}
