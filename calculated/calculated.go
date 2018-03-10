package calculated

import "github.com/pkg/errors"

type QID string

type Expr interface {
	isExpr()
}

func (QID) isExpr()    {}
func (OpExpr) isExpr() {}

// OpDef is a composition of two questions and how to compute them
type OpDef struct {
	op    Op
	left  QID
	right QID
}

type OpExpr struct {
	op    Op
	left  Expr
	right Expr
}

type Op string

const (
	Add      Op = "+"
	Subtract    = "-"
	Multiply    = "*"
	Divide      = "/"
)

var (
	ErrOperandType      = errors.New("invalid operandType")
	ErrRecursiveFormula = errors.New("recursive formula question")
)

func eval(
	expression Expr,
	questions map[QID]float64,
) (float64, error) {
	switch expression := expression.(type) {
	case OpExpr:
		return evalOp(expression, questions)
	case QID:
		return questions[expression], nil
	default:
		return 0, ErrOperandType
	}
}

func evalOp(
	op OpExpr,
	questions map[QID]float64,
) (float64, error) {

	leftVal, err := eval(op.left, questions)
	if err != nil {
		return 0, err
	}
	rightVal, err := eval(op.right, questions)
	if err != nil {
		return 0, err
	}

	switch op.op {
	case Add:
		return leftVal + rightVal, nil
	case Subtract:
		return leftVal - rightVal, nil
	case Multiply:
		return leftVal * rightVal, nil
	case Divide:
		if rightVal == 0.0 {
			return 0, nil
		}

		return leftVal / rightVal, nil
	default:
		return 0, nil
	}
}

func expand(
	list map[QID]OpDef,
) (map[QID]Expr, error) {

	result := map[QID]Expr{}
	for question := range list {
		expr, err := expandOpDef(question, list, map[QID]struct{}{})
		if err != nil {
			return map[QID]Expr{}, err
		}
		result[question] = expr
	}

	return result, nil
}

func expandOpDef(
	id QID,
	list map[QID]OpDef,
	visited map[QID]struct{},
) (Expr, error) {

	found, exists := list[id]
	if !exists {
		return id, nil
	}

	if _, alreadyVisited := visited[id]; alreadyVisited {
		return nil, errors.Wrapf(ErrRecursiveFormula, string(id))
	}

	var err error
	visited[id] = struct{}{}
	result := OpExpr{}
	result.op = found.op
	result.left, err = expandOpDef(found.left, list, visited)
	if err != nil {
		return nil, err
	}

	result.right, err = expandOpDef(found.right, list, visited)
	if err != nil {
		return nil, err
	}

	return result, nil
}
