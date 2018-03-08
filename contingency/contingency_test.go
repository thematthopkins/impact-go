package contingency

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thematthopkins/impact-go/contingency"
)

func TestNoContingency(t *testing.T) {
	enabled := contingency.Enable(
		contingency.Responses{},
		contingency.AnswerDependencies{},
		contingency.AnswerDependencies{},
		[]contingency.QuestionSfid{},
	)

	assert.True(t, enabled)
}

func TestDisablingAnswerValueNotMatched(t *testing.T) {
	enabled := contingency.Enable(
		contingency.Responses{
			"masterQ": contingency.Response{
				Answers: map[contingency.AnswerValueSfid]struct{}{
					"nonMatching": struct{}{},
				},
			},
		},
		contingency.AnswerDependencies{
			"masterQ": "masterV",
		},
		contingency.AnswerDependencies{},
		[]contingency.QuestionSfid{},
	)
	assert.True(t, enabled)
}

func TestDisablingAnswerValueMatched(t *testing.T) {
	enabled := contingency.Enable(
		contingency.Responses{
			"masterQ": contingency.Response{
				Answers: map[contingency.AnswerValueSfid]struct{}{
					"masterV": struct{}{},
				},
			},
		},
		contingency.AnswerDependencies{
			"masterQ": "masterV",
		},
		contingency.AnswerDependencies{},
		[]contingency.QuestionSfid{},
	)
	assert.False(t, enabled)
}

func TestEnablingAnswerValueNotMatched(t *testing.T) {
	enabled := contingency.Enable(
		contingency.Responses{},
		contingency.AnswerDependencies{},
		contingency.AnswerDependencies{
			"masterQ": "masterV",
		},
		[]contingency.QuestionSfid{},
	)
	assert.False(t, enabled)
}

func TestEnablingAnswerValueMatched(t *testing.T) {
	enabled := contingency.Enable(
		contingency.Responses{
			"masterQ": contingency.Response{
				Answers: map[contingency.AnswerValueSfid]struct{}{
					"masterV": struct{}{},
				},
			},
		},
		contingency.AnswerDependencies{},
		contingency.AnswerDependencies{
			"masterQ": "masterV",
		},
		[]contingency.QuestionSfid{},
	)
	assert.True(t, enabled)
}

func TestEnablingQuestionUnmet(t *testing.T) {
	enabled := contingency.Enable(
		contingency.Responses{},
		contingency.AnswerDependencies{},
		contingency.AnswerDependencies{},
		[]contingency.QuestionSfid{
			"masterQ",
			"masterQ2",
		},
	)
	assert.False(t, enabled)
}

func TestEnablingQuestionUnmetPercentage(t *testing.T) {
	enabled := contingency.Enable(
		contingency.Responses{
			"masterQ": contingency.Response{
				ValuePercentage: 80,
			},
		},
		contingency.AnswerDependencies{},
		contingency.AnswerDependencies{},
		[]contingency.QuestionSfid{
			"masterQ",
			"masterQ2",
		},
	)
	assert.False(t, enabled)
}

func TestEnablingQuestionMetPercentage(t *testing.T) {
	enabled := contingency.Enable(
		contingency.Responses{
			"masterQ": contingency.Response{
				ValuePercentage: 100,
			},
		},
		contingency.AnswerDependencies{},
		contingency.AnswerDependencies{},
		[]contingency.QuestionSfid{
			"masterQ",
			"masterQ2",
		},
	)
	assert.True(t, enabled)
}

func TestDisableOverridesEnable(t *testing.T) {
	enabled := contingency.Enable(
		contingency.Responses{
			"masterQ": contingency.Response{
				Answers: map[contingency.AnswerValueSfid]struct{}{
					"masterV": struct{}{},
				},
			},
		},
		contingency.AnswerDependencies{
			"masterQ": "masterV",
		},
		contingency.AnswerDependencies{
			"masterQ": "masterV",
		},
		[]contingency.QuestionSfid{},
	)
	assert.False(t, enabled)
}

func TestDisableOverridesEnableByValuePercentage(t *testing.T) {
	enabled := contingency.Enable(
		contingency.Responses{
			"masterQ": contingency.Response{
				ValuePercentage: 100,
				Answers: map[contingency.AnswerValueSfid]struct{}{
					"masterV": struct{}{},
				},
			},
		},
		contingency.AnswerDependencies{
			"masterQ": "masterV",
		},
		contingency.AnswerDependencies{},
		[]contingency.QuestionSfid{
			"masterQ",
			"masterQ2",
		},
	)
	assert.False(t, enabled)
}

func TestFromGoalDisableByAnswerValue(t *testing.T) {
	q1 := contingency.QuestionSfid("q1")
	a1 := contingency.AnswerValueSfid("a1")
	result := contingency.FromGoal(
		&q1,
		&a1,
		[]contingency.QuestionSfid{
			"goalq1",
			"goalq2",
		},
		contingency.AnswerValueDisables,
	)

	assert.Equal(t,
		map[contingency.QuestionSfid]contingency.QuestionDependencies{
			"goalq1": contingency.QuestionDependencies{
				DisablingAnswerValues: contingency.AnswerDependencies{
					"q1": "a1",
				},
				EnablingAnswerValues: contingency.AnswerDependencies{},
				EnablingQuestions:    []contingency.QuestionSfid{},
			},
			"goalq2": contingency.QuestionDependencies{
				DisablingAnswerValues: contingency.AnswerDependencies{
					"q1": "a1",
				},
				EnablingAnswerValues: contingency.AnswerDependencies{},
				EnablingQuestions:    []contingency.QuestionSfid{},
			},
		},
		result,
	)
}

func TestFromGoalEnabledByAnswerValue(t *testing.T) {
	q1 := contingency.QuestionSfid("q1")
	a1 := contingency.AnswerValueSfid("a1")
	result := contingency.FromGoal(
		&q1,
		&a1,
		[]contingency.QuestionSfid{
			"goalq1",
			"goalq2",
		},
		contingency.AnswerValueEnables,
	)

	assert.Equal(t,
		map[contingency.QuestionSfid]contingency.QuestionDependencies{
			"goalq1": contingency.QuestionDependencies{
				DisablingAnswerValues: contingency.AnswerDependencies{},
				EnablingAnswerValues: contingency.AnswerDependencies{
					"q1": "a1",
				},
				EnablingQuestions: []contingency.QuestionSfid{},
			},
			"goalq2": contingency.QuestionDependencies{
				DisablingAnswerValues: contingency.AnswerDependencies{},
				EnablingAnswerValues: contingency.AnswerDependencies{
					"q1": "a1",
				},
				EnablingQuestions: []contingency.QuestionSfid{},
			},
		},
		result,
	)
}

func TestFromGoalWithNoMasterQ(t *testing.T) {
	result := contingency.FromGoal(
		nil,
		nil,
		[]contingency.QuestionSfid{
			"goalq1",
			"goalq2",
		},
		contingency.AnswerValueEnables,
	)

	assert.Equal(t,
		map[contingency.QuestionSfid]contingency.QuestionDependencies{
			"goalq1": contingency.QuestionDependencies{
				DisablingAnswerValues: contingency.AnswerDependencies{},
				EnablingAnswerValues:  contingency.AnswerDependencies{},
				EnablingQuestions:     []contingency.QuestionSfid{},
			},
			"goalq2": contingency.QuestionDependencies{
				DisablingAnswerValues: contingency.AnswerDependencies{},
				EnablingAnswerValues:  contingency.AnswerDependencies{},
				EnablingQuestions:     []contingency.QuestionSfid{},
			},
		},
		result,
	)

}

func TestFromGoalWithNoMasterQ_Disabled(t *testing.T) {
	result := contingency.FromGoal(
		nil,
		nil,
		[]contingency.QuestionSfid{
			"goalq1",
			"goalq2",
		},
		contingency.AnswerValueDisables,
	)

	assert.Equal(t,
		map[contingency.QuestionSfid]contingency.QuestionDependencies{
			"goalq1": contingency.QuestionDependencies{
				DisablingAnswerValues: contingency.AnswerDependencies{},
				EnablingAnswerValues:  contingency.AnswerDependencies{},
				EnablingQuestions:     []contingency.QuestionSfid{},
			},
			"goalq2": contingency.QuestionDependencies{
				DisablingAnswerValues: contingency.AnswerDependencies{},
				EnablingAnswerValues:  contingency.AnswerDependencies{},
				EnablingQuestions:     []contingency.QuestionSfid{},
			},
		},
		result,
	)

}
