package calculated

import "errors"

type QuestionSfid string

// Formula is a composition of two questions and how to compute them
type Expression struct {
	operator Operator
	operand1 QuestionSfid
	operand2 QuestionSfid
}

type ExpandedExpression struct {
	operator           Operator
	expressionOperand1 *ExpandedExpression
	questionOperand1   *QuestionSfid
	expressionOperand2 *ExpandedExpression
	questionOperand2   *QuestionSfid
}

type Operator string

const (
	Add      Operator = "+"
	Subtract          = "-"
	Multiply          = "*"
	Divide            = "/"
)

var MissingOperand = errors.New("missing operand")

func evaluate(
	expandedExpression ExpandedExpression,
	questions map[QuestionSfid]float64,
) (float64, error) {

	if expandedExpression.expressionOperand1 == nil && expandedExpression.questionOperand1 == nil {
		return 0, MissingOperand
	}

	if expandedExpression.expressionOperand2 == nil && expandedExpression.questionOperand2 == nil {
		return 0, MissingOperand
	}

	operandOneValue := 0.0
	if expandedExpression.expressionOperand1 != nil {
		var err error
		operandOneValue, err = evaluate(*expandedExpression.expressionOperand1, questions)
		if err != nil {
			return 0, err
		}
	} else {
		operandOneValue, _ = questions[*expandedExpression.questionOperand1]
	}

	operandTwoValue := 0.0
	if expandedExpression.expressionOperand2 != nil {
		var err error
		operandTwoValue, err = evaluate(*expandedExpression.expressionOperand2, questions)
		if err != nil {
			return 0, err
		}
	} else {
		operandTwoValue, _ = questions[*expandedExpression.questionOperand2]
	}

	switch expandedExpression.operator {
	case Add:
		return operandOneValue + operandTwoValue, nil
	case Subtract:
		return operandOneValue - operandTwoValue, nil
	case Multiply:
		return operandOneValue * operandTwoValue, nil
	case Divide:
		if operandTwoValue == 0.0 {
			return 0, nil
		}

		return operandOneValue / operandTwoValue, nil
	default:
		return 0, nil
	}
}

func expand(
	list map[QuestionSfid]Expression,
) map[QuestionSfid]ExpandedExpression {

	result := map[QuestionSfid]ExpandedExpression{}
	for question, expression := range list {
		o1 := expression.operand1
		o2 := expression.operand2
		result[question] = ExpandedExpression{
			operator:         expression.operator,
			questionOperand1: &o1,
			questionOperand2: &o2,
		}
	}

	return result
}
