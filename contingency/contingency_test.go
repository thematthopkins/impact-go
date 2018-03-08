package contingency

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoContingency(t *testing.T) {
	enabled := Enable(
		Responses{},
		AnswerDependencies{},
		AnswerDependencies{},
		[]QuestionSfid{},
	)

	assert.True(t, enabled)
}

func TestDisablingAnswerValueNotMatched(t *testing.T) {
	enabled := Enable(
		Responses{
			"masterQ": Response{
				Answers: map[AnswerValueSfid]struct{}{
					"nonMatching": struct{}{},
				},
			},
		},
		AnswerDependencies{
			"masterQ": "masterV",
		},
		AnswerDependencies{},
		[]QuestionSfid{},
	)
	assert.True(t, enabled)
}

func TestDisablingAnswerValueMatched(t *testing.T) {
	enabled := Enable(
		Responses{
			"masterQ": Response{
				Answers: map[AnswerValueSfid]struct{}{
					"masterV": struct{}{},
				},
			},
		},
		AnswerDependencies{
			"masterQ": "masterV",
		},
		AnswerDependencies{},
		[]QuestionSfid{},
	)
	assert.False(t, enabled)
}

func TestEnablingAnswerValueNotMatched(t *testing.T) {
	enabled := Enable(
		Responses{},
		AnswerDependencies{},
		AnswerDependencies{
			"masterQ": "masterV",
		},
		[]QuestionSfid{},
	)
	assert.False(t, enabled)
}

func TestEnablingAnswerValueMatched(t *testing.T) {
	enabled := Enable(
		Responses{
			"masterQ": Response{
				Answers: map[AnswerValueSfid]struct{}{
					"masterV": struct{}{},
				},
			},
		},
		AnswerDependencies{},
		AnswerDependencies{
			"masterQ": "masterV",
		},
		[]QuestionSfid{},
	)
	assert.True(t, enabled)
}

func TestEnablingQuestionUnmet(t *testing.T) {
	enabled := Enable(
		Responses{},
		AnswerDependencies{},
		AnswerDependencies{},
		[]QuestionSfid{
			"masterQ",
			"masterQ2",
		},
	)
	assert.False(t, enabled)
}

func TestEnablingQuestionUnmetPercentage(t *testing.T) {
	enabled := Enable(
		Responses{
			"masterQ": Response{
				ValuePercentage: 80,
			},
		},
		AnswerDependencies{},
		AnswerDependencies{},
		[]QuestionSfid{
			"masterQ",
			"masterQ2",
		},
	)
	assert.False(t, enabled)
}

func TestEnablingQuestionMetPercentage(t *testing.T) {
	enabled := Enable(
		Responses{
			"masterQ": Response{
				ValuePercentage: 100,
			},
		},
		AnswerDependencies{},
		AnswerDependencies{},
		[]QuestionSfid{
			"masterQ",
			"masterQ2",
		},
	)
	assert.True(t, enabled)
}

func TestDisableOverridesEnable(t *testing.T) {
	enabled := Enable(
		Responses{
			"masterQ": Response{
				Answers: map[AnswerValueSfid]struct{}{
					"masterV": struct{}{},
				},
			},
		},
		AnswerDependencies{
			"masterQ": "masterV",
		},
		AnswerDependencies{
			"masterQ": "masterV",
		},
		[]QuestionSfid{},
	)
	assert.False(t, enabled)
}

func TestDisableOverridesEnableByValuePercentage(t *testing.T) {
	enabled := Enable(
		Responses{
			"masterQ": Response{
				ValuePercentage: 100,
				Answers: map[AnswerValueSfid]struct{}{
					"masterV": struct{}{},
				},
			},
		},
		AnswerDependencies{
			"masterQ": "masterV",
		},
		AnswerDependencies{},
		[]QuestionSfid{
			"masterQ",
			"masterQ2",
		},
	)
	assert.False(t, enabled)
}

func TestFromGoalDisableByAnswerValue(t *testing.T) {
	q1 := QuestionSfid("q1")
	a1 := AnswerValueSfid("a1")
	result := FromGoal(
		&q1,
		&a1,
		[]QuestionSfid{
			"goalq1",
			"goalq2",
		},
		AnswerValueDisables,
	)

	assert.Equal(t,
		map[QuestionSfid]QuestionDependencies{
			"goalq1": QuestionDependencies{
				DisablingAnswerValues: AnswerDependencies{
					"q1": "a1",
				},
				EnablingAnswerValues: AnswerDependencies{},
				EnablingQuestions:    []QuestionSfid{},
			},
			"goalq2": QuestionDependencies{
				DisablingAnswerValues: AnswerDependencies{
					"q1": "a1",
				},
				EnablingAnswerValues: AnswerDependencies{},
				EnablingQuestions:    []QuestionSfid{},
			},
		},
		result,
	)
}

func TestFromGoalEnabledByAnswerValue(t *testing.T) {
	q1 := QuestionSfid("q1")
	a1 := AnswerValueSfid("a1")
	result := FromGoal(
		&q1,
		&a1,
		[]QuestionSfid{
			"goalq1",
			"goalq2",
		},
		AnswerValueEnables,
	)

	assert.Equal(t,
		map[QuestionSfid]QuestionDependencies{
			"goalq1": QuestionDependencies{
				DisablingAnswerValues: AnswerDependencies{},
				EnablingAnswerValues: AnswerDependencies{
					"q1": "a1",
				},
				EnablingQuestions: []QuestionSfid{},
			},
			"goalq2": QuestionDependencies{
				DisablingAnswerValues: AnswerDependencies{},
				EnablingAnswerValues: AnswerDependencies{
					"q1": "a1",
				},
				EnablingQuestions: []QuestionSfid{},
			},
		},
		result,
	)
}

func TestFromGoalWithNoMasterQ(t *testing.T) {
	result := FromGoal(
		nil,
		nil,
		[]QuestionSfid{
			"goalq1",
			"goalq2",
		},
		AnswerValueEnables,
	)

	assert.Equal(t,
		map[QuestionSfid]QuestionDependencies{
			"goalq1": QuestionDependencies{
				DisablingAnswerValues: AnswerDependencies{},
				EnablingAnswerValues:  AnswerDependencies{},
				EnablingQuestions:     []QuestionSfid{},
			},
			"goalq2": QuestionDependencies{
				DisablingAnswerValues: AnswerDependencies{},
				EnablingAnswerValues:  AnswerDependencies{},
				EnablingQuestions:     []QuestionSfid{},
			},
		},
		result,
	)

}

func TestFromGoalWithNoMasterQ_Disabled(t *testing.T) {
	result := FromGoal(
		nil,
		nil,
		[]QuestionSfid{
			"goalq1",
			"goalq2",
		},
		AnswerValueDisables,
	)

	assert.Equal(t,
		map[QuestionSfid]QuestionDependencies{
			"goalq1": QuestionDependencies{
				DisablingAnswerValues: AnswerDependencies{},
				EnablingAnswerValues:  AnswerDependencies{},
				EnablingQuestions:     []QuestionSfid{},
			},
			"goalq2": QuestionDependencies{
				DisablingAnswerValues: AnswerDependencies{},
				EnablingAnswerValues:  AnswerDependencies{},
				EnablingQuestions:     []QuestionSfid{},
			},
		},
		result,
	)

}
