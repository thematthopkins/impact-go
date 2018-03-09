package contingency

import "github.com/pkg/errors"

//QuestionSfid Salesforce Id
type QuestionSfid string

//AnswerValueSfid Salesforce Id
type AnswerValueSfid string

//Response supplied for a specific question
type Response struct {
	ValuePercentage int
	Answers         map[AnswerValueSfid]struct{}
}

//Responses is the set of all responses for the assessment
type Responses map[QuestionSfid]Response

//AnswerDependencies question to answer value
type AnswerDependencies map[QuestionSfid]AnswerValueSfid

type questionDependencies struct {
	DisablingAnswerValues AnswerDependencies
	EnablingAnswerValues  AnswerDependencies
	EnablingQuestions     map[QuestionSfid]struct{}
}

type goalQuestionContingencies struct {
	disablingAnswerValues AnswerDependencies
	enablingAnswerValues  AnswerDependencies
}

type answerValueType int

const (
	answerValueEnables answerValueType = iota
	answerValueDisables
)

//ErrCircularContingencies when there is circular contingencies
var ErrCircularContingencies = errors.New("circular contingencies")

func addDescendantContingencies(
	q QuestionSfid,
	result *questionDependencies,
	deps map[QuestionSfid]questionDependencies,
) error {
	questionDeps, ok := deps[q]
	if !ok {
		return nil
	}
	for disablingQuestion, disablingAnswerValue := range questionDeps.DisablingAnswerValues {
		_, alreadyHasDep := result.DisablingAnswerValues[disablingQuestion]

		if alreadyHasDep {
			return errors.Wrapf(ErrCircularContingencies, string(q))
		}

		result.DisablingAnswerValues[disablingQuestion] = disablingAnswerValue

		err := addDescendantContingencies(disablingQuestion, result, deps)
		if err != nil {
			return err
		}
	}

	for enablingQuestion, enablingAnswerValue := range questionDeps.EnablingAnswerValues {
		_, alreadyHasDep := result.EnablingAnswerValues[enablingQuestion]

		if alreadyHasDep {
			return errors.Wrapf(ErrCircularContingencies, string(q))
		}

		result.EnablingAnswerValues[enablingQuestion] = enablingAnswerValue

		err := addDescendantContingencies(enablingQuestion, result, deps)
		if err != nil {
			return err
		}
	}

	for enablingQuestion := range questionDeps.EnablingQuestions {
		_, alreadyHasDep := result.EnablingQuestions[enablingQuestion]
		if alreadyHasDep {
			return errors.Wrapf(ErrCircularContingencies, string(q))
		}

		result.EnablingQuestions[enablingQuestion] = struct{}{}

		err := addDescendantContingencies(enablingQuestion, result, deps)
		if err != nil {
			return err
		}
	}

	return nil
}

func expand(
	contingencies map[QuestionSfid]questionDependencies,
) (map[QuestionSfid]questionDependencies, error) {

	result := map[QuestionSfid]questionDependencies{}
	for q := range contingencies {
		questionDeps := questionDependencies{
			DisablingAnswerValues: AnswerDependencies{},
			EnablingAnswerValues:  AnswerDependencies{},
			EnablingQuestions:     map[QuestionSfid]struct{}{},
		}
		err := addDescendantContingencies(q, &questionDeps, contingencies)
		if err != nil {
			return map[QuestionSfid]questionDependencies{}, err
		}
		result[q] = questionDeps
	}

	return result, nil
}

func fromGoal(
	masterQuestion *QuestionSfid,
	answerValue *AnswerValueSfid,
	goalQuestions []QuestionSfid,
	disabledByAnswerValue answerValueType,
) map[QuestionSfid]questionDependencies {

	result := map[QuestionSfid]questionDependencies{}

	for _, q := range goalQuestions {
		questionDeps := questionDependencies{
			DisablingAnswerValues: AnswerDependencies{},
			EnablingAnswerValues:  AnswerDependencies{},
			EnablingQuestions:     map[QuestionSfid]struct{}{},
		}
		if disabledByAnswerValue == answerValueDisables && masterQuestion != nil {
			questionDeps.DisablingAnswerValues[*masterQuestion] = *answerValue
		}

		if disabledByAnswerValue == answerValueEnables && masterQuestion != nil {
			questionDeps.EnablingAnswerValues[*masterQuestion] = *answerValue
		}
		result[q] = questionDeps
	}

	return result
}

//Enable determines if a question should be hidden
func Enable(
	responses Responses,
	disablingAnswerValues AnswerDependencies,
	enablingAnswerValues AnswerDependencies,
	enablingQuestions map[QuestionSfid]struct{}) bool {

	return enabledByAnswerValue(responses, enablingAnswerValues) &&
		enabledByQuestionValuePercentage(responses, enablingQuestions) &&
		!disabledByAnswerValues(responses, disablingAnswerValues)
}

//disabledByAnswerValues if a question is hidden by answer values
func disabledByAnswerValues(
	responses map[QuestionSfid]Response,
	disablingAnswerValues map[QuestionSfid]AnswerValueSfid) bool {

	for disablingQuestion, disablingAnswer := range disablingAnswerValues {
		if response, hasDisablingResponse := responses[disablingQuestion]; hasDisablingResponse {
			_, hasDisablingAnswer := response.Answers[disablingAnswer]
			if hasDisablingAnswer {
				return true
			}
		}
	}
	return false
}

func enabledByAnswerValue(
	responses map[QuestionSfid]Response,
	enablingAnswerValues map[QuestionSfid]AnswerValueSfid) bool {

	if len(enablingAnswerValues) == 0 {
		return true
	}

	for enablingQuestion, enablingAnswer := range enablingAnswerValues {
		if response, hasEnablingResponse := responses[enablingQuestion]; hasEnablingResponse {
			_, hasEnablingAnswer := response.Answers[enablingAnswer]
			if hasEnablingAnswer {
				return true
			}
		}
	}
	return false
}

func enabledByQuestionValuePercentage(
	responses map[QuestionSfid]Response,
	enablingQuestions map[QuestionSfid]struct{}) bool {

	if len(enablingQuestions) == 0 {
		return true
	}

	for enablingQuestion := range enablingQuestions {
		if response, hasEnablingResponse := responses[enablingQuestion]; hasEnablingResponse {
			if response.ValuePercentage >= 100 {
				return true
			}
		}
	}

	return false
}
