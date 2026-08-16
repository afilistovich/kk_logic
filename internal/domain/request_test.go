package domain_test

import (
	"testing"

	"github.com/afilistovich/kk_logic/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestRequest_CanTransitionTo(t *testing.T) {
	cases := []struct {
		name string
		from domain.RequestStatus
		to   domain.RequestStatus
		want bool
	}{
		{"new to in_progress is allowed", domain.StatusNew, domain.StatusInProgress, true},
		{"new to done is not allowed directly", domain.StatusNew, domain.StatusDone, false},
		{"in_progress to done is allowed", domain.StatusInProgress, domain.StatusDone, true},
		{"in_progress to rejected is allowed", domain.StatusInProgress, domain.StatusRejected, true},
		{"done is a terminal state", domain.StatusDone, domain.StatusInProgress, false},
		{"cancelled is a terminal state", domain.StatusCancelled, domain.StatusInProgress, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &domain.Request{Status: tc.from}
			assert.Equal(t, tc.want, r.CanTransitionTo(tc.to))
		})
	}
}

func TestRequest_TransitionTo(t *testing.T) {
	r := &domain.Request{Status: domain.StatusNew}

	err := r.TransitionTo(domain.StatusInProgress)
	assert.NoError(t, err)
	assert.Equal(t, domain.StatusInProgress, r.Status)

	err = r.TransitionTo(domain.StatusNew) // недопустимый переход назад
	assert.ErrorIs(t, err, domain.ErrInvalidStatusTransition)
}
