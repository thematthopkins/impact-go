package calculated

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	expanded := OpExpr{
		left:  QID("Q1"),
		right: QID("Q2"),
	}

	assessment := map[QID]float64{
		"Q1": 5.0,
		"Q2": 2.5,
	}

	expanded.op = Add
	addResult, addErr := eval(expanded, assessment)
	assert.NoError(t, addErr)
	assert.Equal(t, 7.5, addResult)

	expanded.op = Subtract
	subtractResult, subtractErr := eval(expanded, assessment)
	assert.NoError(t, subtractErr)
	assert.Equal(t, 2.5, subtractResult)

	expanded.op = Multiply
	multiplyResult, multiplyErr := eval(expanded, assessment)
	assert.NoError(t, multiplyErr)
	assert.Equal(t, 12.5, multiplyResult)

	expanded.op = Divide
	divideResult, divideErr := eval(expanded, assessment)
	assert.NoError(t, divideErr)
	assert.Equal(t, 2.0, divideResult)

	expanded.op = "!"
	unknownResult, unknownErr := eval(expanded, assessment)
	assert.NoError(t, unknownErr)
	assert.Equal(t, 0.0, unknownResult)
}

func TestEval_DivideByZero(t *testing.T) {
	divideResult, divideErr := eval(OpExpr{
		op:    Divide,
		left:  QID("Q1"),
		right: QID("Q2"),
	}, map[QID]float64{
		"Q1": 5.0,
		"Q2": 0,
	})
	assert.NoError(t, divideErr)
	assert.Equal(t, 0.0, divideResult)
}

func TestEval_Descendants(t *testing.T) {
	result, err := eval(OpExpr{
		op: Add,
		left: OpExpr{
			op:    Add,
			left:  QID("Q1"),
			right: QID("Q2"),
		},
		right: OpExpr{
			op:    Add,
			left:  QID("Q2"),
			right: QID("Q3"),
		},
	}, map[QID]float64{
		"Q1": 5.0,
		"Q2": 2.5,
		"Q3": 7.25,
	})
	assert.NoError(t, err)
	assert.Equal(t, 17.25, result)
}

func TestEval_DescendentsMissing(t *testing.T) {
	result, err := eval(OpExpr{
		op:    Add,
		left:  QID("Q1"),
		right: QID("Q2"),
	}, map[QID]float64{
		"Q2": 5.3,
	})
	assert.NoError(t, err)
	assert.Equal(t, 5.3, result)
}

func TestEval_InvalidLeft(t *testing.T) {
	_, err := eval(OpExpr{
		op:    Add,
		left:  QID("Q1"),
		right: nil,
	}, map[QID]float64{
		"Q1": 2.5,
	})
	assert.Error(t, err)
}

func TestEval_InvalidRight(t *testing.T) {
	_, err := eval(OpExpr{
		op:    Add,
		left:  nil,
		right: QID("Q1"),
	}, map[QID]float64{
		"Q1": 2.5,
	})
	assert.Error(t, err)
}

func TestExpand_Empty(t *testing.T) {
	result, err := expand(map[QID]OpDef{})

	assert.NoError(t, err)
	assert.Equal(t, map[QID]Expr{}, result)
}

func TestExpand(t *testing.T) {
	result, err := expand(map[QID]OpDef{
		"Q1": OpDef{
			op:    Add,
			left:  "Q2",
			right: "Q3",
		},
		"Q4": OpDef{
			op:    Add,
			left:  "Q5",
			right: "Q6",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, map[QID]Expr{
		"Q1": OpExpr{
			op:    Add,
			left:  QID("Q2"),
			right: QID("Q3"),
		},
		"Q4": OpExpr{
			op:    Add,
			left:  QID("Q5"),
			right: QID("Q6"),
		},
	}, result)
}

func TestExpand_Deep(t *testing.T) {
	result, err := expand(map[QID]OpDef{
		"Q3": OpDef{
			op:    Add,
			left:  "Q1",
			right: "Q2",
		},
		"Q4": OpDef{
			op:    Add,
			left:  "Q2",
			right: "Q1",
		},
		"Q6": OpDef{
			op:    Add,
			left:  "Q3",
			right: "Q4",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, map[QID]Expr{
		"Q3": OpExpr{
			op:    Add,
			left:  QID("Q1"),
			right: QID("Q2"),
		},
		"Q4": OpExpr{
			op:    Add,
			left:  QID("Q2"),
			right: QID("Q1"),
		},
		"Q6": OpExpr{
			op: Add,
			left: OpExpr{
				op:    Add,
				left:  QID("Q1"),
				right: QID("Q2"),
			},
			right: OpExpr{
				op:    Add,
				left:  QID("Q2"),
				right: QID("Q1"),
			},
		},
	}, result)
}

func TestExpanded_Infinite(t *testing.T) {
	_, err := expand(map[QID]OpDef{
		"Q1": OpDef{
			op:    Add,
			left:  "Q1",
			right: "Q2",
		},
	})

	assert.Error(t, err)

	_, err = expand(map[QID]OpDef{
		"Q1": OpDef{
			op:    Add,
			left:  QID("Q2"),
			right: QID("Q1"),
		},
	})

	assert.Error(t, err)
}
