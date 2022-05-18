// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
//

package model

const (
	// RingStateStable is a ring in a stable state and undergoing no changes.
	RingStateStable = "stable"
	// RingStateCreationRequested is a ring in the process of being created.
	RingStateCreationRequested = "creation-requested"
	// RingStateCreationFailed is a ring that failed creation.
	RingStateCreationFailed = "creation-failed"
	// RingStateReleaseRequested is a ring in the process of a release.
	RingStateReleaseRequested = "release-requested"
	// RingStateReleaseFailed is a ring that failed the release.
	RingStateReleaseFailed = "release-failed"
	// RingStateReleaseComplete is a ring that the release is complete.
	RingStateReleaseComplete = "release-complete"
	// RingStateSoakingRequested is a ring that is undergoing soak period.
	RingStateSoakingRequested = "soaking-requested"
	// RingStateSoakingFailed is a ring that is undergoing soak period.
	RingStateSoakingFailed = "soaking-failed"
	// RingStateReleaseRollbackRequested is a ring that the release is rolling back.
	RingStateReleaseRollbackRequested = "release-rollback-requested"
	// RingStateReleaseRollbackComplete is a ring that the release rollback is complete.
	RingStateReleaseRollbackComplete = "release-rollback-complete"
	// RingStateReleaseRollbackFailed is a ring that the release rollback has failed.
	RingStateReleaseRollbackFailed = "release-rollback-failed"
	// RingStateDeletionRequested is a ring in the process of being deleted.
	RingStateDeletionRequested = "deletion-requested"
	// RingStateDeletionFailed is a ring that failed deletion.
	RingStateDeletionFailed = "deletion-failed"
	// RingStateDeleted is a ring that has been deleted
	RingStateDeleted = "deleted"
)

// AllRingStates is a list of all states a ring can be in.
// Warning:
// When creating a new ring state, it must be added to this list.
var AllRingStates = []string{
	RingStateStable,
	RingStateCreationRequested,
	RingStateCreationFailed,
	RingStateReleaseRequested,
	RingStateReleaseFailed,
	RingStateReleaseComplete,
	RingStateSoakingRequested,
	RingStateSoakingFailed,
	RingStateReleaseRollbackRequested,
	RingStateReleaseRollbackFailed,
	RingStateReleaseRollbackComplete,
	RingStateDeletionRequested,
	RingStateDeletionFailed,
	RingStateDeleted,
}

// AllRingStatesPendingWork is a list of all ring states that the supervisor
// will attempt to transition towards stable on the next "tick".
// Warning:
// When creating a new ring state, it must be added to this list if the elrond
// supervisor should perform some action on its next work cycle.
var AllRingStatesPendingWork = []string{
	RingStateCreationRequested,
	RingStateReleaseRequested,
	RingStateSoakingRequested,
	RingStateReleaseRollbackRequested,
	RingStateDeletionRequested,
}

// AllRingRequestStates is a list of all states that a ring can be put in
// via the API.
// Warning:
// When creating a new ring state, it must be added to this list if an API
// endpoint should put the ring in this state.
var AllRingRequestStates = []string{
	RingStateCreationRequested,
	RingStateReleaseRequested,
	RingStateSoakingRequested,
	RingStateReleaseRollbackRequested,
	RingStateDeletionRequested,
}

// ValidTransitionState returns whether a ring can be transitioned into the
// new state or not based on its current state.
func (c *Ring) ValidTransitionState(newState string) bool {
	switch newState {
	case RingStateCreationRequested:
		return validTransitionToRingStateCreationRequested(c.State)
	case RingStateReleaseRequested:
		return validTransitionToRingStateReleaseRequested(c.State)
	case RingStateDeletionRequested:
		return validTransitionToRingStateDeletionRequested(c.State)
	case RingStateSoakingRequested:
		return validTransitionToRingStateSoakingRequested(c.State)
	case RingStateReleaseRollbackRequested:
		return validTransitionToRingStateRollbackRequested(c.State)
	}

	return false
}

func validTransitionToRingStateCreationRequested(currentState string) bool {
	switch currentState {
	case RingStateCreationRequested,
		RingStateCreationFailed:
		return true
	}

	return false
}

func validTransitionToRingStateReleaseRequested(currentState string) bool {
	switch currentState {
	case RingStateStable,
		RingStateReleaseRequested,
		RingStateReleaseFailed:
		return true
	}

	return false
}

func validTransitionToRingStateDeletionRequested(currentState string) bool {
	switch currentState {
	case RingStateStable,
		RingStateCreationRequested,
		RingStateCreationFailed,
		RingStateReleaseFailed,
		RingStateDeletionRequested,
		RingStateDeletionFailed:
		return true
	}

	return false
}

func validTransitionToRingStateSoakingRequested(currentState string) bool {
	switch currentState {
	case RingStateReleaseComplete,
		RingStateSoakingRequested,
		RingStateSoakingFailed:
		return true
	}

	return false
}

func validTransitionToRingStateRollbackRequested(currentState string) bool {
	switch currentState {
	case RingStateSoakingFailed,
		RingStateReleaseFailed,
		RingStateReleaseRollbackFailed:
		return true
	}

	return false
}

// RingStateReport is a report of all ring requests states.
type RingStateReport []StateReportEntry

// StateReportEntry is a report entry of a given request state.
type StateReportEntry struct {
	RequestedState string
	ValidStates    StateList
	InvalidStates  StateList
}

// StateList is a list of states
type StateList []string

// Count provides the number of states in a StateList.
func (sl *StateList) Count() int {
	return len(*sl)
}

// GetRingRequestStateReport returns a RingStateReport based on the current
// model of ring states.
func GetRingRequestStateReport() RingStateReport {
	report := RingStateReport{}

	for _, requestState := range AllRingRequestStates {
		entry := StateReportEntry{
			RequestedState: requestState,
		}

		for _, newState := range AllRingStates {
			c := Ring{State: newState}
			if c.ValidTransitionState(requestState) {
				entry.ValidStates = append(entry.ValidStates, newState)
			} else {
				entry.InvalidStates = append(entry.InvalidStates, newState)
			}
		}

		report = append(report, entry)
	}

	return report
}
