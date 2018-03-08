package contingency

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

//Enable determines if a question should be hidden
func Enable(
	responses Responses,
	disablingAnswerValues AnswerDependencies,
	enablingAnswerValues AnswerDependencies,
	enablingQuestions []QuestionSfid) bool {

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
	enablingQuestions []QuestionSfid) bool {

	if len(enablingQuestions) == 0 {
		return true
	}

	for _, enablingQuestion := range enablingQuestions {
		if response, hasEnablingResponse := responses[enablingQuestion]; hasEnablingResponse {
			if response.ValuePercentage >= 100 {
				return true
			}
		}
	}

	return false
}
