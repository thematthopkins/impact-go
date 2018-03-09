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
		map[QuestionSfid]struct{}{},
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
		map[QuestionSfid]struct{}{},
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
		map[QuestionSfid]struct{}{},
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
		map[QuestionSfid]struct{}{},
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
		map[QuestionSfid]struct{}{},
	)
	assert.True(t, enabled)
}

func TestEnablingQuestionUnmet(t *testing.T) {
	enabled := Enable(
		Responses{},
		AnswerDependencies{},
		AnswerDependencies{},
		map[QuestionSfid]struct{}{
			"masterQ":  struct{}{},
			"masterQ2": struct{}{},
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
		map[QuestionSfid]struct{}{
			"masterQ":  struct{}{},
			"masterQ2": struct{}{},
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
		map[QuestionSfid]struct{}{
			"masterQ":  struct{}{},
			"masterQ2": struct{}{},
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
		map[QuestionSfid]struct{}{},
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
		map[QuestionSfid]struct{}{
			"masterQ":  struct{}{},
			"masterQ2": struct{}{},
		},
	)
	assert.False(t, enabled)
}

func TestFromGoalDisableByAnswerValue(t *testing.T) {
	q1 := QuestionSfid("q1")
	a1 := AnswerValueSfid("a1")
	result := fromGoal(
		&q1,
		&a1,
		[]QuestionSfid{
			"goalq1",
			"goalq2",
		},
		answerValueDisables,
	)

	assert.Equal(t,
		map[QuestionSfid]questionDependencies{
			"goalq1": questionDependencies{
				DisablingAnswerValues: AnswerDependencies{
					"q1": "a1",
				},
				EnablingAnswerValues: AnswerDependencies{},
				EnablingQuestions:    map[QuestionSfid]struct{}{},
			},
			"goalq2": questionDependencies{
				DisablingAnswerValues: AnswerDependencies{
					"q1": "a1",
				},
				EnablingAnswerValues: AnswerDependencies{},
				EnablingQuestions:    map[QuestionSfid]struct{}{},
			},
		},
		result,
	)
}

func TestFromGoalEnabledByAnswerValue(t *testing.T) {
	q1 := QuestionSfid("q1")
	a1 := AnswerValueSfid("a1")
	result := fromGoal(
		&q1,
		&a1,
		[]QuestionSfid{
			"goalq1",
			"goalq2",
		},
		answerValueEnables,
	)

	assert.Equal(t,
		map[QuestionSfid]questionDependencies{
			"goalq1": questionDependencies{
				DisablingAnswerValues: AnswerDependencies{},
				EnablingAnswerValues: AnswerDependencies{
					"q1": "a1",
				},
				EnablingQuestions: map[QuestionSfid]struct{}{},
			},
			"goalq2": questionDependencies{
				DisablingAnswerValues: AnswerDependencies{},
				EnablingAnswerValues: AnswerDependencies{
					"q1": "a1",
				},
				EnablingQuestions: map[QuestionSfid]struct{}{},
			},
		},
		result,
	)
}

func TestFromGoalWithNoMasterQ(t *testing.T) {
	result := fromGoal(
		nil,
		nil,
		[]QuestionSfid{
			"goalq1",
			"goalq2",
		},
		answerValueEnables,
	)

	assert.Equal(t,
		map[QuestionSfid]questionDependencies{
			"goalq1": questionDependencies{
				DisablingAnswerValues: AnswerDependencies{},
				EnablingAnswerValues:  AnswerDependencies{},
				EnablingQuestions:     map[QuestionSfid]struct{}{},
			},
			"goalq2": questionDependencies{
				DisablingAnswerValues: AnswerDependencies{},
				EnablingAnswerValues:  AnswerDependencies{},
				EnablingQuestions:     map[QuestionSfid]struct{}{},
			},
		},
		result,
	)
}

func TestFromGoalWithNoMasterQ_Disabled(t *testing.T) {
	result := fromGoal(
		nil,
		nil,
		[]QuestionSfid{
			"goalq1",
			"goalq2",
		},
		answerValueDisables,
	)

	assert.Equal(t,
		map[QuestionSfid]questionDependencies{
			"goalq1": questionDependencies{
				DisablingAnswerValues: AnswerDependencies{},
				EnablingAnswerValues:  AnswerDependencies{},
				EnablingQuestions:     map[QuestionSfid]struct{}{},
			},
			"goalq2": questionDependencies{
				DisablingAnswerValues: AnswerDependencies{},
				EnablingAnswerValues:  AnswerDependencies{},
				EnablingQuestions:     map[QuestionSfid]struct{}{},
			},
		},
		result,
	)
}

func TestExpand_empty(t *testing.T) {
	result, err := expand(map[QuestionSfid]questionDependencies{})

	assert.NoError(t, err)
	assert.Equal(t,
		map[QuestionSfid]questionDependencies{},
		result,
	)
}

func TestExpand_noCascade(t *testing.T) {
	result, err := expand(map[QuestionSfid]questionDependencies{
		"q1": questionDependencies{
			DisablingAnswerValues: AnswerDependencies{
				"q0": "a1",
			},
		},
		"q2": questionDependencies{
			DisablingAnswerValues: AnswerDependencies{
				"q0": "a1",
			},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t,
		map[QuestionSfid]questionDependencies{
			"q1": questionDependencies{
				DisablingAnswerValues: AnswerDependencies{
					"q0": "a1",
				},
				EnablingAnswerValues: AnswerDependencies{},
				EnablingQuestions:    map[QuestionSfid]struct{}{},
			},
			"q2": questionDependencies{
				DisablingAnswerValues: AnswerDependencies{
					"q0": "a1",
				},
				EnablingAnswerValues: AnswerDependencies{},
				EnablingQuestions:    map[QuestionSfid]struct{}{},
			},
		},
		result,
	)
}

func TestExpand_disablingTwoLevels(t *testing.T) {
	result, err := expand(map[QuestionSfid]questionDependencies{
		"q1": questionDependencies{
			DisablingAnswerValues: AnswerDependencies{
				"q0": "aq0",
			},
		},
		"q2": questionDependencies{
			DisablingAnswerValues: AnswerDependencies{
				"q1": "aq1",
			},
		},
		"q3": questionDependencies{
			DisablingAnswerValues: AnswerDependencies{
				"q2": "aq2",
			},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t,
		map[QuestionSfid]questionDependencies{
			"q1": questionDependencies{
				DisablingAnswerValues: AnswerDependencies{
					"q0": "aq0",
				},
				EnablingAnswerValues: AnswerDependencies{},
				EnablingQuestions:    map[QuestionSfid]struct{}{},
			},
			"q2": questionDependencies{
				DisablingAnswerValues: AnswerDependencies{
					"q0": "aq0",
					"q1": "aq1",
				},
				EnablingAnswerValues: AnswerDependencies{},
				EnablingQuestions:    map[QuestionSfid]struct{}{},
			},
			"q3": questionDependencies{
				DisablingAnswerValues: AnswerDependencies{
					"q0": "aq0",
					"q1": "aq1",
					"q2": "aq2",
				},
				EnablingAnswerValues: AnswerDependencies{},
				EnablingQuestions:    map[QuestionSfid]struct{}{},
			},
		},
		result,
	)
}

func TestExpand_circularDisabled(t *testing.T) {
	_, err := expand(map[QuestionSfid]questionDependencies{
		"q1": questionDependencies{
			DisablingAnswerValues: AnswerDependencies{
				"q2": "aq2",
			},
		},
		"q2": questionDependencies{
			DisablingAnswerValues: AnswerDependencies{
				"q1": "aq1",
			},
		},
	})

	assert.Error(t, err)
}

func TestExpand_enablingTwoLevels(t *testing.T) {
	result, err := expand(map[QuestionSfid]questionDependencies{
		"q1": questionDependencies{
			EnablingAnswerValues: AnswerDependencies{
				"q0": "aq0",
			},
		},
		"q2": questionDependencies{
			EnablingAnswerValues: AnswerDependencies{
				"q1": "aq1",
			},
		},
		"q3": questionDependencies{
			EnablingAnswerValues: AnswerDependencies{
				"q2": "aq2",
			},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t,
		map[QuestionSfid]questionDependencies{
			"q1": questionDependencies{
				EnablingAnswerValues: AnswerDependencies{
					"q0": "aq0",
				},
				DisablingAnswerValues: AnswerDependencies{},
				EnablingQuestions:     map[QuestionSfid]struct{}{},
			},
			"q2": questionDependencies{
				EnablingAnswerValues: AnswerDependencies{
					"q0": "aq0",
					"q1": "aq1",
				},
				DisablingAnswerValues: AnswerDependencies{},
				EnablingQuestions:     map[QuestionSfid]struct{}{},
			},
			"q3": questionDependencies{
				EnablingAnswerValues: AnswerDependencies{
					"q0": "aq0",
					"q1": "aq1",
					"q2": "aq2",
				},
				DisablingAnswerValues: AnswerDependencies{},
				EnablingQuestions:     map[QuestionSfid]struct{}{},
			},
		},
		result,
	)
}

func TestExpand_circularEnabled(t *testing.T) {
	_, err := expand(map[QuestionSfid]questionDependencies{
		"q1": questionDependencies{
			EnablingAnswerValues: AnswerDependencies{
				"q2": "aq2",
			},
		},
		"q2": questionDependencies{
			EnablingAnswerValues: AnswerDependencies{
				"q1": "aq1",
			},
		},
	})

	assert.Error(t, err)
}

func TestExpand_enablingQuestionTwoLevels(t *testing.T) {
	result, err := expand(map[QuestionSfid]questionDependencies{
		"q1": questionDependencies{
			EnablingQuestions: map[QuestionSfid]struct{}{
				"q0": struct{}{},
			},
		},
		"q2": questionDependencies{
			EnablingQuestions: map[QuestionSfid]struct{}{
				"q1": struct{}{},
			},
		},
		"q3": questionDependencies{
			EnablingQuestions: map[QuestionSfid]struct{}{
				"q2": struct{}{},
			},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t,
		map[QuestionSfid]questionDependencies{
			"q1": questionDependencies{
				EnablingAnswerValues:  AnswerDependencies{},
				DisablingAnswerValues: AnswerDependencies{},
				EnablingQuestions: map[QuestionSfid]struct{}{
					"q0": struct{}{},
				},
			},
			"q2": questionDependencies{
				EnablingAnswerValues:  AnswerDependencies{},
				DisablingAnswerValues: AnswerDependencies{},
				EnablingQuestions: map[QuestionSfid]struct{}{
					"q0": struct{}{},
					"q1": struct{}{},
				},
			},
			"q3": questionDependencies{
				EnablingAnswerValues:  AnswerDependencies{},
				DisablingAnswerValues: AnswerDependencies{},
				EnablingQuestions: map[QuestionSfid]struct{}{
					"q0": struct{}{},
					"q1": struct{}{},
					"q2": struct{}{},
				},
			},
		},
		result,
	)
}

func TestExpand_enablingQuestionCircular(t *testing.T) {
	_, err := expand(map[QuestionSfid]questionDependencies{
		"q1": questionDependencies{
			EnablingQuestions: map[QuestionSfid]struct{}{
				"q2": struct{}{},
			},
		},
		"q2": questionDependencies{
			EnablingQuestions: map[QuestionSfid]struct{}{
				"q1": struct{}{},
			},
		},
	})

	assert.Error(t, err)
}
