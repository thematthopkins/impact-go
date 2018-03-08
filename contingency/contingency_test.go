package contingency_test

import (
	"testing"

	"github.com/thematthopkins/impact-go/contingency"

	"github.com/stretchr/testify/assert"
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
