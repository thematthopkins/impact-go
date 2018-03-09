package calculated

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	q1 := QuestionSfid("Q1")
	q2 := QuestionSfid("Q2")
	expanded := ExpandedExpression{
		questionOperand1: &q1,
		questionOperand2: &q2,
	}

	assessment := map[QuestionSfid]float64{
		"Q1": 5.0,
		"Q2": 2.5,
	}

	expanded.operator = Add
	addResult, addErr := evaluate(expanded, assessment)
	assert.NoError(t, addErr)
	assert.Equal(t, 7.5, addResult)

	expanded.operator = Subtract
	subtractResult, subtractErr := evaluate(expanded, assessment)
	assert.NoError(t, subtractErr)
	assert.Equal(t, 2.5, subtractResult)

	expanded.operator = Multiply
	multiplyResult, multiplyErr := evaluate(expanded, assessment)
	assert.NoError(t, multiplyErr)
	assert.Equal(t, 12.5, multiplyResult)

	expanded.operator = Divide
	divideResult, divideErr := evaluate(expanded, assessment)
	assert.NoError(t, divideErr)
	assert.Equal(t, 2.0, divideResult)

	expanded.operator = "!"
	unknownResult, unknownErr := evaluate(expanded, assessment)
	assert.NoError(t, unknownErr)
	assert.Equal(t, 0.0, unknownResult)
}

func TestEvaluate_DivideByZero(t *testing.T) {
	q1 := QuestionSfid("Q1")
	q2 := QuestionSfid("Q2")
	expanded := ExpandedExpression{
		questionOperand1: &q1,
		questionOperand2: &q2,
	}

	assessment := map[QuestionSfid]float64{
		"Q1": 5.0,
		"Q2": 0,
	}

	expanded.operator = Divide
	divideResult, divideErr := evaluate(expanded, assessment)
	assert.NoError(t, divideErr)
	assert.Equal(t, 0.0, divideResult)
}

func TestEvaluate_Descendants(t *testing.T) {
	q1 := QuestionSfid("Q1")
	q2 := QuestionSfid("Q2")
	q3 := QuestionSfid("Q3")
	expandedFirst := ExpandedExpression{
		operator:         Add,
		questionOperand1: &q1,
		questionOperand2: &q2,
	}

	expandedSecond := ExpandedExpression{
		operator:         Add,
		questionOperand1: &q2,
		questionOperand2: &q3,
	}

	expanded := ExpandedExpression{
		operator:           Add,
		expressionOperand1: &expandedFirst,
		expressionOperand2: &expandedSecond,
	}

	assessment := map[QuestionSfid]float64{
		"Q1": 5.0,
		"Q2": 2.5,
		"Q3": 7.25,
	}

	result, err := evaluate(expanded, assessment)
	assert.NoError(t, err)
	assert.Equal(t, 17.25, result)
}

func TestEvaluate_DescendentsMissing(t *testing.T) {
	q1 := QuestionSfid("Q1")
	q2 := QuestionSfid("Q2")
	expanded := ExpandedExpression{
		operator:         Add,
		questionOperand1: &q1,
		questionOperand2: &q2,
	}

	assessment := map[QuestionSfid]float64{
		"Q2": 5.3,
	}

	result, err := evaluate(expanded, assessment)
	assert.NoError(t, err)
	assert.Equal(t, 5.3, result)
}

func TestEvaluate_DescendentsMissingOperand(t *testing.T) {
	q1 := QuestionSfid("Q1")
	expanded := ExpandedExpression{
		operator:         Add,
		questionOperand1: &q1,
		questionOperand2: nil,
	}

	assessment := map[QuestionSfid]float64{
		"Q1": 2.5,
	}

	_, err := evaluate(expanded, assessment)
	assert.Error(t, err)

	expanded.questionOperand1 = nil
	expanded.questionOperand2 = &q1
	_, err = evaluate(expanded, assessment)
	assert.Error(t, err)
}

func TestEvaluate_DescendantsBubbleUp(t *testing.T) {
	q1 := QuestionSfid("Q1")
	q2 := QuestionSfid("Q2")
	expandedFirst := ExpandedExpression{
		operator:         Add,
		questionOperand1: &q1,
		questionOperand2: &q2,
	}

	expandedSecond := ExpandedExpression{
		operator:         Add,
		questionOperand1: &q2,
		questionOperand2: nil,
	}

	expanded := ExpandedExpression{
		operator:           Add,
		expressionOperand1: &expandedFirst,
		expressionOperand2: &expandedSecond,
	}

	assessment := map[QuestionSfid]float64{
		"Q1": 5.0,
		"Q2": 2.5,
	}

	_, err := evaluate(expanded, assessment)
	assert.Error(t, err)

	expandedFirst.questionOperand1 = nil
	_, err = evaluate(expanded, assessment)
	assert.Error(t, err)
}

func TestExpand_Empty(t *testing.T) {
	result, err := expand(map[QuestionSfid]Expression{})

	assert.NoError(t, err)
	assert.Equal(t, map[QuestionSfid]ExpandedExpression{}, result)
}

func TestExpand(t *testing.T) {
	q2 := QuestionSfid("Q2")
	q3 := QuestionSfid("Q3")
	q5 := QuestionSfid("Q5")
	q6 := QuestionSfid("Q6")

	result, err := expand(map[QuestionSfid]Expression{
		"Q1": Expression{
			operator: Add,
			operand1: q2,
			operand2: q3,
		},
		"Q4": Expression{
			operator: Add,
			operand1: q5,
			operand2: q6,
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, map[QuestionSfid]ExpandedExpression{
		"Q1": ExpandedExpression{
			operator:         Add,
			questionOperand1: &q2,
			questionOperand2: &q3,
		},
		"Q4": ExpandedExpression{
			operator:         Add,
			questionOperand1: &q5,
			questionOperand2: &q6,
		},
	}, result)
}

func TestExpand_Deep(t *testing.T) {
	q1 := QuestionSfid("Q1")
	q2 := QuestionSfid("Q2")
	q3 := QuestionSfid("Q3")
	q4 := QuestionSfid("Q4")

	result, err := expand(map[QuestionSfid]Expression{
		"Q3": Expression{
			operator: Add,
			operand1: q1,
			operand2: q2,
		},
		"Q4": Expression{
			operator: Add,
			operand1: q2,
			operand2: q1,
		},
		"Q6": Expression{
			operator: Add,
			operand1: q3,
			operand2: q4,
		},
	})

	qThreeExpanded := ExpandedExpression{
		operator:         Add,
		questionOperand1: &q1,
		questionOperand2: &q2,
	}

	qFourExpanded := ExpandedExpression{
		operator:         Add,
		questionOperand1: &q2,
		questionOperand2: &q1,
	}

	assert.NoError(t, err)
	assert.Equal(t, map[QuestionSfid]ExpandedExpression{
		"Q3": qThreeExpanded,
		"Q4": qFourExpanded,
		"Q6": ExpandedExpression{
			operator:           Add,
			expressionOperand1: &qThreeExpanded,
			expressionOperand2: &qFourExpanded,
		},
	}, result)
}

func TestExpanded_Infinite(t *testing.T) {
	q1 := QuestionSfid("Q1")
	q2 := QuestionSfid("Q2")
	_, err := expand(map[QuestionSfid]Expression{
		"Q1": Expression{
			operator: Add,
			operand1: q1,
			operand2: q2,
		},
	})

	assert.Error(t, err)

	_, err = expand(map[QuestionSfid]Expression{
		"Q1": Expression{
			operator: Add,
			operand1: q2,
			operand2: q1,
		},
	})

	assert.Error(t, err)
}
